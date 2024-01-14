import { Component } from '@angular/core';
import { Reservation } from 'src/app/domains/entity/reservation.model';
import { UserAuth } from 'src/app/domains/entity/user-auth.model';
import { User } from 'src/app/domains/entity/user-profile.model';
import { ReservationService } from 'src/app/services/reservation-service/reservation.service';
import { UserService } from 'src/app/services/user/user.service';

@Component({
  selector: 'app-user-reservations-table',
  templateUrl: './user-reservations-table.component.html',
  styleUrls: ['./user-reservations-table.component.scss']
})
export class UserReservationsTableComponent {
  userReservations!: Reservation[]
  user: UserAuth | null = null;

  constructor(private userService: UserService, private reservationsService: ReservationService) {}


  ngOnInit() {
    this.user = this.userService.getLoggedUser();
    if (this.user) {
      this.reservationsService.getAllReservationsById(this.user.id).subscribe({
        next: (data) => {
          this.userReservations = data.data;
        }, 
        error: (err: Error) => {
          console.log(err)
        }
      })
    }
  }


  cancelReservation(index: number) {
    const reservation = this.userReservations[index];
    if (reservation) {
      this.reservationsService.deleteById(
        reservation.country,
        reservation.id,
        reservation.userID,
        reservation.hostID,
        reservation.accommodationID
      )
      .subscribe(
        (deletedReservation) => {
          // Handle success, if needed
          console.log('Reservation canceled successfully:', deletedReservation);
          // Optionally, you can remove the canceled reservation from the local array
          this.userReservations.splice(index, 1);
        },
        (error) => {
          // Handle error, if needed
          console.error('Error canceling reservation:', error);
        }
      );
    }
  }
}
