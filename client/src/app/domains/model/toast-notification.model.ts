import { ToastNotificationType } from '../enums/toast-notification-type.enum';

export default interface ToastNotification {
  title: string;
  message: string;
  type: ToastNotificationType;
}
