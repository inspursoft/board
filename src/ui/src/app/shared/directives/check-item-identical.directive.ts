import { Directive, forwardRef, Attribute } from '@angular/core';
import { Validator, AbstractControl, NG_VALIDATORS } from '@angular/forms';

@Directive({
    selector: '[validateEqual][formControlName],[validateEqual][formControl],[validateEqual][ngModel]',
    providers: [
        { provide: NG_VALIDATORS, useExisting: forwardRef(() => CheckItemIdenticalDirective), multi: true }
    ]
})
export class CheckItemIdenticalDirective implements Validator {
  
    constructor(
      @Attribute('validateEqual') public validateEqual: string,
      @Attribute('reverse') public reverse: string
    ){}

    private get isReverse() {
        if (!this.reverse) return false;
        return this.reverse === 'true' ? true: false;
    }

    validate(control: AbstractControl): { [key: string]: any } {
        let value = control.value;
        let element = control.root.get(this.validateEqual);

        // value not equal
        if (element && value !== element.value && !this.isReverse) {
            return { 'validateEqual': false };
        }
        // value equal and reverse
        if (element && value === element.value && this.isReverse) {
            delete element.errors['validateEqual'];
            if (!Object.keys(element.errors).length) element.setErrors(null);
        }

        // value not equal and reverse
        if (element && value !== element.value && this.isReverse) {
            element.setErrors({ 'validateEqual': false });
        }

        return null;
    }
}