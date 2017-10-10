/**
 * Created by liyanq on 9/11/17.
 */
import { Component, Input, Output, EventEmitter, OnInit, ViewChild } from "@angular/core"
import { FormControl, FormGroup, ValidationErrors, ValidatorFn, Validators } from "@angular/forms";

export enum CsInputStatus{isView = 0, isEdit = 1}
export enum CsInputType{itWithInput, itWithNoInput, itOnlyWithInput}
export enum CsInputFiledType{iftString, iftNumber}
export type CsInputSupportType = string | number
export class CsInputFiled {
  constructor(public status: CsInputStatus,
              public defaultValue: CsInputSupportType,
              public value: CsInputSupportType) {
  }
}

const InputPatternNumber: RegExp = /^[1-9]\d*$/;
@Component({
  selector: "cs-input",
  templateUrl: "./cs-input.component.html",
  styleUrls: ["./cs-input.component.css"]
})
export class CsInputComponent implements OnInit {
  _isDisabled: boolean = false;
  inputErrors: ValidationErrors;
  inputFormGroup: FormGroup;
  inputControl: FormControl;
  inputValidatorFns: Array<ValidatorFn>;
  inputValidatorMessageParam: string;
  inputField: CsInputFiled;
  @ViewChild("input") inputHtml;
  @Input() inputLabel: string = "";
  @Input() inputFiledType: CsInputFiledType = CsInputFiledType.iftString;
  @Input() inputIsRequired: boolean = false;
  @Input() inputPattern: RegExp;
  @Input() inputMaxlength: number = 0;
  @Input() inputMinlength: number = 0;
  @Input() inputMax: number = 0;
  @Input() inputMin: number = 0;
  @Input() inputType: CsInputType = CsInputType.itWithInput;
  @Input() customerValidatorFns: Array<ValidatorFn>;
  @Input() validatorMessage: Array<{validatorKey: string, validatorMessage: string}>;

  constructor() {
    this.inputValidatorFns = Array<ValidatorFn>();
  }

  ngOnInit() {
    this.inputControl = new FormControl({
      value: this.SimpleFiled,
      disabled: this.isDisabled || this.inputType == CsInputType.itWithNoInput
    });
    this.inputFormGroup = new FormGroup({inputControl: this.inputControl});
    if (this.inputFiledType == CsInputFiledType.iftNumber) {
      this.inputValidatorFns.push(Validators.pattern(InputPatternNumber));
    }
    if (this.inputIsRequired) {
      this.inputValidatorFns.push(Validators.required);
    }
    if (this.inputMaxlength > 0) {
      this.inputValidatorFns.push(Validators.maxLength(this.inputMaxlength));
    }
    if (this.inputMinlength > 0) {
      this.inputValidatorFns.push(Validators.minLength(this.inputMinlength));
    }
    if (this.inputMax > 0) {
      this.inputValidatorFns.push(Validators.max(this.inputMax));
    }
    if (this.inputMin > 0) {
      this.inputValidatorFns.push(Validators.min(this.inputMin));
    }
    if (this.inputPattern) {
      this.inputValidatorFns.push(Validators.pattern(this.inputPattern));
    }
    if (this.customerValidatorFns) {
      this.customerValidatorFns.forEach(value => {
        this.inputValidatorFns.push(value);
      })
    }
  }

  @Input("simpleFiled")
  set SimpleFiled(value: CsInputSupportType) {
    this.inputField = new CsInputFiled(
      CsInputStatus.isView, value, value
    );
  }

  get SimpleFiled(): CsInputSupportType {
    return this.inputField.value;
  }

  @Input("disabled")
  set isDisabled(value: boolean) {
    this._isDisabled = value;
    if (value) {
      this.inputField.status = CsInputStatus.isView;
    }
    if (this.inputControl) {
      this.inputControl.reset({
        value: this.SimpleFiled,
        disabled: value
      });
    }
  }

  get isDisabled() {
    return this._isDisabled;
  }

  public get valid(): boolean {
    let notEmpty = this.inputField.value != 0 && this.inputField.value != "";
    let isInView = this.inputField.status == CsInputStatus.isView;
    if (this.inputIsRequired) {
      return notEmpty && isInView;
    } else {
      return isInView;
    }
  }

  getValidatorMessage(errors: ValidationErrors): string {
    this.inputValidatorMessageParam = "";
    let result: string = "";
    if (this.validatorMessage) {
      this.validatorMessage.forEach(value => {
        if (errors[value.validatorKey]) {
          result = value.validatorMessage;
        }
      });
    }
    if (result == "") {
      if (errors["required"]) {
        result = "ERROR.INPUT_REQUIRED"
      } else if (errors["pattern"] && this.inputFiledType == CsInputFiledType.iftNumber) {
        result = "ERROR.INPUT_ONLY_NUMBER"
      } else if (errors["pattern"] && this.inputFiledType == CsInputFiledType.iftString) {
        result = "ERROR.INPUT_PATTERN"
      } else if (errors["maxlength"]) {
        result = `ERROR.INPUT_MAX_LENGTH`;
        this.inputValidatorMessageParam = `:${this.inputMaxlength}`
      } else if (errors["minlength"]) {
        result = `ERROR.INPUT_MIN_LENGTH`;
        this.inputValidatorMessageParam = `:${this.inputMinlength}`
      } else if (errors["max"]) {
        result = `ERROR.INPUT_MAX_VALUE`;
        this.inputValidatorMessageParam = `:${this.inputMax}`
      } else if (errors["min"]) {
        result = `ERROR.INPUT_MIN_VALUE`;
        this.inputValidatorMessageParam = `:${this.inputMin}`
      }
    }
    return result;
  }

  public checkValueByHost() {
    this.onCheckClick();
  }

  get isPassValidator(): boolean {
    let result = true;
    this.inputErrors = null;
    this.inputValidatorFns.forEach(value => {
      let error = value(this.inputControl);
      if (error) {
        result = false;
        this.inputErrors = error;
      }
    });
    return result;
  }

  onInputKeyPressEvent(event: KeyboardEvent) {
    if (event.keyCode == 13) {
      this.onCheckClick();
    }
  }

  @Output("onEdit") onEditEvent: EventEmitter<any> = new EventEmitter<any>();
  @Output("onCheck") onCheckEvent: EventEmitter<any> = new EventEmitter<any>();
  @Output("onRevert") onRevertEvent: EventEmitter<any> = new EventEmitter<any>();

  onEditClick() {
    if (!this.isDisabled && this.inputType != CsInputType.itWithNoInput) {
      this.inputField.status = CsInputStatus.isEdit;
    } else if (!this.isDisabled && this.inputType == CsInputType.itWithNoInput) {
      this.onEditEvent.emit();
    }
  }

  onCheckClick() {
    if (this.isPassValidator) {
      this.inputField.status = CsInputStatus.isView;
      this.inputField.defaultValue = this.inputField.value;
      this.inputHtml.nativeElement.blur();
      this.onCheckEvent.emit(this.inputField.value);
    }
  }

  onRevertClick() {
    this.inputField.value = this.inputField.defaultValue;
    this.inputField.status = CsInputStatus.isView;
    this.onRevertEvent.emit();
  }
}
