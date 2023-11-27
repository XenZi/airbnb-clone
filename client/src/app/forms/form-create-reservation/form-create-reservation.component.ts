// reservation-form.component.ts
import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ReservationService } from 'src/app/services/reservation-service/reservation.service';

@Component({
  selector: 'app-reservation-form',
  templateUrl: 'form-create-reservation.component.html',
  styleUrls: ['form-create-reservation.component.scss']
})
export class ReservationFormComponent implements OnInit {
  reservationForm!: FormGroup;

  constructor(private fb: FormBuilder, private reservationService: ReservationService) {}

  ngOnInit() {
    this.initializeForm();
  }

  initializeForm() {
    this.reservationForm = this.fb.group({
      startDate: ['', Validators.required],
      endDate: ['', Validators.required],
      username: ['', Validators.required],
      accommodationName: ['', Validators.required],
    });
  }

  submitReservation() {
    if (this.reservationForm.valid) {
      const reservationData = this.reservationForm.value;
      this.reservationService.createReservation(reservationData).subscribe(
        (response) => {
          console.log('Reservation created successfully:', response);
        },
        (error) => {
          console.error('Error creating reservation:', error);
        }
      );
    } else {
    }
  }
}
