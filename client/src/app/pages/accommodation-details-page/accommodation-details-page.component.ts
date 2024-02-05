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

@Component({
  selector: 'app-accommodation-details-page',
  templateUrl: './accommodation-details-page.component.html',
  styleUrls: ['./accommodation-details-page.component.scss'],
})
export class AccommodationDetailsPageComponent implements OnDestroy {
  accommodationID!: string;
  accommodation!: Accommodation;
  userLogged!: UserAuth | null;
  isUserLogged: boolean = false;
  customEventUUID!: string;
  constructor(
    private route: ActivatedRoute,
    private modalService: ModalService,
    private accommodationsService: AccommodationsService,
    private userService: UserService,
    private metricsService: MetricsService,
    private localStorageService: LocalStorageService,
    private previouseRouteService: PreviouseRouteService
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

  sendCreateAt() {
    const currentDate = new Date();
    const formattedDate = currentDate.toISOString().slice(0, 16).replace("T", " ");
    this.metricsService.joinedAt(this.userLogged?.id ? this.userLogged?.id : "not logged in", this.accommodationID, this.customEventUUID, formattedDate)
  }

  checkUserForMetricsAndGetThoseIfHeIsAHost() {
    if (!this.isUserLogged) {
      return;
    }
    if ((this.userLogged?.id as string) !== this.accommodation.userId) {
      return;
    }
    this.metricsService.getMetrics(this.accommodationID, "daily").subscribe({
      next: (data) => {
        console.log("SUCC METRICS", data)
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
        this.checkUserForMetricsAndGetThoseIfHeIsAHost();
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
}
