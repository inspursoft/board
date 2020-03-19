import { QueryList, ViewChildren } from "@angular/core";
import { DropdownExComponent, InputArrayExComponent, InputDropdownNumberComponent, InputExComponent } from "board-components-library";

export class CsComponentBase {
  @ViewChildren(InputExComponent) inputExComponents: QueryList<InputExComponent>;
  @ViewChildren(DropdownExComponent) dropdownExComponents: QueryList<DropdownExComponent>;
  @ViewChildren(InputArrayExComponent) inputArrayExComponents: QueryList<InputArrayExComponent>;
  @ViewChildren(InputDropdownNumberComponent) inputNumberDropdownComponents: QueryList<InputDropdownNumberComponent>;

  verifyInputExValid(): boolean {
    return this.inputExComponents.toArray().every((component: InputExComponent) => {
      if (!component.isValid && !component.inputDisabled) {
        console.log(''+component.isValid + '8' + component.inputDisabled)
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

  verifyInputNumberDropdownValid(): boolean {
    return this.inputNumberDropdownComponents.toArray().every((component: InputDropdownNumberComponent) => {
      if (!component.isValid && !component.disabled) {
        component.checkSelf();
      }
      return component.isValid || component.disabled;
    });
  }
}
