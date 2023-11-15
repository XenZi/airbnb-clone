import { Component } from '@angular/core';
import ToastNotification from 'src/app/domains/model/toast-notification.model';
import { ToastService } from 'src/app/services/toast/toast.service';

@Component({
  selector: 'app-toast-container',
  templateUrl: './toast-container.component.html',
  styleUrls: ['./toast-container.component.scss'],
})
export class ToastContainerComponent {
  toasts: ToastNotification[] = [];

  constructor(private toastService: ToastService) {}

  ngOnInit(): void {
    this.toastService.toast$.subscribe((toasts: ToastNotification[]) => {
      this.toasts = toasts;
    });
  }
}
