window.onload = function() {
    const path = window.location.pathname;
    const split = path.split("/");
    const id = split[split.length - 1];

    const title = document.getElementsByTagName("title")[0]; 
    const name = document.getElementById("name");
    const description = document.getElementById("description");

    function populateData(data) {
        title.innerText = data.name;
        name.innerText = data.name;
        description.innerText = data.description;
    }

    fetch("/api/problems/get/" + id)
    .then(res => res.json())
    .then(populateData)
    .catch(err => console.log(err));
}
