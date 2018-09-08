window.onload = function() {
    const problemsContainer = document.querySelector(".problems");

    function createProblem(id, name) {
        const problem = document.createElement("a");
        problem.href = "/problem/" + id;

        const text = document.createTextNode(name);
        problem.appendChild(text);

        return problem;
    }

    const br = document.createElement("br");
    function populateProblems(problems) {
        for (const p of problems) {
            const el = createProblem(p.id, p.name);
            problemsContainer.appendChild(el);
            problemsContainer.appendChild(br);
        }
    }

    fetch("/api/problems/list")
    .then(res => res.json())
    .then(populateProblems)
    .catch(err => console.log(err));
}
