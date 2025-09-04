'use strict';

const taskContainer = `<main class="container">
    <div class="title">Todo App Zaba</div>
   <form class="form" onsubmit="return false">
        <div class="input-form"> 
        <input class="title-input" type="text" name="" id="title" placeholder="enter your task" >
        <input class="desc-input" type="text" name="" id="dis"placeholder="enter the discription" >
        </div>
        <button type="button" class="add-task">Add Task</button>

    </form>
    <div class="tasks-container">
       
    </div>
</main>
`
const task = `
 <div class="task">
    <div class="info">
        <h3 class="title-Task"></h3>
        <hr>
        <p class="description-Task"></p>
    </div>
    <div class="actions">
         <button type="button" class="update">update</button>
        <button type="button" class="delete">delete</button>
    </div>
</div>
`
document.body.innerHTML = taskContainer
const addTask = document.querySelector(".add-task");
const titleInput = document.querySelector(".title-input");
const descInput = document.querySelector(".desc-input");
const tasksContainer = document.querySelector(".tasks-container");
const Info = document.querySelector(".input-form");
let requared = document.createElement("p")
requared.textContent = "input can not be empty"
requared.style = "font-size: x-small ;color:red;visibility: hidden;"
Info.appendChild(requared)


let id = 0;
let ne = document.createElement("div")
ne.innerHTML = task
let taskNo = ne.firstElementChild
addTask.addEventListener('click', e => {
    e.preventDefault();

    if (titleInput.value === "" || descInput.value === "") {
        titleInput.style.border = "2px solid red";
        descInput.style.border = "2px solid red";
        requared.style.visibility = "visible";
    } else {
        let newTask = taskNo.cloneNode(true); // clone instead of reusing
        newTask.querySelector(".title-Task").textContent = titleInput.value;
        newTask.querySelector(".description-Task").textContent = descInput.value;

        let data = {
            id: id++,
            Title: titleInput.value,
            Description: descInput.value
        };
        console.log(data);
        fetch("http://localhost:8080/Posts", {

            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)

        })

        titleInput.value = ""
        descInput.value = ""

        console.log(titleInput.textContent, descInput.textContent);


        newTask.querySelector(".delete").addEventListener("click", () => {
            tasksContainer.removeChild(newTask);
        });

        tasksContainer.appendChild(newTask);

        titleInput.style.border = "none";
        descInput.style.border = "none";
        requared.style.visibility = "hidden";

    }
});
fetch("http://localhost:8080/Posts")
    .then(response => response.json())
    .then(data => {
        data.forEach(element => {
            let el = document.createElement("div")
            el.innerHTML = task
            let taskEl = el.firstElementChild
            taskEl.setAttribute("data-id", element.id);
            taskEl.querySelector(".title-Task").textContent = element.Title
            taskEl.querySelector(".description-Task").textContent = element.Description
            taskEl.querySelector(".delete").addEventListener("click", () => {
                const id = taskEl.getAttribute("data-id");
                fetch("http://localhost:8080/Posts/" + id, {
                    method: 'DELETE'
                }).then(() => {
                    tasksContainer.removeChild(taskEl);
                });
            });

            tasksContainer.appendChild(taskEl);
        });
    })
    .catch(err => console.error(" problem in fetch  :", err));
