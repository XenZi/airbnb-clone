import { Component, Input } from '@angular/core';
import { FormArray, FormBuilder, FormGroup } from '@angular/forms';
import { Router } from '@angular/router';
import { AccommodationsService } from 'src/app/services/accommodations-service/accommodations.service';

@Component({
  selector: 'app-form-filter',
  templateUrl: './form-filter.component.html',
  styleUrls: ['./form-filter.component.scss']
})
export class FormFilterComponent {

  filterForm: FormGroup;
  convenienceList: string[];

  @Input() cityCopy!: string;
  @Input() countryCopy!: string;
  @Input() numOfVisitorsCopy!: string;
  @Input() startDateCopy!: string;
  @Input() endDateCopy!: string;
  
  constructor(private fb: FormBuilder,private accommodationsService: AccommodationsService,private router: Router) {
    this.filterForm = this.fb.group({
      maxPrice: [''], // Assuming maxPrice is a FormControl
      conveniences: this.fb.array([]),
      distinguished: ['']
    });
  
    // Populate the convenienceList array with your convenience options
    this.convenienceList = ['WiFi', 'Kitchen', 'Air Conditioning', 'Free Parking', 'Pool'];
  
    // Initialize the conveniences form array with checkboxes
    this.initializeConveniencesCheckboxes();
  }
  
  initializeConveniencesCheckboxes() {
    const checkboxes = this.convenienceList.map(() => this.fb.control(false));
    this.filterForm.setControl('conveniences', this.fb.array(checkboxes));
  }
  
  


    onSubmit(e: Event) {
      this.searchAccommodations()
      console.log(this.cityCopy,this.countryCopy,this.numOfVisitorsCopy,this.startDateCopy,this.endDateCopy)
      
    }

    // buildConveniences(): FormArray {
    //   const arr = this.convenienceList.map(() => {
    //     return this.fb.control(false);
    //   });
    //   return this.fb.array(arr);
    // }

    searchAccommodations(){
      
      
      this.router.navigate(['/search'], {
        queryParams: {
          city: this.cityCopy,
          country: this.countryCopy,
          numOfVisitors: this.numOfVisitorsCopy,
          startDate:this.startDateCopy,
          endDate: this.endDateCopy,
          maxPrice:this.filterForm.value.maxPrice as string,
          conveniences:this.fromBooleanToConveniences(),
          distinguished:this.filterForm.value.distinguished as string

        },
      });
      // window.location.reload();
  
    }
  
   

    get convenienceFormArray(): FormArray {
      return this.filterForm.get('conveniences') as FormArray;
    }

    get conveniences(): FormArray {
      return this.filterForm.get('conveniences') as FormArray;
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

}

