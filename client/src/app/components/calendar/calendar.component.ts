import { Component } from '@angular/core';
import {
  CalendarView,
  CalendarEvent,
  CalendarMonthViewDay,
} from 'angular-calendar';
import {
  isSameMonth,
  startOfDay,
  addDays,
  subDays,
  endOfDay,
  addMonths,
  subMonths,
} from 'date-fns';

@Component({
  selector: 'app-calendar',
  templateUrl: './calendar.component.html',
  styleUrls: ['./calendar.component.scss'],
})
export class CalendarComponent {
  // view = CalendarView.Month;
  // viewDate: Date = new Date();
  // events: CalendarEvent[] = [];
  // activeDayIsOpen = false;
  // excludeDays: number[] = [0, 6]; // Sunday and Saturday
  // selectedRangeStart!: Date | null;
  // selectedRangeEnd!: Date | null;
  // constructor() {
  //   // Initialize your events here, including disabled dates
  //   this.initializeEvents();
  // }
  // dayClicked({ date, events }: { date: Date; events: CalendarEvent[] }): void {
  //   if (isSameMonth(date, this.viewDate)) {
  //     if (!this.selectedRangeStart) {
  //       // If no start date is set, set the start date
  //       this.selectedRangeStart = date;
  //     } else if (
  //       this.selectedRangeStart &&
  //       this.selectedRangeEnd &&
  //       (date < this.selectedRangeStart || date > this.selectedRangeEnd)
  //     ) {
  //       // If both start and end dates are set and the clicked date is outside the range,
  //       // reset the selection
  //       this.selectedRangeStart = date;
  //       this.selectedRangeEnd = null;
  //     } else {
  //       // Set the end date if the start date is already set
  //       this.selectedRangeEnd = date;
  //     }
  //     // If both start and end dates are set, perform any necessary action
  //     if (this.selectedRangeStart && this.selectedRangeEnd) {
  //       console.log(
  //         'Selected Range:',
  //         this.selectedRangeStart,
  //         'to',
  //         this.selectedRangeEnd
  //       );
  //       // You can perform additional actions with the selected range here
  //     }
  //   }
  // }
  // eventClicked(event: CalendarEvent): void {
  //   console.log('Event clicked', event);
  // }
  // navigateMonth(change: number): void {
  //   // Navigate to the next or previous month
  //   this.viewDate =
  //     change > 0 ? addMonths(this.viewDate, 1) : subMonths(this.viewDate, 1);
  //   this.initializeEvents(); // Update events for the new month
  // }
  // private initializeEvents(): void {
  //   // Your logic to initialize events, including disabled dates
  //   // Example: Disable specific dates
  //   const disabledDates = [
  //     new Date(2023, 10, 30),
  //     new Date(2023, 11, 1),
  //     new Date(2023, 11, 2),
  //     // Add more dates as needed
  //   ];
  //   // Populate events array with disabled dates
  //   this.events = disabledDates.map((date) => {
  //     return {
  //       start: date,
  //       title: 'Disabled Date',
  //       color: {
  //         primary: '#eee',
  //         secondary: '#FAE3E3',
  //       },
  //       actions: [],
  //     };
  //   });
  // }
  // dayModifier(day: CalendarMonthViewDay): void {
  //   // Highlight the selected range in the calendar
  //   if (
  //     this.selectedRangeStart &&
  //     this.selectedRangeEnd &&
  //     day.date >= startOfDay(this.selectedRangeStart) &&
  //     day.date <= endOfDay(this.selectedRangeEnd)
  //   ) {
  //     day.cssClass = 'cal-day-selected-range';
  //   }
  // }

  public dateValue: Date = new Date();
  dateList: string[] = ['25/12/2023', '26/12/2023', '27/12/2023'];
  constructor() {}
  disabledDate(args: any): void {
    console.log(args.date.getUTCDate());
    console.log(args.date.toLocaleDateString('en-GB'));
    if (this.dateList.includes(args.date.toLocaleDateString('en-GB'))) {
      //set 'true' to disable the weekends
      args.isDisabled = true;
    }
  }
}
