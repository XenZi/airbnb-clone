import { Injectable } from '@angular/core';

import { LocalStorageService } from '../localstorage/local-storage.service';
import { HttpClient } from '@angular/common/http';
import { apiURL } from 'src/app/domains/constants';
import { ModalService } from '../modal/modal.service';
import { ToastService } from '../toast/toast.service';
import { ToastNotificationType } from 'src/app/domains/enums/toast-notification-type.enum';
import { Router } from '@angular/router';
import { Observable } from 'rxjs';
import { User } from 'src/app/domains/entity/user-profile.model';
import { Role } from 'src/app/domains/enums/roles.enum';
import { AuthService } from '../auth-service/auth.service';

@Injectable({
  providedIn: 'root'
})
export class ProfileService {

  constructor(
    private localStorageService: LocalStorageService,
    private http: HttpClient,
    private modalService: ModalService,
    private authService: AuthService,
    private toastSerice: ToastService,
    private router: Router
  ) {}

  create(
    id: string,
    firstName: string,
    lastName: string,
    email: string,
    residence: string,
    role: Role,
    username: string,
    age: number,
    
  ): void {
    this.http
      .post(`${apiURL}/users/`, {
        id,
        firstName,
        lastName,
        email,
        residence,
        role,
        username,
        age
        
      })
      .subscribe({
        next: (data) => {
          this.toastSerice.showToast(
            'Success',
            'Profile Information changed',
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

  delete(
    id:string
  ):void {
    console.log("delete sent to ", `${apiURL}/users/${id}`)
    this.http.delete(`${apiURL}/users/${id}`, {})
    .subscribe({
      next: (data) => {
        this.toastSerice.showToast(
          'Yay',
          'You have been Erased',
          ToastNotificationType.Success
        );
        this.router.navigate(['/'])
        this.authService.logout()
        
        
      },
      error: (err) => {
        this.toastSerice.showToast(
          'Error',
          err.error.error,
          ToastNotificationType.Error
        )
      }
    });
    
  }

  update(
    id: string,
    firstName: string,
    lastName: string,
    email: string,
    residence: string,
    role: Role,
    username: string,
    age: number,

  ): void {
    console.log({
      id,
      firstName,
      lastName,
      email,
      residence,
      role,
      username,
      age
    })
    this.http.put(`${apiURL}/users/${id}`, {
      id,
      firstName,
      lastName,
      email,
      residence,
      role,
      username,
      age
    }).subscribe({
      next: (data) => {
        this.toastSerice.showToast(
          'Congration',
           'you done it!',
            ToastNotificationType.Success
            );
            this.modalService.close();
            window.location.reload();
            this.router.navigate(['/profile', id])      
      },
      error: (err) => {
        this.toastSerice.showToast(
          'Uh Uh Uh, you didnt say the magic word',
          err.error.error,
          ToastNotificationType.Error
        )
      }
    })
  }

  

  public getUserById(id: string): Observable<any> {
    return this.http.get<any>(`${apiURL}/users/${id}`);
  }

}