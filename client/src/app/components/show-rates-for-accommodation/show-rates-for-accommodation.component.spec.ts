import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ShowRatesForAccommodationComponent } from './show-rates-for-accommodation.component';

describe('ShowRatesForAccommodationComponent', () => {
  let component: ShowRatesForAccommodationComponent;
  let fixture: ComponentFixture<ShowRatesForAccommodationComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [ShowRatesForAccommodationComponent]
    });
    fixture = TestBed.createComponent(ShowRatesForAccommodationComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
