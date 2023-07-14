checkCookies();

function checkCookies() {
    const exists = document.cookie.includes('token');

    if (exists) {
        const cookies = document.cookie;
        const time = parseInt(cookies.slice(cookies.indexOf('expire=')).split('expire=').join(''));

        if (time && (time - Date.now()) <= 0) {
            document.cookie = `token=;expires=Thu, 01 Jan 1970 00:00:01 GMT`;
            document.cookie = `expire=;expires=Thu, 01 Jan 1970 00:00:01 GMT`;
            return;
        }
    }
}