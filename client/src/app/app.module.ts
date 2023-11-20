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
import { ButtonComponent } from './components/button/button.component';
import { ReactiveFormsModule } from '@angular/forms';
import { HTTP_INTERCEPTORS, HttpClientModule } from '@angular/common/http';
import { ToastNotificationComponent } from './components/toast/toast-notification/toast-notification.component';
import { ToastContainerComponent } from './components/toast/toast-container/toast-container.component';
import { FormRegisterComponent } from './forms/form-register/form-register.component';
import { RecaptchaFormsModule, RecaptchaModule } from 'ng-recaptcha';
import { ConfirmAccountPageComponent } from './pages/confirm-account-page/confirm-account-page.component';
import { BasePageComponent } from './pages/base-page/base-page.component';
import { ConfirmingAccInfoComponent } from './components/confirming-acc-info/confirming-acc-info.component';
@NgModule({
  declarations: [
    AppComponent,
    ModalComponent,
    FormLoginComponent,
    HeaderComponent,
    DropdownComponent,
    SimpleProfileMenuComponent,
    ButtonComponent,
    ToastNotificationComponent,
    ToastContainerComponent,
    FormRegisterComponent,
    ConfirmAccountPageComponent,
    BasePageComponent,
    ConfirmingAccInfoComponent,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    FontAwesomeModule,
    HttpClientModule,
    ReactiveFormsModule,
    RecaptchaModule,
    RecaptchaFormsModule,
  ],
  providers: [],
  bootstrap: [AppComponent],
})
export class AppModule {}
