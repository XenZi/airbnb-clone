import { Injectable } from '@angular/core';
import { LocalStorageService } from '../localstorage/local-storage.service';
import { HttpClient } from '@angular/common/http';
import { apiURL } from 'src/app/domains/constants';
import { ModalService } from '../modal/modal.service';
import { ToastService } from '../toast/toast.service';
import { ToastNotificationType } from 'src/app/domains/enums/toast-notification-type.enum';
import { Router } from '@angular/router';
import { UserService } from '../user/user.service';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  constructor(
    private localStorageService: LocalStorageService,
    private http: HttpClient,
    private modalService: ModalService,
    private toastSerice: ToastService,
    private router: Router,
    private userService: UserService
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
          console.log(data);
          this.localStorageService.setItem('token', data.data?.token);
          this.localStorageService.setItem(
            'user',
            JSON.stringify(data.data?.user)
          );
          setTimeout(() => {
            this.toastSerice.showToast(
              'You have successfully logged in',
              'You made it',
              ToastNotificationType.Success
            );
            this.modalService.close();
          }, 1000);
          window.location.reload();
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
    username: string,
    age: number
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
        age,
      })
      .subscribe({
        next: (data) => {
          this.toastSerice.showToast(
            'You have successfully registered',
            'You can expect mail for confirmation',
            ToastNotificationType.Success
          );
          this.modalService.close();
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
            'You have successfully requested password reset',
            'Check out email',
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
            'You have successfully reset your password',
            'You have successfully reset your password',
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

  changePassword(
    oldPassword: string,
    password: string,
    confirmedPassword: string
  ) {
    this.http
      .post(`${apiURL}/auth/change-password`, {
        oldPassword,
        password,
        confirmedPassword,
      })
      .subscribe({
        next: (data) => {
          this.toastSerice.showToast(
            'You have successfully changed your password',
            'You have successfully changed your password',
            ToastNotificationType.Success
          );
          this.localStorageService.clear();
          this.router.navigate(['/']);
        },
        error: (err) => {
          this.toastSerice.showToast(
            'Error',
            err.error.error,
            ToastNotificationType.Error
          );
          if (err.error.status == 401) {
            this.localStorageService.clear();
            this.router.navigate(['/']);
          }
        },
      });
  }

  logout(): void {
    this.localStorageService.clear();
    window.location.reload();
    setTimeout(() => {
      this.toastSerice.showToast(
        "You're logged out",
        "You're logged out",
        ToastNotificationType.Info
      );
    }, 1000);
  }

  updateCredentials(email: string, username: string, password: string): void {
    this.http
      .post(`${apiURL}/auth/update-credentials`, {
        email,
        username,
        password,
      })
      .subscribe({
        next: (data: any) => {
          this.toastSerice.showToast(
            'You made it',
            data?.data?.message,
            ToastNotificationType.Success
          );
        },
        error: (err) => {
          this.toastSerice.showToast(
            'Error',
            err?.error?.error,
            ToastNotificationType.Error
          );
          console.log(err);
        },
      });
  }
}
