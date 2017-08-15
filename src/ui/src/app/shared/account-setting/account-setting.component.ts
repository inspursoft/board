
import { Component, EventEmitter, Input, Output } from "@angular/core"

@Component({
  selector:"user-setting",
  templateUrl:"./account-setting.component.html",
  styleUrls:["./account-setting.component.css"]
})
export class AccountSettingComponent{
  _isOpen: boolean = false;
  isAlertClose: boolean = true;
  errMessage: string;
  @Output() isOpenChange: EventEmitter<boolean> = new EventEmitter<boolean>();

  @Input()
  get isOpen() {
    return this._isOpen;
  }

  set isOpen(open: boolean) {
    this._isOpen = open;
    this.isOpenChange.emit(this._isOpen);

  }
}