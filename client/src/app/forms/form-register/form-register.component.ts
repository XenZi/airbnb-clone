import { Component } from '@angular/core';
import {
  AbstractControl,
  FormBuilder,
  FormGroup,
  ValidationErrors,
  ValidatorFn,
  Validators,
} from '@angular/forms';
import { Role } from 'src/app/domains/enums/roles.enum';
import { AuthService } from 'src/app/services/auth-service/auth.service';

import { ToastService } from 'src/app/services/toast/toast.service';
import { ToastNotificationType } from 'src/app/domains/enums/toast-notification-type.enum';
import { RecaptchaErrorParameters } from 'ng-recaptcha';
import { customPasswordStrengthValidator } from 'src/app/utils/validations.utils';
import { formatErrors } from 'src/app/utils/formatter.utils';

@Component({
  selector: 'app-form-register',
  templateUrl: './form-register.component.html',
  styleUrls: ['./form-register.component.scss'],
})
export class FormRegisterComponent {
  registerForm: FormGroup;
  roles: string[] = [Role.Guest, Role.Host];
  errors: string = '';
  isCaptchaValidated: boolean = false;
  constructor(
    private authService: AuthService,
    private formBuilder: FormBuilder,
    private toastService: ToastService
  ) {
    this.registerForm = this.formBuilder.group({
      email: ['', [Validators.required, Validators.email]],
      username: [
        '',
        [
          Validators.required,
          Validators.minLength(3),
          Validators.pattern('^[A-Za-z]+$'),
        ],
      ],
      firstName: ['', [Validators.required, Validators.minLength(2)]],
      lastName: ['', [Validators.required, Validators.minLength(2)]],
      currentPlace: ['', [Validators.required, Validators.minLength(2)]],
      password: ['', [Validators.required, customPasswordStrengthValidator()]],
      role: ['Guest', Validators.required],
    });
  }

  resolved(ecaptchaResponse: string) {
    this.isCaptchaValidated = true;
  }

  onErrorCaptcha(errorDetails: RecaptchaErrorParameters) {}
  onSubmit(e: Event) {
    e.preventDefault();
    if (!this.isCaptchaValidated) {
      this.toastService.showToast(
        'reCaptcha Error',
        'You must validate captcha first',
        ToastNotificationType.Error
      );
      return;
    }
    if (!this.registerForm.valid) {
      Object.keys(this.registerForm.controls).forEach((key) => {
        const controlErrors = this.registerForm.get(key)?.errors;
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
    this.authService.register(
      this.registerForm.value.email,
      this.registerForm.value.firstName,
      this.registerForm.value.lastName,
      this.registerForm.value.currentPlace,
      this.registerForm.value.password,
      this.registerForm.value.role,
      this.registerForm.value.username
    );
  }

  // customPasswordStrengthValidator(): ValidatorFn {
  //   return (control: AbstractControl): ValidationErrors | null => {
  //     const value = control.value;
  //     console.log(value);
  //     if (!value) {
  //       return null;
  //     }

  //     const hasUpperCase = /[A-Z]+/.test(value);
  //     const hasNumeric = /[0-9]+/.test(value);
  //     const hasSpecialChars = /[!@#$%^&*()_+\?><>:';\]\[']/.test(value);
  //     const minLength = (value as string).length >= 8;
  //     const passwordValid =
  //       hasUpperCase && hasNumeric && hasSpecialChars && minLength;
  //     console.log(hasUpperCase);
  //     console.log(hasNumeric);
  //     console.log(hasSpecialChars);
  //     console.log(minLength);
  //     return !passwordValid ? { passwordStrength: true } : null;
  //   };
  // }
}