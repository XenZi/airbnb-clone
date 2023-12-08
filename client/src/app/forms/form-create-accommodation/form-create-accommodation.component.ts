import { Component } from '@angular/core';
import {
  AbstractControl,
  FormArray,
  FormBuilder,
  FormGroup,
  ValidationErrors,
  ValidatorFn,
  Validators,
} from '@angular/forms';
import { UserAuth } from 'src/app/domains/entity/user-auth.model';
import { User } from 'src/app/domains/entity/user-profile.model';
import { AccommodationsService } from 'src/app/services/accommodations-service/accommodations.service';
import { UserService } from 'src/app/services/user/user.service';
import { formatErrors } from 'src/app/utils/formatter.utils';

@Component({
  selector: 'app-form-create-accommodation',
  templateUrl: './form-create-accommodation.component.html',
  styleUrls: ['./form-create-accommodation.component.scss'],
})
export class FormCreateAccommodationComponent {
  createAccommodationForm: FormGroup;
  convenienceList = [
    'WiFi',
    'Kitchen',
    'Air Conditioning',
    'Free Parking',
    'Pool',
  ];
  errors: string = '';
  user: UserAuth | null = null;
  constructor(
    private fb: FormBuilder,
    private accommodationsService: AccommodationsService,
    private userService: UserService
  ) {
    this.createAccommodationForm = this.fb.group({
      name: [''],
      address: [''],
      city: [''],
      country: [''],
      conveniences: this.buildConveniences(),
      maxNumOfVisitors: [''],
      minNumOfVisitors: [''],
      dateAvailabilities: this.fb.array([this.initDateAvailability()]),
    });
  }

  ngOnInit() {
    this.user = this.userService.getLoggedUser();
  }
  initDateAvailability(): FormGroup {
    return this.fb.group({
      startDate: ['', Validators.required],
      endDate: ['', Validators.required],
      price: ['', [Validators.required, Validators.min(0)]],
    });
  }
  buildConveniences(): FormArray {
    const arr = this.convenienceList.map(() => {
      return this.fb.control(false);
    });
    return this.fb.array(arr);
  }

  get convenienceFormArray(): FormArray {
    return this.createAccommodationForm.get('conveniences') as FormArray;
  }

  get conveniences(): FormArray {
    return this.createAccommodationForm.get('conveniences') as FormArray;
  }

  get dateAvailabilities(): FormArray {
    return this.createAccommodationForm.get('dateAvailabilities') as FormArray;
  }

  addDateAvailability(): void {
    this.dateAvailabilities.push(this.initDateAvailability());
  }

  removeLastDateAvailability(): void {
    this.dateAvailabilities.removeAt(this.dateAvailabilities.length - 1);
  }

  fromBooleanToConveniences(): string[] {
    let convArray: string[] = [];
  
    this.conveniences.value.forEach((el: boolean, i: number) => {
      if (el === true && this.convenienceList[i]) {
        convArray.push(this.convenienceList[i]);
      }
    });
    console.log(convArray)
    return convArray;
  }

  createLocationCsv():string{
    const address=this.createAccommodationForm.value.address;
    const city=  this.createAccommodationForm.value.city;
    const country= this.createAccommodationForm.value.country;
    console.log(address + ","+city +","+country)
    return address + ","+city +","+country
  }

  onSubmit(e: Event) {
    e.preventDefault();
    if (!this.createAccommodationForm.valid) {
      this.errors = '';
      Object.keys(this.createAccommodationForm.controls).forEach((key) => {
        const controlErrors = this.createAccommodationForm.get(key)?.errors;
        if (controlErrors) {
          this.errors += formatErrors(key);
        }
      });
      console.log(this.errors);
      return;
    }
    this.accommodationsService.create(
      this.user?.id as string,
      this.user?.username as string,
      this.createAccommodationForm.value.name,
      this.createAccommodationForm.value.address,
      this.createAccommodationForm.value.city,
      this.createAccommodationForm.value.country,
      this.fromBooleanToConveniences(),
      this.createAccommodationForm.value.minNumOfVisitors as number,
      this.createAccommodationForm.value.maxNumOfVisitors as number,
      this.dateAvailabilities.value,
      this.createLocationCsv(),
    );
  }
}
