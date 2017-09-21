/**
 * Created by liyanq on 9/12/17.
 */
import { Component, Input, Output, EventEmitter, OnInit } from "@angular/core"

export enum CsInputArrStatus{iasView, iasEdit}
export enum CsInputArrType{iasString, iasNumber}
export type CsInputArrSupportType = string | number
export class CsInputArrFiled {
  constructor(public status: CsInputArrStatus,
              public defaultValue: CsInputArrSupportType,
              public value: CsInputArrSupportType) {
  }
}

@Component({
  selector: "cs-input-array",
  templateUrl: "./cs-input-array.component.html",
  styleUrls: ["./cs-input-array.component.css"]
})
export class CsInputArrayComponent implements OnInit {
  _sourceArr: Array<CsInputArrSupportType>;
  _isDisable: boolean = false;
  FiledArray: Array<CsInputArrFiled>;
  @Input() type: CsInputArrType = CsInputArrType.iasString;

  constructor() {
    this.FiledArray = Array();
  }

  ngOnInit() {
    this.FiledArray.push(new CsInputArrFiled(CsInputArrStatus.iasView, "", ""));
  }

  @Input("source")
  set sourceArr(value: Array<string>) {
    this._sourceArr = value;
    value.forEach(value => {
      this.FiledArray.push(new CsInputArrFiled(CsInputArrStatus.iasView, value, value));
    })
  }


  @Input() labelText: string = "";
  @Input() inputMaxlength: string;


  @Input()
  set disabled(value) {
    this._isDisable = value;
    if (value) {
      this.FiledArray.forEach(item => item.status = CsInputArrStatus.iasView);
    }
  }

  get disabled() {
    return this._isDisable;
  }

  get inputType(): string {
    return this.type == CsInputArrType.iasString ? "text" : "number";
  }

  changeEditState(item: CsInputArrFiled) {
    if (!this.disabled) {
      item.status = CsInputArrStatus.iasEdit;
    }
  }

  @Output("onEdit") onEditEvent: EventEmitter<any> = new EventEmitter<any>();
  @Output("onCheck") onCheckEvent: EventEmitter<any> = new EventEmitter<any>();
  @Output("onRevert") onRevertEvent: EventEmitter<any> = new EventEmitter<any>();
  @Output("onMinus") onMinusEvent: EventEmitter<any> = new EventEmitter<any>();

  onEditClick(index: number) {
    this.FiledArray[index].status = CsInputArrStatus.iasEdit;
    this.onEditEvent.emit();
  }

  onCheckClick(index: number) {
    this.FiledArray[index].status = CsInputArrStatus.iasView;
    this.FiledArray[index].defaultValue = this.FiledArray[index].value;
    if (index == this.FiledArray.length - 1) {
      this.type == CsInputArrType.iasString ?
        this._sourceArr.push(this.FiledArray[index].value) :
        this._sourceArr.push(Number(this.FiledArray[index].value).valueOf());
    } else {
      this.type == CsInputArrType.iasString ?
        this._sourceArr[index] = this.FiledArray[index].value :
        this._sourceArr[index] = Number(this.FiledArray[index].value).valueOf();
    }
    this.onCheckEvent.emit();
  }

  onRevertClick(index: number) {
    this.FiledArray[index].status = CsInputArrStatus.iasView;
    this.FiledArray[index].value = this.FiledArray[index].defaultValue;
    this.onRevertEvent.emit();
  }

  onPlusClick() {
    this.FiledArray.push(new CsInputArrFiled(CsInputArrStatus.iasView, "", ""));
  }

  onMinusClick(index: number) {
    this._sourceArr.splice(index, 1);
    this.FiledArray.splice(index, 1);
    this.onMinusEvent.emit();
  }

  onInputKeyPress(event:KeyboardEvent,index:number){
    if (event.keyCode == 13){
      this.onCheckClick(index);
    }
  }
}
