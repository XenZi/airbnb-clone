import { Component, Input } from '@angular/core';
import { HostRate } from 'src/app/domains/entity/ratings.model';
import { RatingService } from 'src/app/services/rating/rating.service';

@Component({
  selector: 'app-show-rates-for-user',
  templateUrl: './show-rates-for-user.component.html',
  styleUrls: ['./show-rates-for-user.component.scss'],
})
export class ShowRatesForUserComponent {
  @Input() hostID!: string;
  ratings: HostRate[] = [];
  constructor(private ratingService: RatingService) {}

  ngOnInit() {
    this.ratingService.getAllRatingsForHost(this.hostID).subscribe({
      next: (data) => {
        this.ratings = data.data;
        console.log(this.ratings);
      },
      error: (err) => {
        console.log(err);
      },
    });
  }
}
