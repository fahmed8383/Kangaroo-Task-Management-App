import { Component, OnInit } from '@angular/core';
import { MainService } from '../main.service';
import { faEdit } from "@fortawesome/free-solid-svg-icons";
import { Router } from '@angular/router';

@Component({
  selector: 'app-myaccount',
  templateUrl: './myaccount.component.html',
  styleUrls: ['./myaccount.component.scss']
})
export class MyaccountComponent implements OnInit {

//// variables

  // set imported icon names for component
  faEdit = faEdit;

  // determines whether to show editable inputs or not
  editEmail: boolean = false;
  editPassword: boolean = false;

  // holds the registration info
  regInfo = {
    email: "",
    password: "",
    confirmPassword: "", 
  }

  // checks if there are any errors in the form
  regErr: string = "";

  // checks if the email is still available and not already taken
  emailAvailability: boolean = true;

//// functions

  // checks if the entered email is a valid email format
  validateEmail(email){
    var emailFormat = /\S+@\S+\.\S+/;
    return emailFormat.test(email);
  }

  // checks if the username exists in the database and sets the usernameAvailability variable accordingly
  checkEmailAvailability(){
    // if the entered email is not in a valid format, return
    if(!this.validateEmail(this.regInfo.email)){
      return
    }
    fetch("/api/check-email-availability", {
      method: "post",
      body: JSON.stringify(this.regInfo),
    })
    .then(function(response) {
      return response.json();
    })
    .then((data)=>{
      if(data.msg == "success"){
        this.emailAvailability = false
      }
      else{
        this.emailAvailability = true
      }
    }) 
    .catch(err => {
      console.log(err);
    });
  }

  // responsible for reseting the email for the user
  resetEmail(){

    // if the user entered an unavailable email, return prematurely
    if(!this.emailAvailability){
      return
    }

    // checks that all required fields are filled before sending request to backend
    if(this.regInfo.email == ""){
      this.regErr = "missing field";
    }
    // check if the entered email is a valid email format
    else if(!this.validateEmail(this.regInfo.email)){
      this.regErr = "invalid email";
    }
    else{
      this.regErr = "";

      fetch("/api/reset-email", {
        method: "post",
        body: JSON.stringify(this.regInfo),
      })
      .then(function(response) {
        return response.json();
      })
      .then((data)=>{
        // if successful, set editEmail variable back to false
        if(data.msg = "success"){
          this.editEmail = false
        }
      }) 
      .catch(err => {
        console.log(err);
      });
    }
  }

  // responsible for reseting the password for the user
  resetPassword(){

    // checks that all required fields are filled before sending request to backend
    if(this.regInfo.password == ""){
      this.regErr = "missing field";
      return
    }

    // send request to backend
    fetch("/api/reset-loggedin-password", {
      method: "post",
      body: JSON.stringify(this.regInfo),
    })
    .then(function(response) {
      return response.json();
    })
    .then((data)=>{
      // if successful, set editPassword variable back to false and clear the password fields
      if(data.msg == "success"){
        this.regInfo.password = "";
        this.regInfo.confirmPassword = "";
        this.editPassword = false;
      }
    }) 
    .catch(err => {
      console.log(err);
    });
  }

  constructor(private _mainservice: MainService, private router: Router) { }

  ngOnInit(): void {
    /*// if user is not logged in while trying to access this page, redirect to home
    if(!this._mainservice.loggedIn){
      this.router.navigate(['/']);
    }

    // fetch the users schedule and sorting method if logged in and on correct page
    if(this._mainservice.loggedIn){
      fetch("/api/get-user-info", {
        method: "get",
      })
      .then(function(response) {
        return response.json();
      })
      .then((data)=>{
        this.regInfo = data
      }) 
      .catch(err => {
        console.log(err);
      });
    }*/
  }

}
