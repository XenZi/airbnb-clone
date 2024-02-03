import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FormFilterComponent } from './form-filter.component';

describe('FormFilterComponent', () => {
  let component: FormFilterComponent;
  let fixture: ComponentFixture<FormFilterComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [FormFilterComponent]
    });
    fixture = TestBed.createComponent(FormFilterComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
