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
