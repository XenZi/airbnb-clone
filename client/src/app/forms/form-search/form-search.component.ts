import { Component, HostListener } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { AccommodationsService } from 'src/app/services/accommodations-service/accommodations.service';


@Component({
  selector: 'app-form-search',
  templateUrl: './form-search.component.html',
  styleUrls: ['./form-search.component.scss'],
})
export class FormSearchComponent {
  searchForm: FormGroup;

  constructor(private formBuilder: FormBuilder,private accommodationsService: AccommodationsService) {
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
      this.searchForm.value.city,
      this.searchForm.value.dateRange,
      this.searchForm.value.guestsNumber
    );
    this.searchAccommodations()
    console.log(this.searchAccommodations())
  }
  searchAccommodations(){
    this.accommodationsService.search(this.searchForm.value.city as string,this.searchForm.value.country as string,this.searchForm.value.guestsNumber as string).subscribe({next:(data)=>{console.log(data.data)}})

  }

  
  

  

  


}
