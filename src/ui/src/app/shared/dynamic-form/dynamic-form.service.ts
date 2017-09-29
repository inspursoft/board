import { Injectable } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';

import { FormControlBase } from './form-control-base'

@Injectable()
export class DynamicFormService {
  
  constructor(){}

  getControls(control: FormControlBase<any>) {
    return [control];
  }

  toFormGroup(controls: FormControlBase<any>[]) {
    let group: any = [];
    controls.forEach(control => {
      group[control.key] = control.required ? 
        new FormControl(control.value || '', Validators.required) 
        : new FormControl(control.value || '');
    });
    return new FormGroup(group);
  }
}