import { AfterViewInit, Component, EventEmitter, Input, Output } from "@angular/core";
import { FormControl, FormGroup, ValidationErrors, ValidatorFn, Validators } from "@angular/forms";
import { AbstractControl } from "@angular/forms/src/model";

@Component({
  selector: 'cs-input-dropdown',
  templateUrl: './cs-input-dropdown.component.html',
  styleUrls: ['./cs-input-dropdown.component.css']
})

export class CsInputDropdownComponent implements AfterViewInit {
  @Input() inputPlaceholder = "";
  @Input() inputIsRequired = false;
  @Input() inputMax = 0;
  @Input() inputMin = 0;
  @Input() inputWidth = 100;
  @Input() inputConflictNumbers: Array<number>;
  @Input() validatorMessage: Array<{ validatorKey: string, validatorMessage: string }>;
  @Output() onCommitEvent: EventEmitter<number>;

  @Input() set inputDisabled(value: boolean) {
    value ? this.inputControl.disable() : this.inputControl.enable();
  }

  @Input() set inputValue(value: number) {
    this.inputControl.setValue(value);
  }

  inputFormGroup: FormGroup;
  inputControl: FormControl;
  dropdownMenuList: Array<number>;
  inputValidatorFns: Array<ValidatorFn>;
  inputValidatorMessageParam = '';

  constructor() {
    this.inputConflictNumbers = Array<number>();
    this.inputValidatorFns = Array<ValidatorFn>();
    this.dropdownMenuList = Array<number>();
    this.validatorMessage = Array<{ validatorKey: string, validatorMessage: string }>();
    this.onCommitEvent = new EventEmitter<number>();
    this.inputControl = new FormControl(0, {updateOn: 'change'});
    this.inputFormGroup = new FormGroup({inputControl: this.inputControl}, {updateOn: 'change'});
  }

  ngAfterViewInit() {
    if (this.inputIsRequired) {
      this.inputValidatorFns.push(Validators.required);
    }
    if (this.inputMax > 0) {
      this.inputValidatorFns.push(Validators.max(this.inputMax));
    }
    this.inputValidatorFns.push(this.checkMin.bind(this));
    this.inputValidatorFns.push(this.checkAlreadyUsed.bind(this));
    this.inputControl.setValidators(this.inputValidatorFns);
  }

  get valid(): boolean {
    if ((this.inputIsRequired && this.inputControl.value < this.inputMin) ||
      (!this.inputIsRequired && this.inputControl.value != 0 && this.inputControl.value < this.inputMin)) {
      return false
    }
    return this.inputControl.valid;
  }

  checkAlreadyUsed(c: AbstractControl): ValidationErrors | null {
    if (this.inputConflictNumbers.indexOf(c.value) > -1) {
      return {'inUsed': 'inUsed'}
    } else {
      return null;
    }
  }

  checkMin(c: AbstractControl): ValidationErrors | null {
    if (this.inputIsRequired && c.value < this.inputMin) {
      return {'min': 'min'}
    } else if (!this.inputIsRequired && c.value != 0 && c.value < this.inputMin) {
      return {'min': 'min'}
    } else {
      return null;
    }
  }

  validNumberStart(validNumber: number): string {
    let validNumberStr = `${validNumber}`;
    let find = `${this.inputControl.value}`;
    let index = validNumberStr.indexOf(find);
    return validNumberStr.slice(0, index);
  }

  validNumberEnd(validNumber: number): string {
    let validNumberStr = `${validNumber}`;
    let find = `${this.inputControl.value}`;
    let index = validNumberStr.indexOf(find) + find.length;
    return validNumberStr.slice(index);
  }

  getNextMinValidNumber(baseNumber: number, increase: number): number {
    let result = baseNumber + increase;
    if (result < this.inputMin || result > this.inputMax || !this.inputControl.value) {
      return 0;
    }
    if (this.inputConflictNumbers.find(value => value == result)) {
      return this.getNextMinValidNumber(baseNumber, increase + 1)
    }
    let strResult = `${result}`;
    let find = `${this.inputControl.value}`;
    if (strResult.indexOf(find) == -1) {
      return this.getNextMinValidNumber(baseNumber, increase + 1)
    }
    return result;
  }

  getValidatorMessage(errors: ValidationErrors): string {
    this.inputValidatorMessageParam = "";
    let result: string = "";
    this.validatorMessage.forEach(value => {
      if (errors[value.validatorKey]) {
        result = value.validatorMessage;
      }
    });
    if (result == "") {
      if (errors["required"]) {
        result = "ERROR.INPUT_REQUIRED"
      } else if (errors["max"]) {
        result = `ERROR.INPUT_MAX_VALUE`;
        this.inputValidatorMessageParam = `:${this.inputMax}`
      } else if (errors["min"]) {
        result = `ERROR.INPUT_MIN_VALUE`;
        this.inputValidatorMessageParam = `:${this.inputMin}`
      } else if (Object.keys(errors).length > 0) {
        result = errors[Object.keys(errors)[0]];
      }
    }
    return result;
  }

  onBlurEvent() {
    if (this.inputControl.value == null) {
      this.inputControl.setValue(0);
    }
  }

  onChangeEvent() {
    this.dropdownMenuList.splice(0, this.dropdownMenuList.length);
    for (let i = 0; i < 5; i++) {
      let startMin = this.dropdownMenuList.length > 0 ? this.dropdownMenuList[this.dropdownMenuList.length - 1] : this.inputMin - 1;
      let validValue = this.getNextMinValidNumber(startMin, 1);
      if (validValue > 0) {
        this.dropdownMenuList.push(validValue);
      }
    }
    if (this.inputControl.valid) {
      this.onCommitEvent.emit(this.inputControl.value);
    }
  }

  onDropdownItemCLick(selectNumber: number) {
    this.inputControl.setValue(selectNumber);
    this.onChangeEvent();
  }

  checkInputSelf() {
    if (this.inputControl.enabled && (this.inputControl.touched || this.inputIsRequired)) {
      this.inputControl.markAsTouched({onlySelf: true});
      this.inputControl.updateValueAndValidity();
    }
  }
}