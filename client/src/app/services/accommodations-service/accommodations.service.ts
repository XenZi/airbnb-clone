import { Injectable } from '@angular/core';

import { LocalStorageService } from '../localstorage/local-storage.service';
import { HttpClient } from '@angular/common/http';
import { apiURL } from 'src/app/domains/constants';
import { ModalService } from '../modal/modal.service';
import { ToastService } from '../toast/toast.service';
import { ToastNotificationType } from 'src/app/domains/enums/toast-notification-type.enum';
import { Router } from '@angular/router';
import { Observable } from 'rxjs';
import { Accommodation } from 'src/app/domains/entity/accommodation-model';

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

  create(
    name: string,
    location: string,
    conveniences: string,
    minNumOfVisitors: string,
    maxNumOfVisitors: string
  ): void {
    this.http
      .post(`${apiURL}/accommodations/`, {
        name,
        location,
        conveniences,
        minNumOfVisitors,
        maxNumOfVisitors,
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

  public loadAccommodations(): Observable<Accommodation[]> {
    return this.http.get<Accommodation[]>(`${apiURL}/accommodations/`);
  }

  public getAccommodationById(id: string): Observable<Accommodation> {
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
        window.location.reload();
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
    location: string,
    conveniences: string,
    minNumOfVisitors: number,
    maxNumOfVisitors: number
  ): void {
    this.http
      .put(`${apiURL}/accommodations/${id}`, {
        name,
        location,
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
}