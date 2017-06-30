import { Directive, Input, OnChanges, SimpleChanges } from '@angular/core';
import { NG_VALIDATORS, Validator, Validators, ValidatorFn, AbstractControl } from '@angular/forms';

@Directive({
  selector: '[checkItemExisting]',
  providers: [
    {provide: NG_VALIDATORS, useExisting: CheckItemExistingDirective, multi: true }
  ]
})
export class CheckItemExistingDirective implements Validator, OnChanges {

  @Input() checkItemExisting: string;
  @Input() targetName: string;

  valFn: ValidatorFn = Validators.nullValidator;

  ngOnChanges(changes: SimpleChanges): void {
    const change = changes['checkItemExisting'];
    if(change) {
      const val: string = change.currentValue;
      console.log(val);
      this.valFn = CheckItemExisting(val);
    } else {
      this.valFn = Validators.nullValidator;
    }
  }
  validate(control: AbstractControl): {[key: string]: any} {
    return this.valFn(control);
  }
}

export function CheckItemExisting(inputVal: string): ValidatorFn {
  return (control: AbstractControl): {[key: string]: any} => {
    const controlValue = control.value;
    return (controlValue === 'admin')? {checkItemExisting: {controlValue}}: null;
  }
}