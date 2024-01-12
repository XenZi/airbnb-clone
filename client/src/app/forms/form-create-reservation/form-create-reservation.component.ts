// reservation-form.component.ts
import { Component, Input, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Accommodation } from 'src/app/domains/entity/accommodation-model';
import { UserAuth } from 'src/app/domains/entity/user-auth.model';
import { ReservationService } from 'src/app/services/reservation-service/reservation.service';
import { UserService } from 'src/app/services/user/user.service';

@Component({
  selector: 'app-reservation-form',
  templateUrl: 'form-create-reservation.component.html',
  styleUrls: ['form-create-reservation.component.scss'],
})
export class ReservationFormComponent implements OnInit {
  reservationForm!: FormGroup;
  @Input() accommodation!: Accommodation;
  @Input() accommodationID!: string
  user: UserAuth | null = null;
  availabilityData: any[] = [];
  constructor(
    private fb: FormBuilder,
    private reservationService: ReservationService,
    private userService: UserService
  ) {}

  ngOnInit() {
    let accommodationID: string = this.accommodationID
    this.initializeForm();
    this.user = this.userService.getLoggedUser();
    this.reservationService.getAvailability(accommodationID).subscribe({next: (data) => {
      console.log(data)
      this.availabilityData = data
    },error: (err) => {
      console.log(err)
    }})
    
  }

  initializeForm() {
    this.reservationForm = this.fb.group({
      range: ['', [Validators.required]],
    });
  }

  handleDateChange(rangeDates: Date[]) {
    console.log(rangeDates);
    const formattedRange = rangeDates.map((date) => this.formatDate(date));
    this.reservationForm.get('range')?.setValue(formattedRange);
  }
  formatDate(date: Date) {
    let day = date.getDate();
    let month = date.getMonth() + 1;
    let year = date.getFullYear();

    let newDay = (day < 10 ? '0' + day : day) as string;
    let newMonth = month < 10 ? '0' + month : month;

    return `${year}-${newMonth}-${newDay}`;
  }

  submitReservation() {
    if (!this.reservationForm.valid) {
      console.log('not valid');
      return;
    }
    let userID: string = this.user?.id as string;
    let accommodationID: string = this.accommodationID
    let startDate: string = this.reservationForm.value.range[0] 
    let endDate: string = this.reservationForm.value.range[this.reservationForm.value.range.length - 1]
    let username: string = this.user?.username as string;
    let accommodationName: string = this.accommodation.name;
    let location: string = "bb,Belgrade,Serbia";
    let price: number = 50;
    let numOfDays: number = this.reservationForm.value.range.length;
    let dateRange: string[] = this.reservationForm.value.range;
    let hostID: string = this.accommodation.userId
    console.log(this.accommodation)
    let reservationData = {"userID": userID,
    "accommodationID": accommodationID,
    "startDate": startDate,
    "endDate": endDate,
    "username": username,
    "accommodationName": accommodationName,
    "location": location,
    "price": price,
    "numOfDays": numOfDays,
    "dateRange": dateRange,
    "hostID" : hostID
  }
  console.log(reservationData)
  this.reservationService.createReservation(reservationData)

    console.log(
      userID,
      accommodationID,
      startDate,
      endDate,
      username,
      accommodationName,
      location,
      price,
      numOfDays,
      dateRange
    );
  }
}
