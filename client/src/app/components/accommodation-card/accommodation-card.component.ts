import { Component, Input } from '@angular/core';
import { Route, Router } from '@angular/router';
import { Observable } from 'rxjs';
import { Accommodation } from 'src/app/domains/entity/accommodation-model';
import { AccommodationsService } from 'src/app/services/accommodations-service/accommodations.service';

@Component({
  selector: 'app-accommodation-card',
  templateUrl: './accommodation-card.component.html',
  styleUrls: ['./accommodation-card.component.scss'],
})
export class AccommodationCardComponent {
  @Input() accommodation!: Accommodation;

  constructor(private router: Router) {}

  ngOnInit() {}
  clickOnCard() {
    this.router.navigate([`/accommodations/${this.accommodation.id}`]);
  }
}
