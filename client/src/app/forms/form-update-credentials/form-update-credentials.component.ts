import { Component } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { AuthService } from 'src/app/services/auth-service/auth.service';
import { UserService } from 'src/app/services/user/user.service';
import { formatErrors } from 'src/app/utils/formatter.utils';

@Component({
  selector: 'app-form-update-credentials',
  templateUrl: './form-update-credentials.component.html',
  styleUrls: ['./form-update-credentials.component.scss'],
})
export class FormUpdateCredentialsComponent {
  formUpdateCredentials: FormGroup;
  errors: string[] = [];
  constructor(
    private formBuilder: FormBuilder,
    private authService: AuthService,
    private userService: UserService
  ) {
    this.formUpdateCredentials = this.formBuilder.group({
      email: ['', [Validators.required, Validators.email]],
      username: [
        '',
        [
          Validators.required,
          Validators.minLength(3),
          Validators.pattern('^[A-Za-z]+$'),
        ],
      ],
      password: ['', [Validators.required]],
    });
  }

  ngOnInit() {
    let foundUser = this.userService.getLoggedUser();
    if (foundUser != null) {
      this.formUpdateCredentials = this.formBuilder.group({
        email: [foundUser.email, [Validators.required, Validators.email]],
        username: [
          foundUser.username,
          [
            Validators.required,
            Validators.minLength(3),
            Validators.pattern('^[A-Za-z]+$'),
          ],
        ],
        password: ['', [Validators.required]],
      });
    }
  }

  onSubmit(e: Event) {
    e.preventDefault();
    this.errors = [];
    if (!this.formUpdateCredentials.valid) {
      Object.keys(this.formUpdateCredentials.controls).forEach((key) => {
        const controlErrors = this.formUpdateCredentials.get(key)?.errors;
        if (controlErrors) {
          this.errors.push(formatErrors(key));
        }
      });
      return;
    }
    this.authService.updateCredentials(
      this.formUpdateCredentials.value.email,
      this.formUpdateCredentials.value.username,
      this.formUpdateCredentials.value.password
    );
  }
}
