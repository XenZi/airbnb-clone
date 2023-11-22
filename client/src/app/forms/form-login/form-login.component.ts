import { Component } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { AuthService } from 'src/app/services/auth-service/auth.service';
import { ModalService } from 'src/app/services/modal/modal.service';
import { FormRequestResetPasswordComponent } from '../form-request-reset-password/form-request-reset-password.component';

@Component({
  selector: 'app-form-login',
  templateUrl: './form-login.component.html',
  styleUrls: ['./form-login.component.scss'],
})
export class FormLoginComponent {
  loginForm: FormGroup;
  constructor(
    private authService: AuthService,
    private formBuilder: FormBuilder,
    private modalService: ModalService
  ) {
    this.loginForm = this.formBuilder.group({
      email: ['', [Validators.required, Validators.email]],
      password: ['', [Validators.required]],
    });
  }

  onSubmit(e: Event) {
    e.preventDefault();
    if (!this.loginForm.valid) {
      return;
    }
    this.authService.login(
      this.loginForm.value.email,
      this.loginForm.value.password
    );
  }

  forgotPasswordClick() {
    this.modalService.close();
    this.modalService.open(
      FormRequestResetPasswordComponent,
      'Request Password Reset'
    );
  }
}
