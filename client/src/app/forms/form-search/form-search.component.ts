import { Component, HostListener } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';

@Component({
  selector: 'app-form-search',
  templateUrl: './form-search.component.html',
  styleUrls: ['./form-search.component.scss'],
})
export class FormSearchComponent {
  searchForm: FormGroup;

  constructor(private formBuilder: FormBuilder) {
    this.searchForm = this.formBuilder.group({
      city: [''],
      country: [''],
      dateRange: [''],
      guestsNumber: [''],
    });
  }
  handleDateChange(rangeDates: Date[]) {
    const formattedRange = rangeDates.map((date) => date.toUTCString());
    this.searchForm.get('dateRange')?.setValue(formattedRange);
  }

  onSubmit(e: Event) {
    e.preventDefault();
    console.log(
      this.searchForm.value.city,
      this.searchForm.value.country,
      this.searchForm.value.dateRange,
      this.searchForm.value.guestsNumber
    );
  }
}
