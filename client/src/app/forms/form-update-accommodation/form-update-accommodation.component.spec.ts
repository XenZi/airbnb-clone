import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FormUpdateAccommodationComponent } from './form-update-accommodation.component';

describe('FormUpdateAccommodationComponent', () => {
  let component: FormUpdateAccommodationComponent;
  let fixture: ComponentFixture<FormUpdateAccommodationComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [FormUpdateAccommodationComponent]
    });
    fixture = TestBed.createComponent(FormUpdateAccommodationComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
