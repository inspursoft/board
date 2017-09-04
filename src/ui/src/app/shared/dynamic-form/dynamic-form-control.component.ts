import { Component, Input, OnInit } from '@angular/core';
import { FormGroup } from '@angular/forms';

import { FormControlBase } from './form-control-base';

@Component({
  selector: 'df-control',
  templateUrl: './dynamic-form-control.component.html'
})
export class DynamicFormControlComponent implements OnInit {
  
  @Input() control: FormControlBase<any>;
  @Input() form: FormGroup;

  editable: boolean = true;
  revertable: boolean = false;
  original: string;

  optionKey: string;
  subOptions: any;
  
  compoundValue: {};

  ngOnInit(): void {
    if(this.control.controlType === 'co-dropdown') {
      this._getSubOptions(this.control.value, 0);
      this.compoundValue = {
        '1': this.control.value,
        '2': this.subOptions[0]
      };
      this.control.value = this.compoundValue;
    }
    if(this.control.controlType === 'toggle-text') {
      this.original = this.control.value;
    }
  }

  _getSubOptions(value: string, index: number) {
    let option = this.control['options'][index];
    if(option.key === value) {
      this.subOptions = option.value;
    }
  }

  get isValid() {
    return this.form.controls[this.control.key].valid;
  }

  onChange(value) {
    switch(this.control.controlType){
    case 'co-dropdown':
      this.optionKey = value;
      for(let i in this.control['options']) {
        this._getSubOptions(value, +i);
      } 
      this.compoundValue = {
        '1': value,
        '2': this.subOptions[0]
      };
      this.control.value = this.compoundValue;
      break;
    case 'toggle-text':
      this.revertable = true;
    default:
      this.control.value = value;
    }
  }

  onSubChange(value) {
    this.compoundValue = {
      '1': this.compoundValue['1'],
      '2': value
    };
    this.control.value = this.compoundValue;
  }

  revert() {
    this.control.value = this.original;
    this.revertable = false;
  }
}