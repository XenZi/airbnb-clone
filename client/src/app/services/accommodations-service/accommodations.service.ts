import { Injectable } from '@angular/core';

import { LocalStorageService } from '../localstorage/local-storage.service';
import { HttpClient } from '@angular/common/http';
import { apiURL } from 'src/app/domains/constants';
import { ModalService } from '../modal/modal.service';
import { ToastService } from '../toast/toast.service';
import { ToastNotificationType } from 'src/app/domains/enums/toast-notification-type.enum';
import { Router } from '@angular/router';
import { Observable, map } from 'rxjs';
import { Accommodation } from 'src/app/domains/entity/accommodation-model';
import { DateAvailability } from 'src/app/domains/entity/date-availability.model';

@Injectable({
  providedIn: 'root',
})
export class AccommodationsService {
  constructor(
    private localStorageService: LocalStorageService,
    private http: HttpClient,
    private modalService: ModalService,
    private toastSerice: ToastService,
    private router: Router
  ) {}

  username = localStorage.getItem('username');

  create(
    userId: string,
    username: string,
    name: string,
    address: string,
    city: string,
    country: string,
    conveniences: string[],
    minNumOfVisitors: number,
    maxNumOfVisitors: number,
    availableAccommodationDates: DateAvailability[],
    location: string
  ): void {
    this.http
      .post(`${apiURL}/accommodations/`, {
        userId,
        username,
        name,
        address,
        city,
        country,
        conveniences,
        minNumOfVisitors: Number(minNumOfVisitors),
        maxNumOfVisitors: Number(maxNumOfVisitors),
        availableAccommodationDates,
        location,
      })
      .subscribe({
        next: (data) => {
          this.toastSerice.showToast(
            'Success',
            'Accommodation created!',
            ToastNotificationType.Success
          );
          this.modalService.close();
          this.router.navigate(['/']);
          window.location.reload();
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
    // this.router.navigate(['/']);
    // window.location.reload();
  }

  public loadAccommodations(): Observable<any> {
    return this.http.get<any>(`${apiURL}/accommodations/`);
  }
  public getAccommodationById(id: string): Observable<any> {
    return this.http.get<Accommodation>(`${apiURL}/accommodations/${id}`);
  }

  deleteById(id: string): void {
    this.http.delete(`${apiURL}/accommodations/${id}`, {}).subscribe({
      next: (data) => {
        this.toastSerice.showToast(
          'Success',
          'Accommodation updated!',
          ToastNotificationType.Success
        );
        this.router.navigate(['/']);
        // window.location.reload();
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

  update(
    id: string,
    name: string,
    address: string,
    city: string,
    conveniences: string[],
    minNumOfVisitors: number,
    maxNumOfVisitors: number
  ): void {
    this.http
      .put(`${apiURL}/accommodations/${id}`, {
        name,
        address,
        city,
        conveniences,
        minNumOfVisitors,
        maxNumOfVisitors,
      })
      .subscribe({
        next: (data) => {
          this.toastSerice.showToast(
            'Success',
            'Accommodation updated!',
            ToastNotificationType.Success
          );
          this.modalService.close();
          window.location.reload();
          this.router.navigate(['/accommodations', id]);
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

  public search(
    city: string,
    country: string,
    numOfVisitors: string,
    startDate: string,
    endDate: string
  ): Observable<any> {
    this.router.navigate(['/search'], {
      queryParams: {
        city: city,
        country: country,
        numOfVisitors: numOfVisitors,
        startDate: startDate,
        endDate: endDate,
      },
    });
    console.log('pocetni datum je', startDate);
    return this.http.get<any>(
      `${apiURL}/accommodations/search?city=${city}&country=${country}&numOfVisitors=${numOfVisitors}&startDate=${startDate}&endDate=${endDate}`
    );
  }
}
