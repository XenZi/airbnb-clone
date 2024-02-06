import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FormUpdateAvailabilityComponent } from './form-update-availability.component';

describe('FormUpdateAvailabilityComponent', () => {
  let component: FormUpdateAvailabilityComponent;
  let fixture: ComponentFixture<FormUpdateAvailabilityComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [FormUpdateAvailabilityComponent]
    });
    fixture = TestBed.createComponent(FormUpdateAvailabilityComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
