
import { Directive, Attribute, HostBinding, HostListener, Input,ElementRef } from '@angular/core';
import { ValidatorFn, AbstractControl, Validator, NG_VALIDATORS } from '@angular/forms';

export function validatorEmailFunction(): ValidatorFn {
  return (control: AbstractControl): { [key: string]: any } => {
    const emailPattern = /^(([^<>()[\]\.,;:\s@\"]+(\.[^<>()[\]\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
    if (control.value) {
      const regExp = new RegExp(emailPattern);
      if (!regExp.test(control.value)) {
        return { "emailFormat": true };
      }
    }
    return null;
  }
}

@Directive({
  selector: 'input[check-item-pattern]',
  providers: [{
    provide: NG_VALIDATORS,
    useExisting: CheckItemPatternDirective,
    multi: true
  }]
})

export class CheckItemPatternDirective implements Validator {

  @Input() checkItemType: string = "email"
  getValidatorFn(): ValidatorFn {
    if (this.checkItemType == "email") {
      return validatorEmailFunction();
    }
  }
 
  validate(control: AbstractControl): { [key: string]: any } {
    let validatorFn: ValidatorFn = this.getValidatorFn();
    return validatorFn(control);
  }
}