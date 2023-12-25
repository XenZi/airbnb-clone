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
          this.userReservations = data;
        }, 
        error: (err: Error) => {
          console.log(err)
        }
      })
    }
  }


  cancelReservation(index: number) {
    this.reservationsService.deleteById(this.userReservations[index].id, this.userReservations[index].country)
  }
}
