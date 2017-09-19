/**
 * Created by liyanq on 9/1/17.
 */


import { Component, Input, Output, EventEmitter } from "@angular/core"

@Component({
  selector: "volume-mounts",
  templateUrl: "./volume-mounts.component.html",
  styleUrls: ["./volume-mounts.component.css"]
})
export class VolumeMountsComponent {
  _isOpen: boolean = false;
  isAlertOpen: boolean = false;
  volumeErrMsg: string = "";

  @Input()
  get isOpen() {
    return this._isOpen;
  }

  set isOpen(open: boolean) {
    this._isOpen = open;
    this.isOpenChange.emit(this._isOpen);
  }

  @Input() volumeData: {container_dir: string, target_storagename: string, target_dir: string};
  @Output() isOpenChange: EventEmitter<boolean> = new EventEmitter<boolean>();

  confirmVolumeInfo() {
    this.isOpen = false;
  }
}