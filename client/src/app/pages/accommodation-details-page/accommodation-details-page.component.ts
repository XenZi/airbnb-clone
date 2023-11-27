import { Component } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Accommodation } from 'src/app/domains/entity/accommodation-model';
import { AccommodationsService } from 'src/app/services/accommodations-service/accommodations.service';

@Component({
  selector: 'app-accommodation-details-page',
  templateUrl: './accommodation-details-page.component.html',
  styleUrls: ['./accommodation-details-page.component.scss']
})
export class AccommodationDetailsPageComponent {
  accommodationID:string | undefined 
  accommodation!:Accommodation;

  constructor(
    private route: ActivatedRoute,
    
    private accommodationsService: AccommodationsService,
    
    
    
  ) {}

  ngOnInit(){

    this.getAccommodationID();
    this.getAccommodationById();
    
  
  }
  getAccommodationID() {
    this.route.paramMap.subscribe((params) => {
      this.accommodationID = String(params.get('id'));
    });
  }

  getAccommodationById() {
    this.accommodationsService.getAccommodationById(this.accommodationID as string).subscribe((data) => {
      this.accommodation= data;
    });
  }

}
