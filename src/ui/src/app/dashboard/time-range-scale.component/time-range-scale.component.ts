import { Component, Input, Output, EventEmitter } from "@angular/core"

export interface scaleOption {
  readonly description: string;
  readonly value: string;
  readonly id: number;
}
interface StandardKeyValue<T> {
  [key: string]: T
}

@Component({
  selector: "time-range-scale",
  templateUrl: "./time-range-scale.component.html",
  styleUrls: ["./time-range-scale.component.css"]
})
export class TimeRangeScale {
  @Input() options: Array<scaleOption>;
  @Output() changeScale: EventEmitter<scaleOption> = new EventEmitter<scaleOption>();
  _activeIndex: number = 0;

  changeBlock(index: number, data: scaleOption): void {
    this._activeIndex = index;
    this.changeScale.emit(data);
  }

  getClassByIndex(index: number): StandardKeyValue<boolean> {
    return {
      "normal-block": true,
      "left-block": index == 0,
      "right-block": index == this.options.length - 1,
      "middle-block": index > 0 && index < this.options.length - 1,
      "active": this._activeIndex == index
    }
  }
}