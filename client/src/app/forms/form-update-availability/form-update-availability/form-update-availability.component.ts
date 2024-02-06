import { Component, Input, OnInit } from '@angular/core';

import { UserAuth } from 'src/app/domains/entity/user-auth.model';
import { ReservationService } from 'src/app/services/reservation-service/reservation.service';
import {
  AbstractControl,
  FormArray,
  FormBuilder,
  FormGroup,
  ValidationErrors,
  ValidatorFn,
  Validators,
} from '@angular/forms';
import { formatErrors } from 'src/app/utils/formatter.utils';
import { th } from 'date-fns/locale';
import { DateAvailability } from 'src/app/domains/entity/date-availability.model';


@Component({
  selector: 'app-form-update-availability',
  templateUrl: './form-update-availability.component.html',
  styleUrls: ['./form-update-availability.component.scss']
})
export class FormUpdateAvailabilityComponent implements OnInit{
  user: UserAuth | null = null;
  updateAvailabilityForm: FormGroup;
  errors: string = '';
  @Input()accommodationID!: string;
  @Input()id!: string;
  @Input()country!: string;
  @Input()price!: number;
  @Input()location!:string;

  


  constructor(
    private fb: FormBuilder,
    private reservationService: ReservationService
  ){
    this.updateAvailabilityForm = this.fb.group({
      dateAvailabilities: this.fb.array([this.ngOnInit()]),

    });
  }

  ngOnInit(): void {
    this.updateAvailabilityForm = this.fb.group({
      startDate: ['', Validators.required],
      endDate: ['', Validators.required],
      price: ['', Validators.required]
    });
  }




  onSubmit(){
 
  if(!this.updateAvailabilityForm.valid){
  this.errors = '';
  Object.keys(this.updateAvailabilityForm.controls).forEach((key) => {
    const controlErrors = this.updateAvailabilityForm.get(key)?.errors;
    if (controlErrors) {
          this.errors += formatErrors(key);
        }
  });
  return;
  }

 
  const availabilites=this.processDateAvailabilities()

  this.reservationService.update(
    this.accommodationID,this.id,this.country,this.price,this.location,[availabilites]
  )

  }
   processDateAvailabilities(): DateAvailability {
    let processedData: DateAvailability = {dateRange: [],price:0,id: '',accommodationID: '',location: ''};
  
   
      const startDate = new Date(this.updateAvailabilityForm.get('startDate')?.value as string);
      const endDate = new Date(this.updateAvailabilityForm.get('endDate')?.value as string);
      const price = this.updateAvailabilityForm.get('price')?.value
  
      // Generate date range
      const currentDates = this.getDatesBetween(startDate, endDate).map(date => date.toISOString().split('T')[0]);
      console.log(currentDates)
  
      // Do something with startDate, endDate, and price
      console.log(`Start Date: ${startDate.toISOString().split('T')[0]}, End Date: ${endDate.toISOString().split('T')[0]}, Price: ${price}`);
      
      // Add currentDates and price to the processedData array
      processedData =  {dateRange: currentDates, price: price};
    
  
    // The processedData array now contains objects with dateRange (dates only) and price for every entry
    return processedData;
  }
  
  // Helper function to get dates between start and end dates
  getDatesBetween(startDate: Date, endDate: Date): Date[] {
    const dates = [];
    let currentDate = new Date(startDate);
  
    while (currentDate <= endDate) {
      dates.push(new Date(currentDate));
      currentDate.setDate(currentDate.getDate() + 1);
    }
  
    return dates;
  }


}
