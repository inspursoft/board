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
  FiledArray: Array<CsInputArrFiled>;

  constructor() {
    this.FiledArray = Array();
  }

  ngOnInit() {
    this.type == CsInputArrType.iasString ?
      this.FiledArray.push(new CsInputArrFiled(CsInputArrStatus.iasView, "", "")) :
      this.FiledArray.push(new CsInputArrFiled(CsInputArrStatus.iasView, 0, 0))
  }

  @Input("source")
  set sourceArr(value: Array<string>) {
    this._sourceArr = value;
    value.forEach(value => {
      this.FiledArray.push(new CsInputArrFiled(CsInputArrStatus.iasView, value, value));
    })
  }

  @Input() type: CsInputArrType = CsInputArrType.iasString;
  @Input() labelText: string = "";
  @Input() inputMaxlength: string;
  @Input() disabled: boolean;

  get inputType(): string {
    return this.type == CsInputArrType.iasString ? "text" : "number";
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
      this._sourceArr.push(this.FiledArray[index].value);
      this.type == CsInputArrType.iasString ?
        this.FiledArray.push(new CsInputArrFiled(CsInputArrStatus.iasView, "", "")) :
        this.FiledArray.push(new CsInputArrFiled(CsInputArrStatus.iasView, 0, 0))
    } else {
      this._sourceArr[index] = this.FiledArray[index].value;
    }
    this.onCheckEvent.emit();
  }

  onRevertClick(index: number) {
    this.FiledArray[index].status = CsInputArrStatus.iasView;
    this.FiledArray[index].value = this.FiledArray[index].defaultValue;
    this.onRevertEvent.emit();
  }

  onMinusClick(index: number) {
    this._sourceArr.splice(index, 1);
    this.FiledArray.splice(index, 1);
    this.onMinusEvent.emit();
  }
}
