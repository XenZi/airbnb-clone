import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FormRateAccommodationComponent } from './form-rate-accommodation.component';

describe('FormRateAccommodationComponent', () => {
  let component: FormRateAccommodationComponent;
  let fixture: ComponentFixture<FormRateAccommodationComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [FormRateAccommodationComponent]
    });
    fixture = TestBed.createComponent(FormRateAccommodationComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
