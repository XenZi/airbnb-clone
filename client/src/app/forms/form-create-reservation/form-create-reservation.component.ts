import { Component, Input, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Accommodation } from 'src/app/domains/entity/accommodation-model';
import { ReservationService } from 'src/app/services/reservation-service/reservation.service';
import { UserService } from 'src/app/services/user/user.service';

@Component({
  selector: 'app-reservation-form',
  templateUrl: 'form-create-reservation.component.html',
  styleUrls: ['form-create-reservation.component.scss']
})
export class ReservationFormComponent implements OnInit {
  reservationForm!: FormGroup;
  @Input() accommodation!: Accommodation;
  selectedDateRange: Date[] = [];

  constructor(
    private fb: FormBuilder,
    private reservationService: ReservationService

  ) { this.reservationForm = this.fb.group({
      dateRange: ['', Validators.required],
    });}

  ngOnInit() {
    this.initializeForm();
  }

  initializeForm() {
    this.reservationForm = this.fb.group({
      dateRange: ['', Validators.required],
    });
  }

  submitReservation() {
    if (this.reservationForm.valid) {
      const [startDate, endDate] = this.selectedDateRange;

      const reservationData = {
        startDate,
        endDate,
        accommodationID: this.accommodation.id,
        location: this.accommodation.location,
        accommodationName: this.accommodation.name,
        userID: this.reservationService.getLoggedUserId(),
      };

      this.reservationService.createReservation(reservationData).subscribe(
        (response: any) => {
          console.log('Reservation created successfully:', response);
        },
        (error: any) => {
          console.error('Error creating reservation:', error);
        }
      );
    } else {
      console.log('Form is not valid. Please check the input fields.');
    }
  }

  
}
