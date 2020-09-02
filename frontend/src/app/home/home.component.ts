import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})
export class HomeComponent implements OnInit {

//// variables

  // holds all carousel images and captions
  carousel = [{img: "../../assets/images/1.PNG", caption: "Organize by category based cards"}, {img: "../../assets/images/2.PNG", caption: "Add tasks and subtasks"}, {img: "../../assets/images/3.PNG", caption: "Track progress"}, {img: "../../assets/images/4.PNG", caption: "Arrange as you please"}];

  // current slide index, changed when user clicks on next or prev
  currentSlide = 0;

//// functions

  // decrease carousel index by 1
  prevSlide(){
    // if it is on the first slide, move to the last slide in the carousel
    if(this.currentSlide == 0){
      this.currentSlide = this.carousel.length-1;
      return
    }

    // else decrease by 1
    this.currentSlide--;
  }

  // increase carousel index by 1
  nextSlide(){
    // if it is on the last slide, move to the first slide in the carousel
    if(this.currentSlide == this.carousel.length-1){
      this.currentSlide = 0;
      return
    }

    // else increase by 1
    this.currentSlide++;
  }

  constructor() { }

  ngOnInit(): void {
  }

}
