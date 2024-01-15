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
  
  getAvailability(accommodationID: string): Observable<any> {
    const url = `${this.apiURL}/reservations/${accommodationID}/availability`;
    console.log(url)
    return this.http.get(url);
  }

  createReservation(reservationData: any): void {
    console.log("USLO");
    this.http.post(`${apiURL}/reservations/`, reservationData)
      .subscribe(
        (data) => {
          console.log("USLO1");
          this.toastService.showToast(
            'Success',
            'Reservation created!',
            ToastNotificationType.Success
          );
        },
        (err) => {
          this.toastService.showToast(
            'Error',
            err.error.error,
            ToastNotificationType.Error
          );
          console.error(err); // Log the error for debugging
        }
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

  deleteById(country: string, id: string, userID: string, hostID: string, accommodationID: string): Observable<any> {
    const url = `${this.apiURL}/${country}/${id}/${userID}/${hostID}/${accommodationID}`;
    return this.http.delete(url);
  }
  
  getAllReservationsById(id: string): Observable<any> {
    return this.http.get(`${apiURL}/reservations/user/guest/${id}`)
  }
}
