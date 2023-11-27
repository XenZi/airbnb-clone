import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { apiURL } from 'src/app/domains/constants';

@Injectable({
  providedIn: 'root'
})
export class ReservationService {
  private apiURL = apiURL;

  constructor(private http: HttpClient) {}

  createReservation(reservationData: any): Observable<any> {
    return this.http.post(`${this.apiURL}/reservations`, reservationData);
  }
}