import { Component, Input, OnInit, Output, EventEmitter } from '@angular/core';
import { FormGroup } from '@angular/forms';

import { FormControlBase } from './form-control-base';
import { TextInput } from './custom-control/text-input';
import { DynamicFormService } from './dynamic-form.service';

@Component({
  selector: 'dynamic-form',
  templateUrl: './dynamic-form.component.html'
})
export class DynamicFormComponent implements OnInit {
  
  maximumCount = 5;
  minimumCount = 1;

  @Input() formControls: FormControlBase<any>[] = [];

  @Output() addItem: EventEmitter<FormGroup> = new EventEmitter<FormGroup>();
  @Output() removeItem: EventEmitter<FormGroup> = new EventEmitter<FormGroup>();

  form: FormGroup;
  
  constructor(private dynamicFormService: DynamicFormService){}  

  ngOnInit() {
    this.form = this.dynamicFormService.toFormGroup(this.formControls);
  }

  addControl() {
    this.addItem.emit(this.form);
  }

  removeControl() {
    this.removeItem.emit(this.form);
  }

}