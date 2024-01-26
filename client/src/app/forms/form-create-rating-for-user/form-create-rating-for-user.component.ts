import { Component, Input } from '@angular/core';
import { UserAuth } from 'src/app/domains/entity/user-auth.model';
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
          console.log(data);
        },
        error: (err) => {
          console.log(err);
        },
      });
  }
}
