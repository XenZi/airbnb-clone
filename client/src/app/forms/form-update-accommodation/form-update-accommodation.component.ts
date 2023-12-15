import { Component, Input } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { timeout } from 'rxjs';
import { Accommodation } from 'src/app/domains/entity/accommodation-model';
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
  @Input() accommodationID!:string
  
  
  updateAccommodationForm: FormGroup;
  accommodation!:Accommodation;
  errors: string = '';
  isCaptchaValidated: boolean = false;
  
  constructor(
    private accommodationsService:AccommodationsService,
    private formBuilder: FormBuilder,
    private toastService: ToastService
  ) {
    this.updateAccommodationForm = this.formBuilder.group({
      name: [
        ''],
      location: [''],
      conveniences: [''],
      minNumOfVisitors: [''],
      maxNumOfVisitors: [''],
    });
  }

  ngOnInit(){

    
    
    this.getAccommodationById();
    setTimeout(() => {
      this.updateAccommodationForm = this.formBuilder.group({
        name: [this.accommodation.name, Validators.required],
        location:[this.accommodation.city, Validators.required],
        conveniences:[this.accommodation.conveniences, Validators.required],
        minNumOfVisitors:[this.accommodation.minNumOfVisitors, Validators.required],
        maxNumOfVisitors:[this.accommodation.maxNumOfVisitors, Validators.required],
        
      
      });
     }, 200);
    
    
  }

 
  getAccommodationById() {
    this.accommodationsService.getAccommodationById(this.accommodationID as string).subscribe((data) => {
      this.accommodation= data;
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
    console.log(this.updateAccommodationForm.value.minNumOfVisitors as number)
    this.accommodationsService.update(

      this.accommodationID,
      this.updateAccommodationForm.value.name,
      this.updateAccommodationForm.value.address,
      this.updateAccommodationForm.value.city,
      this.updateAccommodationForm.value.country,
      this.updateAccommodationForm.value.conveniences,
      this.updateAccommodationForm.value.minNumOfVisitors as number,
      this.updateAccommodationForm.value.maxNumOfVisitors as number

    );

    
  }

}
