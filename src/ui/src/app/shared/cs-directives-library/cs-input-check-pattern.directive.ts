import { Directive, Input } from '@angular/core';
import { AbstractControl, Validators } from '@angular/forms';
import { CsInputComponent } from "../cs-components-library/cs-input/cs-input.component";

const emailPattern = /^(([^<>()[\]\.,;:\s@\"]+(\.[^<>()[\]\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
const passwordPattern = /^(?=.*\d)(?=.*[a-z])(?=.*[A-Z])(?!.*\s).{8,20}$/;
const usernamePattern = /^[0-9a-z_]{4,40}$/;

@Directive({
  selector: '[checkItemPatternEx]'
})
export class CsInputCheckPatternDirective {
  @Input() checkItemPatternEx = "";

  constructor(private csInputComponent: CsInputComponent) {
    this.csInputComponent.inputControl.setValidators(this.validateAction.bind(this))
  }

  validateAction(control: AbstractControl): {[key: string]: any} {
    const value = control.value;
    switch (this.checkItemPatternEx) {
      case "email":
        return emailPattern.test(value) ? Validators.nullValidator : {"checkItemPattern": "ACCOUNT.EMAIL_IS_ILLEGAL"};
      case "password":
        return passwordPattern.test(value) ? Validators.nullValidator : {"passwordPattern": "ACCOUNT.PASSWORD_FORMAT"};
      case "username":
        return usernamePattern.test(value) ? Validators.nullValidator : {"checkItemPattern": "ACCOUNT.USERNAME_ARE_NOT_IDENTICAL"};
    }
  }
}