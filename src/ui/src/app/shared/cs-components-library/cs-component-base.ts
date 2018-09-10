import { QueryList, ViewChildren } from "@angular/core";
import { CsInputComponent } from "./cs-input/cs-input.component";
import { CsInputArrayComponent } from "./cs-input-array/cs-input-array.component";

export class CsComponentBase {
  @ViewChildren(CsInputComponent) inputComponents: QueryList<CsInputComponent>;
  @ViewChildren(CsInputArrayComponent) inputArrayComponents: QueryList<CsInputArrayComponent>;

  verifyInputValid(): boolean {
    return this.inputComponents.toArray().every((component: CsInputComponent) => {
      if (!component.valid) {
        component.checkInputSelf();
      }
      return component.valid || component.inputControl.disabled;
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