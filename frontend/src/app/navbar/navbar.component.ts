import { Component, HostListener, OnInit } from '@angular/core';
import { faTimes } from "@fortawesome/free-solid-svg-icons";
import { faBars } from "@fortawesome/free-solid-svg-icons";

@Component({
  selector: 'app-navbar',
  templateUrl: './navbar.component.html',
  styleUrls: ['./navbar.component.scss']
})
export class NavbarComponent implements OnInit {

//// variables

  // set imported icon names for component
  faTimes = faTimes;
  faBars = faBars;

  // set base mobile nav to false
  showMobileNav: boolean = false;

  // determines the transparency of the navbar
  navbarTransparent: boolean = true;

  // determines whether the modals connected to the navbar should be displayed or not
  displayRegModal: boolean = false;
  displayVerificationModal: boolean = false;
  displayLoginModal: boolean = false;

  // holds the registration info
  regInfo = {
    userName: "",
    email: "",
    password: "",
    confirmPassword: "",
    captcha: "",
    verificationCode: "" 
  }

  // holds the response from registration api call
  regResp: string = "";

  // checks if there are any errors in the registration form
  regErr: string = "";

  // checks if the username is still available and not already taken
  usernameAvailability: boolean = true;

  // checks if the email is still available and not already taken
  emailAvailability: boolean = true;

//// functions

  // change mobilenav value to the opposite on click
  mobileNavClick(){
    this.showMobileNav = !this.showMobileNav;
  }

  // on scroll listener to change the navabar color after a certain threshold
  @HostListener("window:scroll", [])
  onWindowScroll() {

    const number = window.pageYOffset || document.documentElement.scrollTop || document.body.scrollTop || 0;
    if (number > 100) {
      this.navbarTransparent = false;
    } else {
      this.navbarTransparent = true;
    }

  }

  // displays the modals and turns scrolling off
  openRegModal(){
    this.regErr = "";
    this.displayRegModal = true;
    document.body.style.overflow = 'hidden';
  }

  openLoginModal(){
    this.regErr = "";
    this.displayLoginModal = true;
    document.body.style.overflow = 'hidden';
  }

  // closes the modals if the correct element is clicked, clears all entered data, and enables scrolling again
  closeRegModal(e){
    if(e.target.className == "modal" || e.target.className == "close"){
      this.displayRegModal = false;
      document.body.style.overflow = 'auto';
      this.regInfo = {
        userName: "",
        email: "",
        password: "",
        confirmPassword: "",
        captcha: "",
        verificationCode: "" 
      }
    }
  }

  closeLoginModal(e){
    if(e.target.className == "modal" || e.target.className == "close"){
      this.displayLoginModal = false;
      document.body.style.overflow = 'auto';
      this.regInfo = {
        userName: "",
        email: "",
        password: "",
        confirmPassword: "",
        captcha: "",
        verificationCode: "" 
      }
    }
  }

  // closes the email verification modal and sends a request to backend to delete registration entry for the user
  closeVerificationModal(){
    fetch("/api/cancel-registration", {
      method: "post",
      body: JSON.stringify(this.regInfo),
    })
    .then(function(response) {
      return response.json();
    })
    .then((data)=>{

      // if successful, close the modal and reset the form
      if(data.msg == "success"){
        this.displayVerificationModal = false;
        document.body.style.overflow = 'auto';
        this.regInfo = {
          userName: "",
          email: "",
          password: "",
          confirmPassword: "",
          captcha: "",
          verificationCode: "" 
        }
      }
    }) 
    .catch(err => {
      console.log(err);
    });
  }

  // responsible for registering user without email confimation. The email confirmation is required to fully register the account. Sends an email with the code.
  sendEmailVerification(){

    // if the user entered an unavailable email or username, return prematurely
    if(!this.emailAvailability || !this.usernameAvailability){
      return
    }

    // checks that all required fields and the captcha are filled before sending request to backend
    if(this.regInfo.userName == ""||this.regInfo.email == ""||this.regInfo.password==""||this.regInfo.password!=this.regInfo.confirmPassword||this.regInfo.captcha==""){
      this.regErr = "missing field";
    }
    // check if the entered email is a valid email format
    else if(!this.validateEmail(this.regInfo.email)){
      this.regErr = "invalid email";
    }
    else{
      this.regErr = "";

      fetch("/api/send-email-verification", {
        method: "post",
        body: JSON.stringify(this.regInfo),
      })
      .then(function(response) {
        return response.json();
      })
      .then((data)=>{
        // if succesful, close the modal and open the email verification modal
        if(data.msg = "success"){
          this.displayRegModal = false;
          this.displayVerificationModal = true;
        }
      }) 
      .catch(err => {
        console.log(err);
      });
    }
  }

  // responsible for sending resending email confirmation code.
  resendEmailVerification(){

    // if the user entered an unavailable email or username, return prematurely
    if(!this.emailAvailability){
      return
    }

    // checks that the email field is filled before sending request to backend
    if(this.regInfo.email == ""){
      this.regErr = "missing field";
    }
    // check if the entered email is a valid email format
    else if(!this.validateEmail(this.regInfo.email)){
      this.regErr = "invalid email";
    }
    else{
      this.regErr = "";

      fetch("/api/resend-email-verification", {
        method: "post",
        body: JSON.stringify(this.regInfo),
      })
      .then(function(response) {
        return response.json();
      })
      .then((data)=>{
      }) 
      .catch(err => {
        console.log(err);
      });
    }
  }

  validateVerificationCode(){
    // if the user entered an unavailable email or username, return prematurely
    if(!this.emailAvailability){
      return
    }

    // checks that all required fields and the captcha are filled before sending request to backend
    if(this.regInfo.userName == ""||this.regInfo.email == ""||this.regInfo.verificationCode==""){
      this.regErr = "missing field";
      return
    }
    // check if the entered email is a valid email format
    if(!this.validateEmail(this.regInfo.email)){
      this.regErr = "invalid email";
      return
    }

    fetch("/api/validate-verification-code", {
      method: "post",
      body: JSON.stringify(this.regInfo),
    })
    .then(function(response) {
      return response.json();
    })
    .then((data)=>{

      // if successful, close the modal and reset the form
      if(data.msg == "success"){
        this.regErr = "";
        this.displayVerificationModal = false;
        document.body.style.overflow = 'auto';
        this.regInfo = {
          userName: "",
          email: "",
          password: "",
          confirmPassword: "",
          captcha: "",
          verificationCode: "" 
        }
      }
      else{
        this.regErr = "invalid authorization code";
      }
    }) 
    .catch(err => {
      console.log(err);
    });
  }

  // checks if the entered email is a valid email format
  validateEmail(email){
    var emailFormat = /\S+@\S+\.\S+/;
    return emailFormat.test(email);
  }

  // checks if the username exists in the database and sets the usernameAvailability variable accordingly
  checkUsernameAvailability(){
    fetch("/api/check-username-availability", {
      method: "post",
      body: JSON.stringify(this.regInfo),
    })
    .then(function(response) {
      return response.json();
    })
    .then((data)=>{
      if(data.msg == "success"){
        this.usernameAvailability = false
      }
      else{
        this.usernameAvailability = true
      }
    }) 
    .catch(err => {
      console.log(err);
    });
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

  login(){

    // checks that all required fields and the captcha are filled before sending request to backend
    if(this.regInfo.userName == ""||this.regInfo.password == ""){
      this.regErr = "missing field";
      return
    }

    // send request to backend
    fetch("/api/login", {
      method: "post",
      body: JSON.stringify(this.regInfo),
    })
    .then(function(response) {
      return response.json();
    })
    .then((data)=>{
      if(data.msg == "success"){
        this.regErr = ""
      }
      else{
        this.regErr = "invalid login credentials"
      }
    }) 
    .catch(err => {
      console.log(err);
    });
  }

  constructor() { }

  ngOnInit(): void {
  }

}