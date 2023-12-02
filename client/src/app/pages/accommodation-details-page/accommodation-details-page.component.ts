import { ModalService } from './../../services/modal/modal.service';
import { Component } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Accommodation } from 'src/app/domains/entity/accommodation-model';
import { User } from 'src/app/domains/entity/user-profile.model';
import { FormUpdateAccommodationComponent } from 'src/app/forms/form-update-accommodation/form-update-accommodation.component';
import { AccommodationsService } from 'src/app/services/accommodations-service/accommodations.service';

@Component({
  selector: 'app-accommodation-details-page',
  templateUrl: './accommodation-details-page.component.html',
  styleUrls: ['./accommodation-details-page.component.scss'],
})
export class AccommodationDetailsPageComponent {
  accommodationID!: string;
  accommodation!: Accommodation;
  

  constructor(
    private route: ActivatedRoute,
    private modalService: ModalService,
    private accommodationsService: AccommodationsService
  ) {}

  ngOnInit() {
    this.getAccommodationID();
    this.getAccommodationById();
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
  
  getUsernameFromLocal(){
    const userData = localStorage.getItem('user');

    if(userData) {
      const parsedUserData = JSON.parse(userData);
      const username = parsedUserData.username;
      console.log(username); // This will log the value of the "username" key
    } else {
      console.log('No user data found in localStorage');
    }
    
  }

  getAccommodationById() {
    this.accommodationsService
      .getAccommodationById(this.accommodationID as string)
      .subscribe((data) => {
        this.accommodation = data;
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
