import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ConfirmAccountPageComponent } from './confirm-account-page.component';

describe('ConfirmAccountPageComponent', () => {
  let component: ConfirmAccountPageComponent;
  let fixture: ComponentFixture<ConfirmAccountPageComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [ConfirmAccountPageComponent]
    });
    fixture = TestBed.createComponent(ConfirmAccountPageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
