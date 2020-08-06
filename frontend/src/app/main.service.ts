import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class MainService {

  constructor() { }

//// variables

  // states whether user is logged in or not by seeing if session id cookie exists
  loggedIn: boolean = false

  // holds the users schedule
  schedule = [
    {title: "Test 1", tasks: [{title: "Machine Learning Course", description: "", status: "incomplete", type: "subtasks", dueDate: "2020-07-26", createdDate: "2020-07-26", progress: "80", subTasks: []}, {title: "Machine Learning Course", description: "", status: "incomplete", type: "subtasks", dueDate: "2020-07-25", createdDate: "2020-07-25", progress: "80", subTasks: []}]},
    {title: "Test 2", tasks: [{title: "Machine Learning Course", description: "", status: "incomplete", type: "checkbox", dueDate: "2020-06-23", createdDate: "2020-06-23", progress: "80", subTasks: []}]},
  ]

  // holds the sorting method
  sortMethod = "desc-due";

//// functions

  // checks if a cookie by input name exists
  // modified from https://www.w3schools.com/js/js_cookies.asp
  cookieExists(cname) {
    var name = cname + "=";
    var decodedCookie = decodeURIComponent(document.cookie);
    var ca = decodedCookie.split(';');
    for(var i = 0; i <ca.length; i++) {
      var c = ca[i];
      while (c.charAt(0) == ' ') {
        c = c.substring(1);
      }
      if (c.indexOf(name) == 0) {
        return true
      }
    }
    return false
  }
}
