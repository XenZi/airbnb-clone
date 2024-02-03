import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-filter',
  templateUrl: './filter.component.html',
  styleUrls: ['./filter.component.scss']
})
export class FilterComponent {

  isClicked: boolean = false;
  @Input() cityCopy!: string;
  @Input() countryCopy!: string;
  @Input() numOfVisitorsCopy!: string;
  @Input() startDateCopy!: string;
  @Input() endDateCopy!: string;

  

  onFilterClick() {
    this.isClicked = !this.isClicked;
  }

}
