/**
 * Dropdown Component
 * v1.0
 * Created by liyanq on 9/4/17.
 */

import { Component, Input, Output, EventEmitter, HostListener } from "@angular/core"

@Component({
  selector: "cs-dropdown",
  templateUrl: "./cs-dropdown.component.html",
  styleUrls: ["./cs-dropdown.component.css"]
})
export class CsDropdownComponent {
  dropdownList: Array<Object>;
  @Input("width") dropdownWidth: number = 100;
  @Input("defaultText") dropdownText: string = "";
  @Input("disabled") dropdownDisabled: boolean = false;

  @Input("list")
  get list() {
    return this.dropdownList;
  }

  set list(list: Array<Object>) {
    this.dropdownList = list;
  }

  @Input("textKey") dropdownListTextKey: string;
  @Output("onChange") dropdownChange: EventEmitter<Object>;

  constructor() {
    this.dropdownList = Array<Object>();
    this.dropdownChange = new EventEmitter<Object>();
  }

  changeSelect(item: Object) {
    this.dropdownText = item[this.dropdownListTextKey];
    this.dropdownChange.emit(item);
  }
}