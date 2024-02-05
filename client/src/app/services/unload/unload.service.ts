import { Injectable } from '@angular/core';
import { BehaviorSubject, filter } from 'rxjs';
import { MetricsService } from '../metrics/metrics.service';
import { HttpClient } from '@angular/common/http';
import { apiURL } from 'src/app/domains/constants';
import { LocalStorageService } from '../localstorage/local-storage.service';
import { NavigationEnd, NavigationStart, Router } from '@angular/router';

@Injectable({
  providedIn: 'root'
})
export class UnloadService {
  private routeHistory: string[] = [];

  constructor(private router: Router) {
    this.router.events
      .pipe(filter((event) => event instanceof NavigationEnd))
      .subscribe((event: any) => {
        this.routeHistory.push(event.urlAfterRedirects);
        console.log(this.routeHistory)
      });
  }

  getRouteHistory(): string[] {
    return this.routeHistory;
  }

}
