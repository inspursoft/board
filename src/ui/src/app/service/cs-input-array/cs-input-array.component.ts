/**
 * Created by liyanq on 9/12/17.
 */
import { Component, Input, Output, EventEmitter, OnInit } from "@angular/core"
import { FormControl, FormGroup, ValidationErrors, ValidatorFn, Validators } from "@angular/forms";

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
@Component({
  selector: "cs-input-array",
  templateUrl: "./cs-input-array.component.html",
  styleUrls: ["./cs-input-array.component.css"]
})
export class CsInputArrayComponent implements OnInit {
  _isDisable: boolean = false;
  inputArrayFileds: Array<CsInputArrFiled>;
  inputArrayErrors: Map<string, ValidationErrors>;
  inputArrayFormGroup: FormGroup;
  inputArrayValidatorFns: Array<ValidatorFn>;
  inputArrayValidatorMessageParam: string = "";
  @Input() inputArrayType: CsInputArrType = CsInputArrType.iasString;
  @Input() customerArrayValidatorFns: Array<ValidatorFn>;
  @Input() inputArrayIsRequired: boolean = false;
  @Input() inputArrayPattern: RegExp;
  @Input() inputArrayMaxlength: number = 0;
  @Input() inputArrayMinlength: number = 0;
  @Input() validatorMessage: Array<{validatorKey: string, validatorMessage: string}>;
  @Input() inputArraySource: Array<CsInputArrSupportType>;
  @Input() labelText: string = "";

  constructor() {
    this.inputArrayFileds = Array<CsInputArrFiled>();
    this.inputArrayValidatorFns = Array<ValidatorFn>();
    this.inputArrayErrors = new Map<string, ValidationErrors>();
  }

  ngOnInit() {
    if (this.inputArraySource.length == 0) {
      this.inputArrayFormGroup = new FormGroup({
        inputControl_0: new FormControl("")
      });
      this.inputArrayFileds.push(new CsInputArrFiled(CsInputArrStatus.iasView, "", ""));
    } else {
      this.inputArrayFormGroup = new FormGroup({});
      this.inputArraySource.forEach((value, index) => {
        let formControlName = `inputControl_` + index;
        this.inputArrayFormGroup.addControl(formControlName, new FormControl(value));
        this.inputArrayFileds.push(new CsInputArrFiled(CsInputArrStatus.iasView, value.toString(), value.toString()));
      })
    }
    if (this.inputArrayType == CsInputArrType.iasNumber) {
      this.inputArrayValidatorFns.push(Validators.pattern(InputArrayPatternNumber));
    }
    if (this.inputArrayIsRequired) {
      this.inputArrayValidatorFns.push(Validators.required);
    }
    if (this.inputArrayMaxlength > 0) {
      this.inputArrayValidatorFns.push(Validators.maxLength(this.inputArrayMaxlength));
    }
    if (this.inputArrayMinlength > 0) {
      this.inputArrayValidatorFns.push(Validators.minLength(this.inputArrayMinlength));
    }
    if (this.inputArrayPattern) {
      this.inputArrayValidatorFns.push(Validators.pattern(this.inputArrayPattern));
    }
    if (this.customerArrayValidatorFns) {
      this.customerArrayValidatorFns.forEach(value => {
        this.inputArrayValidatorFns.push(value);
      })
    }
  }

  @Input()
  set disabled(value) {
    this._isDisable = value;
    if (value) {
      this.inputArrayFileds.forEach(item => item.status = CsInputArrStatus.iasView);
    }
  }

  getControlName(index: number): string {
    if (index < this.inputArrayFileds.length) {
      return 'inputControl_' + index;
    }
  }

  get plusEnable(): boolean {
    return this.inputArrayFileds[this.inputArrayFileds.length - 1].value != "";
  }

  get disabled() {
    return this._isDisable;
  }

  validatorArrayErrors(controlName: string): boolean {
    let result = true;
    this.inputArrayErrors.set(controlName, null);
    this.inputArrayValidatorFns.forEach(value => {
      let errors = value(this.inputArrayFormGroup.controls[controlName]);
      if (errors) {
        result = false;
        this.inputArrayErrors.set(controlName, errors);
      }
    });
    return result;
  }

  getArrayErrors(index: number) {
    let formControlName = this.getControlName(index);
    return this.inputArrayErrors.get(formControlName);
  }

  getValidatorMessage(index: number) {
    this.inputArrayValidatorMessageParam = "";
    let formControlName = this.getControlName(index);
    let result: string = "";
    let errors = this.inputArrayErrors.get(formControlName);
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
      } else if (errors["pattern"] && this.inputArrayType == CsInputArrType.iasNumber) {
        result = "ERROR.INPUT_ONLY_NUMBER"
      } else if (errors["pattern"] && this.inputArrayType == CsInputArrType.iasString) {
        result = "ERROR.INPUT_PATTERN"
      } else if (errors["maxlength"]) {
        result = `ERROR.INPUT_MAX_LENGTH`;
        this.inputArrayValidatorMessageParam = `:${this.inputArrayMaxlength}`
      } else if (errors["minlength"]) {
        result = `ERROR.INPUT_MIN_LENGTH`;
        this.inputArrayValidatorMessageParam = `:${this.inputArrayMinlength}`
      }
    }
    return result;
  }

  arrayIsValid(): boolean {
    let result: boolean = true;
    this.inputArrayFileds.forEach(value => {
      if (value.status == CsInputArrStatus.iasEdit) {
        result = false;
      }
    });
    return result;
  }

  @Output("onEdit") onEditEvent: EventEmitter<any> = new EventEmitter<any>();
  @Output("onCheck") onCheckEvent: EventEmitter<any> = new EventEmitter<any>();
  @Output("onRevert") onRevertEvent: EventEmitter<any> = new EventEmitter<any>();
  @Output("onMinus") onMinusEvent: EventEmitter<any> = new EventEmitter<any>();

  onEditClick(index: number) {
    if (!this.disabled) {
      this.inputArrayFileds[index].status = CsInputArrStatus.iasEdit;
      this.onEditEvent.emit();
    }
  }

  onCheckClick(index: number) {
    let controlName = this.getControlName(index);
    if ((this.inputArrayFormGroup.touched || this.inputArrayFormGroup.dirty) && this.validatorArrayErrors(controlName)) {
      this.inputArrayFileds[index].status = CsInputArrStatus.iasView;
      this.inputArrayFileds[index].defaultValue = this.inputArrayFileds[index].value;
      if (this.inputArraySource.length == 0) {
        this.inputArrayType == CsInputArrType.iasString ?
          this.inputArraySource.push(this.inputArrayFileds[index].value) :
          this.inputArraySource.push(Number(this.inputArrayFileds[index].value).valueOf());
      } else {
        this.inputArrayType == CsInputArrType.iasString ?
          this.inputArraySource[index] = this.inputArrayFileds[index].value :
          this.inputArraySource[index] = Number(this.inputArrayFileds[index].value).valueOf();
      }
      this.onCheckEvent.emit();
    }
  }

  onRevertClick(index: number) {
    this.inputArrayFileds[index].status = CsInputArrStatus.iasView;
    this.inputArrayFileds[index].value = this.inputArrayFileds[index].defaultValue;
    this.onRevertEvent.emit();
  }

  onPlusClick() {
    if (this.plusEnable) {
      this.inputArrayFileds.push(new CsInputArrFiled(CsInputArrStatus.iasView, "", ""));
      let inputControl = new FormControl("");
      let inputControlName = this.getControlName(this.inputArrayFileds.length - 1);
      this.inputArrayFormGroup.addControl(inputControlName, inputControl);
      this.inputArrayType == CsInputArrType.iasString ?
        this.inputArraySource.push("") :
        this.inputArraySource.push(0);
    }
  }

  onMinusClick(index: number) {
    this.inputArraySource.splice(index, 1);
    this.inputArrayFileds.splice(index, 1);
    this.onMinusEvent.emit();
  }

  onInputKeyPressEvent(event: KeyboardEvent, index: number) {
    if (event.keyCode == 13) {
      this.onCheckClick(index);
    }
  }
}
