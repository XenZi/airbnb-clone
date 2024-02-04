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
import { FormArray } from '@angular/forms';

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

 

  getCountries(): Observable<any[]> {
    return this.http.get<any[]>('assets/countries.json');
  }



  private formatFormData(
    userId: string,
    username: string,
    email:string,
    name: string,
    address: string,
    city: string,
    country: string,
    conveniences: string[],
    minNumOfVisitors: string,
    maxNumOfVisitors: string,
    availableAccommodationDates: DateAvailability[],
    location: string,
    images:FormArray,
    paying:string
  ): FormData {
    let formData: FormData = new FormData();
    formData.append('userId', userId);
    formData.append('username', username);
    formData.append('email', email);
    formData.append('name', name);
    formData.append('address', address);
    formData.append('city', city);
    formData.append('country', country);
    conveniences.forEach(conv => formData.append('conveniences', conv));
    formData.append('minNumOfVisitors',minNumOfVisitors);
    formData.append('maxNumOfVisitors', maxNumOfVisitors);
    formData.append("availableAccommodationDates", JSON.stringify(availableAccommodationDates));
    formData.append('location', location);

    if (Array.from(images as unknown as Array<any>).length !== 0) {
      Array.from(images as unknown as Array<any>).forEach((file) => {
        formData.append('images', file);
      });
    }
    formData.append('paying', paying);



    return formData;
  }
  





  create(
    userId: string,
    username: string,
    email:string,
    name: string,
    address: string,
    city: string,
    country: string,
    conveniences: string[],
    minNumOfVisitors: string,
    maxNumOfVisitors: string,
    availableAccommodationDates: DateAvailability[],
    location: string,
    images:FormArray,
    paying:string,
  ): void {
    console.log(email)
    this.http
      .post(`${apiURL}/accommodations/`, this.formatFormData(
        userId,
        username,
        email,
        name,
        address,
        city,
        country,
        conveniences,
        minNumOfVisitors,
        maxNumOfVisitors,
        availableAccommodationDates,
        location,
        images,
        paying
      ))
      .subscribe({
        next: (data) => {
          this.toastSerice.showToast(
            'Success',
            'Accommodation created!',
            ToastNotificationType.Success
          );
          this.modalService.close();
          // this.router.navigate(['/']);
          // window.location.reload();
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
    this.router.navigate(['/']);
    window.location.reload();
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
    endDate: string,
    maxPrice:string,
    conveniences:string[],
    isDistinguished:string
  ): Observable<any> {
    
    
    console.log('pocetni datum je', isDistinguished);
    return this.http.get<any>(
      `${apiURL}/accommodations/search?city=${city}&country=${country}&numOfVisitors=${numOfVisitors}&startDate=${startDate}&endDate=${endDate}&maxPrice=${maxPrice}&conveniences=${conveniences}&isDistinguished=${isDistinguished}`
    );
    
  }
}
