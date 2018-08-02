/**
 * Created by liyanq on 9/11/17.
 */
import { Component, ElementRef, EventEmitter, Input, OnInit, Output, ViewChild } from "@angular/core"
import { AsyncValidatorFn, FormControl, FormGroup, ValidationErrors, ValidatorFn, Validators } from "@angular/forms";

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
  isInValidatorWIP:boolean = false;
  inputFormGroup: FormGroup;
  inputValidatorFns: Array<ValidatorFn>;
  inputValidatorMessageParam: string;
  inputField: CsInputFiled;
  isCheckInputOnKeyPress: boolean = false;
  isAlreadyChecked: boolean = false;
  public inputControl: FormControl;
  @ViewChild("input") inputHtml: ElementRef;
  @ViewChild("container") containerHtml: ElementRef;
  @Input() inputLabel: string = "";
  @Input() inputLabelMinWidth: string = "180px";
  @Input() inputFiledType: CsInputFiledType = CsInputFiledType.iftString;
  @Input() inputIsRequired: boolean = false;
  @Input() inputPattern: RegExp;
  @Input() inputMaxlength: number = 0;
  @Input() inputMinlength: number = 0;
  @Input() inputMax: number = 0;
  @Input() inputMin: number = 0;
  @Input() inputType: CsInputType = CsInputType.itWithInput;
  @Input() customerValidatorFns: Array<ValidatorFn>;
  @Input() customerValidatorAsyncFunc: AsyncValidatorFn;
  @Input() validatorMessage: Array<{validatorKey: string, validatorMessage: string}>;
  @Input() inputPlaceholder: string = "";

  constructor() {
    this.inputValidatorFns = Array<ValidatorFn>();
    this.inputControl = new FormControl({value: "", disabled: false});
  }

  ngOnInit() {
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
    this.inputControl.statusChanges.subscribe(() => {
      if (this.inputControl.valid && this.isInValidatorWIP) {
        this.inputControl.clearValidators();
        this.inputControl.clearAsyncValidators();
        this.isInValidatorWIP = false;
        this.inputField.status = CsInputStatus.isView;
        this.inputField.defaultValue = this.inputField.value;
        this.inputHtml.nativeElement.blur();
        this.onCheckEvent.emit(this.inputField.value);
        if (this.isCheckInputOnKeyPress) {
          let nextInputElement: Element = this.containerHtml.nativeElement.parentElement.nextElementSibling;
          if (nextInputElement) {
            let nextLabelElement: NodeListOf<HTMLLabelElement> = nextInputElement.getElementsByClassName("cs-input-label") as NodeListOf<HTMLLabelElement>;
            if (nextLabelElement && nextLabelElement.length > 0) {
              nextLabelElement[0].click();
            }
          }
        }
      } else if (this.inputControl.invalid && this.isInValidatorWIP) {
        this.inputControl.clearValidators();
        this.inputControl.clearAsyncValidators();
        this.isInValidatorWIP = false;
        this.inputField.status = CsInputStatus.isEdit;
        this.inputHtml.nativeElement.focus();
      }
    });
  }

  @Input("simpleFiled")
  set SimpleFiled(value: CsInputSupportType) {
    this.inputField = new CsInputFiled(
      CsInputStatus.isView, value, value
    );
    this.inputControl.setValue(this.inputField.value);
  }

  get SimpleFiled(): CsInputSupportType {
    return this.inputField.value;
  }

  @Input("disabled")
  set isDisabled(value: boolean) {
    if (value) {
      this.inputField.status = CsInputStatus.isView;
    }
    this.inputControl.reset({value: this.SimpleFiled, disabled: value});
  }

  public get valid(): boolean {
    return this.inputControl.valid && this.isAlreadyChecked && this.inputField.status == CsInputStatus.isView;
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
        result = "ERROR.INPUT_PATTERN";
        this.inputValidatorMessageParam = `:${errors["pattern"]["requiredPattern"]}`
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
      } else if (Object.keys(errors).length > 0){
        result = errors[Object.keys(errors)[0]];
      }
    }
    return result;
  }

  onInputKeyPressEvent(event: KeyboardEvent) {
    if (event.keyCode == 13) {
      this.isCheckInputOnKeyPress = true;
      this.checkInputSelf();
    }
  }

  onInputBlur() {
    if (this.inputControl.enabled && this.inputField.status == CsInputStatus.isEdit && this.inputType != CsInputType.itWithNoInput) {
      this.isCheckInputOnKeyPress = false;
      this.checkInputSelf();
    }
  }

  @Output("onEdit") onEditEvent: EventEmitter<any> = new EventEmitter<any>();
  @Output("onCheck") onCheckEvent: EventEmitter<any> = new EventEmitter<any>();
  @Output("onRevert") onRevertEvent: EventEmitter<any> = new EventEmitter<any>();

  onEditClick() {
    if (this.inputControl.enabled && this.inputField.status == CsInputStatus.isView && this.inputType != CsInputType.itWithNoInput) {
      this.inputField.status = CsInputStatus.isEdit;
      this.inputHtml.nativeElement.focus();
    } else if (this.inputControl.enabled && this.inputType == CsInputType.itWithNoInput) {
      this.inputHtml.nativeElement.blur();
      this.onEditEvent.emit();
    }
  }

  onCheckClick(): void {
    this.isCheckInputOnKeyPress = false;
    this.checkInputSelf();
  }

  onRevertClick() {
    this.inputField.value = this.inputField.defaultValue;
    this.inputField.status = CsInputStatus.isView;
    this.onRevertEvent.emit();
  }

  public checkInputSelf() {
    if (this.inputControl.enabled) {
      this.isInValidatorWIP = true;
      this.isAlreadyChecked = true;
      this.inputControl.setValidators(this.inputValidatorFns);
      this.inputControl.setAsyncValidators(this.customerValidatorAsyncFunc);
      this.inputControl.markAsTouched({onlySelf: true});
      this.inputControl.updateValueAndValidity({onlySelf: false, emitEvent: true});
    }
  }
}
