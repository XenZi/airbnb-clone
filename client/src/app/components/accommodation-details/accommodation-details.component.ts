import { Component, Input } from '@angular/core';
import { Accommodation } from 'src/app/domains/entity/accommodation-model';

@Component({
  selector: 'app-accommodation-details',
  templateUrl: './accommodation-details.component.html',
  styleUrls: ['./accommodation-details.component.scss'],
})
export class AccommodationDetailsComponent {
  @Input() accommodation!: Accommodation;
}
