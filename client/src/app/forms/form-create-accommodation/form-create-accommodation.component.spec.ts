import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FormCreateAccommodationComponent } from './form-create-accommodation.component';

describe('FormCreateAccommodationComponent', () => {
  let component: FormCreateAccommodationComponent;
  let fixture: ComponentFixture<FormCreateAccommodationComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [FormCreateAccommodationComponent]
    });
    fixture = TestBed.createComponent(FormCreateAccommodationComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
