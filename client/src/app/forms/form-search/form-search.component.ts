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
  startDate:string | undefined
  endDate:string|undefined
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
    console.log(this.searchForm.value.dateRange[0])
    
    
    console.log(this.startDate,this.endDate)
  }
  searchAccommodations(){
    this.startDate=this.formatingStartDate()
    this.endDate=this.formatingEndDate()

    if (this.startDate=="NaN-NaN-NaN"){
      this.startDate=""
    }
    if (this.endDate=="NaN-NaN-NaN"){
      this.endDate=""
    }
    
    this.accommodationsService.search(this.searchForm.value.city as string,this.searchForm.value.country as string,this.searchForm.value.guestsNumber as string,this.startDate as string,this.endDate as string).subscribe({next:(data)=>{console.log(data.data)}})

  }

  formatingStartDate():string{
    
    const startDateString: string = this.searchForm.value.dateRange[0];
    const date = new Date(startDateString);

    // Extract year, month, and day
    const year = date.getFullYear();
    // getMonth() returns 0 for January, 1 for February, ..., 11 for December
    const month = (date.getMonth() + 1).toString().padStart(2, '0'); // Adjust month by adding 1 and pad with '0' if necessary
    const day = date.getDate().toString().padStart(2, '0'); // Pad day with '0' if necessary

    // Form the yyyy-mm-dd format
    const formattedDate = `${year}-${month}-${day}`;

    console.log(formattedDate); // Output: "2023-12-15"
    return formattedDate
  }

  formatingEndDate():string{
    
    const startDateString: string = this.searchForm.value.dateRange[this.searchForm.value.dateRange.length - 1];
    const date = new Date(startDateString);

    // Extract year, month, and day
    const year = date.getFullYear();
    // getMonth() returns 0 for January, 1 for February, ..., 11 for December
    const month = (date.getMonth() + 1).toString().padStart(2, '0'); // Adjust month by adding 1 and pad with '0' if necessary
    const day = date.getDate().toString().padStart(2, '0'); // Pad day with '0' if necessary

    // Form the yyyy-mm-dd format
    const formattedDate = `${year}-${month}-${day}`;

    console.log(formattedDate); // Output: "2023-12-15"
    return formattedDate
  }



  
  

  

  


}
