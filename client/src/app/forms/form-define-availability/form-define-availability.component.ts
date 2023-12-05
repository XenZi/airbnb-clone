import { Component, Input } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { DateAvailability } from 'src/app/domains/entity/date-availability.model';

@Component({
  selector: 'app-form-define-availability',
  templateUrl: './form-define-availability.component.html',
  styleUrls: ['./form-define-availability.component.scss']
})
export class FormDefineAvailabilityComponent {
  @Input() accommodationID!:string
  availableDatesAndPrice: DateAvailability[] = [];
  availabilityForm!: FormGroup ;
  availabilityPriceInstances: any[] = [];
  constructor(private formBuilder: FormBuilder) {
    this.createForm();
  }

  createForm() {
    this.availabilityForm = this.formBuilder.group({
      startDate: ['', Validators.required],
      endDate: ['', Validators.required],
      price: ['', Validators.required],
    });
  }

  onSubmit() {
    

      this.availabilityPriceInstances.push({});
      console.log(this.availabilityPriceInstances)

     
      
    
      
  }

}
