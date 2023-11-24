import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FormRequestResetPasswordComponent } from './form-request-reset-password.component';

describe('FormRequestResetPasswordComponent', () => {
  let component: FormRequestResetPasswordComponent;
  let fixture: ComponentFixture<FormRequestResetPasswordComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [FormRequestResetPasswordComponent]
    });
    fixture = TestBed.createComponent(FormRequestResetPasswordComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
