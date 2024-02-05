import { TestBed } from '@angular/core/testing';

import { PreviouseRouteService } from './previouse-route.service';

describe('PreviouseRouteService', () => {
  let service: PreviouseRouteService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(PreviouseRouteService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
