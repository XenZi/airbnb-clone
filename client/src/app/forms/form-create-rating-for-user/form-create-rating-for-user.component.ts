import { Component, Input } from '@angular/core';
import { HostRate } from 'src/app/domains/entity/ratings.model';
import { UserAuth } from 'src/app/domains/entity/user-auth.model';
import { User } from 'src/app/domains/entity/user-profile.model';
import { RatingService } from 'src/app/services/rating/rating.service';
import { ReservationService } from 'src/app/services/reservation-service/reservation.service';
import { UserService } from 'src/app/services/user/user.service';

@Component({
  selector: 'app-form-create-rating-for-user',
  templateUrl: './form-create-rating-for-user.component.html',
  styleUrls: ['./form-create-rating-for-user.component.scss'],
})
export class FormCreateRatingForUserComponent {
  stars: number[] = [1, 2, 3, 4, 5];
  rating: number = 0;
  @Input() hostID!: string;
  @Input() userFromParent!: User;
  @Input() hostEmail!: string;
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
      .getPastReservationsByUserForGuest(this.hostID, this.user.id)
      .subscribe({
        next: (data) => {
          this.doesUserHaveReservedThisInPast = data ? true : false;
        },
        error: (err) => {
          console.log(err);
        },
      });
    this.ratingService
      .getRatingForHostByUser(this.user.id, this.hostID)
      .subscribe({
        next: (data) => {
          this.doesUserHaveRatedThisPreviously = data.data ? true : false;
          this.rating = data.data.rate;
        },
        error: (err) => {
          console.log(err);
        },
      });
  }

  rateHost() {
    let rateHost: HostRate = {
      rate: this.rating,
      guest: {
        id: this.user.id,
        username: this.user.username,
        email: this.user.email,
      },
      host: {
        id: this.hostID,
        username: this.userFromParent.username,
        email: this.userFromParent.email,
      },
    };
    if (this.doesUserHaveRatedThisPreviously) {
      this.ratingService.updateRateHost(rateHost);
      return;
    }
    this.ratingService.rateHost(rateHost);
  }

  onRatingChanged(newRating: number) {
    this.rating = newRating;
  }

  deleteRating() {
    this.ratingService.deleteRateHost(this.hostID, this.user.id);
  }
}
