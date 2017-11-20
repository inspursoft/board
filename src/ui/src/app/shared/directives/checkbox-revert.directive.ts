/**
 * Created by liyanq on 16/11/2017.
 */

import { Directive, HostListener, Input } from "@angular/core"

@Directive({
  selector: "[checkbox-revert]"
})
export class CheckboxRevert {
  private _elem: HTMLInputElement;

  @Input("checkbox-revert")
  set restInfo(info: {isNeeded: boolean, value: boolean}) {
    if (info && this._elem && info.isNeeded) {
      this._elem.checked = info.value;
      info.isNeeded = false;
    }
  }

  @HostListener("change", ["$event"])
  listenerChange(event: Event) {
    this._elem = event.target as HTMLInputElement;
  }
}