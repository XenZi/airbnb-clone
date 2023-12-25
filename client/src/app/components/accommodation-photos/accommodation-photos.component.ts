import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-accommodation-photos',
  templateUrl: './accommodation-photos.component.html',
  styleUrls: ['./accommodation-photos.component.scss']
})
export class AccommodationPhotosComponent {
@Input()imageIds!: string[]


}
