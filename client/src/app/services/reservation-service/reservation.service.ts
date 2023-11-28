import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, catchError, tap } from 'rxjs';
import { apiURL } from 'src/app/domains/constants';
import { ToastService } from '../toast/toast.service';
import { Router } from '@angular/router';
import { ToastNotificationType } from 'src/app/domains/enums/toast-notification-type.enum';
import { UserService } from '../user/user.service';
import { UserAuth } from 'src/app/domains/entity/user-auth.model';

@Injectable({
  providedIn: 'root'
})
export class ReservationService {
  private apiURL = apiURL;
 

  constructor(
    private http: HttpClient,
    private toastService: ToastService,
    private router: Router,
    private userService: UserService 
  ) {}

  createReservation(reservationData: any): Observable<any> {
    return this.http.post(`https://localhost:443/api/reservations/`, reservationData)
      .pipe(
        tap((data) => {
          this.toastService.showToast(
            'Success',
            'Reservation created!',
            ToastNotificationType.Success
          );
        }),
        catchError((err) => {
          this.toastService.showToast(
            'Error',
            err.error.error,
            ToastNotificationType.Error
          );
          throw err;
        })
      );
  }

  getLoggedUserId(): string | null {
    const loggedUser = this.userService.getLoggedUser();
  
    if (loggedUser) {
      return loggedUser.id; 
    } else {
    
      return null;
    }
  }

  deleteById(
    id:string,
    
    
  ): void {
    this.http
      .delete(`${apiURL}${this.getLoggedUserId}/reservations/${id}`, {
       
      })
      .subscribe({
        next: (data) => {
          this.toastService.showToast(
            'Success',
            'Reservation deleted!',
            ToastNotificationType.Success
            
          );
          
          
          
        },
        error: (err) => {
          this.toastService.showToast(
            'Error',
            err.error.error,
            ToastNotificationType.Error
            );
          },
        });
        this.router.navigate(['/']);
  }
  
}
