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
    console.log("JOINED AT EVENT")
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
    console.log("LEFT AT EVENT")
    this.http.post(`${apiURL}/metrics/leftAt`, {
        userID, accommodationID, customUUID, leftAt
      }).subscribe({
        next: (data) => {
          console.log(data);
        },
        error: (err) => {
          console.log(err);
        }
      });
  }
  rated(userID: string, accommodationID: string, ratedAt: string) {
    console.log("RATED EVENT")
    this.http.post(`${apiURL}/metrics/rated`, {
      userID, accommodationID, ratedAt
    }).subscribe({
      next: (data) => {
        console.log(data);
      },
      error: (err) => {
        console.log(err);
      }
    });
  }
  reserved(userID: string, accommodationID: string, reservedAt: string) {
    console.log("RESERVED EVENT")
    this.http.post(`${apiURL}/metrics/reserved`, {
      userID, accommodationID, reservedAt
    }).subscribe({
      next: (data) => {
        console.log(data);
      },
      error: (err) => {
        console.log(err);
      }
    });
  }
}