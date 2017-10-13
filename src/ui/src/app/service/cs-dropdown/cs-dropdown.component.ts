/**
 * Dropdown Component
 * v1.0
 * Created by liyanq on 9/4/17.
 */

import { Component, Input, Output, EventEmitter } from "@angular/core"

export const ONLY_FOR_CLICK = "OnlyClick";
export type EnableSelectCallBack = (item: Object) => boolean;
@Component({
  selector: "cs-dropdown",
  templateUrl: "./cs-dropdown.component.html",
  styleUrls: ["./cs-dropdown.component.css"]
})
export class CsDropdownComponent {
  isShowDefaultText: boolean = true;
  dropdownText: string;

  @Input() dropdownCanSelect: EnableSelectCallBack;
  @Input() dropdownWidth: number = 100;
  @Input() dropdownDefaultText;
  @Input() dropdownDisabled: boolean = false;
  @Input() dropdownList: Array<any>;
  @Input() dropdownListTextKey;
  @Input() dropdownTitleFontSize:number = 14;
  @Output("onChange") dropdownChange: EventEmitter<any>;
  @Output("onOnlyClickItem") dropdownClick: EventEmitter<any>;

  constructor() {
    this.dropdownChange = new EventEmitter<any>();
    this.dropdownClick = new EventEmitter<any>();
  }

  getItemClass(item: any) {
    return {
      'special': (typeof item == "object") && item['isSpecial'],
      'active': this.dropdownText == this.getItemDescription(item) ||
      this.dropdownDefaultText == this.getItemDescription(item)
    }
  }

  getItemDescription(item: any): string {
    if (typeof item == "object") {
      return item[this.dropdownListTextKey];
    }
    return item.toString();
  }

  changeSelect(item: any) {
    if (item[ONLY_FOR_CLICK]) {
      this.dropdownClick.emit(item);
    } else {
      if (this.dropdownCanSelect && !this.dropdownCanSelect(item)) {
        return;
      }
      this.isShowDefaultText = false;
      if (this.dropdownText != this.getItemDescription(item)) {
        this.dropdownText = this.getItemDescription(item);
        this.dropdownChange.emit(item);
      }
    }
  }
}