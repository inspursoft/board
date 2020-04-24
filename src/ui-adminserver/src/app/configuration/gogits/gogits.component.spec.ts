import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { GogitsComponent } from './gogits.component';

describe('GogitsComponent', () => {
  let component: GogitsComponent;
  let fixture: ComponentFixture<GogitsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ GogitsComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(GogitsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
