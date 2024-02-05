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
  providedIn: 'root',
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
    return this.http.get(url);
  }

  createReservation(reservationData: any): void {
    this.http.post(`${apiURL}/reservations/`, reservationData).subscribe(
      (data) => {
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

  deleteById(
    country: string,
    id: string,
    userID: string,
    hostID: string,
    accommodationID: string,
    endDate: string
  ): Observable<any> {
    const url = `${this.apiURL}/reservations/${country}/${id}/${userID}/${hostID}/${accommodationID}/${endDate}`;
    return this.http.delete(url);
  }

  getAllReservationsById(id: string): Observable<any> {
    return this.http.get(`${apiURL}/reservations/user/guest/${id}`);
  }


  getPastReservations(
    accommodationID: string,
    userID: string
  ): Observable<any> {
    return this.http.get(`${apiURL}/reservations/${accommodationID}/${userID}`);
  }

  getPastReservationsByUserForGuest(
    hostID: string,
    userID: string
  ): Observable<any> {
    console.log(hostID);
    console.log(userID);
    return this.http.get(`${apiURL}/reservations/host/${hostID}/${userID}`);
  }
  getAllReservationsByHost(hostId: string): Observable<any> {
    return this.http.get(`${apiURL}/reservations/user/host/${hostId}`);
  }

 
}
