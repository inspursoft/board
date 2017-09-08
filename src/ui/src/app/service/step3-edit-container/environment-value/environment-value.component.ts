/**
 * Created by liyanq on 9/4/17.
 */

import { Component, EventEmitter, Input, Output } from "@angular/core"

@Component({
  selector:"environment-value",
  templateUrl:"./environment-value.component.html",
  styleUrls:["./environment-value.component.css"]
})
export class EnvironmentValueComponent{
  _isOpen: boolean = false;
  isAlertOpen: boolean = false;
  afterCommitErr: string = "";

  @Input()
  get isOpen() {
    return this._isOpen;
  }

  set isOpen(open: boolean) {
    this._isOpen = open;
    this.isOpenChange.emit(this._isOpen);
  }

  @Output() isOpenChange: EventEmitter<boolean> = new EventEmitter<boolean>();
}