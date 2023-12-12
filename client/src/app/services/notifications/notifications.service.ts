import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { apiURL } from 'src/app/domains/constants';
import { Notifications } from 'src/app/domains/entity/notification.model';

@Injectable({
  providedIn: 'root',
})
export class NotificationsService {
  constructor(private http: HttpClient) {}

  getAllNotificationsForUser(id: string): Observable<any> {
    return this.http.get<any>(`${apiURL}/notifications/${id}`);
  }

  makeAllNotificationsReader(
    id: string,
    notification: Notifications
  ): Observable<any> {
    return this.http.put<any>(`${apiURL}/notifications/${id}`, {
      notification,
    });
  }
}
