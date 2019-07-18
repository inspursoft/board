import { QueryList, ViewChildren } from "@angular/core";
import { CsInputComponent } from "./cs-input/cs-input.component";
import { CsInputArrayComponent } from "./cs-input-array/cs-input-array.component";
import { CsDropdownComponent } from "./cs-dropdown/cs-dropdown.component";
import { CsInputDropdownComponent } from "./cs-input-dropdown/cs-input-dropdown.component";
import { DropdownExComponent, InputArrayExComponent, InputExComponent } from "board-components-library";

export class CsComponentBase {
  @ViewChildren(InputExComponent) inputExComponents: QueryList<InputExComponent>;
  @ViewChildren(DropdownExComponent) dropdownExComponents: QueryList<DropdownExComponent>;
  @ViewChildren(InputArrayExComponent) inputArrayExComponents: QueryList<InputArrayExComponent>;

  verifyInputExValid(): boolean {
    return this.inputExComponents.toArray().every((component: InputExComponent) => {
      if (!component.isValid && !component.inputDisabled) {
        component.checkSelf();
      }
      return component.isValid || component.inputDisabled;
    });
  }

  verifyInputArrayExValid(): boolean {
    return this.inputArrayExComponents.toArray().every((component: InputArrayExComponent) => {
      if (!component.isValid && !component.inputDisabled) {
        component.checkSelf();
      }
      return component.isValid || component.inputDisabled;
    });
  }

  verifyDropdownExValid(): boolean {
    return this.dropdownExComponents.toArray().every((component: DropdownExComponent) => {
      if (!component.isValid && !component.dropdownDisabled) {
        component.checkSelf();
      }
      return component.isValid || component.dropdownDisabled;
    });
  }


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
