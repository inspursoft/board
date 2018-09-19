/**
 * Created by liyanq on 9/1/17.
 */
import { Component, Input, Output, EventEmitter} from "@angular/core"
import { CsComponentBase } from "../../../shared/cs-components-library/cs-component-base";

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
export class VolumeMountsComponent extends CsComponentBase{
  _isOpen: boolean = false;
  volumeDataOrigin:VolumeOutPut;
  patternName: RegExp = /^[a-z0-9A-Z_]+$/;
  patternMountPath: RegExp = /^[a-z0-9A-Z_/]+$/;
  patternPath: RegExp = /^[a-z0-9A-Z_/.:]+$/;
  @Input() volumeData:VolumeOutPut;

  @Input()
  get isOpen() {
    return this._isOpen;
  }

  set isOpen(open: boolean) {
    this._isOpen = open;
    this.volumeDataOrigin = Object.create(this.volumeData);
    this.isOpenChange.emit(this._isOpen);
  }

  @Output() isOpenChange: EventEmitter<boolean> = new EventEmitter<boolean>();
  @Output() onConfirmEvent: EventEmitter<VolumeOutPut> = new EventEmitter<VolumeOutPut>();

  confirmVolumeInfo() {
    if (this.verifyInputValid()){
      this.onConfirmEvent.emit(this.volumeDataOrigin);
      this.isOpen = false;
    }
  }
}