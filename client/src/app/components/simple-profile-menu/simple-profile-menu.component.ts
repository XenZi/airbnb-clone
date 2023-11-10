import { Component, HostListener, Input } from '@angular/core';
import { faL } from '@fortawesome/free-solid-svg-icons';
import { SimpleProfileMenuItem } from 'src/app/domains/model/simple-profile-menu-item.model';

@Component({
  selector: 'app-simple-profile-menu',
  templateUrl: './simple-profile-menu.component.html',
  styleUrls: ['./simple-profile-menu.component.scss'],
})
export class SimpleProfileMenuComponent {
  isClicked: boolean = false;
  items: SimpleProfileMenuItem[] = [
    {
      icon: 'fa-solid fa-right-to-bracket',
      title: 'Log in',
      action: () => {
        this.callLogin();
      },
    },
    {
      icon: 'fa-solid fa-user',
      title: 'Register',
      action: () => {
        this.callRegister();
      },
    },
  ];

  constructor() {}

  @HostListener('document:click', ['$event'])
  onClickOutside(event: any) {
    if (this.isClicked) {
      const list = [...event.target.classList];
      const filteredList = list.filter((x: string) =>
        x.includes('simple-profile-menu')
      );
      if (filteredList.length == 0) {
        this.isClicked = false;
      }
    }
  }

  clickBox() {
    this.isClicked = !this.isClicked;
  }

  callRegister() {}

  callLogin() {}
}
