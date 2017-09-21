/**
 * Created by liyanq on 9/12/17.
 */
import { Component, Input, Output, EventEmitter, OnInit } from "@angular/core"
import { AbstractControl, FormControl, FormGroup, ValidationErrors, ValidatorFn } from "@angular/forms";

export enum CsInputArrStatus{iasView, iasEdit}
export enum CsInputArrType{iasString, iasNumber}
export type CsInputArrSupportType = string | number
export class CsInputArrFiled {
  constructor(public status: CsInputArrStatus,
              public defaultValue: string,
              public value: string) {
  }
}

const InputArrayPatternNumber: RegExp = /^[1-9]\d*$/;
const InputArrayPatternString: RegExp = null;
@Component({
  selector: "cs-input-array",
  templateUrl: "./cs-input-array.component.html",
  styleUrls: ["./cs-input-array.component.css"]
})
export class CsInputArrayComponent implements OnInit {
  _isDisable: boolean = false;
  FiledArray: Array<CsInputArrFiled>;
  inputArrayFormGroup: FormGroup;
  @Input() inputArrayType: CsInputArrType = CsInputArrType.iasString;
  @Input() validatorFns: Array<ValidatorFn>;
  @Input() validatorMessage: Array<{validatorKey: string, validatorMessage: string}>;
  @Input() labelText: string = "";
  @Input() inputMaxlength: string;
  @Input() inputArraySource: Array<CsInputArrSupportType>;

  constructor() {
    this.FiledArray = Array();
  }

  ngOnInit() {
    if (this.inputArraySource.length == 0) {
      this.inputArrayFormGroup = new FormGroup({
        inputControl_0: new FormControl("", this.validatorFns)
      });
      this.FiledArray.push(new CsInputArrFiled(CsInputArrStatus.iasView, "", ""));
    } else {
      this.inputArrayFormGroup = new FormGroup({});
      this.inputArraySource.forEach((value, index) => {
        let formControlName = `inputControl_` + index;
        this.inputArrayFormGroup.controls[formControlName] = new FormControl(value, this.validatorFns);
        this.FiledArray.push(new CsInputArrFiled(CsInputArrStatus.iasView, value.toString(), value.toString()));
      })
    }
  }

  @Input()
  set disabled(value) {
    this._isDisable = value;
    if (value) {
      this.FiledArray.forEach(item => item.status = CsInputArrStatus.iasView);
    }
  }

  getControlName(index: number): string {
    if (index < this.FiledArray.length) {
      return 'inputControl_' + index;
    }
  }

  get plusDisabled(): boolean {
    let lastControl = this.getFormControl(this.FiledArray.length - 1);
    return this.FiledArray[this.FiledArray.length - 1].value == "" || lastControl.invalid;
  }

  get disabled() {
    return this._isDisable;
  }

  get inputArrayPattern() {
    return this.inputArrayType == CsInputArrType.iasString ? InputArrayPatternString : InputArrayPatternNumber;
  }

  getFormControl(index: number): AbstractControl {
    let formControlName = this.getControlName(index);
    return this.inputArrayFormGroup.controls[formControlName]
  }

  getValidatorMessage(index: number) {
    let formControlName = this.getControlName(index);
    let errors: ValidationErrors = this.inputArrayFormGroup.controls[formControlName].errors;
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

  @Output("onEdit") onEditEvent: EventEmitter<any> = new EventEmitter<any>();
  @Output("onCheck") onCheckEvent: EventEmitter<any> = new EventEmitter<any>();
  @Output("onRevert") onRevertEvent: EventEmitter<any> = new EventEmitter<any>();
  @Output("onMinus") onMinusEvent: EventEmitter<any> = new EventEmitter<any>();

  onEditClick(index: number) {
    if (!this.disabled) {
      this.FiledArray[index].status = CsInputArrStatus.iasEdit;
      this.onEditEvent.emit();
    }
  }

  onCheckClick(index: number) {
    if (this.getFormControl(index).valid) {
      this.FiledArray[index].status = CsInputArrStatus.iasView;
      this.FiledArray[index].defaultValue = this.FiledArray[index].value;
      if (this.inputArraySource.length == 0) {
        this.inputArrayType == CsInputArrType.iasString ?
          this.inputArraySource.push(this.FiledArray[index].value) :
          this.inputArraySource.push(Number(this.FiledArray[index].value).valueOf());
      } else {
        this.inputArrayType == CsInputArrType.iasString ?
          this.inputArraySource[index] = this.FiledArray[index].value :
          this.inputArraySource[index] = Number(this.FiledArray[index].value).valueOf();
      }
      this.onCheckEvent.emit();
    }
  }

  onRevertClick(index: number) {
    this.FiledArray[index].status = CsInputArrStatus.iasView;
    this.FiledArray[index].value = this.FiledArray[index].defaultValue;
    this.onRevertEvent.emit();
  }

  onPlusClick() {
    if (!this.plusDisabled) {
      this.FiledArray.push(new CsInputArrFiled(CsInputArrStatus.iasView, "", ""));
      let inputControl = new FormControl("", this.validatorFns);
      let inputControlName = `inputControl_${this.FiledArray.length - 1}`;
      this.inputArrayFormGroup.controls[inputControlName] = inputControl;
      this.inputArrayType == CsInputArrType.iasString ?
        this.inputArraySource.push("") :
        this.inputArraySource.push(0);
    }
  }

  onMinusClick(index: number) {
    this.inputArraySource.splice(index, 1);
    this.FiledArray.splice(index, 1);
    this.onMinusEvent.emit();
  }

  onInputKeyPress(event: KeyboardEvent, index: number) {
    if (event.keyCode == 13) {
      this.onCheckClick(index);
    }
  }
}
