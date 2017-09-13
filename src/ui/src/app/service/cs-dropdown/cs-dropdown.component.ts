/**
 * Dropdown Component
 * v1.0
 * Created by liyanq on 9/4/17.
 */

import { Component, Input, Output, EventEmitter } from "@angular/core"

export const ONLY_FOR_CLICK = "OnlyClick";
@Component({
  selector: "cs-dropdown",
  templateUrl: "./cs-dropdown.component.html",
  styleUrls: ["./cs-dropdown.component.css"]
})
export class CsDropdownComponent {
  isShowDefaultText: boolean = true;
  dropdownText: string;
  @Input("width") dropdownWidth: number = 100;
  @Input("defaultText") dropdownDefaultText;
  @Input("disabled") dropdownDisabled: boolean = false;
  @Input("list") dropdownList: Array<Object>;
  @Input("textKey") dropdownListTextKey: string;
  @Output("onChange") dropdownChange: EventEmitter<Object>;
  @Output("onOnlyClickItem") dropdownClick: EventEmitter<Object>;

  constructor() {
    this.dropdownChange = new EventEmitter<Object>();
    this.dropdownClick = new EventEmitter<Object>();
  }

  changeSelect(item: Object) {
    if (item[ONLY_FOR_CLICK]) {
      this.dropdownClick.emit(item);
    } else {
      this.isShowDefaultText = false;
      this.dropdownText = item[this.dropdownListTextKey];
      this.dropdownChange.emit(item);
    }
  }
}