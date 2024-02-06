import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ShowMetricsComponent } from './show-metrics.component';

describe('ShowMetricsComponent', () => {
  let component: ShowMetricsComponent;
  let fixture: ComponentFixture<ShowMetricsComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [ShowMetricsComponent]
    });
    fixture = TestBed.createComponent(ShowMetricsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
