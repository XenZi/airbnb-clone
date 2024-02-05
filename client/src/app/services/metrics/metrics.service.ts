import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { apiURL } from 'src/app/domains/constants';

@Injectable({
  providedIn: 'root'
})
export class MetricsService {

  constructor(  private http: HttpClient,
    ) { }

  

  public getMetrics(id: string, period: string): Observable<any> {
      return this.http.get(`${apiURL}/metrics_get/get/${id}/${period}`)
  }

     
  public joinedAt(userID: string, accommodationID: string, customUUID: string, joinedAt: string) {
    console.log("TEST VALUES: ", {
      userID,
      accommodationID,
      joinedAt,
      customUUID
    })
    this.http.post(`${apiURL}/metrics/joinedAt`, {
      userID,
      accommodationID,
      joinedAt,
      customUUID,
    }).subscribe({
      next: (data) => {
        console.log(data);
      },
      error: (err) => {
        console.log(err);
      }
    })
  }

  leftAt(userID: string, accommodationID: string,customUUID: string, leftAt: string) {
      this.http.post(`${apiURL}/metrics/leftAt`, {
        userID,
        accommodationID,
        customUUID,
        leftAt
      }).subscribe({
        next: (data) => {
          console.log(data)
        },
        error: (err) => {
          console.log(err);
        }
      })
  }
}
