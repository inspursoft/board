import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { MyInputTemplateComponent } from './my-input-template.component';

describe('MyInputTemplateComponent', () => {
  let component: MyInputTemplateComponent;
  let fixture: ComponentFixture<MyInputTemplateComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ MyInputTemplateComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(MyInputTemplateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
