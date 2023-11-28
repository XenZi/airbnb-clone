import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FormCreateReservationComponent } from './form-create-reservation.component';

describe('FormCreateReservationComponent', () => {
  let component: FormCreateReservationComponent;
  let fixture: ComponentFixture<FormCreateReservationComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [FormCreateReservationComponent]
    });
    fixture = TestBed.createComponent(FormCreateReservationComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
