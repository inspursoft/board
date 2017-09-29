/**
 * Created by liyanq on 9/1/17.
 */


import { Component, Input, Output, EventEmitter, ViewChildren, QueryList } from "@angular/core"
import { CsInputComponent } from "../../cs-input/cs-input.component";

@Component({
  selector: "volume-mounts",
  templateUrl: "./volume-mounts.component.html",
  styleUrls: ["./volume-mounts.component.css"]
})
export class VolumeMountsComponent {
  _isOpen: boolean = false;
  patternVolumeName: RegExp = /^[a-zA-Z_]+$/;
  patternContainerDir: RegExp = /^[a-zA-Z_/.]+$/;
  patternTargetDir: RegExp = /^[a-zA-Z_/.]+$/;
  isAlertOpen: boolean = false;
  volumeErrMsg: string = "";
  @ViewChildren(CsInputComponent) inputList: QueryList<CsInputComponent>;
  @Input() volumeData: {container_dir: string, target_storagename: string, target_storageServer, target_dir: string};

  @Input()
  get isOpen() {
    return this._isOpen;
  }

  set isOpen(open: boolean) {
    this._isOpen = open;
    this.isOpenChange.emit(this._isOpen);
  }

  @Output() isOpenChange: EventEmitter<boolean> = new EventEmitter<boolean>();
  @Output() onConfirmEvent: EventEmitter<Object> = new EventEmitter<Object>();

  get isConfirmEnabled(): boolean {
    let result = true;
    if (this.inputList) {
      this.inputList.forEach(value => {
        if (!value.valid) {
          result = false;
        }
      });
    }
    return result;
  }

  confirmVolumeInfo() {
    this.inputList.forEach(value => value.checkValueByHost());
    this.onConfirmEvent.emit(this.volumeData);
    this.isOpen = false;
  }
}