import { Component } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { AuthService } from 'src/app/services/auth-service/auth.service';

@Component({
  selector: 'app-form-request-reset-password',
  templateUrl: './form-request-reset-password.component.html',
  styleUrls: ['./form-request-reset-password.component.scss'],
})
export class FormRequestResetPasswordComponent {
  requestResetPassword: FormGroup;

  constructor(
    private formBuilder: FormBuilder,
    private authService: AuthService
  ) {
    this.requestResetPassword = this.formBuilder.group({
      email: ['', [Validators.required, Validators.email]],
    });
  }

  onSubmit(e: Event) {
    e.preventDefault();
    if (!this.requestResetPassword.valid) {
      return;
    }
    this.authService.requestPasswordReset(
      this.requestResetPassword.value.email
    );
  }
}
