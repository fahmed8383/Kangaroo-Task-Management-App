import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { HomeComponent } from './home/home.component';
import { PasswordResetComponent } from './password-reset/password-reset.component';
import { MyplannerComponent } from './myplanner/myplanner.component';


const routes: Routes = [
  { path: '', component: HomeComponent },
  { path: 'password-reset/:user/:token', component: PasswordResetComponent},
  { path: 'demo', component: MyplannerComponent},
  { path: 'my-planner', component: MyplannerComponent}
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
