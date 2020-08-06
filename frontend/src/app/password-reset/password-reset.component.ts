import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';


@Component({
  selector: 'app-password-reset',
  templateUrl: './password-reset.component.html',
  styleUrls: ['./password-reset.component.scss']
})
export class PasswordResetComponent implements OnInit {

//// variables

  // holds the password reset info
  resetInfo = {
    userName: "",
    token: "",
    password: "",
    confirmPassword: "",
  }

  // holds the response from api call
  apiResp: string = "";

  // checks if there are any errors in the api call
  apiErr: string = "";

//// functions

  resetPassword(){

    // checks that all required fields are filled before sending request to backend
    if(this.resetInfo.password == "" || this.resetInfo.userName == "" || this.resetInfo.token == ""){
      this.apiErr = "missing field";
      return
    }

    // send request to backend
    fetch("/api/reset-user-password", {
      method: "post",
      body: JSON.stringify(this.resetInfo),
    })
    .then(function(response) {
      return response.json();
    })
    .then((data)=>{
      if(data.msg == "success"){
        this.router.navigate(['/']);
      }
      else{
        this.apiErr = "invalid request"
        this.router.navigate(['/']);
      }
    }) 
    .catch(err => {
      console.log(err);
    });
  }

  constructor(private route: ActivatedRoute, private router: Router) { }

  ngOnInit(): void {
    // get the variables from the url
    this.resetInfo.userName = this.route.snapshot.paramMap.get('user');
    this.resetInfo.token = this.route.snapshot.paramMap.get('token'); 
  }

}
