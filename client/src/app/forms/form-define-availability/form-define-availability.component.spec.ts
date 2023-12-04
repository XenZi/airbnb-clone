import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FormDefineAvailabilityComponent } from './form-define-availability.component';

describe('FormDefineAvailabilityComponent', () => {
  let component: FormDefineAvailabilityComponent;
  let fixture: ComponentFixture<FormDefineAvailabilityComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [FormDefineAvailabilityComponent]
    });
    fixture = TestBed.createComponent(FormDefineAvailabilityComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
