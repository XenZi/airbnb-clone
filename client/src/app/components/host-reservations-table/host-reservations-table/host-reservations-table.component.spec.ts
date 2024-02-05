import { ComponentFixture, TestBed } from '@angular/core/testing';

import { HostReservationsTableComponent } from './host-reservations-table.component';

describe('HostReservationsTableComponent', () => {
  let component: HostReservationsTableComponent;
  let fixture: ComponentFixture<HostReservationsTableComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [HostReservationsTableComponent]
    });
    fixture = TestBed.createComponent(HostReservationsTableComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
