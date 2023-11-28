import { TestBed } from '@angular/core/testing';

import { AccommodationsService } from './accommodations.service';

describe('AccommodationsService', () => {
  let service: AccommodationsService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(AccommodationsService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
