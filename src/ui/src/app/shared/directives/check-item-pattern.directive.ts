
import { Directive, Input,Attribute } from '@angular/core';
import { ValidatorFn, AbstractControl, Validator, NG_VALIDATORS } from '@angular/forms';

function validatorEmailFunction(): ValidatorFn {
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

function validatorPasswordFunction(): ValidatorFn {
  return (control: AbstractControl): { [key: string]: any } => {
    const passwordPattern = /^(?=.*\d)(?=.*[a-z])(?=.*[A-Z])(?!.*\s).{8,20}$/;
    if (control.value) {
      const regExp = new RegExp(passwordPattern);
      if (!regExp.test(control.value)) {
        return { "passwordFormat": true };
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

  @Input('check-item-pattern') checkItemType: string = "email";
  getValidatorFn(): ValidatorFn {
    if (this.checkItemType == "email") {
      return validatorEmailFunction();
    }
    else if (this.checkItemType == "password"){
      return validatorPasswordFunction();
    }
  }
 
  validate(control: AbstractControl): { [key: string]: any } {
    let validatorFn: ValidatorFn = this.getValidatorFn();
    return validatorFn(control);
  }
}