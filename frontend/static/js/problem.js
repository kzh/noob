window.onload = function() {
    const path = window.location.pathname;
    const split = path.split("/");
    const id = split[split.length - 2];

    function handleError(res) {
        if (res.error == undefined) {
            return res;
        }

        throw new Error(res.error);
    }

    const title = document.getElementsByTagName("title")[0]; 
    const msg = document.querySelector(".message");
    function error(err) {
        if (err instanceof SyntaxError) {
            err = new Error("Service unavailable.");
        }

        msg.classList.remove("hidden");
        msg.innerText = err;
        title.innerText = "Error";
    }

    function populateData(data) {
        const name = document.getElementById("name");
        const description = document.getElementById("description");

        title.innerText = "Problem: " + data.name;
        name.innerText = data.name;
        description.textContent = data.description;

        const content = document.querySelector(".content");
        content.classList.remove("hidden");

        msg.classList.add("hidden"); 
    }

    fetch("/api/problems/get/" + id)
    .then(res => res.json())
    .then(handleError)
    .then(populateData)
    .catch(error);

    const editBtn = document.getElementById("edit");
    if (editBtn) {
        editBtn.addEventListener("click", function() {
            window.location.href += "edit/";
        });
    }

    const idEl = document.querySelector("input[name=\"id\"]");
    if (idEl) {
        idEl.value = id;
    }

    function submitCallback(res) {
        if (res.message != undefined) {
            msg.classList.remove("hidden"); 
            msg.innerText = res.message;
        }
    }

    const submitBtn = document.getElementById("submit");
    submitBtn.addEventListener("click", function() {
        const code = document.getElementById("code");

        const submission = {
            id: id,
            code: code.value,
        };

        let formBody = [];
        for (const property in submission) {
              const encodedKey = encodeURIComponent(property);
              const encodedValue = encodeURIComponent(submission[property]);
              formBody.push(encodedKey + "=" + encodedValue);
        }
        formBody = formBody.join("&");
        
        fetch("/api/submissions/submit", {
            method: "POST",
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: formBody,
        })
        .then(res => res.json())
        .then(handleError)
        .then(submitCallback)
        .catch(error);
    });
}
