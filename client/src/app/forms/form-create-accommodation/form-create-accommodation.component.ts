
import { Component } from '@angular/core';
import {
  AbstractControl,
  FormBuilder,
  FormGroup,
  ValidationErrors,
  ValidatorFn,
  Validators,
} from '@angular/forms';

import { ToastService } from 'src/app/services/toast/toast.service';
import { ToastNotificationType } from 'src/app/domains/enums/toast-notification-type.enum';
import { AccommodationsService } from 'src/app/services/accommodations-service/accommodations.service';
import { formatErrors } from 'src/app/utils/formatter.utils';
@Component({
  selector: 'app-form-create-accommodation',
  templateUrl: './form-create-accommodation.component.html',
  styleUrls: ['./form-create-accommodation.component.scss']
})
export class FormCreateAccommodationComponent {

  createAccommodationForm: FormGroup;
  
  errors: string = '';
  isCaptchaValidated: boolean = false;
  userId!:string ;
  username!:string;
  
  constructor(
    private accommodationsService:AccommodationsService,
    private formBuilder: FormBuilder,
    private toastService: ToastService
  ) {
    this.createAccommodationForm = this.formBuilder.group({
      name: [
        '',
        [
          Validators.required,
          Validators.minLength(3),
          
        ],
      ],
      location: ['', [Validators.required]],
        wiFi: [false],
        kitchen: [false],
        airConditioning: [false],
        freeParking: [false],
        pool: [false],
      minNumOfVisitors: ['', [Validators.required]],
      maxNumOfVisitors: ['', [Validators.required]],
    });
  }

  ngOnInit(){
    this.getUsernameFromLocal()
  }

  getUsernameFromLocal(){
    const userData = localStorage.getItem('user');

    if(userData) {
      const parsedUserData = JSON.parse(userData);
      this.userId=parsedUserData.id
      this.username = parsedUserData.username;
      console.log(this.username); // This will log the value of the "username" key
    } else {
      console.log('No user data found in localStorage');
    }
    
  }
 
 

  
  onSubmit(e: Event) {
    e.preventDefault();
    
    if (!this.createAccommodationForm.valid) {
      console.log("bilo sta1")
      Object.keys(this.createAccommodationForm.controls).forEach((key) => {
        const controlErrors = this.createAccommodationForm.get(key)?.errors;
        if (controlErrors) {
          this.errors += formatErrors(key);
        }
      });
      this.toastService.showToast(
        'Error',
        this.errors,
        ToastNotificationType.Error
      );
      this.errors = '';
      return;
    }
    const conveniencesForm = this.createAccommodationForm.get('conveniences');

    
  
      
  
      const wiFi = this.createAccommodationForm.value.wiFi ? 'Wi-Fi' : '';
      const kitchen = this.createAccommodationForm.value.kitchen ? 'Kitchen' : '';
      const airConditioning = this.createAccommodationForm.value.airConditioning ? 'Air Conditioning' : '';
      const freeParking = this.createAccommodationForm.value.freeParking ? 'Free Parking' : '';
      const pool = this.createAccommodationForm.value.pool ? 'Pool' : '';

// Concatenate values into a CSV string
const conveniencesCsv = [wiFi, kitchen, airConditioning, freeParking, pool].filter(Boolean).join(', ');

console.log('CSV string:', conveniencesCsv);
    

    this.accommodationsService.create(
      this.userId,
      this.username,
      this.createAccommodationForm.value.name,
      this.createAccommodationForm.value.location,
      conveniencesCsv,
      this.createAccommodationForm.value.minNumOfVisitors,
      this.createAccommodationForm.value.maxNumOfVisitors

    );
    
  }

}
