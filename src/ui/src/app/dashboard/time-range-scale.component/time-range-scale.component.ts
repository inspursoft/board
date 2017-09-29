import { Component, Input, Output, EventEmitter, OnChanges, SimpleChanges } from "@angular/core"

export interface scaleOption {
  readonly description: string;
  readonly value: string;
  readonly valueOfSecond: number;
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
export class TimeRangeScale implements OnChanges {
  @Input() options: Array<scaleOption>;
  @Input() curScale: scaleOption;
  @Output() scaleChange: EventEmitter<scaleOption> = new EventEmitter<scaleOption>();
  _activeIndex: number = 0;

  ngOnChanges(changes: SimpleChanges) {
    if (changes["curScale"]) {
      for (let i = 0; i < this.options.length; i++) {
        if (this.options[i].id == changes["curScale"].currentValue["id"]) {
          this._activeIndex = i;
        }
      }
    }
  }

  changeBlock(index: number, data: scaleOption): void {
    this._activeIndex = index;
    this.scaleChange.emit(data);
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