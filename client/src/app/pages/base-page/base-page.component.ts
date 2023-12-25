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
  accommodations!: Accommodation[];

  constructor(private accommodationService: AccommodationsService) {}

  ngOnInit() {
    this.loadAccommodations();
  }

  private loadAccommodations(): void {
    this.accommodationService.loadAccommodations().subscribe({
      next: (response) => {
        this.accommodations = response.data;
        console.log(this.accommodations)
        
      },
      error: (error) => {
        console.log(error);
      },
    });
  }
}
