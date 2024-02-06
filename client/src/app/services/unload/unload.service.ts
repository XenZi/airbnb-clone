import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable, filter } from 'rxjs';
import { MetricsService } from '../metrics/metrics.service';
import { HttpClient } from '@angular/common/http';
import { apiURL } from 'src/app/domains/constants';
import { LocalStorageService } from '../localstorage/local-storage.service';
import { NavigationEnd, NavigationStart, Router } from '@angular/router';

@Injectable({
  providedIn: 'root'
})
export class UnloadService {
  private previousUrl: BehaviorSubject<string> = new BehaviorSubject<string>("");
  public previousUrl$: Observable<string> = this.previousUrl.asObservable();

  constructor(private metricsService: MetricsService) { }

  setPreviousUrl(previousUrl: string) {
    if (previousUrl.includes("/accommodations")) {
      let currentStateOfLocalStorage = JSON.parse(sessionStorage.getItem("joins") as string);
      console.log(currentStateOfLocalStorage[currentStateOfLocalStorage.length -1])
      this.metricsService.leftAt(currentStateOfLocalStorage[currentStateOfLocalStorage.length -1].userID, currentStateOfLocalStorage[currentStateOfLocalStorage.length -1].accommodationID, currentStateOfLocalStorage[currentStateOfLocalStorage.length -1].customUUID,currentStateOfLocalStorage[currentStateOfLocalStorage.length -1].leftAt)
      currentStateOfLocalStorage = currentStateOfLocalStorage.splice(-1)
      sessionStorage.setItem('joins', JSON.stringify(currentStateOfLocalStorage))
    }
    this.previousUrl.next(previousUrl);
  }


}
