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

  //set imported icon names for component
  faTimes = faTimes;
  faBars = faBars;

  //set base mobile nav to false
  showMobileNav: boolean = false;

  //determines the transparency of the navbar
  navbarTransparent: boolean = true;

  //determines whether the modals connected to the navbar should be displayed or not
  displayRegModal: boolean = false;

  regInfo = {
    userName: "",
    email: "",
    password: "",
    confirmPassword: "" 
  }

//// functions

  //change mobilenav value to the opposite on click
  mobileNavClick(){
    this.showMobileNav = !this.showMobileNav;
  }

  //on scroll listener to change the navabar color after a certain threshold
  @HostListener("window:scroll", [])
  onWindowScroll() {

    const number = window.pageYOffset || document.documentElement.scrollTop || document.body.scrollTop || 0;
    if (number > 100) {
      this.navbarTransparent = false;
    } else {
      this.navbarTransparent = true;
    }

  }

  //displays the modals and turns scrolling off
  openRegModal(){
    this.displayRegModal = true;
    document.body.style.overflow = 'hidden';
  }

  //closes the modals if thje correct element is clicked and enables scrolling again
  closeRegModal(e){
    if(e.target.className == "modal" || e.target.className == "close"){
      this.displayRegModal = false;
      document.body.style.overflow = 'auto';
    }
  }

  sendEmailVerification(){
    console.log(this.regInfo)
  }

  constructor() { }

  ngOnInit(): void {
  }

}

/*
import { Component, Renderer2, ElementRef, ViewChild, HostListener, OnInit, AfterViewInit } from '@angular/core';
import { faTimes } from "@fortawesome/free-solid-svg-icons";
import { faBars } from "@fortawesome/free-solid-svg-icons";

@Component({
  selector: 'app-navbar',
  templateUrl: './navbar.component.html',
  styleUrls: ['./navbar.component.scss']
})
export class NavbarComponent implements OnInit, AfterViewInit {

//// variables

  //set imported icon names for component
  faTimes = faTimes;
  faBars = faBars;

  //set base mobile nav to false
  showMobileNav: boolean = false;

  //determines the transparency of the navbar
  navbarTransparent: boolean = true;

  //determines whether the modals connected to the navbar should be displayed or not
  displayRegModal: boolean = false;

//// functions

  //change mobilenav value to the opposite on click
  mobileNavClick(){
    this.showMobileNav = !this.showMobileNav;
  }

  //on scroll listener to change the navabar color after a certain threshold
  @HostListener("window:scroll", [])
  onWindowScroll() {

    const number = window.pageYOffset || document.documentElement.scrollTop || document.body.scrollTop || 0;
    if (number > 100) {
      this.navbarTransparent = false;
    } else {
      this.navbarTransparent = true;
    }

  }

  //functions responsible for opening and closing the modal
  openRegModal(){
    this.displayRegModal = true;
  }

  //get element with the following directives in the dom
  @ViewChild('regModal') regModal: ElementRef;

  //add renderer to constructor 
  constructor(private renderer: Renderer2) {
  }

  ngOnInit(): void {
  }

  //these functions witll run after view has be initialized
  ngAfterViewInit(){
    //renderer will listen for click event
    this.renderer.listen('window', 'click',(e:Event)=>{
      if(e.target == this.regModal.nativeElement){
        this.displayRegModal = false;
      }
    });
  }

}
*/
