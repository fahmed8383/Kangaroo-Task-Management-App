import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { MyplannerComponent } from './myplanner.component';

describe('MyplannerComponent', () => {
  let component: MyplannerComponent;
  let fixture: ComponentFixture<MyplannerComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ MyplannerComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(MyplannerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
