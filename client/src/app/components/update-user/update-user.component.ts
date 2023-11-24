import { Component } from '@angular/core';
import { UpdateUserListItem } from 'src/app/domains/model/update-user-list-item.model';

@Component({
  selector: 'app-update-user',
  templateUrl: './update-user.component.html',
  styleUrls: ['./update-user.component.scss'],
})
export class UpdateUserComponent {
  currActiveIndex: number = 0;
  updateUserListOptions: UpdateUserListItem[] = [
    {
      text: 'Update password',
      action: () => {
        this.currActiveIndex = 0;
        this.changeCurrentActiveTo(0);
      },
    },
  ];

  constructor() {}

  isIndexActive(number: number): boolean {
    return this.currActiveIndex === number;
  }

  changeCurrentActiveTo(index: number) {}
}
