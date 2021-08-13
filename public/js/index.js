const server = {
    async filter() {
        const response = await fetch('/api/filters', {
            method: 'GET',
            mode: 'cors',
            credentials: 'same-origin'
        });
        return await response.json();
    },
    async rows(data = {}) {
        const response = await fetch('/api/vacancies', {
            method: 'GET',
            mode: 'cors',
            credentials: 'same-origin',
            header: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        });
        return await response.json();
    },
    async search(id = '', data = {}) {
        const response = await fetch(`/api/search/${id}`, {
            method: 'GET',
            mode: 'cors',
            credentials: 'same-origin',
            header: {
                'Content-Type': 'application/x-www-form-urlencoded'
            },
            body: data
        });
        return await response.json();
    },
    async response(id = '') {
        return await fetch(`/api/vacancies/${id}`, {
            method: 'POST',
            mode: 'cors',
            credentials: 'same-origin'
        });
    },
    async toggle(id = '') {
        return await fetch(`/api/vacancies/${id}`, {
            method: 'PATCH',
            mode: 'cors',
            credentials: 'same-origin'
        });
    }
};

// Получение id вакансии.
function getId(row = {}, parentStep = 0) {
    for (i = 0; i < parentStep; i++) {
        row = $(row).parent()
    }
    return parseInt($(row).attr('data-index'));
};

// Возвращает фильтр.
function getFilter() {
    let filter = {
        salary: '',
        older_than: '',
        excluded_sites: [],
        excluded_groups: [],
        positions: [],
        companies: [],
        areas: []
    };

    $('#site .checkbox > input').each(function() {
        if (!$(this).prop('checked')) {
            filter.excluded_sites.push($(this).parent().text().trim());
        }
    });

    $('#group .checkbox > input').each(function() {
        if (!$(this).prop('checked')) {
            if ($(this).parent().text().trim() != 'Нет данных') {
                filter.excluded_groups.push($(this).parent().text().trim());
            }
        }
    });

    $('#position .dropdown__menu span').each(function() {
        if ($(this).text().trim() != 'Нет данных') {
            filter.positions.push($(this).text().trim());
        }
    });

    const salary = $('#salary input').val().toString();
    if (Number(salary) > 0) {
        filter.salary = salary;
    }

    $('#company .dropdown__menu span').each(function() {
        if ($(this).text().trim() != 'Нет данных') {
            filter.companies.push($(this).text().trim());
        }
    });

    $('#area .dropdown__menu span').each(function() {
        if ($(this).text().trim() != 'Нет данных') {
            filter.areas.push($(this).text().trim());
        }
    });

    return filter;
};

const table = $('.vacancies');
const head = $(table).find('thead');
const body = $(table).find('tbody');
const stepInterval = 5000;

var isLoading = false;
var oldFilter = {};

// Обновление всех строк после
// изменения фильтров.
function updateRows() {
    const filter = getFilter();

    if (oldFilter.toString() != filter.toString()) {
        server.rows(filter)
            .then((data) => {
                $(body).find('tr').remove();
                data.forEach(appendRow);
                oldFilter = filter;
            })
            .catch((err) => {
                console.log(err);
            });
    }
};

// Добавляет строку на страницу.
function appendRow(data = {}) {
    const row = composeRow(data);
    $(body).append(row);
};

// Возвращает элемент строки.
function composeRow(data = {}) {
    const tmpl = $('#vacancy__template').prop('content');
    const row = $(tmpl).clone(true);
    const cells = $(row).find('td');

    $(row).find('tr').attr('data-index', data.vacancy_id);
    $(row).find('.logo').addClass(data.site);
    cells[1].innerText = data.group;
    var link = $(cells[2]).find('a');
    $(link).text(data.vacancy_name);
    $(link).prop('href', data.vacancy_link);
    cells[3].innerText = data.salary;
    cells[4].innerText = data.company;
    cells[5].innerText = data.area;
    cells[6].innerText = data.description;
    cells[7].innerText = data.at_published;
    $(cells[8]).find('input').attr('checked', data.selected);

    return row;
};

// Возвращает элемент группы.
function composeGroup(data = '') {
    const tmpl = $('#group__template').prop('content');
    var group = $(tmpl).clone(true);
    var li = $(group).find('li label');
    const content = $(li).html();
    $(li).html(data + content);
    $(li).find('input').prop('checked', true)
    return group;
};

// Возвращает элемент выбранного поля.
function composeLabel(data = '') {
    const tmpl = $('#label__template').prop('content');
    var label = $(tmpl).clone(true);
    $(label).find('span').text(data);
    return label;
};

// Добавляет элемент сайта в шапку таблицы.
function appendSite(data = {}) {
    const site = composeGroup(data);
    $(head).find('#site ul').append(site);
};

// Добавляет элемент группы в шапку таблицы.
function appendGroup(data = {}) {
    const group = composeGroup(data);
    $(head).find('#group ul').append(group);
};

// Добавляет элемент позиции в шапку таблицы.
function appendPosition(data = '') {
    const position = composeLabel(data);
    $(head).find('#position ul:last-child').append(position);
};

// Добавляет минимальную ЗП в шапку таблицы.
function appendSalary(data = '') {
    $(head).find('#salary input').val(data);
};

// Добавляет элемент компании в шапку таблицы.
function appendCompany(data = '') {
    const company = composeLabel(data);
    $(head).find('#company ul:last-child').append(company);
};

// Добавляет элемент региона в шапку таблицы.
function appendArea(data = '') {
    let area = composeLabel(data);
    $(head).find('#area ul:last-child').append(area);
};

// Добавляет в начало таблицы новую строку.
function prependRow(data = {}, index = 0) {
    setTimeout(function() {
        let row = composeRow(data);
        $(body).prepend(row);
        row = $(body).find('tr:first-child');

        $(row)
            .hide()
            .stop(true, true)
            .fadeIn({ duration: 1000, queue: false })
            .css('display', 'none')
            .slideDown({ duration: 'slow', easing: 'swing' })
    }, 500 * index);
};

// Удаляет со страницы список с найденными значениями.
function removeFoundList(element = {}) {
    const list = $(element).find('.found__list');

    $(list).find('li').remove();
    $(list).addClass('hidden');

    if ($(element).hasClass('dropdown__menu--input')) {
        $(element).find('.not__found__data').addClass('hidden');
    }

    if ($(element).find('.input__container').hasClass('input__container--button')) {
        $(element).find('input').val('');
    }
};

// Загрузка контента (строк)
// при прокрутке таблицы вниз.
$(window).scroll(function() {
    const throttleFetch = throttle(checkPosition, 50);
    throttleFetch();
});

// Загрузка контента (строк)
// при изменении размеров экрана.
$(window).resize(function() {
    const throttleFetch = throttle(checkPosition, 50);
    throttleFetch();
});

// Тротлинг для отправления запросов
// при прокрутке таблицы вниз.
function throttle(callee, timeout) {
    var timer = null;

    return function perform(...args) {
        if (timer) return;

        timer = setTimeout(() => {
            callee(...args);
            clearTimeout(timer);
            timer = null;
        }, timeout);
    };
};

// Определение позиции по вертикали.
// Если позиция наиболее близка к концу таблицы,
// то выполняется запрос новых строк с пагинацией.
async function checkPosition() {
    const height = document.body.offsetHeight;
    const screenHeight = window.innerHeight;
    const scrolled = window.scrollY;

    const threshold = height - screenHeight / 4;
    const position = scrolled + screenHeight;

    if (position >= threshold) {
        await fetchRows();
    }
};

// Загрузка контента при прокрутке таблицы вниз.
async function fetchRows() {
    if (isLoading) return;
    isLoading = true;

    const id = getId($(body).find('tr').last());
    var filter = getFilter();
    filter.append('older_than', id);

    const rows = await server.rows(filter);
    rows.forEach(appendRow);

    isLoading = false;
};

// Если таблица пуста, то происходит
// попытка раз в 5 секунд запросить данные.
$(document).ready(function() {
    const interval = () => {
        if ($(body).is(':empty')) {
            const filter = getFilter();
            server.rows(filter)
                .then((data) => {
                    data.forEach(appendRow);
                    return true;
                })
                .catch((err) => {
                    console.log(err);
                });
        }
        return false;
    };

    if (!interval()) {
        setInterval(interval, stepInterval);
    }
});

// Загрузка фильтров при загрузке страницы.
$(document).ready(function() {
    const notFound = (element) => {
        $(element).find('.check__list').addClass('hidden');
        $(element).find('.not__found__data').removeClass('hidden');
    };

    server.filter()
        .then((data) => {
            console.log(data)
            if (data.sites != null) {
                if (data.sites.length > 0) {
                    data.sites.forEach(appendSite);
                } else {
                    notFound($('#site'));
                }
            } else {
                notFound($('#site'));
            }

            if (data.groups != null) {
                if (data.groups.length > 0) {
                    data.groups.forEach(appendGroup);
                } else {
                    notFound($('#group'));
                }
            } else {
                notFound($('#group'));
            }

            // data.positions.forEach(appendPosition);
            // appendSalary(data.salary);
            // data.companies.forEach(appendCompany);
            // data.areas.forEach(appendArea);
        })
        .catch((err) => {
            console.log(err);
        });
});

// Если таблица не пуста, то каждые 20 секунд
// будут запрашиваться новые данные.
// $(document).ready(function() {
//     setInterval(function() {
//         if (!$(body).is(':empty')) {
//             const id = getId($(body).find('tr').first());
//             let filter = getFilter();
//             filter.append('newer_than', id);

//             server.rows(filter)
//                 .then((data) => {
//                     data.forEach(prependRow);
//                 })
//                 .catch((err) => {
//                     console.log(err);
//                 });
//         }
//     }, stepInterval * 4);
// });

// `Click` вне фильтров.
$(document).mouseup(function(e) {
    const dropdown = $('.dropdown');
    if (!$(dropdown).is(e.target) && $(dropdown).has(e.target).length === 0) {
        $('.dropdown .dropdown__menu').each(function() {
            $(this).addClass('hidden');
            $(this).parent().removeClass('active');
            removeFoundList(this);  
        });

        updateRows();
    }
});

// `Click` по выпадающему списку.
$(document).ready(function() {
    oldFilter = getFilter();

    var toggleScroll = (element) => {
        const count = $(element).children().length;
        if (count > 4) {
            $(element).addClass('scroll');
        } else {
            $(element).removeClass('scroll');
        }
    };

    var toggleHidden = (element) => {
        if ($.trim($(element).html()) === '') {
            $(element).addClass('hidden');
        } else {
            $(element).removeClass('hidden');
        }
    };

    $('.dropdown__item').click(function() {
        const menu = $(this).parent().find('.dropdown__menu');

        if ($(menu).is(':hidden')) {
            $(menu).removeClass('hidden');
            $(menu).parent().addClass('active');

            if ($(menu).prop('tagName') === 'UL') {
                toggleScroll(menu);
            } else {
                const ul = $(menu).find('.selected__list');
                toggleScroll(ul);
                toggleHidden(ul);
            }
        } else {
            $(menu).addClass('hidden');
            $(menu).parent().removeClass('active');
            removeFoundList(menu);
            updateRows();
        }

        $('.dropdown__menu').each(function() {
            if (!$(this).is(menu)) {
                const list = $(this).find('.check__list');
                const notFound = $(this).find('.not__found__data');

                if ($(list).children().length === 0) {
                    $(list).addClass('hidden');
                    $(notFound).removeClass('hidden');
                } else {
                    $(list).removeClass('hidden');
                    $(notFound).addClass('hidden');
                }

                $(this).addClass('hidden');
                $(this).parent().removeClass('active');

                removeFoundList(this);     
            }
        });
    });
});

// `Click` по кнопке поиска.
$(document).ready(function() {
    $('.input__container').on('click', 'button', function() {
        const id = $(this).parent().parent().parent().attr('id');
        const value = $(this).parent().find('input').val();
        const menu = $(this).parent().parent();

        if ($.trim(value).length > 0) {
            let params = new URLSearchParams('value=' + value);

            server.search(id, params)
                .then((data) => {
                    const list = $(menu).find('.found__list');
                    const notFound = $(menu).find('.not__found__data');

                    $(list).find('li').remove();

                    if (data.length === 0) {
                        $(list).addClass('hidden');
                        $(notFound).removeClass('hidden');
                    } else {
                        $(list).removeClass('hidden');
                        $(notFound).addClass('hidden');

                        if (data.length > 8) {
                            $(list).addClass('scroll');
                            $(list).addClass('scroll__8')
                        }
                        data.forEach(function(e) {
                            $(list).append('<li>' + e + '</li');
                        });
                    }
                })
                .catch((err) => {
                    console.log(err);
                })
        }
    });
});

// `Click` по кнопке очистить.
$(document).ready(function() {
    $('.input__container input[type="search"]').on('search', function() {
        const search = $(this).parent().parent();
        removeFoundList(search);
    });
});

// `Click` по выбранному элементу.
$(document).ready(function() {
    $('.found__list').on('click', 'li', function() {
        const menu = $(this).parent().parent().parent();
        const foundList = $(menu).find('.found__list');
        const selectedList = $(menu).find('.selected__list');
        let data = composeLabel($(this).text());
        const uniq = new Set();

        $(selectedList).find('li').each(function() {
            uniq.add($(this).text().trim());
        });

        if (!uniq.has($(this).text().trim())) {
            $(this).remove();
            const count = $(foundList).children().length;

            if (count === 0) {
                removeFoundList($(foundList).parent().parent());
            } else if (count <= 8) {
                $(foundList).removeClass('scroll__8');
                $(foundList).removeClass('scroll');
            }

            if ($(selectedList).is(':hidden')) {
                $(selectedList).removeClass('hidden');
            }

            $(selectedList).append(data);

            if ($(selectedList).children().length > 4) {
                $(selectedList).addClass('scroll');
            }
        }
    });
});

// `Click` по кнопке удалить.
$(document).ready(function() {
    $('.selected__list').on('click', '.remove__button', function() {
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

// `Click` по переключателю слайдера.
$(document).ready(function() {
    $('.slidebar th:last-child, .slidebar__button label').click(function() {
        $(head).toggleClass('close');
    });
});

// `Click` по кнопке отклик.
$(document).ready(function() {
    $(table).on('click', '.response .response__button', function(event) {
        const id = getId(event.target, 4);
        server.response(id)
            .then((response) => {
                if (response.ok) {
                    $(this).parent().parent().parent().remove();
                }
            })
            .catch((err) => {
                console.log(err);
            });
    });
});

// `Check` по чекбоксу с вакансией.
$(document).ready(function() {
    $(table).on('click', '.response input', function(event) {
        const id = getId(event.target, 6);

        server.toggle(id)
            .then((response) => {
                if (!response.ok) {
                    console.log(response.status);
                }
            })
            .catch((err) => {
                console.log(err);
            });
    });
});
