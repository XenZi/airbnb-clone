import { Component } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ToastNotificationType } from 'src/app/domains/enums/toast-notification-type.enum';
import { AccommodationsService } from 'src/app/services/accommodations-service/accommodations.service';
import { ToastService } from 'src/app/services/toast/toast.service';
import { formatErrors } from 'src/app/utils/formatter.utils';

@Component({
  selector: 'app-form-update-accommodation',
  templateUrl: './form-update-accommodation.component.html',
  styleUrls: ['./form-update-accommodation.component.scss']
})
export class FormUpdateAccommodationComponent {

  
  
  updateAccommodationForm: FormGroup;
  
  errors: string = '';
  isCaptchaValidated: boolean = false;
  
  constructor(
    private accommodationsService:AccommodationsService,
    private formBuilder: FormBuilder,
    private toastService: ToastService
  ) {
    this.updateAccommodationForm = this.formBuilder.group({
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
    
    if (!this.updateAccommodationForm.valid) {
      console.log("bilo sta1")
      Object.keys(this.updateAccommodationForm.controls).forEach((key) => {
        const controlErrors = this.updateAccommodationForm.get(key)?.errors;
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
      
      this.updateAccommodationForm.value.name,
      this.updateAccommodationForm.value.location,
      this.updateAccommodationForm.value.conveniences,
      this.updateAccommodationForm.value.minNumOfVisitors,
      this.updateAccommodationForm.value.maxNumOfVisitors

    );
    
  }

}
