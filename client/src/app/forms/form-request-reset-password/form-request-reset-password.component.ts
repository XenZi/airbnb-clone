import { Component } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { UserAuth } from 'src/app/domains/entity/user-auth.model';
import { AuthService } from 'src/app/services/auth-service/auth.service';
import { UserService } from 'src/app/services/user/user.service';

@Component({
  selector: 'app-form-request-reset-password',
  templateUrl: './form-request-reset-password.component.html',
  styleUrls: ['./form-request-reset-password.component.scss'],
})
export class FormRequestResetPasswordComponent {
  requestResetPassword: FormGroup;
  constructor(
    private formBuilder: FormBuilder,
    private authService: AuthService,
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
      this.requestResetPassword.value.email,
    );
  }
}
