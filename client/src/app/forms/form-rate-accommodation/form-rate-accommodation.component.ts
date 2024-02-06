import { Component, Host, Input } from '@angular/core';
import { Guest } from 'src/app/domains/entity/ratings.model';
import { UserAuth } from 'src/app/domains/entity/user-auth.model';
import { AccommodationRate } from 'src/app/domains/entity/ratings.model';
import { RatingService } from 'src/app/services/rating/rating.service';
import { UserService } from 'src/app/services/user/user.service';
import { ReservationService } from 'src/app/services/reservation-service/reservation.service';

@Component({
  selector: 'app-form-rate-accommodation',
  templateUrl: './form-rate-accommodation.component.html',
  styleUrls: ['./form-rate-accommodation.component.scss'],
})
export class FormRateAccommodationComponent {
  stars: number[] = [1, 2, 3, 4, 5];
  rating: number = 0;
  @Input() accommodationID!: string;
  user!: UserAuth;
  doesUserHaveReservedThisInPast: boolean = false;
  doesUserHaveRatedThisPreviously: boolean = false;
  constructor(
    private userService: UserService,
    private ratingService: RatingService,
    private reservationsService: ReservationService
  ) {}

  ngOnInit() {
    this.user = this.userService.getLoggedUser() as UserAuth;
    this.reservationsService
      .getPastReservations(this.accommodationID, this.user.id)
      .subscribe({
        next: (data) => {
          this.doesUserHaveReservedThisInPast = data.data ? true : false;
        },
        error: (err) => {
          console.log(err);
        },
      });
    this.ratingService
      .getRatingForAccommodationByGuest(this.accommodationID, this.user.id)
      .subscribe({
        next: (data) => {
          this.doesUserHaveRatedThisPreviously = true;
          this.rating = data.data.rate;
        },
        error: (err) => {
        },
      });
  }

  onRatingChanged(newRating: number) {
    this.rating = newRating;
  }

  rateAccommodation() {
    let rateAccommodation: AccommodationRate = {
      accommodationID: this.accommodationID,
      guest: {
        username: this.user.username,
        id: this.user.id,
        email: this.user.email,
      },
      rate: this.rating,
    };
    if (this.doesUserHaveRatedThisPreviously) {
      this.ratingService.updateRateForAccommodation(rateAccommodation);
      return;
    }
    this.ratingService.rateAccommodation(rateAccommodation);
  }

  deleteRating() {
    this.ratingService.deleteRateForAccomodation(
      this.accommodationID,
      this.user.id
    );
  }
}
