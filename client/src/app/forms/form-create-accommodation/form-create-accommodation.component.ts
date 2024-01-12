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
import { ICountry } from 'ngx-country-picker';
import { Observable } from 'rxjs';


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
  countries: any[] = [];
  selectedCountry: string = '';
  country!:string
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
      pictures: this.fb.array([]),
    });
  }
  

  
  
  ngOnInit() {
    this.user = this.userService.getLoggedUser();
    this.getCountriesData();
    console.log(this.getCountriesData)
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
  getCountriesData(): void {
    this.accommodationsService.getCountries()
      .subscribe((data: any[]) => {
        this.countries = data;
        console.log(this.countries); // Output the countries data to console
      });
  }

  onCountrySelection(event: Event): void {
    const selectedCountryName = (event.target as HTMLSelectElement).value;
    // Handle the selected country code here
    console.log('Selected Country Code:', selectedCountryName);
    this.country=selectedCountryName
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

  onFileSelected(event: any) {
    const files = event.target.files;
    const fileArray = this.createAccommodationForm.get('pictures') as FormArray;
    
    fileArray.clear();

    // Add each file to the form array
    for (let i = 0; i < 1; i++) {
      const fileControl = this.fb.control(files[i]);
      fileArray.push(fileControl);
      console.log(files[i])
      console.log(fileControl)
    }
    
    
  }

  createLocationCsv():string{
    const address=this.createAccommodationForm.value.address;
    const city=  this.createAccommodationForm.value.city;
    const countryCSV= this.country;
    console.log(address + ","+city +","+countryCSV)
    return address + ","+city +","+countryCSV
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
    const images = this.createAccommodationForm.get('pictures');

    console.log(images?.value)
    console.log(this.dateAvailabilities.value)
    this.accommodationsService.create(
      this.user?.id as string,
      this.user?.username as string,
      this.createAccommodationForm.value.name,
      this.createAccommodationForm.value.address,
      this.createAccommodationForm.value.city,
      this.country,
      this.fromBooleanToConveniences(),
      this.createAccommodationForm.value.minNumOfVisitors,
      this.createAccommodationForm.value.maxNumOfVisitors,
      this.dateAvailabilities.value,
      this.createLocationCsv(),
      images?.value,
    );
  }
}
