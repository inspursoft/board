/**
 * Created by liyanq on 9/12/17.
 */
import { Component, EventEmitter, Input, OnInit, Output, QueryList, ViewChildren } from "@angular/core"
import { AsyncValidatorFn, ValidationErrors, Validators } from "@angular/forms";
import { CsInputComponent } from "../cs-input/cs-input.component";
import { AbstractControl } from "@angular/forms/src/model";

export enum CsInputArrType{iasString, iasNumber}
export type CsInputArrSupportType = string | number
@Component({
  selector: "cs-input-array",
  templateUrl: "./cs-input-array.component.html",
  styleUrls: ["./cs-input-array.component.css"]
})
export class CsInputArrayComponent implements OnInit {
  @ViewChildren(CsInputComponent) inputList: QueryList<CsInputComponent>;
  @Input() inputArrayFixedSource: Array<CsInputArrSupportType>;
  @Input() inputArraySource: Array<CsInputArrSupportType>;
  @Input() inputArrayType: CsInputArrType = CsInputArrType.iasString;
  @Input() inputArrayDisabled: boolean = false;
  @Input() inputArrayPattern: RegExp;
  @Input() inputArrayMaxlength: number = 0;
  @Input() inputArrayMinlength: number = 0;
  @Input() inputArrayMax: number = 0;
  @Input() inputArrayMin: number = 0;
  @Input() inputArrayLabelText: string = "";
  @Input() inputArrayIsRequired: boolean = false;
  @Input() inputArrayLabelMinWidth: string = "180";
  @Input() validatorMessage: Array<{validatorKey: string, validatorMessage: string}>;
  @Input() customerValidatorAsyncFunc: AsyncValidatorFn;
  @Output("onMinus") onMinusEvent: EventEmitter<any> = new EventEmitter<any>();
  @Output("onCheck") onCheckEvent: EventEmitter<any> = new EventEmitter<any>();
  @Output("onRevert") onRevertEvent: EventEmitter<any> = new EventEmitter<any>();

  constructor() {
    this.inputArrayFixedSource = Array<CsInputArrSupportType>();
    this.inputArraySource = Array<CsInputArrSupportType>();
    this.validatorMessage = Array<{validatorKey: string, validatorMessage: string}>();
  }

  ngOnInit() {
    this.validatorMessage.push({validatorKey: "notRepeat", validatorMessage: "ERROR.INPUT_NOT_REPEAT"})
  }

  public get valid(): boolean {
    return this.inputList.toArray().every(value => value.valid == true);
  }

  public checkInputSelf() {
    this.inputList.toArray().forEach(value => value.checkInputSelf());
  }

  checkRepeatAction(c: AbstractControl): ValidationErrors | null {
    if (this.inputList) {
      let ctr = this.inputList.toArray().find(value => {
        if (this.inputArrayType == CsInputArrType.iasString) {
          return value.inputControl != c && (value.inputControl.value as string).trim() === (c.value as string).trim();
        } else {
          return value.inputControl != c && value.inputControl.value === c.value;
        }
      });
      if (ctr) {
        return {notRepeat: "ERROR.INPUT_NOT_REPEAT"};
      } else {
        return Validators.nullValidator;
      }
    } else {
      return Validators.nullValidator;
    }
  }

  get selfObject() {
    return this;
  }

  onMinusClick(index: number) {
    this.inputArraySource.splice(index, 1);
    this.onMinusEvent.emit();
  }

  onPlusClick() {
    this.inputArraySource.push(this.inputArrayType == CsInputArrType.iasString ? "" : 0);
  }

  onCheckClick(index: number, value: CsInputArrSupportType) {
    this.inputArraySource[index] = value;
    this.onCheckEvent.emit();
  }

  onRevertClick(index: number, value: CsInputArrSupportType) {
    this.inputArraySource[index] = value;
    this.onRevertEvent.emit();
  }
}
