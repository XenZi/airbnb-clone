
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
      conveniences: [''],
      minNumOfVisitors: ['', [Validators.required]],
      maxNumOfVisitors: ['', [Validators.required]],
    });
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
    this.accommodationsService.create(
      
      this.createAccommodationForm.value.name,
      this.createAccommodationForm.value.location,
      this.createAccommodationForm.value.conveniences,
      this.createAccommodationForm.value.minNumOfVisitors,
      this.createAccommodationForm.value.maxNumOfVisitors

    );
    
  }

}