import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { UserService } from '../user/user.service';
import { Observable } from 'rxjs';
import { apiURL } from 'src/app/domains/constants';
import {
  AccommodationRate,
  HostRate,
} from 'src/app/domains/entity/ratings.model';

@Injectable({
  providedIn: 'root',
})
export class RatingService {
  constructor(private http: HttpClient, private userService: UserService) {}

  getAllRatingsForAccommodation(accommodationID: string): Observable<any> {
    return this.http.get<any>(
      `${apiURL}/recommendations/rating/accommodation/${accommodationID}`
    );
  }

  rateAccommodation(rateAccommodation: AccommodationRate) {
    this.http
      .post<AccommodationRate>(
        `${apiURL}/recommendations/rating/accommodation`,
        rateAccommodation
      )
      .subscribe({
        next: (data) => {
          console.log(data);
        },
        error: (err) => {
          console.log(err);
        },
      });
  }

  getRatingForAccommodationByGuest(
    accommodationID: string,
    guestID: string
  ): Observable<any> {
    return this.http.get(
      `${apiURL}/recommendations/rating/${accommodationID}/${guestID}`
    );
  }

  updateRateForAccommodation(rateAccommodation: AccommodationRate) {
    this.http
      .put<AccommodationRate>(
        `${apiURL}/recommendations/rating/accommodation`,
        rateAccommodation
      )
      .subscribe({
        next: (data) => {
          console.log(data);
        },
        error: (err) => {
          console.log(err);
        },
      });
  }

  deleteRateForAccomodation(accommodationID: string, guestID: string) {
    this.http
      .delete(
        `${apiURL}/recommendations/rating/accommodation/${accommodationID}/${guestID}`
      )
      .subscribe({
        next: (data) => {
          console.log(data);
        },
        error: (err) => {
          console.log(err);
        },
      });
  }

  getAllRatingsForHost(hostID: string): Observable<any> {
    return this.http.get<any>(
      `${apiURL}/recommendations/rating/host/${hostID}`
    );
  }
}
