import { Component, OnInit } from '@angular/core';
import { MainService } from '../main.service';
import { faPlus, faChevronDown, faArrowRight, faChevronRight, faChevronLeft, faTrash } from "@fortawesome/free-solid-svg-icons";
import { Router } from '@angular/router';

@Component({
  selector: 'app-myplanner',
  templateUrl: './myplanner.component.html',
  styleUrls: ['./myplanner.component.scss']
})
export class MyplannerComponent implements OnInit {

//// variables

  // set imported icon names for component
  faPlus = faPlus;
  arrowDown = faChevronDown;
  arrowRight = faArrowRight;
  chevRight = faChevronRight;
  chevLeft = faChevronLeft;
  faTrash = faTrash;

  // checks if there are any errors in the api call
  apiErr: string = "";

  // determines whether the modals should be displayed or not
  displayTopicModal: boolean = false;
  displayTaskModal: boolean = false;
  displayTaskInfoModal: boolean = false;
  displayTaskEditModal: boolean = false;
  displaySidebar: boolean = false;
  displayDeleteTopicModal: boolean = false;

  // determine if topic is editable
  editTopic: boolean = false;

  // template for creating new topics
  topicTemp = {title: "", tasks: []}

  // template for creating new tasks
  taskTemp = {title: "", description: "", status: "incomplete", type: "", dueDate: "", createdDate: "", progress: "", subTasks: []}

  // template for creating new subtasks
  subTaskTemp = {title: "", status: false}

  // holds sorting method for ng-model
  sortingMethod: string = "";

  // determines whether completed tasks will be shown
  displayCompleted: boolean = false;

  // holds the temporary topic and task index at which a function was called
  tempI = 0;
  tempJ = 0;

//// functions

  getSchedule(){
    return this._mainservice.schedule
  }

  // displays the modals and turns scrolling off 
  openTopicModal(){
    this.apiErr = "";
    this.displaySidebar = false;
    this.displayTopicModal = true;
    document.body.style.overflow = 'hidden';
  }

  openTaskModal(i){
    this.apiErr = "";
    this.displayTaskModal = true;
    document.body.style.overflow = 'hidden';

    // set the topic index for the call
    this.tempI = i;
  }

  openTaskInfoModal(i, j){
    this.apiErr = "";
    this.displayTaskInfoModal = true;
    document.body.style.overflow = 'hidden';

    // set the topic index for the call
    this.tempI = i;
    this.tempJ = j;
  }

  openTaskEditModal(){
    this.displayTaskInfoModal = false;
    this.displayTaskEditModal = true;
  }

  openSidebar() {
    this.displaySidebar = true;
    document.body.style.overflow = 'hidden';
  }

  openDeleteTopicModal(i) {
    this.tempI = i;
    this.displayDeleteTopicModal = true;
    document.body.style.overflow = 'hidden';
  }

  // determines the background color for the task modal
  getBackgroundColor(task){
    if(task.status == "completed"){
      return "green";
    }
    if(task.status == "onHold"){
      return "grey";
    }
    if(task.type == 'checkbox'){
      return "#7575FF";
    }
    if(task.type == 'date'){
      return "red";
    }
    else{
      return "";
    }
  }

  // closes the modals if the correct element is clicked, clears all entered data, and enables scrolling again
  closeModal(e){
    if(e.target.className == "modal" || e.target.className == "close" || e.target.className == "close2"){
      if(this.displayTopicModal){
        this.displayTopicModal = false;
        this.topicTemp = {title: "", tasks: []}
      }
      if(this.displayTaskModal){
        this.displayTaskModal = false;
        this.taskTemp = {title: "", description: "", status: "", type: "", dueDate: "", createdDate: "", progress: "", subTasks: [] }
      }
      if(this.displayTaskInfoModal){
        this.displayTaskInfoModal = false;
      }
      if(this.displayTaskEditModal){
        this.displayTaskEditModal = false;
      }
      if(this.displaySidebar){
        this.displaySidebar = false;
      }
      if(this.displayDeleteTopicModal){
        this.displayDeleteTopicModal = false;
      }
      document.body.style.overflow = 'auto';
    }
  }

  // adds topic to schedule array and sends the modified schedule to the backend
  addTopic(){

    // checks that all required fields are filled
    if(this.topicTemp.title == ""){
      this.apiErr = "missing field";
      return
    }

    this._mainservice.schedule.push(this.topicTemp)
    this.displayTopicModal = false;
    document.body.style.overflow = 'auto';
    this.topicTemp = {title: "", tasks: []}

    // if on the my-planner page, send requests to backend
    if(this.router.url == "/my-planner"){
      let body = {
        schedule: this._mainservice.schedule,
        sort: this._mainservice.sortMethod
      }

      fetch("/api/my-schedule", {
        method: "post",
        body: JSON.stringify(body),
      })
      .then(function(response) {
        return response.json();
      })
      .then((data)=>{
        if(data.msg == "success"){
        }
        if(data.msg == "logged-out"){
          this._mainservice.loggedIn = false;
          this.router.navigate(['/']);
        }
      }) 
      .catch(err => {
        console.log(err);
      });
    }
  }

  // sets edit mode on for topics
  editTopics(){
    this.displaySidebar = false;
    this.editTopic = true;
  }

  // removes topic from the schedule array, waits until the user clicks save before sending schedule to the backend
  deleteTopic(){
    this._mainservice.schedule.splice(this.tempI, 1);
    this.displayDeleteTopicModal = false;
    document.body.style.overflow = 'auto';
  }

  // moves topic left and right in the schedule array, waits until the user clicks save before sending schedule to the backend
  moveTopicLeft(i){
    // if i is not the first element in the array
    if(i!=0){
      var temp = this._mainservice.schedule[i-1];
      this._mainservice.schedule[i-1] = this._mainservice.schedule[i];
      this._mainservice.schedule[i] = temp;
    }
    // if i is the first element in the array
    else{
      var temp = this._mainservice.schedule.shift();
      this._mainservice.schedule.push(temp);
    }
  }

  moveTopicRight(i){
    // if i is not the last element in the array
    if(i!=this._mainservice.schedule.length-1){
      var temp = this._mainservice.schedule[i+1];
      this._mainservice.schedule[i+1] = this._mainservice.schedule[i];
      this._mainservice.schedule[i] = temp;
    }
    // if i is the last element in the array
    else{
      var temp = this._mainservice.schedule.pop();
      this._mainservice.schedule.unshift(temp);
    }
  }

  // saves topics and sends updated schedule to the backend
  saveTopics(){
    this.editTopic = false;

    // if on the my-planner page, send requests to backend
    if(this.router.url == "/my-planner"){
      let body = {
        schedule: this._mainservice.schedule,
        sort: this._mainservice.sortMethod
      }

     fetch("/api/my-schedule", {
        method: "post",
        body: JSON.stringify(body),
      })
      .then(function(response) {
        return response.json();
      })
      .then((data)=>{
        if(data.msg == "success"){
        }
        if(data.msg == "logged-out"){
          this._mainservice.loggedIn = false;
          this.router.navigate(['/']);
        }
      }) 
      .catch(err => {
        console.log(err);
      });
    }
  }

  // adds task to schedule array and sends the modified schedule to the backend
  addTask(){

    // checks that all required fields are filled
    if(this.taskTemp.title == ""||this.taskTemp.type == ""){
      this.apiErr = "missing field";
      return
    }

    this.taskTemp.createdDate = new Date().toISOString().slice(0, 10)

    // use the index set in the opentaskModal() function to append the task array
    this._mainservice.schedule[this.tempI].tasks.push(this.taskTemp)
    this.displayTaskModal = false;
    document.body.style.overflow = 'auto';
    this.taskTemp = {title: "", description: "", status: "incomplete", type: "", dueDate: "", createdDate: "", progress: "", subTasks: [] }

    // resort schedule after adding new task
    this.sortSchedule();

    // if on the my-planner page, send requests to backend
    if(this.router.url == "/my-planner"){
      let body = {
        schedule: this._mainservice.schedule,
        sort: this._mainservice.sortMethod
      }

     fetch("/api/my-schedule", {
        method: "post",
        body: JSON.stringify(body),
      })
      .then(function(response) {
        return response.json();
      })
      .then((data)=>{
        if(data.msg == "success"){
        }
        if(data.msg == "logged-out"){
          this._mainservice.loggedIn = false;
          this.router.navigate(['/']);
        }
      }) 
      .catch(err => {
        console.log(err);
      });
    }
  }

  // closes task edit modal and sends updated info to backend
  saveTask(){

    // close modal
    this.displayTaskEditModal = false;
    document.body.style.overflow = 'auto';

    // if on the my-planner page, send requests to backend
    if(this.router.url == "/my-planner"){
      let body = {
        schedule: this._mainservice.schedule,
        sort: this._mainservice.sortMethod
      }

     fetch("/api/my-schedule", {
        method: "post",
        body: JSON.stringify(body),
      })
      .then(function(response) {
        return response.json();
      })
      .then((data)=>{
        if(data.msg == "success"){
        }
        if(data.msg == "logged-out"){
          this._mainservice.loggedIn = false;
          this.router.navigate(['/']);
        }
      }) 
      .catch(err => {
        console.log(err);
      });
    }
  }

  // removes task from schedule array and sends the modified schedule to the backend
  deleteTask(){

    // close modal
    this.displayTaskInfoModal = false;
    document.body.style.overflow = 'auto';

    // remove task from appropriate index
    this._mainservice.schedule[this.tempI].tasks.splice(this.tempJ, 1);

    // if on the my-planner page, send requests to backend
    if(this.router.url == "/my-planner"){
      let body = {
        schedule: this._mainservice.schedule,
        sort: this._mainservice.sortMethod
      }

     fetch("/api/my-schedule", {
        method: "post",
        body: JSON.stringify(body),
      })
      .then(function(response) {
        return response.json();
      })
      .then((data)=>{
        if(data.msg == "success"){
        }
        if(data.msg == "logged-out"){
          this._mainservice.loggedIn = false;
          this.router.navigate(['/']);
        }
      }) 
      .catch(err => {
        console.log(err);
      });
    }
  }

  // adds subtask to schedule array and sends the modified schedule to the backend
  addSubTask(){
    
    // checks that all required fields are filled
    if(this.subTaskTemp.title == ""){
      return
    }

    // add subtask to appropriate task
    this._mainservice.schedule[this.tempI].tasks[this.tempJ].subTasks.push(this.subTaskTemp);

    // reset subTaskTemp array
    this.subTaskTemp = {title: "", status: false}

    // recalculate progress for task
    this.calculateProgress();

    // if on the my-planner page, send requests to backend
    if(this.router.url == "/my-planner"){
      let body = {
        schedule: this._mainservice.schedule,
        sort: this._mainservice.sortMethod
      }

     fetch("/api/my-schedule", {
        method: "post",
        body: JSON.stringify(body),
      })
      .then(function(response) {
        return response.json();
      })
      .then((data)=>{
        if(data.msg == "success"){
        }
        if(data.msg == "logged-out"){
          this._mainservice.loggedIn = false;
          this.router.navigate(['/']);
        }
      }) 
      .catch(err => {
        console.log(err);
      });
    }
  }

  // removes subtask from schedule array and sends the modified schedule to the backend
  deleteSubtask(k){

    // remove subtask from appropriate index
    this._mainservice.schedule[this.tempI].tasks[this.tempJ].subTasks.splice(k, 1);

    // recalculate progress for task
    this.calculateProgress();

    // if on the my-planner page, send requests to backend
    if(this.router.url == "/my-planner"){
      let body = {
        schedule: this._mainservice.schedule,
        sort: this._mainservice.sortMethod
      }

     fetch("/api/my-schedule", {
        method: "post",
        body: JSON.stringify(body),
      })
      .then(function(response) {
        return response.json();
      })
      .then((data)=>{
        if(data.msg == "success"){
        }
        if(data.msg == "logged-out"){
          this._mainservice.loggedIn = false;
          this.router.navigate(['/']);
        }
      }) 
      .catch(err => {
        console.log(err);
      });
    }
  }

  // calculates percentage of task completed for the progressbar, uses number of subtasks completed to calculate
  calculateProgress(){

    // if no subtasks exist currently set percentage to 0
    if(this._mainservice.schedule[this.tempI].tasks[this.tempJ].subTasks.length == 0){
      this._mainservice.schedule[this.tempI].tasks[this.tempJ].progress = "0";
    }

    let completed = 0;

    // loop through subtasks and increase completed by 1 if subtask is marked as completed
    for(let k=0;k<this._mainservice.schedule[this.tempI].tasks[this.tempJ].subTasks.length; k++){
      if(this._mainservice.schedule[this.tempI].tasks[this.tempJ].subTasks[k].status){
        completed++;
      }
    }

    // set percentage to percent completed
    this._mainservice.schedule[this.tempI].tasks[this.tempJ].progress = (completed/this._mainservice.schedule[this.tempI].tasks[this.tempJ].subTasks.length*100).toString();

    // if progress is a 100%, set the task to completed
    if((completed/this._mainservice.schedule[this.tempI].tasks[this.tempJ].subTasks.length*100).toString() == "100"){
      this._mainservice.schedule[this.tempI].tasks[this.tempJ].status = "completed";
    }
  }

  // set sorting method for service and send data to backend
  setSortingMethod(){
    this._mainservice.sortMethod = this.sortingMethod;
    this.sortSchedule();

    // if on the my-planner page, send requests to backend
    if(this.router.url == "/my-planner"){
      let body = {
        schedule: this._mainservice.schedule,
        sort: this._mainservice.sortMethod
      }

     fetch("/api/my-schedule", {
        method: "post",
        body: JSON.stringify(body),
      })
      .then(function(response) {
        return response.json();
      })
      .then((data)=>{
        if(data.msg == "success"){
        }
        if(data.msg == "logged-out"){
          this._mainservice.loggedIn = false;
          this.router.navigate(['/']);
        }
      }) 
      .catch(err => {
        console.log(err);
      });
    }
  }

  // sorts the schedule based on the format chosen by the user
  sortSchedule(){
    if(this._mainservice.sortMethod == "desc-due"){
      this._mainservice.schedule[this.tempI].tasks.sort(function(a,b){
        // if a does not have a due date, return a after b
        if(a.dueDate == ""){
          return 1;
        }
        // if b does not have a due date, return b after a
        if(b.dueDate == ""){
          return -1;
        }
        return new Date(a.dueDate).getTime() - new Date(b.dueDate).getTime();
      });
      return
    }
    if(this._mainservice.sortMethod == "asc-due"){
      this._mainservice.schedule[this.tempI].tasks.sort(function(a,b){
        // if a does not have a due date, return a after b
        if(a.dueDate == ""){
          return 1;
        }
        // if b does not have a due date, return b after a
        if(b.dueDate == ""){
          return -1;
        }
        return new Date(b.dueDate).getTime() - new Date(a.dueDate).getTime();
      });
      return
    }
    if(this._mainservice.sortMethod == "desc-created"){
      this._mainservice.schedule[this.tempI].tasks.sort(function(a,b){
        return new Date(a.createdDate).getTime() - new Date(b.createdDate).getTime();
      });
      return
    }
    if(this._mainservice.sortMethod == "asc-created"){
      this._mainservice.schedule[this.tempI].tasks.sort(function(a,b){
        return new Date(b.createdDate).getTime() - new Date(a.createdDate).getTime();
      });
      return
    }
  }

  // changes the session storage variable for whether completed tasks would be shown
  changeDisplayCompleted(){
    sessionStorage.setItem("displayCompleted", this.displayCompleted.toString());
  }

  constructor(private _mainservice: MainService, private router: Router) { }

  ngOnInit(): void {

    // if user is not logged in while trying to access the "my-planner" version of this page, redirect to home
    if(this.router.url == "/my-planner"){
      if(!this._mainservice.loggedIn){
        this.router.navigate(['/']);
      }

      // fetch the users schedule and sorting method if logged in and on correct page
      if(this._mainservice.loggedIn){
        fetch("/api/my-schedule", {
          method: "get",
        })
        .then(function(response) {
          return response.json();
        })
        .then((data)=>{
          this._mainservice.schedule = data.schedule;
          this._mainservice.sortMethod = data.sort;
        }) 
        .catch(err => {
          console.log(err);
        });
      }
    }

    this.sortingMethod = this._mainservice.sortMethod;

    // get the displayCompleted variable from session storage
    this.displayCompleted = sessionStorage.getItem("displayCompleted") == "true";
  }

}
