import { ModalService } from './../../services/modal/modal.service';
import { Component, HostListener, OnDestroy } from '@angular/core';
import { ActivatedRoute, NavigationStart, Router } from '@angular/router';
import { Accommodation } from 'src/app/domains/entity/accommodation-model';
import { UserAuth } from 'src/app/domains/entity/user-auth.model';
import { v4 as uuidv4 } from 'uuid';

import { FormUpdateAccommodationComponent } from 'src/app/forms/form-update-accommodation/form-update-accommodation.component';
import { AccommodationsService } from 'src/app/services/accommodations-service/accommodations.service';
import { MetricsService } from 'src/app/services/metrics/metrics.service';
import { UserService } from 'src/app/services/user/user.service';
import { UnloadService } from 'src/app/services/unload/unload.service';
import { Subscription } from 'rxjs';
import { LocalStorageService } from 'src/app/services/localstorage/local-storage.service';
import { PreviouseRouteService } from 'src/app/services/previous-route/previouse-route.service';
import { Metrics } from 'src/app/domains/entity/metrics.model';
import { ReservationService } from 'src/app/services/reservation-service/reservation.service';
import { FormUpdateAvailabilityComponent } from 'src/app/forms/form-update-availability/form-update-availability/form-update-availability.component';
import { th } from 'date-fns/locale';
import { DateAvailability } from 'src/app/domains/entity/date-availability.model';

@Component({
  selector: 'app-accommodation-details-page',
  templateUrl: './accommodation-details-page.component.html',
  styleUrls: ['./accommodation-details-page.component.scss'],
})
export class AccommodationDetailsPageComponent implements OnDestroy {
  accommodationID!: string;
  accommodation!: Accommodation;
  availabilityData: DateAvailability [] = [];
  userLogged!: UserAuth | null;
  isUserLogged: boolean = false;
  customEventUUID!: string;
  currentStateOfLookingForMetrics: number = 0;
  currentStateOfMetrics: Metrics = {
    numberOfRatings: 0,
    numberOfReservations: 0,
    numberOfVisits: 0,
    onScreenTime: 0
  };
  constructor(
    private route: ActivatedRoute,
    private modalService: ModalService,
    private accommodationsService: AccommodationsService,
    private userService: UserService,
    private metricsService: MetricsService,
    private localStorageService: LocalStorageService,
    private previouseRouteService: PreviouseRouteService,
    private reservationService: ReservationService
  ) {
  }

  ngOnDestroy(): void {
    const currentDate = new Date();
    const formattedDate = currentDate.toISOString().slice(0, 16).replace("T", " ");
    let currentStateOfLocalStorage = JSON.parse(sessionStorage.getItem("joins") as string);
    if (!currentStateOfLocalStorage) {
      currentStateOfLocalStorage = [{
        userID: this.userLogged?.id ? this.userLogged?.id : "not logged in",
        accommodationID: this.accommodationID,
        customUUID: this.customEventUUID,
        leftAt: formattedDate
      }]
    } else {
      currentStateOfLocalStorage = [...currentStateOfLocalStorage,  {
        userID: this.userLogged?.id ? this.userLogged?.id : "not logged in",
        accommodationID: this.accommodationID,
        customUUID: this.customEventUUID,
        leftAt: formattedDate
      }]
    }
    sessionStorage.setItem("joins", JSON.stringify(currentStateOfLocalStorage));
  }


  ngOnInit() {
    this.getAccommodationID();
    this.getAccommodationById();
    this.userLogged = this.userService.getLoggedUser();
    this.isUserLogged = this.userLogged ? true : false;
    this.customEventUUID = uuidv4();
    this.sendCreateAt();
    console.log(this.previouseRouteService.getPreviousUrl());
    this.reservationService.getAvailability(this.accommodationID).subscribe({next: (data) => {
      console.log(data)
      this.availabilityData = data
    },error: (err) => {
      console.log(err)
    }})
  }

  @HostListener('window:beforeunload', ['$event'])
  testFunc() {
    const currentDate = new Date();
    const formattedDate = currentDate.toISOString().slice(0, 16).replace("T", " ");
    let currentStateOfLocalStorage = JSON.parse(sessionStorage.getItem("joins") as string);
    if (!currentStateOfLocalStorage) {
      currentStateOfLocalStorage = [{
        userID: this.userLogged?.id ? this.userLogged?.id : "not logged in",
        accommodationID: this.accommodationID,
        customUUID: this.customEventUUID,
        leftAt: formattedDate
      }]
    } else {
      currentStateOfLocalStorage = [...currentStateOfLocalStorage,  {
        userID: this.userLogged?.id ? this.userLogged?.id : "not logged in",
        accommodationID: this.accommodationID,
        customUUID: this.customEventUUID,
        leftAt: formattedDate
      }]
    }
    sessionStorage.setItem("joins", JSON.stringify(currentStateOfLocalStorage));
  }

  setNewCurrentStateOfLookingForMetrics(num: number) {
    this.currentStateOfLookingForMetrics = num;

  }

  handleValueChangeOfStateOfMetrics(num: number) {
    console.log(num);
    this.checkUserForMetricsAndGetThoseIfHeIsAHost(this.currentStateOfLookingForMetrics);
  }

  sendCreateAt() {
    const currentDate = new Date();
    const formattedDate = currentDate.toISOString().slice(0, 16).replace("T", " ");
    this.metricsService.joinedAt(this.userLogged?.id ? this.userLogged?.id : "not logged in", this.accommodationID, this.customEventUUID, formattedDate)
  }

  checkUserForMetricsAndGetThoseIfHeIsAHost(numb: number) {
    if (!this.isUserLogged) {
      return;
    }
    if ((this.userLogged?.id as string) !== this.accommodation.userId) {
      return;
    }
    this.metricsService.getMetrics(this.accommodationID, numb ? "monthly" : "daily").subscribe({
      next: (data) => {
        console.log("SUCC METRICS", data)
        this.currentStateOfMetrics = data.data;
      },
      error: (err) => {
        console.log(err);
      }
    })
  }

  updateClick() {
    this.callUpdateAccommodation();
    console.log('Uslo');
  }

  deleteClick() {
    this.callDeleteAccommodation();
    console.log('Uslo');
  }

  getAccommodationID() {
    this.route.paramMap.subscribe((params) => {
      this.accommodationID = String(params.get('id'));
    });
  }

  getAccommodationById() {
    this.accommodationsService
      .getAccommodationById(this.accommodationID as string)
      .subscribe((data) => {
        console.log(data);
        this.accommodation = data.data;
        this.checkUserForMetricsAndGetThoseIfHeIsAHost(0);
      });
  }

  callUpdateAccommodation() {
    this.modalService.open(
      FormUpdateAccommodationComponent,
      'Update accommodation',
      {
        accommodationID: this.accommodationID,
      }
    );
  }

  callDeleteAccommodation() {
    this.accommodationsService.deleteById(this.accommodationID as string);
  }
  
  callUpdateAvailability(index: number){
    this.modalService.open(
      FormUpdateAvailabilityComponent,
      'Update availability',
      {
        accommodationID:this.accommodationID,
        id:this.availabilityData[index].id,
        country: this.accommodation.country,
        price: this.availabilityData[index].price,
        location: this.createLocationCsv()
        
      }
    )
  }
  createLocationCsv():string{
    const address=this.accommodation.address;
    const city=  this.accommodation.city;
    const countryCSV= this.accommodation.country;
    console.log(address + ","+city +","+countryCSV)
    return address + ","+city +","+countryCSV
  }

}
