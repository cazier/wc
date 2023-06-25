function loadGroupStage(group) {
    const request = new Request(
        `${window.location.origin}/api/match/group/${group}`, {
        method: "GET",
    });

    fetch(request)
        .then((response) => {
            if (response.status === 200) {
                return response.json();
            } else {
                throw new Error("Whoops");
            }
        })
        .then((response) => {
            console.log(response);
        });
}

function getCookie(name) {
    const cookies = document.cookie;

    for (var cookie in cookies.split(';')) {
        var [key, value] = cookie.split('=');
        if (key == name) {
            return value;
        }
    }
    return null;
}
