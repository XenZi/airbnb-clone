import { Injectable } from '@angular/core';
import { LocalStorageService } from '../localstorage/local-storage.service';
import { HttpClient } from '@angular/common/http';
import { apiURL } from 'src/app/domains/constants';
import { ModalService } from '../modal/modal.service';
import { ToastService } from '../toast/toast.service';
import { ToastNotificationType } from 'src/app/domains/enums/toast-notification-type.enum';
import { Router } from '@angular/router';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  constructor(
    private localStorageService: LocalStorageService,
    private http: HttpClient,
    private modalService: ModalService,
    private toastSerice: ToastService,
    private router: Router
  ) {}

  login(email: string, password: string): void {
    this.http
      .post(
        `${apiURL}/auth/login`,
        {
          email,
          password,
        },
        {
          headers: {
            'Content-Type': 'application/json',
          },
        }
      )
      .subscribe({
        next: (data: any) => {
          this.localStorageService.setItem('token', data.data?.Token);
          this.localStorageService.setItem(
            'user',
            JSON.stringify(data.data?.User)
          );
          this.toastSerice.showToast(
            'You have successfully logged in',
            'You made it',
            ToastNotificationType.Success
          );
          this.modalService.close();
        },
        error: (err) => {
          console.log(err);
          this.toastSerice.showToast(
            'Error',
            err.error.error,
            ToastNotificationType.Error
          );
          this.modalService.close();
        },
      });
  }

  register(
    email: string,
    firstName: string,
    lastName: string,
    currentPlace: string,
    password: string,
    role: string,
    username: string
  ): void {
    this.http
      .post(`${apiURL}/auth/register`, {
        email,
        firstName,
        lastName,
        currentPlace,
        password,
        role,
        username,
      })
      .subscribe({
        next: (data) => {
          this.toastSerice.showToast(
            'You have successfully registered',
            'You can expect mail for confirmation',
            ToastNotificationType.Success
          );
        },
        error: (err) => {
          this.toastSerice.showToast(
            'Error',
            err.error.error,
            ToastNotificationType.Error
          );
        },
      });
  }

  confirmAccount(token: string): void {
    this.http.post(`${apiURL}/auth/confirm-account/${token}`, {}).subscribe({
      next: (data) => {
        this.toastSerice.showToast(
          'You have successfully confirmed your account',
          'You have successfully confirmed your account',
          ToastNotificationType.Success
        );
        this.router.navigate(['/']);
      },
      error: (err) => {
        this.toastSerice.showToast(
          'Error',
          err.error.error,
          ToastNotificationType.Error
        );
      },
    });
  }

  requestPasswordReset(email: string) {
    this.http
      .post(`${apiURL}/auth/request-reset-password`, {
        email,
      })
      .subscribe({
        next: (data) => {
          console.log(data);
          this.toastSerice.showToast(
            'You have successfully confirmed your account',
            'You have successfully confirmed your account',
            ToastNotificationType.Success
          );
          this.router.navigate(['/']);
        },
        error: (err) => {
          console.log(err);
          this.toastSerice.showToast(
            'Error',
            err.error.error,
            ToastNotificationType.Error
          );
        },
      });
  }

  resetPassword(
    password: string,
    confirmedPassword: string,
    token: string
  ): void {
    this.http
      .post(`${apiURL}/auth/reset-password/${token}`, {
        password,
        confirmedPassword,
      })
      .subscribe({
        next: (data) => {
          this.toastSerice.showToast(
            'You have successfully reseted your password',
            'You have successfully reseted your password',
            ToastNotificationType.Success
          );
          this.router.navigate(['/']);
        },
        error: (err) => {
          console.log(err);
          this.toastSerice.showToast(
            'Error',
            err.error.error,
            ToastNotificationType.Error
          );
        },
      });
  }
}
