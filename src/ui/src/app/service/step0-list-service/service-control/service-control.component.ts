/**
 * Created by liyanq on 04/12/2017.
 */

import { Component, EventEmitter, Input, Output, OnInit} from "@angular/core"

@Component({
  selector: "service-control",
  styleUrls: ["./service-control.component.css"],
  templateUrl: "./service-control.component.html"
})
export class ServiceControlComponent implements OnInit{
  _isOpen: boolean = false;
  dropDownListNum: Array<number>;

  constructor() {
    this.dropDownListNum = Array<number>();
  }

  ngOnInit() {
    for (let i = 1; i <= 100; i++) {
      this.dropDownListNum.push(i)
    }
  }

  @Output() isOpenChange: EventEmitter<boolean> = new EventEmitter<boolean>();

  @Input() get isOpen() {
    return this._isOpen;
  }

  set isOpen(open: boolean) {
    this._isOpen = open;
    this.isOpenChange.emit(this._isOpen);
  }


}