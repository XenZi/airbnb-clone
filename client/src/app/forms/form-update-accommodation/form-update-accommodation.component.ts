import { Component, Input } from '@angular/core';
import { FormArray, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { timeout } from 'rxjs';
import { Accommodation } from 'src/app/domains/entity/accommodation-model';
import { ToastNotificationType } from 'src/app/domains/enums/toast-notification-type.enum';
import { AccommodationsService } from 'src/app/services/accommodations-service/accommodations.service';
import { ToastService } from 'src/app/services/toast/toast.service';
import { formatErrors } from 'src/app/utils/formatter.utils';

@Component({
  selector: 'app-form-update-accommodation',
  templateUrl: './form-update-accommodation.component.html',
  styleUrls: ['./form-update-accommodation.component.scss'],
})
export class FormUpdateAccommodationComponent {
  @Input() accommodationID!: string;
  convenienceList = [
    'WiFi',
    'Kitchen',
    'Air Conditioning',
    'Free Parking',
    'Pool',
  ];
  updateAccommodationForm: FormGroup;
  accommodation!: Accommodation;
  errors: string = '';
  isCaptchaValidated: boolean = false;

  constructor(
    private accommodationsService: AccommodationsService,
    private formBuilder: FormBuilder,
    private toastService: ToastService
  ) {
    this.updateAccommodationForm = this.formBuilder.group({
      name: [''],
      address: [''],
      city: [''],
      conveniences: this.buildConveniences([]),
      minNumOfVisitors: [''],
      maxNumOfVisitors: [''],
    });
  }

  ngOnInit() {
    this.accommodationsService
      .getAccommodationById(this.accommodationID as string)
      .subscribe((data) => {
        this.accommodation = data.data;
        console.log(this.accommodation);
        this.updateAccommodationForm = this.formBuilder.group({
          name: [this.accommodation.name, Validators.required],
          address: [this.accommodation.address, Validators.required],
          city: [this.accommodation.city, Validators.required],
          conveniences: this.buildConveniences(this.accommodation.conveniences),
          minNumOfVisitors: [
            this.accommodation.minNumOfVisitors,
            Validators.required,
          ],
          maxNumOfVisitors: [
            this.accommodation.maxNumOfVisitors,
            Validators.required,
          ],
        });
      });
  }

  buildConveniences(selectedConveniences: string[]): FormArray {
    const arr = this.convenienceList.map((convenience) => {
      return this.formBuilder.control(
        selectedConveniences.includes(convenience)
      );
    });
    return this.formBuilder.array(arr);
  }
  fromBooleanToConveniences(): string[] {
    let convArray: string[] = [];

    this.convenienceFormArray.value.forEach((el: boolean, i: number) => {
      if (el === true && this.convenienceList[i]) {
        convArray.push(this.convenienceList[i]);
      }
    });
    console.log(convArray);
    return convArray;
  }
  get convenienceFormArray(): FormArray {
    return this.updateAccommodationForm.get('conveniences') as FormArray;
  }

  onSubmit(e: Event) {
    e.preventDefault();

    if (!this.updateAccommodationForm.valid) {
      console.log('bilo sta1');
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
    console.log(this.updateAccommodationForm.value.conveniences);
    this.accommodationsService.update(
      this.accommodationID,
      this.updateAccommodationForm.value.name,
      this.updateAccommodationForm.value.address,
      this.updateAccommodationForm.value.city,
      this.fromBooleanToConveniences(),
      this.updateAccommodationForm.value.minNumOfVisitors as number,
      this.updateAccommodationForm.value.maxNumOfVisitors as number
    );
  }
}
