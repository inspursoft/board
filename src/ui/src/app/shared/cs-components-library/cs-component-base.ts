import {QueryList, ViewChildren} from "@angular/core";
import {CsInputComponent} from "./cs-input/cs-input.component";
import {CsInputArrayComponent} from "./cs-input-array/cs-input-array.component";
import {CsDropdownComponent} from "./cs-dropdown/cs-dropdown.component";
import {CsInputDropdownComponent} from "./cs-input-dropdown/cs-input-dropdown.component";

export class CsComponentBase {
  @ViewChildren(CsInputComponent) inputComponents: QueryList<CsInputComponent>;
  @ViewChildren(CsInputArrayComponent) inputArrayComponents: QueryList<CsInputArrayComponent>;
  @ViewChildren(CsDropdownComponent) dropdownComponents: QueryList<CsDropdownComponent>;
  @ViewChildren(CsInputDropdownComponent) inputDropdownComponents: QueryList<CsInputDropdownComponent>;

  verifyInputValid(): boolean {
    return this.inputComponents.toArray().every((component: CsInputComponent) => {
      if (!component.valid) {
        component.checkInputSelf();
      }
      return component.valid || component.inputControl.disabled;
    });
  }

  verifyInputDropdownValid(): boolean {
    return this.inputDropdownComponents.toArray().every((component: CsInputDropdownComponent) => {
      if (!component.valid) {
        component.checkInputSelf();
      }
      return component.valid || component.inputControl.disabled;
    });
  }

  verifyDropdownValid(): boolean {
    return this.dropdownComponents.toArray().every((component: CsDropdownComponent) => {
      if (!component.valid) {
        component.checkInputSelf();
      }
      return component.valid;
    });
  }

  verifyInputArrayValid(): boolean {
    return this.inputArrayComponents.toArray().every((component: CsInputArrayComponent) => {
      if (!component.valid) {
        component.checkInputSelf();
      }
      return component.valid || component.inputArrayDisabled;
    });
  }
}