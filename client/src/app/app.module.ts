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
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { HTTP_INTERCEPTORS, HttpClientModule } from '@angular/common/http';
import { ToastNotificationComponent } from './components/toast/toast-notification/toast-notification.component';
import { ToastContainerComponent } from './components/toast/toast-container/toast-container.component';
import { FormRegisterComponent } from './forms/form-register/form-register.component';
import { RecaptchaFormsModule, RecaptchaModule } from 'ng-recaptcha';
import { ConfirmAccountPageComponent } from './pages/confirm-account-page/confirm-account-page.component';
import { BasePageComponent } from './pages/base-page/base-page.component';
import { ConfirmingAccInfoComponent } from './components/confirming-acc-info/confirming-acc-info.component';
import { ResetPasswordPageComponent } from './pages/reset-password-page/reset-password-page.component';
import { FormResetPasswordComponent } from './forms/form-reset-password/form-reset-password.component';
import { FormRequestResetPasswordComponent } from './forms/form-request-reset-password/form-request-reset-password.component';
import { RoleBasedPageComponent } from './pages/role-based-page/role-based-page.component';
import { UpdateUserComponent } from './components/update-user/update-user.component';
import { FormChangePasswordComponent } from './forms/form-change-password/form-change-password.component';
import { TokenInterceptor } from './interceptors/token.interceptor';
import { ReservationFormComponent } from './forms/form-create-reservation/form-create-reservation.component';
import { FormCreateAccommodationComponent } from './forms/form-create-accommodation/form-create-accommodation.component';
import { AccommodationCardComponent } from './components/accommodation-card/accommodation-card.component';
import { AccommodationDetailsPageComponent } from './pages/accommodation-details-page/accommodation-details-page.component';
import { TopLevelInfoComponent } from './components/top-level-info/top-level-info.component';
import { UserProfilePageComponent } from './pages/user-profile-page/user-profile-page.component';
import { FormUpdateUserProfileComponent } from './forms/form-update-user-profile/form-update-user-profile.component';
import { FormUpdateAccommodationComponent } from './forms/form-update-accommodation/form-update-accommodation.component';
import { AccommodationPhotosComponent } from './components/accommodation-photos/accommodation-photos.component';
import { AccommodationDetailsComponent } from './components/accommodation-details/accommodation-details.component';
import { ReservationComponent } from './components/reservation/reservation.component';
import { CalendarComponent } from './components/calendar/calendar.component';
import { CalendarModule } from 'primeng/calendar';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { FormUpdateCredentialsComponent } from './forms/form-update-credentials/form-update-credentials.component';
import { NotificationsComponent } from './components/notifications/notifications.component';
import { SearchComponent } from './components/search/search.component';
import { FormSearchComponent } from './forms/form-search/form-search.component';
import { SearchPageComponent } from './pages/search-page/search-page.component';
import { UserReservationsTableComponent } from './components/user-reservations-table/user-reservations-table.component';
import { FormRateAccommodationComponent } from './forms/form-rate-accommodation/form-rate-accommodation.component';
import { StarRatingComponent } from './components/star-rating/star-rating.component';
import { ShowRatesForAccommodationComponent } from './components/show-rates-for-accommodation/show-rates-for-accommodation.component';
import { ShowRatesForUserComponent } from './components/show-rates-for-user/show-rates-for-user.component';
import { FormCreateRatingForUserComponent } from './forms/form-create-rating-for-user/form-create-rating-for-user.component';
import { FilterComponent } from './components/filter/filter.component';
import { FormFilterComponent } from './forms/form-filter/form-filter.component';
import { UnloadService } from './services/unload/unload.service';
import { HostReservationsTableComponent } from './components/host-reservations-table/host-reservations-table/host-reservations-table.component';
import { ShowMetricsComponent } from './components/show-metrics/show-metrics.component';
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
    ResetPasswordPageComponent,
    FormResetPasswordComponent,
    FormRequestResetPasswordComponent,
    RoleBasedPageComponent,
    UpdateUserComponent,
    FormChangePasswordComponent,
    ReservationFormComponent,
    FormCreateAccommodationComponent,
    AccommodationCardComponent,
    AccommodationDetailsPageComponent,
    TopLevelInfoComponent,
    UserProfilePageComponent,
    FormUpdateUserProfileComponent,
    FormUpdateAccommodationComponent,
    AccommodationPhotosComponent,
    AccommodationDetailsComponent,
    ReservationComponent,
    CalendarComponent,
    FormUpdateCredentialsComponent,
    NotificationsComponent,
    SearchComponent,
    FormSearchComponent,
    SearchPageComponent,
    UserReservationsTableComponent,
    FormRateAccommodationComponent,
    StarRatingComponent,
    ShowRatesForAccommodationComponent,
    ShowRatesForUserComponent,
    FormCreateRatingForUserComponent,
    FilterComponent,
    FormFilterComponent,
    HostReservationsTableComponent,
    ShowMetricsComponent,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    FontAwesomeModule,
    HttpClientModule,
    ReactiveFormsModule,
    RecaptchaModule,
    RecaptchaFormsModule,
    FormsModule,
    CalendarModule,
    BrowserAnimationsModule,
  ],
  providers: [
    {
      provide: HTTP_INTERCEPTORS,
      useClass: TokenInterceptor,
      multi: true,
    },
    UnloadService
  ],
  bootstrap: [AppComponent],
})
export class AppModule {}
