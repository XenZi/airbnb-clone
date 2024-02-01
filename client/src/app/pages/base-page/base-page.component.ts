import { Component } from '@angular/core';
import { Observable } from 'rxjs';
import { Accommodation } from 'src/app/domains/entity/accommodation-model';
import { UserAuth } from 'src/app/domains/entity/user-auth.model';
import { AccommodationsService } from 'src/app/services/accommodations-service/accommodations.service';
import { RecommendationService } from 'src/app/services/recommendation/recommendation.service';
import { UserService } from 'src/app/services/user/user.service';

@Component({
  selector: 'app-base-page',
  templateUrl: './base-page.component.html',
  styleUrls: ['./base-page.component.scss'],
})
export class BasePageComponent {
  accommodations!: Accommodation[];
  recommendedAccommodations!: Accommodation[];
  userAuth: UserAuth | null = null;
  constructor(
    private accommodationService: AccommodationsService,
    private recommendationService: RecommendationService,
    private userService: UserService
  ) {}

  ngOnInit() {
    this.loadAccommodations();
    this.userAuth = this.userService.getLoggedUser();
    this.userAuth
      ? this.recommendationService
          .getAllRecommendedByUserID(this.userAuth.id)
          .subscribe({
            next: (data) => {
              this.recommendedAccommodations = data.data;
            },
            error: (err) => {
              console.log(err);
            },
          })
      : this.recommendationService.getAllRecommendedByRating().subscribe({
          next: (data) => {
            this.recommendedAccommodations = data.data;
          },
          error: (err) => {
            console.log(err);
          },
        });
  }

  private loadAccommodations(): void {
    this.accommodationService.loadAccommodations().subscribe({
      next: (response) => {
        this.accommodations = response.data;
        console.log(this.accommodations);
      },
      error: (error) => {
        console.log(error);
      },
    });
  }
}
