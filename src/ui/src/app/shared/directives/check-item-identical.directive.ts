import { Directive, Input, OnChanges, SimpleChanges } from '@angular/core';
import { NG_VALIDATORS, Validator, Validators, ValidatorFn,  AbstractControl } from '@angular/forms';

@Directive({
  selector: '[checkItemIdentical]',
  providers: [{ provide: NG_VALIDATORS, useExisting: CheckItemIdenticalDirective, multi: true }]
})
export class CheckItemIdenticalDirective implements Validator, OnChanges {
  @Input() comparison: string;
  valFn:ValidatorFn = Validators.nullValidator;

  ngOnChanges(changes: SimpleChanges): void {
    const change = changes['comparison'];
    if(change) {
      const value: string = change.currentValue;
      this.valFn = CheckItemIdenticalValidator(value);  
    } else {
      this.valFn = Validators.nullValidator;
    }
  }

  validate(control: AbstractControl): {[key: string]: any} {
    return this.valFn(control);
  }
}

export function CheckItemIdenticalValidator(comparison: string): ValidatorFn {
  return (control: AbstractControl): {[key: string]: any} => {
    const inputValue = control.value;
    return (comparison && comparison === inputValue) ? null : {'checkItemIdentical':  {inputValue}};
  };
}