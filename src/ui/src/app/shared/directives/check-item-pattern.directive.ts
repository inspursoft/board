import { Directive, Input } from '@angular/core';
import { ValidatorFn, AbstractControl, Validator, NG_VALIDATORS, Validators } from '@angular/forms';

const emailPattern = /^(([^<>()[\]\.,;:\s@\"]+(\.[^<>()[\]\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
const passwordPattern = /^(?=.*\d)(?=.*[a-z])(?=.*[A-Z])(?!.*\s).{8,20}$/;
const usernamePattern = /^[0-9a-z_]{4,35}$/;

@Directive({
  selector: '[checkItemPattern]',
  providers: [{
    provide: NG_VALIDATORS,
    useExisting: CheckItemPatternDirective,
    multi: true
  }]
})
export class CheckItemPatternDirective implements Validator {

  @Input() checkItemPattern;

  validate(control: AbstractControl): {[key: string]: any} {
    const value = control.value;
    switch (this.checkItemPattern) {
      case "email":
        return emailPattern.test(value) ? Validators.nullValidator : {'checkItemPattern': value};
      case "password":
        return passwordPattern.test(value) ? Validators.nullValidator : {'checkItemPattern': value};
      case "username":
        return usernamePattern.test(value) ? Validators.nullValidator : {'checkItemPattern': value};
    }
  }
}
