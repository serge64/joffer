const loginForm = {
    form: document.getElementById('login'),
    email: document.getElementById('login-email'),
    password: document.getElementById('login-password')
};

const registerForm = {
    form: document.getElementById('register'),
    email: document.getElementById('register-email'),
    password: document.getElementById('register-password'),
    password2: document.getElementById('register-password2')
};

const errFieldBlank = 'Поле не должно быть пустым';
const errEmailNotValid = 'Электронная почта не действительна';
const errFieldNotValid = 'Пароль должен содержать не менее 6 символов';
const errPasswordNotConfirm = 'Пароль не совпадает';

function setErrorFor(input = {}, message = '') {
    input.classList.add('form__input--error');
    input.parentElement.querySelector('.form__input-error-message').textContent = message;
}

function setSuccessFor(input = {}) {
    input.classList.remove('form__input--error');
    input.parentElement.querySelector('.form__input-error-message').textContent = '';
}

function setErrorForMessage(messageObject = {}, type = '', message = '') {
    let msg = message.slice(0, 1).toUpperCase() + message.slice(1).toLowerCase();
    messageObject.textContent = msg;
    messageObject.classList.remove('form__message--success', 'form__message--error');
    messageObject.classList.add(`form__message--${type}`);
}

function setSuccessForMessage(messageObject = {}) {
    messageObject.textContent = '';
    messageObject.classList.remove('form__message--success', 'form__message--error');
}

function isEmail(email = '') {
    const re = /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
    return re.test(email);
}

function checkLoginInputs() {
    const emailValue = loginForm.email.value.trim();
    const passwordValue = loginForm.password.value.trim();
    var check = true;

    if (emailValue === '') {
        setErrorFor(loginForm.email, errFieldBlank);
        check = false;
    } else if (!isEmail(emailValue)) {
        setErrorFor(loginForm.email, errEmailNotValid);
        check = false;
    } else {
        setSuccessFor(loginForm.email);
    }

    if (passwordValue === '') {
        setErrorFor(loginForm.password, errFieldBlank);
        check = false;
    } else {
        setSuccessFor(loginForm.password);
    }

    setSuccessForMessage(loginForm.form.querySelector('.form__message'))

    return check;
}

function checkRegisterInputs() {
    const emailValue = registerForm.email.value.trim();
    const passwordValue = registerForm.password.value.trim();
    const password2Value = registerForm.password2.value.trim();
    var check = true;

    if (emailValue === '') {
        setErrorFor(registerForm.email, errFieldBlank);
        check = false;
    } else if (!isEmail(emailValue)) {
        setErrorFor(registerForm.email, errEmailNotValid);
        check = false;
    } else {
        setSuccessFor(registerForm.email);
    }

    if (passwordValue === '') {
        setErrorFor(registerForm.password, errFieldBlank);
        check = false;
    } else if (passwordValue.length > 0 && passwordValue.length < 6) {
        setErrorFor(registerForm.password, errFieldNotValid);
        check = false;
    } else {
        setSuccessFor(registerForm.password);

        if (passwordValue != password2Value) {
            setErrorFor(registerForm.password2, errPasswordNotConfirm);
            check = false;
        } else {
            setSuccessFor(registerForm.password2);
        }
    }

    setSuccessForMessage(registerForm.form.querySelector('.form__message'))

    return check;
}

function handleResponse(messageObject = {}, response = {}) {
    if (response.redirected) {
        window.location.href = response.url;
    } else {
        responseObject.json()
            .then((data) => {
                setErrorForMessage(messageObject, 'error', data.message);
            });
    }
}

async function handleRequest(url = '', data = {}) {
    const response = await fetch(url, {
        method: 'POST',
        // mode: 'cors',
        // credentials: 'same-origin',
        // redirect: 'follow',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    });

    return await response;
}

loginForm.form.querySelectorAll('.form__input').forEach(e => {
    e.addEventListener('input', () => {
        setSuccessFor(e);
    });
});

registerForm.form.querySelectorAll('.form__input').forEach(e => {
    e.addEventListener('input', () => {
        setSuccessFor(e);
    });
});

loginForm.form.querySelector('.form__link').addEventListener('click', e => {
    e.preventDefault();
    loginForm.form.classList.add('form--hidden');
    registerForm.form.classList.remove('form--hidden');
});

registerForm.form.querySelector('.form__link').addEventListener('click', e => {
    e.preventDefault();
    loginForm.form.classList.remove('form--hidden');
    registerForm.form.classList.add('form--hidden');
});

loginForm.form.querySelector('.form__button').addEventListener('click', e => {
    e.preventDefault();
    if (checkLoginInputs()) {
        let data = {
            email: loginForm.email.value.trim(),
            password: loginForm.password.value.trim()
        };

        handleRequest('/auth/signin', data)
                .then((response) => {
                    handleResponse(loginForm.form.querySelector('.form__message'), response);
                });
    }
});

registerForm.form.querySelector('.form__button').addEventListener('click', e => {
    e.preventDefault();
    if (checkRegisterInputs()) {
        let data = {
            email: registerForm.email.value.trim(),
            password: registerForm.password.value.trim()
        };

        handleRequest('/auth/signup', data)
            .then((response) => {
                handleResponse(registerForm.form.querySelector('.form__message'), response);
            });
    }
});