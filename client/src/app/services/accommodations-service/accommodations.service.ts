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
  providedIn: 'root'
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
    maxNumOfVisitors: string,
    
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

  public loadAccommodations():Observable<Accommodation[]>{
    return this.http.get<Accommodation[]>(`${apiURL}/accommodations/`)
  }

  public getAccommodationById(id: string): Observable<Accommodation> {
    return this.http.get<Accommodation>(`${apiURL}/accommodations/${id}`);
  }

}
