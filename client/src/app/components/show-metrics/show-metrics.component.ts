import { Component, EventEmitter, Input, Output } from '@angular/core';

@Component({
  selector: 'app-show-metrics',
  templateUrl: './show-metrics.component.html',
  styleUrls: ['./show-metrics.component.scss']
})
export class ShowMetricsComponent {
  @Input() numberOfRatings!: number;
  @Input() numberOfReservations!: number;
  @Input() numberOfVisits!: number;
  @Input() onScreenTime!: number;
  @Input() changeValueOfUpper!: Function;
  @Output() handleChangeValue: EventEmitter<any> = new EventEmitter();
  currentActiveSelection: number = 0;

  changeCurrentActiveSelection(num: number) {
    this.handleChangeValue.emit(num);
    this.currentActiveSelection = num;
    this.changeValueOfUpper(num);
  }
}
