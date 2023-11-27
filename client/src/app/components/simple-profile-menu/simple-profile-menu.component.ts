import { Component, HostListener, Input } from '@angular/core';
import { faL } from '@fortawesome/free-solid-svg-icons';
import { SimpleProfileMenuItem } from 'src/app/domains/model/simple-profile-menu-item.model';
import { FormLoginComponent } from 'src/app/forms/form-login/form-login.component';
import { FormRegisterComponent } from 'src/app/forms/form-register/form-register.component';
import { ModalService } from 'src/app/services/modal/modal.service';
import { UpdateUserComponent } from '../update-user/update-user.component';
import { UserService } from 'src/app/services/user/user.service';
import { FormCreateAccommodationComponent } from 'src/app/forms/form-create-accommodation/form-create-accommodation.component';
import { Role } from 'src/app/domains/enums/roles.enum';
import { AuthService } from 'src/app/services/auth-service/auth.service';

@Component({
  selector: 'app-simple-profile-menu',
  templateUrl: './simple-profile-menu.component.html',
  styleUrls: ['./simple-profile-menu.component.scss'],
})
export class SimpleProfileMenuComponent {
  isClicked: boolean = false;
  isUserLogged: boolean = false;
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
  loggedItems: SimpleProfileMenuItem[] = [
    {
      icon: 'fa-solid fa-user',
      title: 'Update profile',
      action: () => {
        this.callUpdateProfile();
      },
    },{
      icon: 'fa-solid fa-user',
      title: 'Add accommodation',
      action: () => {
        this.callNewAccommodation();
      },
    },
    {
      icon: 'fa-solid fa-door-open',
      title: 'Log out',
      action: () => {
        this.callLogout();
      },
    },
  ];
  hostItems: SimpleProfileMenuItem[] = [
    {
      icon: 'fa-solid fa-user',
      title: 'GGGGG',
      action: () => {
        this.callUpdateProfile();
      },
    },
  ];
  guestItems: SimpleProfileMenuItem[] = [
    {
      icon: 'fa-solid fa-user',
      title: 'HHHHHH',
      action: () => {
        console.log('SADSAD');
      },
    },
  ];
  constructor(
    private modalService: ModalService,
    private userService: UserService,
    private authService: AuthService
  ) {}

  ngOnInit() {
    this.isUserLogged = this.userService.getLoggedUser() == null ? false : true;
    this.items = this.isUserLogged
      ? this.userService.getLoggedUser()?.role == Role.Host
        ? [...this.loggedItems, ...this.hostItems]
        : [...this.loggedItems, ...this.guestItems]
      : [...this.items];
  }

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

  callUpdateProfile() {
    this.modalService.open(UpdateUserComponent, 'Update your profile');
  }
  callNewAccommodation(){
    this.modalService.open(FormCreateAccommodationComponent, 'Create accommodation');

  }
  callLogout() {
    this.authService.logout();
  } 
}
