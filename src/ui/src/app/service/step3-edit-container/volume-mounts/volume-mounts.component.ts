/**
 * Created by liyanq on 9/1/17.
 */
import { Component, Input, Output, EventEmitter, ViewChildren, QueryList } from "@angular/core"
import { CsInputComponent } from "../../../shared/cs-components-library/cs-input/cs-input.component";

export interface VolumeOutPut {
  out_name: string;     //old=>target_storagename || deployment_yaml.volume_list.volume_name
  out_mountPath: string;//old=>container_dir
  out_path: string;     //old=>target_dir || deployment_yaml.volume_list.volume_path
  out_medium: string;   //old=>target_storageServer
}

@Component({
  selector: "volume-mounts",
  templateUrl: "./volume-mounts.component.html",
  styleUrls: ["./volume-mounts.component.css"]
})
export class VolumeMountsComponent {
  _isOpen: boolean = false;
  patternName: RegExp = /^[a-z0-9A-Z_]+$/;
  patternMountPath: RegExp = /^[a-z0-9A-Z_/]+$/;
  patternPath: RegExp = /^[a-z0-9A-Z_/]+$/;
  isAlertOpen: boolean = false;
  volumeErrMsg: string = "";
  @ViewChildren(CsInputComponent) inputList: QueryList<CsInputComponent>;
  @Input() volumeData:VolumeOutPut;

  @Input()
  get isOpen() {
    return this._isOpen;
  }

  set isOpen(open: boolean) {
    this._isOpen = open;
    this.isOpenChange.emit(this._isOpen);
  }

  @Output() isOpenChange: EventEmitter<boolean> = new EventEmitter<boolean>();
  @Output() onConfirmEvent: EventEmitter<VolumeOutPut> = new EventEmitter<VolumeOutPut>();

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