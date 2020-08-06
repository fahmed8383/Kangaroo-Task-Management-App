import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { NgCircleProgressModule } from 'ng-circle-progress';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { FormsModule } from '@angular/forms';
import { RecaptchaModule, RecaptchaFormsModule } from 'ng-recaptcha'; 
import { NavbarComponent } from './navbar/navbar.component';
import { HomeComponent } from './home/home.component';
import { PasswordResetComponent } from './password-reset/password-reset.component';
import { MainService } from './main.service';
import { MyplannerComponent } from './myplanner/myplanner.component';

@NgModule({
  declarations: [
    AppComponent,
    NavbarComponent,
    HomeComponent,
    PasswordResetComponent,
    MyplannerComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    FontAwesomeModule,
    FormsModule,
    RecaptchaModule,
    RecaptchaFormsModule,
    NgCircleProgressModule.forRoot({
      radius: 30,
      outerStrokeWidth: 4,
      showInnerStroke: false,
      outerStrokeColor: "#f7f7f7",
      unitsColor: "#f7f7f7",
      titleColor: "#f7f7f7", 
      responsive: true,
      animationDuration: 300,
    })
  ],
  providers: [MainService],
  bootstrap: [AppComponent]
})
export class AppModule { }
