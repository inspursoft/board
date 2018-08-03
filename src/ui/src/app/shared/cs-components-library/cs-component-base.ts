import { QueryList, ViewChildren } from "@angular/core";
import { CsInputComponent } from "./cs-input/cs-input.component";

export class CsComponentBase {
  @ViewChildren(CsInputComponent) inputComponents: QueryList<CsInputComponent>;

  verifyInputValid(): boolean {
    return this.inputComponents.toArray().every((component: CsInputComponent) => {
      if (!component.valid) {
        component.checkInputSelf();
      }
      return component.valid || component.inputControl.disabled;
    });
  }
}