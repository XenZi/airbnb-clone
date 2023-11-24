import { Component } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute } from '@angular/router';
import { ToastNotificationType } from 'src/app/domains/enums/toast-notification-type.enum';
import { AuthService } from 'src/app/services/auth-service/auth.service';
import { ToastService } from 'src/app/services/toast/toast.service';
import { formatErrors } from 'src/app/utils/formatter.utils';
import { customPasswordStrengthValidator } from 'src/app/utils/validations.utils';

@Component({
  selector: 'app-form-change-password',
  templateUrl: './form-change-password.component.html',
  styleUrls: ['./form-change-password.component.scss'],
})
export class FormChangePasswordComponent {
  changePasswordForm: FormGroup;
  errors: string = '';
  constructor(
    private formBuilder: FormBuilder,
    private toastService: ToastService,
    private authService: AuthService,
    private route: ActivatedRoute
  ) {
    this.changePasswordForm = this.formBuilder.group({
      oldPassword: [
        '',
        [Validators.required, customPasswordStrengthValidator()],
      ],
      password: ['', [Validators.required, customPasswordStrengthValidator()]],
      confirmedPassword: [
        '',
        [Validators.required, customPasswordStrengthValidator()],
      ],
    });
  }

  onSubmit(e: Event) {
    e.preventDefault();
    if (!this.changePasswordForm.valid) {
      Object.keys(this.changePasswordForm.controls).forEach((key) => {
        const controlErrors = this.changePasswordForm.get(key)?.errors;
        if (controlErrors) {
          this.errors += formatErrors(key);
        }
      });
      this.toastService.showToast(
        'Error',
        this.errors,
        ToastNotificationType.Error
      );
      this.errors = '';
      return;
    }
    if (
      this.changePasswordForm.value.password !==
      this.changePasswordForm.value.confirmedPassword
    ) {
      this.toastService.showToast(
        'Error',
        'Password and confirmed password are not the same',
        ToastNotificationType.Error
      );
      return;
    }
    this.authService.changePassword(
      this.changePasswordForm.value.oldPassword,
      this.changePasswordForm.value.password,
      this.changePasswordForm.value.confirmedPassword
    );
  }
}
