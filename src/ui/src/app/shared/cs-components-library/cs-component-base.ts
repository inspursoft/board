import { QueryList, ViewChildren } from "@angular/core";
import { CsInputComponent } from "./cs-input/cs-input.component";
import { CsInputArrayComponent } from "./cs-input-array/cs-input-array.component";
import { CsDropdownComponent } from "./cs-dropdown/cs-dropdown.component";
import { CsInputDropdownComponent } from "./cs-input-dropdown/cs-input-dropdown.component";
import { InputExComponent } from "board-components-library";

export class CsComponentBase {
  @ViewChildren(CsInputComponent) inputComponents: QueryList<CsInputComponent>;
  @ViewChildren(InputExComponent) inputExComponents: QueryList<InputExComponent>;
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

  verifyInputExValid(): boolean {
    return this.inputExComponents.toArray().every((component: InputExComponent) => {
      if (!component.isValid) {
        component.checkSelf();
      }
      return component.isValid;
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
