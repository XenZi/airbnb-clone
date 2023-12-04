import { Component } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute } from '@angular/router';
import { ToastNotificationType } from 'src/app/domains/enums/toast-notification-type.enum';
import { AuthService } from 'src/app/services/auth-service/auth.service';
import { ToastService } from 'src/app/services/toast/toast.service';
import { formatErrors } from 'src/app/utils/formatter.utils';
import { customPasswordStrengthValidator } from 'src/app/utils/validations.utils';

@Component({
  selector: 'app-form-reset-password',
  templateUrl: './form-reset-password.component.html',
  styleUrls: ['./form-reset-password.component.scss'],
})
export class FormResetPasswordComponent {
  resetPasswordForm: FormGroup;
  errors: string[] = [];
  token!: string;
  constructor(
    private formBuilder: FormBuilder,
    private toastService: ToastService,
    private authService: AuthService,
    private route: ActivatedRoute
  ) {
    this.resetPasswordForm = this.formBuilder.group({
      password: ['', [Validators.required, customPasswordStrengthValidator()]],
      confirmedPassword: [
        '',
        [Validators.required, customPasswordStrengthValidator()],
      ],
    });
  }

  ngOnInit() {
    this.route.paramMap.subscribe((params) => {
      this.token = String(params.get('token')) + '=';
    });
  }

  onSubmit(e: Event) {
    e.preventDefault();
    this.errors = [];
    if (!this.resetPasswordForm.valid) {
      Object.keys(this.resetPasswordForm.controls).forEach((key) => {
        const controlErrors = this.resetPasswordForm.get(key)?.errors;
        if (controlErrors) {
          this.errors.push(formatErrors(key));
        }
      });
      return;
    }
    if (
      this.resetPasswordForm.value.password !==
      this.resetPasswordForm.value.confirmedPassword
    ) {
      this.toastService.showToast(
        'Error',
        'Password and confirmed password are not the same',
        ToastNotificationType.Error
      );
      return;
    }
    this.authService.resetPassword(
      this.resetPasswordForm.value.password,
      this.resetPasswordForm.value.confirmedPassword,
      this.token
    );
  }
}
