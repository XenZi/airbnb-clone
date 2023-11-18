import { Component, HostListener, Input } from '@angular/core';
import { faL } from '@fortawesome/free-solid-svg-icons';
import { SimpleProfileMenuItem } from 'src/app/domains/model/simple-profile-menu-item.model';
import { FormLoginComponent } from 'src/app/forms/form-login/form-login.component';
import { FormRegisterComponent } from 'src/app/forms/form-register/form-register.component';
import { ModalService } from 'src/app/services/modal/modal.service';

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

  constructor(private modalService: ModalService) {}

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

  callRegister() {
    this.modalService.open(FormRegisterComponent, 'Register');
  }

  callLogin() {
    this.modalService.open(FormLoginComponent, 'Login');
  }
}
