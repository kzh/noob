window.onload = function() {
    const problemsContainer = document.querySelector(".problems");

    function createProblem(id, name) {
        const problem = document.createElement("a");
        problem.href = "/problem/" + id;

        const text = document.createTextNode(name);
        problem.appendChild(text);

        return problem;
    }

    function handleError(res) {
        if (res.error == undefined) {
            return res;
        }

        throw new Error(res.error);
    }

    function populateProblems(problems) {
        for (const p of problems) {
            const el = createProblem(p.id, p.name);
            problemsContainer.appendChild(el);

            const br = document.createElement("br");
            problemsContainer.appendChild(br);
        }
    }

    fetch("/api/problems/list")
    .then(res => res.json())
    .then(handleError)
    .then(populateProblems)
    .catch(err => alert(err));
}
