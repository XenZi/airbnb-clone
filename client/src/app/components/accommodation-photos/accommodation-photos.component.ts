import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-accommodation-photos',
  templateUrl: './accommodation-photos.component.html',
  styleUrls: ['./accommodation-photos.component.scss']
})
export class AccommodationPhotosComponent {
  @Input()imageIds!: string[]

  startedImagesArray!: string[];

  ngOnInit() {
    this.startedImagesArray = this.startedImagesArray = this.imageIds.slice(1);
  }
}
