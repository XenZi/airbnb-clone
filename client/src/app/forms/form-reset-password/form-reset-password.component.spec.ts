import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FormResetPasswordComponent } from './form-reset-password.component';

describe('FormResetPasswordComponent', () => {
  let component: FormResetPasswordComponent;
  let fixture: ComponentFixture<FormResetPasswordComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [FormResetPasswordComponent]
    });
    fixture = TestBed.createComponent(FormResetPasswordComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
