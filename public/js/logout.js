$(document).ready(function() {
    $('#logout').click(function() {
        fetchLogout()
            .then((response) => {
                if (response.redirected) {
                    window.location.href = response.url;
                }
            })
            .catch((err) => {
                console.log(err);
            })
    });
});

async function fetchLogout() {
    const response = await fetch('/auth/logout', {
        method: 'POST',
        mode: 'cors',
        credentials: 'same-origin',
        redirect: 'follow'
    });
    return await response;
}