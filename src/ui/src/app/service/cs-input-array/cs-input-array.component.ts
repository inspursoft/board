/**
 * Created by liyanq on 9/12/17.
 */
import { Component, Input, Output, EventEmitter, OnInit } from "@angular/core"

export enum CsInputArrStatus{iasView, iasEdit}
export class CsInputArrFiled {
  constructor(public status: CsInputArrStatus,
              public defaultValue: string,
              public value: string) {
  }
}

@Component({
  selector: "cs-input-array",
  templateUrl: "./cs-input-array.component.html",
  styleUrls: ["./cs-input-array.component.css"]
})
export class CsInputArrayComponent implements OnInit {
  _isDisabled: boolean = false;
  _sourceArr: Array<string>;
  FiledArray: Array<CsInputArrFiled>;

  constructor() {
    this.FiledArray = Array<CsInputArrFiled>();
  }

  ngOnInit() {
    this.FiledArray.push(new CsInputArrFiled(CsInputArrStatus.iasView, "", ""))
  }

  @Input("Source")
  set sourceArr(value: Array<string>) {
    this._sourceArr = value;
    value.forEach(value => {
      this.FiledArray.push(new CsInputArrFiled(CsInputArrStatus.iasView, value, value));
    })
  }

  @Input("Label") labelText: string = "";
  @Input() inputMaxlength: string;

  @Input("disabled")
  set isDisabled(value: boolean) {
    this._isDisabled = value;

  }

  get isDisabled() {
    return this._isDisabled;
  }

  @Output("OnEdit") onEditEvent: EventEmitter<any> = new EventEmitter<any>();
  @Output("OnCheck") onCheckEvent: EventEmitter<any> = new EventEmitter<any>();
  @Output("OnRevert") onRevertEvent: EventEmitter<any> = new EventEmitter<any>();
  @Output("OnMinus") onMinusEvent: EventEmitter<any> = new EventEmitter<any>();

  onEditClick(index: number) {
    this.FiledArray[index].status = CsInputArrStatus.iasEdit;
    this.onEditEvent.emit();
  }

  onCheckClick(index: number) {
    this.FiledArray[index].status = CsInputArrStatus.iasView;
    this.FiledArray[index].defaultValue = this.FiledArray[index].value;
    if (index == this.FiledArray.length - 1) {
      this._sourceArr.push(this.FiledArray[index].value);
      this.FiledArray.push(new CsInputArrFiled(CsInputArrStatus.iasView, "", ""));
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
