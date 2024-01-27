import { Component, EventEmitter, Input, Output } from '@angular/core';

@Component({
  selector: 'app-star-rating',
  templateUrl: './star-rating.component.html',
  styleUrls: ['./star-rating.component.scss'],
})
export class StarRatingComponent {
  ratings: number[] = [1, 2, 3, 4, 5];
  @Input() rating: number = 0;
  @Output() ratingChange = new EventEmitter<number>();
  stars: string[] = [];
  constructor() {}

  onStarClick(event: MouseEvent, starIndex: number): void {
    this.rating = starIndex;
    this.ratingChange.emit(this.rating);
  }
}
