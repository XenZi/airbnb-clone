import { Component } from '@angular/core';
import { CalendarModule } from 'primeng/calendar';
@Component({
  selector: 'app-calendar',
  templateUrl: './calendar.component.html',
  styleUrls: ['./calendar.component.scss'],
})
export class CalendarComponent {
  rangeDates: Date[] | undefined;
  invalidDates!: Array<Date>;

  constructor() {
    let invalidDate = new Date();
    invalidDate.setDate(invalidDate.getDate() - 1);
    this.invalidDates = [new Date(), invalidDate];
  }
}
