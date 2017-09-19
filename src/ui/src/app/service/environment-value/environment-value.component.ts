/**
 * Created by liyanq on 9/4/17.
 */

import { Component, EventEmitter, Input, Output, OnInit } from "@angular/core"

export class EnvType {
  constructor(public envName: string,
              public envValue: string) {
  }
}
@Component({
  selector: "environment-value",
  templateUrl: "./environment-value.component.html",
  styleUrls: ["./environment-value.component.css"]
})
export class EnvironmentValueComponent implements OnInit {
  _isOpen: boolean = false;
  envAlertMessage: string;
  envsData: Array<EnvType>;
  isAlertOpen: boolean = false;
  afterCommitErr: string = "";
  @Input() inputEnvsData: Array<EnvType>;

  constructor() {
    this.envsData = Array<EnvType>();
  }

  ngOnInit() {
    if (this.inputEnvsData && this.inputEnvsData.length > 0) {
      this.envsData = this.envsData.concat(this.inputEnvsData);
    }
  }

  @Input()
  get isOpen() {
    return this._isOpen;
  }

  set isOpen(open: boolean) {
    this._isOpen = open;
    this.isOpenChange.emit(this._isOpen);
  }


  @Output() isOpenChange: EventEmitter<boolean> = new EventEmitter<boolean>();
  @Output() onConfirm: EventEmitter<Array<EnvType>> = new EventEmitter<Array<EnvType>>();

  addNewEnv() {
    this.envsData.push(new EnvType("", ""));
  }

  confirmEnvInfo() {
    this.onConfirm.emit(this.envsData);
    this.isOpen = false;
  }

  envMinusClick(index: number) {
    this.envsData.splice(index, 1);
  }

}