/**
 * Created by liyanq on 9/11/17.
 */
import { Component, Input, Output, EventEmitter, ViewChild, OnInit } from "@angular/core"
import { FormControl, FormGroup, ValidatorFn } from "@angular/forms";

export enum CsInputStatus{isView, isEdit}
export enum CsInputType{itWithInput, itWithNoInput}
export type CsInputTypeFiled = CsInputFiled<string> | CsInputFiled<number>;
export class CsInputFiled<T> {
  constructor(public status: CsInputStatus,
              public defaultValue: T,
              public value: T) {
  }
}

@Component({
  selector: "cs-input",
  templateUrl: "./cs-input.component.html",
  styleUrls: ["./cs-input.component.css"]
})
export class CsInputComponent implements OnInit {
  _isDisabled: boolean = false;
  inputFormGroup: FormGroup;
  @ViewChild("input") Input;
  @Input("Label") labelText: string = "";
  @Input("Field") curField: CsInputTypeFiled;
  @Input("Type") curFieldType: CsInputType = CsInputType.itWithInput;
  @Input() validatorFns: Array<ValidatorFn>;
  @Input() inputMaxlength: string;

  ngOnInit() {
    this.inputFormGroup = new FormGroup({
      inputControl: new FormControl("", this.validatorFns)
    })
  }

  @Input("disabled")
  set isDisabled(value: boolean) {
    this._isDisabled = value;
    if (value) {
      this.curField.status = CsInputStatus.isView;
    }
  }

  get isDisabled() {
    return this._isDisabled;
  }

  @Input("SimpleFiled")
  set SimpleFiled(value: string) {
    this.curField = new CsInputFiled(
      CsInputStatus.isView, value, value
    );
  }

  get typeName(): string {
    return typeof this.curField.value;
  }

  @Output("OnEdit") onEditEvent: EventEmitter<any> = new EventEmitter<any>();
  @Output("OnCheck") onCheckEvent: EventEmitter<any> = new EventEmitter<any>();
  @Output("OnRevert") onRevertEvent: EventEmitter<any> = new EventEmitter<any>();

  onEditClick() {
    if (this.curFieldType == CsInputType.itWithInput) {
      this.curField.status = 1;
    }
    this.onEditEvent.emit();
  }

  onCheckClick() {
    if (this.inputFormGroup.valid){
      this.curField.status = 0;
      this.curField.defaultValue = this.curField.value;
      this.onCheckEvent.emit(this.curField.value);
    }
  }

  onRevertClick() {
    this.curField.value = this.curField.defaultValue;
    this.curField.status = 0;
    this.onRevertEvent.emit();
  }
}
