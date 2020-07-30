import { Directive, Input } from '@angular/core';
import { AbstractControl } from '@angular/forms';
import { InputExComponent } from 'board-components-library';
import { TranslateService } from '@ngx-translate/core';

const emailPattern = /^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$/;
const passwordPattern = /^(?=.*\d)(?=.*[a-z])(?=.*[A-Z])(?!.*\s).{8,20}$/;
const usernamePattern = /^[0-9a-z_]{4,35}$/;

@Directive({
  selector: '[appLibCheckPatternEx]'
})
export class LibCheckPatternExDirective {
  @Input() appLibCheckPatternEx = '';
  emailErrorMsg = '';
  passwordErrorMsg = '';
  usernameErrorMsg = '';

  constructor(private inputComponent: InputExComponent,
              private translateService: TranslateService) {
    this.inputComponent.inputValidatorFns.push(this.validateAction.bind(this));
    this.translateService.get('ACCOUNT.EMAIL_IS_ILLEGAL').subscribe(res => this.emailErrorMsg = res);
    this.translateService.get('ACCOUNT.PASSWORD_FORMAT').subscribe(res => this.passwordErrorMsg = res);
    this.translateService.get('ACCOUNT.USERNAME_ARE_NOT_IDENTICAL').subscribe(res => this.usernameErrorMsg = res);
  }

  validateAction(control: AbstractControl): { [key: string]: any } {
    const value = control.value;
    switch (this.appLibCheckPatternEx) {
      case 'email':
        return emailPattern.test(value) ? null : {checkItemPattern: this.emailErrorMsg};
      case 'password':
        return passwordPattern.test(value) ? null : {passwordPattern: this.passwordErrorMsg};
      case 'username':
        return usernamePattern.test(value) ? null : {checkItemPattern: this.usernameErrorMsg};
    }
  }
}
