import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { ModalComponent } from './components/modal/modal.component';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { FormLoginComponent } from './forms/form-login/form-login.component';
import { HeaderComponent } from './components/header/header.component';
import { DropdownComponent } from './components/dropdown/dropdown.component';
import { SimpleProfileMenuComponent } from './components/simple-profile-menu/simple-profile-menu.component';

@NgModule({
  declarations: [AppComponent, ModalComponent, FormLoginComponent, HeaderComponent, DropdownComponent, SimpleProfileMenuComponent],
  imports: [BrowserModule, AppRoutingModule, FontAwesomeModule],
  providers: [],
  bootstrap: [AppComponent],
})
export class AppModule {}
