window.onload = function() {
    const path = window.location.pathname;
    const split = path.split("/");
    const id = split[split.length - 3];

    const idEl = document.querySelector("input[name=\"id\"]");
    idEl.value = id;

    function handleError(res) {
        if (res.error == undefined) {
            return res;
        }

        throw new Error(res.error);
    }

    const msg = document.querySelector(".message");
    function error(err) {
        if (err instanceof SyntaxError) {
            err = new Error("Service unavailable.");
        }
        
        msg.innerText = err;
    }

    let ready1 = false;
    let ready2 = false;

    function reveal() {
        if (!ready1 || !ready2) {
            return;
        }

        const content = document.querySelector(".content");
        content.classList.remove("hidden");

        msg.classList.add("hidden"); 
    }

    function populateProblem(data) {
        const name = document.querySelector("input[name=\"name\"]");
        const description = document.querySelector("textarea[name=\"description\"]");
        name.value = data.name;
        description.value = data.description;

        ready1 = true;
        reveal();
    }

    function populateIO(data) {
        const inputs = document.querySelector("textarea[name=\"inputs\"]");
        const outputs = document.querySelector("textarea[name=\"outputs\"]");
        inputs.value = data.inputs;
        outputs.value = data.outputs;

        ready2 = true;
        reveal();
    }

    fetch("/api/problems/get/" + id)
    .then(res => res.json())
    .then(handleError)
    .then(populateProblem)
    .catch(error);

    fetch("/api/problems/get/" + id + "/io")
    .then(res => res.json())
    .then(handleError)
    .then(populateIO)
    .catch(error);
}
