// reservation-form.component.ts
import { Component, Input, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Accommodation } from 'src/app/domains/entity/accommodation-model';
import { ReservationService } from 'src/app/services/reservation-service/reservation.service';

@Component({
  selector: 'app-reservation-form',
  templateUrl: 'form-create-reservation.component.html',
  styleUrls: ['form-create-reservation.component.scss']
})
export class ReservationFormComponent implements OnInit {
  reservationForm!: FormGroup;
  @Input() accommodation!:Accommodation

  constructor(private fb: FormBuilder, private reservationService: ReservationService) {}

  ngOnInit() {
    this.initializeForm();
  }

  initializeForm() {
    this.reservationForm = this.fb.group({
      startDate: ['', Validators.required],
      endDate: ['', Validators.required],
    });
  }

  submitReservation() {
    if (this.reservationForm.valid) {
      const reservationData = this.reservationForm.value;
      reservationData.accommodationID = this.accommodation.id
      reservationData.location = this.accommodation.city
      reservationData.accommodationName = this.accommodation.name
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
