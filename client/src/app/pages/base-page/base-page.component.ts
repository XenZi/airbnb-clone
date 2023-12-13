import { Component } from '@angular/core';
import { Observable } from 'rxjs';
import { Accommodation } from 'src/app/domains/entity/accommodation-model';
import { AccommodationsService } from 'src/app/services/accommodations-service/accommodations.service';

@Component({
  selector: 'app-base-page',
  templateUrl: './base-page.component.html',
  styleUrls: ['./base-page.component.scss'],
})
export class BasePageComponent {
  accommodations!: Observable<Accommodation[]>;

  constructor(private accommodationService: AccommodationsService) {}

  ngOnInit() {
    this.loadAccommodations();
    console.log(localStorage.getItem('user'))
    console.log(this.accommodations)
    
  }

  private loadAccommodations(): void {
    this.accommodations = this.accommodationService.loadAccommodations();
  }

  
}
