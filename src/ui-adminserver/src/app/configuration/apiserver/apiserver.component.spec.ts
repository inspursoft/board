import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ApiserverComponent } from './apiserver.component';

describe('ApiserverComponent', () => {
  let component: ApiserverComponent;
  let fixture: ComponentFixture<ApiserverComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ApiserverComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ApiserverComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
