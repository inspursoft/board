import { Component, EventEmitter, Input, Output } from "@angular/core";
import { CsModalChildBase } from "../../../../shared/cs-modal-base/cs-modal-child-base";
import { MessageService } from "../../../../shared.service/message.service";

@Component({
  selector: 'pv-monitors',
  templateUrl: './monitors.component.html',
  styleUrls: ['./monitors.component.css']
})
export class MonitorsComponent extends CsModalChildBase {
  private _isOpen = false;
  patternIp: RegExp = /^((?:(?:25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(?:25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d))))$/;
  monitorsArray: Array<string>;
  @Output() isOpenChange: EventEmitter<boolean>;
  @Output() onCommitMonitorText: EventEmitter<string>;

  @Input() set monitorText(text: string) {
    if (text != '') {
      this.monitorsArray = text.split(`;`);
    }
  }

  get monitorText(): string {
    return this.monitorsArray.join(`;`)
  }

  @Input()
  get isOpen() {
    return this._isOpen
  }

  set isOpen(open: boolean) {
    this._isOpen = open;
    this.isOpenChange.emit(this._isOpen);
  }

  constructor(private messageService: MessageService) {
    super();
    this.isOpenChange = new EventEmitter<boolean>();
    this.onCommitMonitorText = new EventEmitter<string>();
    this.monitorsArray = Array<string>();
  }

  changeIp(ip: string, index: number) {
    let monitor = this.monitorsArray[index];
    if (this.monitorsArray.find((value, oldIndex) => value.startsWith(ip) && index != oldIndex)) {
      this.messageService.showAlert(`STORAGE.PV_CONFIG_MONITORS_IP`, {alertType: "warning", view: this.alertView});
      this.monitorsArray[index] = `${monitor} `
    } else {
      let port = this.getPort(monitor);
      this.monitorsArray[index] = `${ip}:${port}`;
    }
  }

  changePort(port: string, index: number) {
    let monitor = this.monitorsArray[index];
    let ip = this.getIp(monitor);
    this.monitorsArray[index] = `${ip}:${port}`;
  }

  getIp(monitor: string): string {
    return monitor.substr(0, monitor.indexOf(':')).trim();
  }

  getPort(monitor: string): string {
    return monitor.substr(monitor.indexOf(':') + 1).trim();
  }

  addNewMonitor() {
    if (this.monitorsArray.find(value => value.startsWith('127.0.0.1'))) {
      this.messageService.showAlert(`STORAGE.PV_CONFIG_MONITORS_IP`, {alertType: "warning", view: this.alertView});
    } else {
      this.monitorsArray.push('127.0.0.1:6789')
    }
  }

  deleteMonitor(index: number) {
    this.monitorsArray.splice(index, 1);
  }

  confirmMonitors() {
    if (this.verifyInputExValid()) {
      this.onCommitMonitorText.emit(this.monitorText);
      this.isOpen = false;
    }
  }
}
