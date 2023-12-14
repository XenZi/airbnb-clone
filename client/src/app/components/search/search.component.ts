import { Component, HostListener } from '@angular/core';

@Component({
  selector: 'app-search',
  templateUrl: './search.component.html',
  styleUrls: ['./search.component.scss'],
})
export class SearchComponent {
  isClicked: boolean = false;

  onSearchClick() {
    this.isClicked = !this.isClicked;
  }
}
