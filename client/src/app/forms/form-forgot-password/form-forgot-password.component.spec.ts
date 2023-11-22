import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FormForgotPasswordComponent } from './form-forgot-password.component';

describe('FormForgotPasswordComponent', () => {
  let component: FormForgotPasswordComponent;
  let fixture: ComponentFixture<FormForgotPasswordComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [FormForgotPasswordComponent]
    });
    fixture = TestBed.createComponent(FormForgotPasswordComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
