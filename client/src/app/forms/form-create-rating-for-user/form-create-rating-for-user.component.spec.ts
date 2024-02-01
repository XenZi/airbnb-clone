import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FormCreateRatingForUserComponent } from './form-create-rating-for-user.component';

describe('FormCreateRatingForUserComponent', () => {
  let component: FormCreateRatingForUserComponent;
  let fixture: ComponentFixture<FormCreateRatingForUserComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [FormCreateRatingForUserComponent]
    });
    fixture = TestBed.createComponent(FormCreateRatingForUserComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
