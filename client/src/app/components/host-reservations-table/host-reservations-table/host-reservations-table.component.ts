import { Component } from '@angular/core';
import { Reservation } from 'src/app/domains/entity/reservation.model';
import { UserAuth } from 'src/app/domains/entity/user-auth.model';
import { ReservationService } from 'src/app/services/reservation-service/reservation.service';
import { UserService } from 'src/app/services/user/user.service';

@Component({
  selector: 'app-host-reservations-table',
  templateUrl: './host-reservations-table.component.html',
  styleUrls: ['./host-reservations-table.component.scss']
})
export class HostReservationsTableComponent {
  hostReservations!: Reservation[]
  user: UserAuth | null = null;

  constructor(private userService: UserService, private reservationsService: ReservationService) {}

  ngOnInit() {
    this.user = this.userService.getLoggedUser();
    if (this.user) {
      this.reservationsService.getAllReservationsByHost(this.user.id).subscribe({
        next: (data) => {
          console.log(data)
          this.hostReservations = data;
          console.log(this.hostReservations)
        }, 
        error: (err: Error) => {
          console.log(err)
        }
      })
    }
  }
}
