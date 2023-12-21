import { ModalService } from './../../services/modal/modal.service';
import { Component } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Accommodation } from 'src/app/domains/entity/accommodation-model';
import { UserAuth } from 'src/app/domains/entity/user-auth.model';

import { FormUpdateAccommodationComponent } from 'src/app/forms/form-update-accommodation/form-update-accommodation.component';
import { AccommodationsService } from 'src/app/services/accommodations-service/accommodations.service';
import { UserService } from 'src/app/services/user/user.service';

@Component({
  selector: 'app-accommodation-details-page',
  templateUrl: './accommodation-details-page.component.html',
  styleUrls: ['./accommodation-details-page.component.scss'],
})
export class AccommodationDetailsPageComponent {
  accommodationID!: string;
  accommodation!: Accommodation;
  userLogged!: UserAuth | null;
  isUserLogged: boolean = false;
  constructor(
    private route: ActivatedRoute,
    private modalService: ModalService,
    private accommodationsService: AccommodationsService,
    private userService: UserService
  ) {}

  ngOnInit() {
    this.getAccommodationID();
    this.getAccommodationById();
    this.userLogged = this.userService.getLoggedUser();
    this.isUserLogged = this.userLogged ? true : false;
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
        this.accommodation = data.data;
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
