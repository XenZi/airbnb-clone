import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ShowRatesForUserComponent } from './show-rates-for-user.component';

describe('ShowRatesForUserComponent', () => {
  let component: ShowRatesForUserComponent;
  let fixture: ComponentFixture<ShowRatesForUserComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [ShowRatesForUserComponent]
    });
    fixture = TestBed.createComponent(ShowRatesForUserComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
