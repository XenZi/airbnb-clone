import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { ToastNotificationType } from 'src/app/domains/enums/toast-notification-type.enum';
import ToastNotification from 'src/app/domains/model/toast-notification.model';
@Injectable({
  providedIn: 'root',
})
export class ToastService {
  private toastSubject: BehaviorSubject<ToastNotification[]> =
    new BehaviorSubject<ToastNotification[]>([]);
  public toast$: Observable<ToastNotification[]> =
    this.toastSubject.asObservable();

  constructor() {}
  public showToast(
    title: string,
    message: string,
    type: ToastNotificationType
  ): void {
    const currentToasts = this.toastSubject.getValue();
    const newToast: ToastNotification = {
      title: title,
      message: message,
      type: type,
    };
    this.toastSubject.next([...currentToasts, newToast]);
  }

  public hideToast(toast: ToastNotification): void {
    const currentToasts = this.toastSubject.getValue();
    const updatedToasts = currentToasts.filter((t) => t !== toast);
    this.toastSubject.next(updatedToasts);
  }
}
