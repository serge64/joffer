var { resumeStore, positionStore } = [];
var groupCounter = 0;

const groupsMap = new Map();
const letterMap = new Map();

const server = {
    async profile() {
        const response = await fetch('/api/profile', {
            method: 'GET',
            mode: 'cors',
            credentials: 'same-origin'
        });
        return await response.json();
    },
    async sendProfile() {
        const response = await fetch('/api/profile', {
            method: 'POST',
            mode: 'cors',
            credentials: 'same-origin'
        });
        return await response.json();
    },
    async deleteProfile() {
        return await fetch('/api/profile', {
            method: 'DELETE',
            mode: 'cors',
            credentials: 'same-origin'
        });
    },
    async sendGroup(group = {}) {
        const response = await fetch('/api/groups', {
            method: 'POST',
            mode: 'cors',
            credentials: 'same-origin',
            header: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(group)
        });
        return await response.json();
    },
    async saveGroup(id = '', group = {}) {
        return await fetch(`/api/groups/${id}`, {
            method: 'PATCH',
            mode: 'cors',
            credentials: 'same-origin',
            header: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(group)
        });
    },
    async deleteGroup(id = '') {
        return await fetch(`/api/groups/${id}`, {
            method: 'DELETE',
            mode: 'cors',
            credentials: 'same-origin'
        });
    },
    async sendLetter(letter = {}) {
        const response = await fetch('/api/letters', {
            method: 'POST',
            mode: 'cors',
            credentials: 'same-origin',
            header: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(letter)
        });
        return await response.json();
    },
    async saveLetter(id = '', letter = {}) {
        return await fetch(`/api/letters/${id}`, {
            method: 'PATCH',
            mode: 'cors',
            credentials: 'same-origin',
            header: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(letter)
        });
    },
    async deleteLetter(id = '') {
        return await fetch(`/api/letters/${id}`, {
            method: 'DELETE',
            mode: 'cors',
            credentials: 'same-origin'
        });
    },
    async response(id = '') {
        return await fetch(`/api/groups/${id}`, {
            method: 'POST',
            mode: 'cors',
            credentials: 'same-origin'
        });
    }
};

// Запрос профиля после загрузки страницы.
$(document).ready(function() {
    server.profile()
        .then((data) => {
            appendProfile(data);
        })
        .catch((err) => {
            console.log(err);
        });
});

// Добавляет `profile` на страницу.
function appendProfile(data = {}) {
    if (Object.keys(data).length != 0) {   
        const { name, email, resumes, groups, letters } = data;

        resumeStore = resumes;

        $('.profile__area aside').find('p').text(`${name} <${email}>`);
        $('.groups').removeClass('hidden');

        groups.forEach(appendGroup);
        letters.forEach((e) => letterMap.set(e.name, e));

        appendResumes(resumeStore);
        appendLetters();
        toggleProfile();
    }
};

// Добавить/Удалить профиль. 
function profileSwitchMode() {
    toggleProfile();

    const active = $('#profile__button').hasClass('active');

    if (active) {
        server.sendProfile()
            .then((data) => {
                window.location.href = data;
            })
            .catch((err) => {
                console.log(err);
            });
    } else {
        const title = 'Добавьте профиль для дальнейшей работы';
        resumeStore = [];

        $('.profile__area aside').find('p').text(title);
        $('.groups').addClass('hidden');
        $('.groups').find('article').each(function() {
            if ($(this).attr('id') != 'add__card') {
                $(this).remove();
            }
        });

        server.deleteProfile()
            .catch((err) => {
                console.log(err);
            })
    }
};

// Переключение режима профиля.
function toggleProfile() {
    $('#profile__button').toggleClass('active');
    $('.title').toggleClass('active');
    $('nav').toggleClass('active');
}

// Добавляет группу на страницу.
function appendGroup(data = {}) {
    const tmpl = $('#article__template').prop('content');
    const node = $(tmpl).clone(true);
    const article = $(node).find('article');
    const items = $(article).children();

    $(article).attr('data-index', data.id);
    $(items[0]).find('input').val(data.name);

    const resume = $(items[1]).find('span');

    $(resume).text(data.resume);
    $(resume).removeClass('placeholder');

    const letter = $(items[2]).find('span');

    $(letter).text(data.letter);
    $(letter).removeClass('placeholder');

    const list = $(items[3]).find('ul');

    if (data.positions != null) {
        if (data.positions.length > 0) {
            data.positions.forEach(function(e) {
                $(list).append(composeSelected(e));
            });
            $(list).removeClass('hidden');
        }
    }

    const button = $(items[4]).find('.response__button');
    const span = $(button).find('span');

    if (data.count === 0) {
        $(button).addClass('disable');
        $(span).addClass('hidden');
    } else {
        $(button).removeClass('disable');
        $(span).removeClass('hidden');
        $(span).text(`(${data.count})`);
    }

    groupsMap.set(data.id.toString(), getGroupData(article));

    $('#add__card').before(article);
};

// Добавляет список резюме на страницу.
function appendResumes(data = []) {
    $('.groups').find('article .resumes').each(function() {
        const list = $(this).find('.option__list');
        $(list).find('li').remove();

        if (data.length > 4) {
            $(list).addClass('scroll');
        }

        data.forEach(function(e) {
            $(list).append(`<li>${e}</li>`)
        });      
    });
};

// Добавляет список сопроводительных
// писем на страницу.
function appendLetters() {
    $('article').find('.letters').each(function() {
        const list = $(this).find('.option__list');
        $(list).find('li').remove();

        if (letterMap.size > 4) {
            $(list).addClass('scroll');
        }

        letterMap.forEach(function(value) {
            $(list).append(composeLetter(value));
        });      
    });
};

// Возвращает `letter`.
function composeLetter(data = {}) {
    const tmpl = $('#letter__template').prop('content');
    const letter = $(tmpl).clone(true);
    $(letter).find('li').attr('data-index', data.id);
    $(letter).find('span').text(data.name);
    return letter;
};

// Возвращает выбранную позицию.
function composeSelected(data = '') {
    const tmpl = $('#selected__template').prop('content');
    const label = $(tmpl).clone(true);
    $(label).find('span').text(data);
    return label;
};

// Возвращает данные по группе.
function getGroupData(article = {}) {
    const fields = $(article).children();
    var data = {
        id: "",
        name: "",
        resume: "",
        letter: "",
        positions: []
    }

    data.id = $(article).attr('data-index');
    data.name = $(fields[0]).find('input').val().trim();
    data.resume = $(fields[1]).find('span').text().trim();
    data.letter = $(fields[2]).find('.select__trigger span').text().trim();

    $(fields[3]).find('li').each(function() {
        const value = $(this).find('span');
        data.positions.push($(value).text());
    });

    return data;
};

// Возвращает новую группу.
function composeEmptyGroup() {
    const tmpl = $('#article__template').prop('content');
    const node = $(tmpl).clone(true);
    const id = '_' + ++groupCounter;
    const article = $(node).find('article');
    $(article).attr('data-index', id);
    groupsMap.set(id, getGroupData(article));
    return article;
};

// Обработчик изменений в группе.
function differenceHandler(article = {}) {
    const isEqual = (obj1, obj2) => {
        const obj1Keys = Object.keys(obj1);
        const obj2Keys = Object.keys(obj2);

        if(obj1Keys.length !== obj2Keys.length) {
            return false;
        }

        for (let objKey of obj1Keys) {
            if (obj1[objKey] !== obj2[objKey]) {
                if(typeof obj1[objKey] == "object" && typeof obj2[objKey] == "object") {
                    if(!isEqual(obj1[objKey], obj2[objKey])) {
                        return false;
                    }
                } 
                else {
                    return false;
                }
            }
        }
    
        return true;
    };

    if ($(article).attr('id') == undefined) {
        const change = $(article).find('.change');
        const id = $(article).attr('data-index').toString();
        const data = getGroupData(article);

        if (groupsMap.get(id) != undefined) {
            if (!isEqual(data, groupsMap.get(id))) {
                if (valid(data)) {
                    if ($(change).hasClass('close')) {
                        slideUp(change);
                    }
                }
            } else {
                slideDown(change);
            }
        }
    }
}

// Скрывает кнопку сохранить изменения.
function slideDown(element = {}) {
    $(element).addClass('close');
    $(element)
        .stop(true, true)
        .slideUp({ duration: 'fast', easing: 'swing' })
        .fadeOut({ duration: 500, queue: false })
}

// Показывает кнопку сохранить изменения.
function slideUp(element = {}) {
    $(element).removeClass('close');
    $(element)
        .hide()
        .stop(true, true)
        .slideDown({ duration: 'fast', easing: 'swing' })
        .css('display', 'flex')
        .fadeIn({ duration: 500, queue: false });
}

// `Click` по кнопке `Добавить/Удалить профиль`.
$(document).ready(function() {
    $('#profile__button').click(profileSwitchMode);
});

// `Click` по выпадающему списку.
$(document).ready(function() {
    $('.groups').on('click', '.select__input', function() {
        const select = $(this).parent();
        const list = $(select).find('.option__list');

        if ($(list).children().length > 4) {
            $(list).addClass('scroll');
        } else {
            $(list).removeClass('scroll');
        }

        $(select).toggleClass('open');

        $('.select__container').each(function() {
            if (!$(this).is(select)) {
                $(this).removeClass('open');
            }
        });
    });
});

// `Click` по элементу в выпадающем списке.
$(document).ready(function() {
    $('.groups').on('click', '.option__list li', function(e) {
        const container = $(this).parent().parent().parent();
        const input = $(this).parent().parent().parent().find('.select__trigger span');
        const element = $(e.target).prop('tagName');

        if (element === 'LI' || element === 'SPAN') {
            $(input).removeClass('placeholder');
            $(input).text($(this).text());
            $(container).removeClass('open');
        }       
    });
});

// `Click` по кнопке создать письмо.
$(document).ready(function() {
    $('.groups').on('click', '.new__letter', function() {
        const container = $(this).parent().parent();
        $(container).find('.letter__title').removeClass('disable');
        showLetterEditor(container, false);
    });
});

// Показать редактор соповодительного письма.
function showLetterEditor(letter = {}, isEdit = false, data = {}) {
    $(letter).find('.option__list').addClass('hidden');
    $(letter).find('.letter__content').addClass('hidden');
    $(letter).find('.letter__container').removeClass('hidden');

    if (isEdit) {
        $(letter).find('.add__button').addClass('hidden');
        $(letter).find('.save__button').removeClass('hidden');
        $(letter).find('input').val(data.name);
        $(letter).find('textarea').val(data.body);
    } else {
        $(letter).find('.add__button').removeClass('hidden');
        $(letter).find('.save__button').addClass('hidden');
    }
};

// `Click` по кнопке отмена редактирования письма.
$(document).ready(function() {
    $('.groups').on('click', '.cancel__button', function() {
        const container = $(this).parent().parent().parent().parent();
        resetDropdown(container);
    });
});

// `Click` по кнопке добавить письмо.
$(document).ready(function() {
    const validLetter = (element) => {
        const input = $(element).find('input');
        const textarea = $(element).find('textarea');
    
        if ($(input).val().trim().length === 0) {
            $(input).addClass('error');
            return false;
        }
    
        if (letterMap.has($(input).val().trim())) {
            $(input).addClass('error');
            return false;
        }
    
        if ($(textarea).val().trim().length === 0) {
            $(textarea).addClass('error');
            return false;
        }
    
        return true;
    };

    $('.groups').on('click', '.add__button', function() {
        const container = $(this).parent().parent().parent().parent();
        const letter = $(this).parent().parent();

        if (validLetter(letter)) {
            var field = {
                name: $(letter).find('input').val().trim(),
                body: $(letter).find('textarea').val()
            }

            server.sendLetter(field)
                .then((data) => {
                    field.id = data;
                    letterMap.set(field.name, field);
                    appendLetters();
                })
                .catch((err) => {
                    console.log(err);
                });

            resetDropdown(container);
        }
    });
});

// `Click` по кнопке сохранить изменения в письме.
$(document).ready(function() {
    const validLetter = (element) => {
        const textarea = $(element).find('textarea');

        if ($(textarea).val().trim().length === 0) {
            $(textarea).addClass('error');
            return false;
        }

        return true;
    };

    $('.groups').on('click', '.save__button', function() {
        const container = $(this).parent().parent().parent().parent();
        const letter = $(this).parent().parent();

        if (validLetter(letter)) {
            var field = {
                name: $(letter).find('input').val().trim(),
                body: $(letter).find('textarea').val()
            };

            const id = $(letter).attr('data-index');

            server.saveLetter(id, field)
                .then((response) => {
                    if (response.ok) {
                        field.id = id;
                        letterMap.set(field.name, field);
                        appendLetters();
                    }
                })
                .catch((err) => {
                    console.log(err);
                });

            resetDropdown(container);
        }
    });
});

// Сбрасывает выпадающий список.
function resetDropdown(container = {}) {
    const list = $(container).find('.option__list');
    const button = $(container).find('.letter__content');
    const letter = $(container).find('.letter__container');
    const input = $(letter).find('input');
    const textarea = $(letter).find('textarea');

    $(list).removeClass('hidden');
    $(button).removeClass('hidden');
    $(letter).addClass('hidden');
    $(input).val('');
    $(input).removeClass('error');
    $(textarea).val('');
    $(textarea).removeClass('error');
};

// `Click` по кнопке редактировать письмо.
$(document).ready(function() {
    $('.groups').on('click', '.letter__action .edit__button', function() {
        const container = $(this).parent().parent().parent().parent();
        const letter = $(this).parent().parent();
        const name = $(letter).find('span').text().trim();
        const data = letterMap.get(name);
        const id = $(letter).attr('data-index');
        $(container).find('.letter__title').addClass('disable');
        $(container).find('.letter__container').attr('data-index', id);
        showLetterEditor(container, true, data);
    });
});

// `Click` по кнопке удалить письмо.
$(document).ready(function() {
    $('.groups').on('click', '.letter__action .remove__button', function() {
        const letter = $(this).parent().parent();
        const name = $(letter).find('span').text().trim();
        const id = $(letter).attr('data-index');

        $('.select__trigger span').each(function() {
            if ($(this).text().trim() === name) {
                $(this).text($(this).attr('data-value'));
                $(this).addClass('placeholder');
            }
        });

        server.deleteLetter(id)
            .then((response) => {
                if (response.ok) {
                    letterMap.delete(name);
                    appendLetters();
                }
            })
            .catch((err) => {
                console.log(err);
            });
    });
});

// `Click` вне выпадающего списка.
$(document).mouseup(function(e) {
    const select = $('.select__container');
    if (!$(select).is(e.target) && $(select).has(e.target).length === 0) {
        $('.select__container').each(function() {
            $(this).removeClass('open');
        });
    }
});

// `Click` по кнопке добавить позицию.
$(document).ready(function() {
    $('.groups').on('click', '.input__container--button button', function() {
        submitNewPosition(this);
    });
});

$(document).ready(function() {
    $('.groups').on('keypress', '.input__container--button input', function(e) {
        if (e.which === 13) {
            submitNewPosition(this);
        }
    });
});

function submitNewPosition(element = {}) {
    const list = $(element).parent().parent().find('.selected__list');
    const input = $(element).parent().find('input');
    const uniq = new Set();

    if ($(input).val().trim().length > 0) {
        $(list).find('li').each(function() {
            uniq.add($(element).text().trim());
        });

        if (!uniq.has($(input).val().trim())) {
            $(list).removeClass('hidden');

            let data = composeSelected($(input).val());
            $(list).append(data);

            if ($(list).children().length > 4) {
                $(list).addClass('scroll');
            }

            $(input).val('');
        }
    }
};

// `Click` по кнопке удалить позицию.
$(document).ready(function() {
    $('.groups').on('click', '.selected__list .remove__button', function() {
        const list = $(this).parent().parent();
        $(this).parent().remove();
        const count = $(list).children().length;

        if (count <= 4) {
            $(list).removeClass('scroll');
        }

        if (count === 0) {
            $(list).addClass('hidden');
        }
    });
});

// `Click` по кнопке удалить группу.
$(document).ready(function() {
    const remove = (article, id) => {
        $(article).remove();
        groupsMap.delete(id);
    };

    $('.groups').on('click', '.delete__button', function() {
        const article = $(this).parent().parent().parent();
        const id = $(article).attr('data-index').toString();

        if (id.startsWith('_')) {
            remove(article, id);
        } else {
            server.deleteGroup(id)
            .then((response) => {
                if (response.ok) {
                    remove(article, id);
                }
            })
            .catch((err) => {
                console.log(err);
            });
        }
    });
});

// `Click` по кнопке отклик.
$(document).ready(function() {
    $('.groups').on('click', 'article .response__button', function() {
        const article = $(this).parent().parent().parent();
        const button = $(article).find('.response__button');

        if (!$(button).hasClass('disable')) {
            const data = getGroupData(article);

            if (valid(data)) {
                server.response(data.id)
                    .catch((err) => {
                        console.log(err);
                    });
            }
        }
    });
});

// Проверка на валидность значений.
function valid(data = {}) {
    if (data.name.trim() === '') return false;
    if (data.resume.trim() === '' || data.resume === 'Выберите') return false;
    if (data.letter.trim() === '' || data.letter === 'Выберите') return false;
    return true;
}

// `Click` по кнопке создать группу.
$(document).ready(function() {
    const button = $('#add__card');

    $(button).click(function() {
        let data = composeEmptyGroup();
        $(button).before(data);
        appendResumes(resumeStore);
        appendLetters();
    });
});

// `Change` по группе.
$(document).ready(function() {
    $('.groups').on('click', 'article', function() {
        differenceHandler(this);   
    });
});

// `Keyup` по группе.
$(document).ready(function() {
    $('.groups').on('keyup', 'article', function() {
        differenceHandler(this);
    });
});

// `Click` по кнопке сохранить изменения.
$(document).ready(function() {
    $('.groups').on('click', '.change__question button', function() {
        const article = $(this).parent().parent().parent();
        const group = getGroupData(article);

        if (valid(group)) {
            const id = $(article).attr('data-index').toString();
            const change = $(this).parent().parent();

            var body = {
                name: group.name,
                resume: group.resume,
                letter: group.letter,
                positions: []
            };

            group.positions.forEach(function(value) {
                body.positions.push(value);
            });

            if (id.startsWith('_')) {
                server.sendGroup(body)
                    .then((data) => {
                        console.log(data)
                        groupsMap.delete(id);
                        groupsMap.set(data.id, group);
                        $(article).attr('data-index', data.id);
                        slideDown(change);
                    })
                    .catch((err) => {
                        console.log(err);
                    })
            } else {
                server.saveGroup(id, body)
                    .then((response) => {
                        if (response.ok) {
                            groupsMap.delete(id);
                            groupsMap.set(id, group);
                            slideDown(change);
                        }
                    })
                    .catch((err) => {
                        console.log(err);
                    });  
            }    
        }
    });
});
