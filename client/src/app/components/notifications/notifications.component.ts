import { Component, HostListener } from '@angular/core';
import { Notifications } from 'src/app/domains/entity/notification.model';
import { User } from 'src/app/domains/entity/user-profile.model';
import { NotificationsService } from 'src/app/services/notifications/notifications.service';
import { UserService } from 'src/app/services/user/user.service';

@Component({
  selector: 'app-notifications',
  templateUrl: './notifications.component.html',
  styleUrls: ['./notifications.component.scss'],
})
export class NotificationsComponent {
  unreadedNotifications: boolean = true;
  isUserLogged: boolean = false;
  isClicked: boolean = false;
  notifications!: Notifications;
  user: User | null = null;
  constructor(
    private userService: UserService,
    private notificationsService: NotificationsService
  ) {}

  ngOnInit() {
    this.isUserLogged = this.userService.getLoggedUser() == null ? false : true;
    this.user = (this.userService.getLoggedUser() as unknown as User) ?? null;
    if (this.user) {
      this.notificationsService
        .getAllNotificationsForUser(this.user.id)
        .subscribe({
          next: (data) => {
            console.log(data);
            this.notifications = data.data;
          },
        });
    }
  }

  @HostListener('document:click', ['$event'])
  onClickOutside(event: any) {
    if (this.isClicked) {
      const list = [...event.target.classList];
      const filteredList = list.filter((x: string) =>
        x.includes('notifications')
      );
      if (filteredList.length == 0) {
        this.isClicked = false;
      }
    }
  }

  clickBox() {
    console.log(this.user);
    this.isClicked = !this.isClicked;
    if (this.isClicked) {
      this.notificationsService.makeAllNotificationsReader(
        this.user?.id as string,
        this.notifications
      );
    }
    console.log(this.isClicked);
  }
}
