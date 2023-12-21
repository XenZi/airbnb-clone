export interface Notifications {
  id: string;
  notifications: Notification[];
}

export interface Notification {
  text: string;
  createdAt: string;
  isOpened: boolean;
}
