import { Component } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';

@Component({
  selector: 'app-form-forgot-password',
  templateUrl: './form-forgot-password.component.html',
  styleUrls: ['./form-forgot-password.component.scss'],
})
export class FormForgotPasswordComponent {
  resetPasswordForm: FormGroup;

  constructor(private formBuilder: FormBuilder) {
    this.resetPasswordForm = this.formBuilder.group({
      password: ['', [Validators.required]],
      confirmedPassword: ['', [Validators.required]],
    });
  }
}
