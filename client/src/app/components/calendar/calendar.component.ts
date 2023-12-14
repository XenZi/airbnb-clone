import { Component, EventEmitter, Output } from '@angular/core';
import { CalendarModule } from 'primeng/calendar';
@Component({
  selector: 'app-calendar',
  templateUrl: './calendar.component.html',
  styleUrls: ['./calendar.component.scss'],
})
export class CalendarComponent {
  rangeDates: Date[] | undefined;
  @Output() datesChanged = new EventEmitter<Date[]>();
  constructor() {}

  onDatesChange(): void {
    if (this.rangeDates && this.rangeDates.length === 2) {
      const fullRange = this.getFullDateRange(
        this.rangeDates[0],
        this.rangeDates[1]
      );
      this.datesChanged.emit(fullRange);
    }
  }

  private getFullDateRange(startDate: Date, endDate: Date): Date[] {
    let dates = [];
    let currentDate = new Date(startDate);
    while (currentDate <= endDate) {
      dates.push(new Date(currentDate));
      currentDate.setDate(currentDate.getDate() + 1);
    }
    return dates;
  }
}
