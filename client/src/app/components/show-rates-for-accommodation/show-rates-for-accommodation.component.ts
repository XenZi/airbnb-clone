import { Component, Input } from '@angular/core';
import { AccommodationRate } from 'src/app/domains/entity/ratings.model';
import { RatingService } from 'src/app/services/rating/rating.service';

@Component({
  selector: 'app-show-rates-for-accommodation',
  templateUrl: './show-rates-for-accommodation.component.html',
  styleUrls: ['./show-rates-for-accommodation.component.scss'],
})
export class ShowRatesForAccommodationComponent {
  @Input() accommodationID!: string;
  rates: AccommodationRate[] = [];
  constructor(private ratingService: RatingService) {}

  ngOnInit() {
    this.ratingService
      .getAllRatingsForAccommodation(this.accommodationID)
      .subscribe({
        next: (data) => {
          this.rates = data.data;
        },
        error: (err) => {
          console.log(err);
        },
      });
  }
}
