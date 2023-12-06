import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AccommodationPhotosComponent } from './accommodation-photos.component';

describe('AccommodationPhotosComponent', () => {
  let component: AccommodationPhotosComponent;
  let fixture: ComponentFixture<AccommodationPhotosComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [AccommodationPhotosComponent]
    });
    fixture = TestBed.createComponent(AccommodationPhotosComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
