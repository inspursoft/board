/**
 * Created by liyanq on 9/11/17.
 */
import { Component, Input, Output, EventEmitter, OnInit } from "@angular/core"
import { FormControl, FormGroup, ValidationErrors, ValidatorFn } from "@angular/forms";

export enum CsInputStatus{isView, isEdit}
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
const InputPatternString: RegExp = null;
@Component({
  selector: "cs-input",
  templateUrl: "./cs-input.component.html",
  styleUrls: ["./cs-input.component.css"]
})
export class CsInputComponent implements OnInit {
  _isDisabled: boolean = false;
  inputFormGroup: FormGroup;
  @Input() inputLabel: string = "";
  @Input() inputFiledType: CsInputFiledType = CsInputFiledType.iftString;
  @Input() inputField: CsInputFiled;
  @Input() inputType: CsInputType = CsInputType.itWithInput;
  @Input() inputMaxlength: string;
  @Input() validatorFns: Array<ValidatorFn>;
  @Input() validatorMessage: Array<{validatorKey: string, validatorMessage: string}>;

  ngOnInit() {
    this.inputFormGroup = new FormGroup({
      inputControl: new FormControl(this.SimpleFiled, this.validatorFns)
    });
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
  }

  get inputPattern(): RegExp {
    return this.inputFiledType == CsInputFiledType.iftString ?
      InputPatternString :
      InputPatternNumber;
  }

  get isDisabled() {
    return this._isDisabled;
  }

  get typeName(): string {
    return typeof this.inputField.value;
  }

  get valid(): boolean {
    return this.inputFormGroup.valid;
  }

  getValidatorMessage(errors: ValidationErrors): string {
    let result: string = "";
    if (errors["pattern"]) {
      result = "只能输入数字"
    }
    else if (this.validatorMessage) {
      this.validatorMessage.forEach(value => {
        if (errors[value.validatorKey]) {
          result = value.validatorMessage;
        }
      })
    }
    return result;
  }

  public checkValueByHost() {
    this.onCheckClick();
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
    if (this.inputFormGroup.valid) {
      this.inputField.status = CsInputStatus.isView;
      this.inputField.defaultValue = this.inputField.value;
      this.onCheckEvent.emit(this.inputField.value);
    }
  }

  onRevertClick() {
    this.inputField.value = this.inputField.defaultValue;
    this.inputField.status = CsInputStatus.isView;
    this.onRevertEvent.emit();
  }

  onInputKeyPress(event: KeyboardEvent) {
    if (event.keyCode == 13) {
      this.onCheckClick();
    }
  }
}
