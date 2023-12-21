import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AccommodationCardComponent } from './accommodation-card.component';

describe('AccommodationCardComponent', () => {
  let component: AccommodationCardComponent;
  let fixture: ComponentFixture<AccommodationCardComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [AccommodationCardComponent]
    });
    fixture = TestBed.createComponent(AccommodationCardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
