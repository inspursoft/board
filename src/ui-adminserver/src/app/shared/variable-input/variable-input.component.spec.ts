import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { VariableInputComponent } from './variable-input.component';

describe('VariableInputComponent', () => {
  let component: VariableInputComponent;
  let fixture: ComponentFixture<VariableInputComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ VariableInputComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(VariableInputComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
