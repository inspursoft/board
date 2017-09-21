/**
 * Created by liyanq on 9/1/17.
 */


import { Component, Input, Output, EventEmitter, ViewChildren, QueryList, OnInit } from "@angular/core"
import { CsInputComponent } from "../../cs-input/cs-input.component";

@Component({
  selector: "volume-mounts",
  templateUrl: "./volume-mounts.component.html",
  styleUrls: ["./volume-mounts.component.css"]
})
export class VolumeMountsComponent implements OnInit {
  _isOpen: boolean = false;
  isAlertOpen: boolean = false;
  volumeErrMsg: string = "";
  volumeData: {
    container_dir: string,
    target_storagename: string,
    target_storageServer,
    target_dir: string
  };
  @ViewChildren(CsInputComponent) inputList: QueryList<CsInputComponent>;
  @Output() isOpenChange: EventEmitter<boolean> = new EventEmitter<boolean>();
  @Output() onConfirmEvent: EventEmitter<Object> = new EventEmitter<Object>();
  ngOnInit() {
    this.volumeData = {
      container_dir: "",
      target_storagename: "",
      target_storageServer: "",
      target_dir: ""
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

  confirmVolumeInfo() {
    this.inputList.forEach(value => value.checkValueByHost());
    this.onConfirmEvent.emit(this.volumeData);
    this.isOpen = false;
  }
}