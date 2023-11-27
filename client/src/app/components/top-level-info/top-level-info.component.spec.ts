import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TopLevelInfoComponent } from './top-level-info.component';

describe('TopLevelInfoComponent', () => {
  let component: TopLevelInfoComponent;
  let fixture: ComponentFixture<TopLevelInfoComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [TopLevelInfoComponent]
    });
    fixture = TestBed.createComponent(TopLevelInfoComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
