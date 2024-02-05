import { Component } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute } from '@angular/router';
import { UserAuth } from 'src/app/domains/entity/user-auth.model';
import { ToastNotificationType } from 'src/app/domains/enums/toast-notification-type.enum';
import { AuthService } from 'src/app/services/auth-service/auth.service';
import { ToastService } from 'src/app/services/toast/toast.service';
import { UserService } from 'src/app/services/user/user.service';
import { formatErrors } from 'src/app/utils/formatter.utils';
import { customPasswordStrengthValidator } from 'src/app/utils/validations.utils';

@Component({
  selector: 'app-form-change-password',
  templateUrl: './form-change-password.component.html',
  styleUrls: ['./form-change-password.component.scss'],
})
export class FormChangePasswordComponent {
  changePasswordForm: FormGroup;
  errors: string[] = [];
  userID!: string;
  constructor(
    private formBuilder: FormBuilder,
    private toastService: ToastService,
    private authService: AuthService,
    private route: ActivatedRoute,
    private userService: UserService,
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

  ngOnInit() {
    this.userID = (this.userService.getLoggedUser() as UserAuth).id;
  }
  onSubmit(e: Event) {
    e.preventDefault();
    this.errors = [];
    if (!this.changePasswordForm.valid) {
      Object.keys(this.changePasswordForm.controls).forEach((key) => {
        const controlErrors = this.changePasswordForm.get(key)?.errors;
        if (controlErrors) {
          this.errors.push(formatErrors(key));
        }
      });
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
      this.changePasswordForm.value.confirmedPassword,
      this.userID
    );
  }
}
