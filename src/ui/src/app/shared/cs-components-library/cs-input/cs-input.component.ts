/**
 * Created by liyanq on 9/11/17.
 */
import { AfterViewInit, Component, ElementRef, EventEmitter, Input, OnInit, Output, ViewChild } from "@angular/core"
import { AsyncValidatorFn, FormControl, FormGroup, ValidationErrors, ValidatorFn, Validators } from "@angular/forms";
import { AbstractControl } from "@angular/forms/src/model";

export enum CsInputFiledType {iftString, iftNumber, iftPassword, iftEmail}
export enum CsInputType{itWithInput, itWithNoInput, itOnlyWithInput}
export enum CsInputStatus {isView = 0, isEdit = 1}
export type CsInputSupportType = string | number

class CsInputFiled {
  constructor(public status: CsInputStatus,
              public defaultValue: CsInputSupportType,
              public value: CsInputSupportType) {
  }
}

class CustomValidator {
  static passwordValidate(c: AbstractControl): ValidationErrors | null {
    let sourceElement: HTMLCollectionOf<Element> = document.getElementsByClassName("source-password");
    let verifyElement: HTMLCollectionOf<Element> = document.getElementsByClassName("verify-password");
    if (sourceElement && sourceElement.length > 0 && verifyElement && verifyElement.length > 0) {
      let source = (sourceElement.item(0) as HTMLInputElement).value;
      let verify = (verifyElement.item(0) as HTMLInputElement).value;
      if (source != "" && verify != "") {
        return source == verify ? null : {verifyPassword: "verify-password"}
      } else {
        return null;
      }
    } else {
      return null;
    }
  }
}

const PATTERN_Number: RegExp = /^[0-9]\d*$/;
@Component({
  selector: "cs-input",
  templateUrl: "./cs-input.component.html",
  styleUrls: ["./cs-input.component.css"]
})
export class CsInputComponent implements OnInit, AfterViewInit {
  inputFormGroup: FormGroup;
  inputValidatorFns: Array<ValidatorFn>;
  inputValidatorMessageParam: string;
  inputField: CsInputFiled;
  isCheckInputOnKeyPress: boolean = false;
  _isDisabled = false;
  private isAfterViewInit = false;
  public inputControl: FormControl;
  @ViewChild("input") inputHtml: ElementRef;
  @ViewChild("container") containerHtml: ElementRef;
  @Input() inputLabel: string = "";
  @Input() inputLabelMinWidth: string = "180";
  @Input() inputFiledType: CsInputFiledType = CsInputFiledType.iftString;
  @Input() inputIsRequired: boolean = false;
  @Input() inputPattern: RegExp;
  @Input() inputMaxlength: number = 0;
  @Input() inputMinlength: number = 0;
  @Input() inputMax: number = 0;
  @Input() inputMin: number = 0;
  @Input() inputType: CsInputType = CsInputType.itWithInput;
  @Input() customerValidatorAsyncFunc: AsyncValidatorFn;
  @Input() customerValidatorFunc: ValidatorFn;
  @Input() validatorMessage: Array<{validatorKey: string, validatorMessage: string}>;
  @Input() inputPlaceholder: string = "";
  @Input() sourcePassword: boolean = false;
  @Input() verifyPassword: boolean = false;


  constructor() {
    this.inputValidatorFns = Array<ValidatorFn>();
    this.inputControl = new FormControl("", {updateOn: 'blur'});
    this.inputFormGroup = new FormGroup({inputControl: this.inputControl}, {updateOn: 'blur'});
  }

  ngAfterViewInit() {
    this.isAfterViewInit = true;
    if (this.customerValidatorAsyncFunc) {
      this.inputControl.setAsyncValidators(this.customerValidatorAsyncFunc);
    }
  }

  ngOnInit() {
    if (this.inputControl.validator) {
      this.inputValidatorFns.push(this.inputControl.validator);
    }
    if (this.customerValidatorFunc) {
      this.inputValidatorFns.push(this.customerValidatorFunc);
    }
    if (this.inputFiledType == CsInputFiledType.iftNumber) {
      this.inputValidatorFns.push(Validators.pattern(PATTERN_Number));
    }
    if (this.inputIsRequired && this.inputType != CsInputType.itWithNoInput) {
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
    if (this.verifyPassword || this.sourcePassword) {
      this.inputValidatorFns.push(CustomValidator.passwordValidate);
    }
    this.inputControl.setValidators(this.inputValidatorFns);
    this.inputControl.statusChanges.subscribe(() => {
      if (this.inputControl.valid) {
        this.inputField.status = CsInputStatus.isView;
        this.inputField.defaultValue = this.inputControl.value;
        this.onCheckEvent.emit(this.inputFiledType == CsInputFiledType.iftNumber ? Number(this.inputControl.value) : this.inputControl.value);
        if (this.isCheckInputOnKeyPress) {
          this.isCheckInputOnKeyPress = false;
          let nextInputElement: Element = this.containerHtml.nativeElement.parentElement.nextElementSibling;
          if (nextInputElement) {
            let nextLabelElement: NodeListOf<HTMLLabelElement> = nextInputElement.getElementsByClassName("cs-input-label") as NodeListOf<HTMLLabelElement>;
            if (nextLabelElement && nextLabelElement.length > 0) {
              nextLabelElement[0].click();
            }
          }
        }
      }
    });
  }

  get inputFieldTypeName(): string {
    switch (this.inputFiledType) {
      case CsInputFiledType.iftNumber:
        return "number";
      case CsInputFiledType.iftPassword:
        return "password";
      case CsInputFiledType.iftEmail:
        return "email";
      default:
        return "text";
    }
  }

  @Input("simpleFiled") set SimpleFiled(value: CsInputSupportType) {
    if (this.inputField) {
      this.inputField.value = value;
      this.inputField.defaultValue = value;
    } else {
      this.inputField = new CsInputFiled(
        CsInputStatus.isView, value, value
      );
    }
    this.inputControl.setValue(value);
  }

  @Input("disabled")
  set isDisabled(value: boolean) {
    this._isDisabled = value;
    if (value) {
      this.inputControl.clearAsyncValidators();
      this.inputControl.disable();
      this.inputField.status = CsInputStatus.isView;
    } else if (this.isAfterViewInit) {
      this.inputControl.enable();
      this.inputControl.setAsyncValidators(this.customerValidatorAsyncFunc);
    }
  }

  get isDisabled(): boolean {
    return this._isDisabled;
  }

  public get valid(): boolean {
    return this.inputIsRequired ? this.inputControl.valid && this.inputField.status == CsInputStatus.isView :
      this.inputControl.valid && this.inputField.status == CsInputStatus.isView;
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
      } else if (errors["passwordPattern"]) {
        result = `ACCOUNT.PASSWORD_FORMAT`;
      } else if (errors["verifyPassword"]) {
        result = `ACCOUNT.PASSWORDS_ARE_NOT_IDENTICAL`;
      } else if (Object.keys(errors).length > 0){
        result = errors[Object.keys(errors)[0]];
      }
    }
    return result;
  }

  onInputKeyPressEvent(event: KeyboardEvent) {
    if (event.keyCode == 13) {
      this.isCheckInputOnKeyPress = true;
      (this.inputHtml.nativeElement as HTMLElement).blur();
    }
  }

  onInputBlur() {
    if (this.inputControl.valid && this.inputField.status == CsInputStatus.isEdit && this.inputType != CsInputType.itWithNoInput) {
      this.inputField.status = CsInputStatus.isView;
    }
  }

  onInputFocus() {
    if (this.inputControl.enabled && this.inputField.status == CsInputStatus.isView && this.inputType != CsInputType.itWithNoInput) {
      this.inputField.status = CsInputStatus.isEdit;
    }
  }

  @Output("onEdit") onEditEvent: EventEmitter<any> = new EventEmitter<any>();
  @Output("onCheck") onCheckEvent: EventEmitter<any> = new EventEmitter<any>();
  @Output("onRevert") onRevertEvent: EventEmitter<any> = new EventEmitter<any>();

  onEditClick() {
    if (this.inputControl.enabled && this.inputField.status == CsInputStatus.isView && this.inputType != CsInputType.itWithNoInput) {
      this.inputHtml.nativeElement.focus();
      if (document.activeElement == this.inputHtml.nativeElement){
        this.inputField.status = CsInputStatus.isEdit;
        this.onEditEvent.emit();
      }
    } else if (this.inputControl.enabled && this.inputType == CsInputType.itWithNoInput) {
      this.inputHtml.nativeElement.blur();
      this.onEditEvent.emit();
    }
  }

  onCheckClick(): void {
    this.checkInputSelf();
  }

  onRevertClick() {
    this.inputField.status = CsInputStatus.isView;
    this.inputControl.reset(this.inputField.defaultValue);
    this.onRevertEvent.emit(this.inputControl.value);
  }

  public checkInputSelf() {
    if (this.inputControl.enabled && (this.inputControl.touched || this.inputIsRequired)) {
      (this.inputHtml.nativeElement as HTMLElement).focus();
      this.inputControl.markAsTouched({onlySelf: true});
      this.inputControl.updateValueAndValidity();
    }
  }
}
