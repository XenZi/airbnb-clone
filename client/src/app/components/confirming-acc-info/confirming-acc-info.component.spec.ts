import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ConfirmingAccInfoComponent } from './confirming-acc-info.component';

describe('ConfirmingAccInfoComponent', () => {
  let component: ConfirmingAccInfoComponent;
  let fixture: ComponentFixture<ConfirmingAccInfoComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [ConfirmingAccInfoComponent]
    });
    fixture = TestBed.createComponent(ConfirmingAccInfoComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
