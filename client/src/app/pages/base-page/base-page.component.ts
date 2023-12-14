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
    this.accommodationService.loadAccommodations().subscribe({
      next: (response) => {
        // Check the structure of the response here
        // For example, if response has a 'data' property containing accommodations
        // Modify the following line accordingly based on the actual response structure
        this.accommodations = response.data; // Update this line based on the actual response structure
      },
      error: (error) => {
        console.log(error)
      }
    });
  }

  
}
