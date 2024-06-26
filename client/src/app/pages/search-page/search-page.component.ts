import { Component } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { Accommodation } from 'src/app/domains/entity/accommodation-model';
import { AccommodationsService } from 'src/app/services/accommodations-service/accommodations.service';

@Component({
  selector: 'app-search-page',
  templateUrl: './search-page.component.html',
  styleUrls: ['./search-page.component.scss']
})
export class SearchPageComponent {
  accommodations!: Accommodation[];
  city!: string;
  country!: string;
  numOfVisitors!: string;
  startDate!:string;
  endDate!:string;
  maxPrice!:string;
  conveniences!:string[]
  distinguished!:string

  cityCopy!: string;
  countryCopy!: string;
  numOfVisitorsCopy!: string;
  startDateCopy!:string;
  endDateCopy!:string

  constructor(private accommodationService: AccommodationsService,private router: Router,
    private route: ActivatedRoute,
    ) {}

  ngOnInit() {
    this.route.queryParams.subscribe(params => {
      this.city = params['city'] || ''; // Retrieve the 'city' parameter value or set an empty string as default
      this.country = params['country'] || ''; // Retrieve the 'country' parameter value or set an empty string as default
      this.numOfVisitors = params['numOfVisitors'] || "1"; // Retrieve the 'numOfVisitors' parameter value as a number or set 0 as default
      this.startDate=params['startDate'] || ""
      this.endDate=params['endDate']||""
      this.maxPrice=params['maxPrice']||"0"
      this.conveniences=params['conveniences']||""
      this.distinguished=params['distinguished']||""
      // Once you have the query parameters, you can use them to perform the search
      this.cityCopy=this.city
      this.countryCopy=this.country
      this.numOfVisitorsCopy=this.numOfVisitors
      this.startDateCopy=this.startDateCopy
      this.endDateCopy=this.endDateCopy
    });
    this.loadSearchedAccommodations();
  }



  private loadSearchedAccommodations(): void {
    console.log(this.city,this.country,this.numOfVisitors)
    this.accommodationService.search(this.city as string,this.country as string,this.numOfVisitors as string,this.startDate as string,this.endDate as string,this.maxPrice as string,this.conveniences,this.distinguished).subscribe({
      next: (response) => {
        this.accommodations = response.data;
        console.log(this.accommodations)
        
      },
      error: (error) => {
        console.log(error);
      },
    });
    
  }

}
