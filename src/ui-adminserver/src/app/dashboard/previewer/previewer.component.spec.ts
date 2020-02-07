import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PreviewerComponent } from './previewer.component';

describe('PreviewerComponent', () => {
  let component: PreviewerComponent;
  let fixture: ComponentFixture<PreviewerComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ PreviewerComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PreviewerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
