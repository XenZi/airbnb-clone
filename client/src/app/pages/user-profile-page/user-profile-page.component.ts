import { Component } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { UserAuth } from 'src/app/domains/entity/user-auth.model';
import { User } from 'src/app/domains/entity/user-profile.model';
import { Role } from 'src/app/domains/enums/roles.enum';
import { FormUpdateUserProfileComponent } from 'src/app/forms/form-update-user-profile/form-update-user-profile.component';
import { ModalService } from 'src/app/services/modal/modal.service';
import { ProfileService } from 'src/app/services/profile/profile.service';
import { UserService } from 'src/app/services/user/user.service';

@Component({
  selector: 'app-user-profile-page',
  templateUrl: './user-profile-page.component.html',
  styleUrls: ['./user-profile-page.component.scss'],
})
export class UserProfilePageComponent {
  profileID!: string;
  user!: User;
  loggedUser!: UserAuth;
  constructor(
    private route: ActivatedRoute,
    private profileService: ProfileService,
    private modalService: ModalService,
    private userService: UserService,

  ) {}

  ngOnInit() {
    this.getUserID();
    this.profileService
      .getUserById(this.profileID as string)
      .subscribe((data) => {
        this.user = data.data;
      });
    if (this.user === undefined) {
      this.user = {
        id: 'id ssl mock',
        firstName: 'ime ssl mock',
        lastName: 'prezime ssl mock',
        email: 'mail ssl mock',
        residence: 'rezidencija ssl mock',
        role: Role.Host,
        username: 'username ssl mock',
        age: 25,
        distinguished: true,
        rating: 4.8,
      };
    }

    this.loggedUser = this.userService.getLoggedUser() as UserAuth;
  }

  updateClick() {
    this.callUpdateProfile();
  }

  deleteClick() {
    this.callDeleteProfile();
  }

  async callUpdateProfile() {
    // let foundUser = await this.user;
    // await this.modalService.open(
    //   FormUpdateUserProfileComponent,
    //   'Update Profile',
    //   {
    //     user: foundUser,
    //   }
    // );
    this.profileService
      .getUserById(this.profileID as string)
      .subscribe((data) => {
        this.modalService.open(
          FormUpdateUserProfileComponent,
          'Update Profile',
          {
            user: data.data,
          }
        );
      });
  }

  callDeleteProfile() {
    this.profileService.delete(this.profileID as string);
  }

  getUserID() {
    this.route.paramMap.subscribe((params) => {
      this.profileID = String(params.get('id'));
    });
  }

  getUserById() {
    this.profileService
      .getUserById(this.profileID as string)
      .subscribe((data) => {
        this.user = data.data;
      });
  }

}
