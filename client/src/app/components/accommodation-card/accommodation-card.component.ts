import { Component } from '@angular/core';
import { Observable } from 'rxjs';
import { Accommodation } from 'src/app/domains/entity/accommodation-model';
import { AccommodationsService } from 'src/app/services/accommodations-service/accommodations.service';

@Component({
  selector: 'app-accommodation-card',
  templateUrl: './accommodation-card.component.html',
  styleUrls: ['./accommodation-card.component.scss']
})
export class AccommodationCardComponent {
  
  accommodations!:Observable<Accommodation[]>

  constructor(private accommodationService:AccommodationsService){}
  ngOnInit():void{
    this.loadAccommodations();
    console.log(this.accommodations)
  }

  private loadAccommodations(){
    this.accommodations=this.accommodationService.loadAccommodations();
  }



}
