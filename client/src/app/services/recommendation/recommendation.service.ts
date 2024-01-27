import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { apiURL } from 'src/app/domains/constants';

@Injectable({
  providedIn: 'root',
})
export class RecommendationService {
  constructor(private http: HttpClient) {}

  getAllRecommendedByRating(): Observable<any> {
    return this.http.get<any>(`${apiURL}/recommendations/top-rated`);
  }

  getAllRecommendedByUserID(userID: string): Observable<any> {
    return this.http.get<any>(`${apiURL}/recommendations/${userID}`);
  }
}
