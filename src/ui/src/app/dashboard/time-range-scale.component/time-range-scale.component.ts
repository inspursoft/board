import { Component, Input, Output, EventEmitter, OnChanges, SimpleChanges } from '@angular/core';
import { ScaleOption } from '../dashboard.types';

interface StandardKeyValue<T> {
  [key: string]: T;
}

@Component({
  selector: 'app-time-range-scale',
  templateUrl: './time-range-scale.component.html',
  styleUrls: ['./time-range-scale.component.css']
})
export class TimeRangeScaleComponent implements OnChanges {
  @Input() options: Array<ScaleOption>;
  @Input() curScale: ScaleOption;
  @Input() disabled = false;
  @Output() scaleChange: EventEmitter<ScaleOption> = new EventEmitter<ScaleOption>();
  activeIndex = 0;

  ngOnChanges(changes: SimpleChanges) {
    if (Reflect.has(changes, 'curScale')) {
      for (let i = 0; i < this.options.length; i++) {
        if (this.options[i].id === Reflect.get(Reflect.get(changes, 'curScale').currentValue, 'id')) {
          this.activeIndex = i;
        }
      }
    }
  }

  changeBlock(index: number, data: ScaleOption): void {
    if (this.activeIndex !== index && !this.disabled) {
      this.activeIndex = index;
      this.scaleChange.emit(data);
    }
  }

  getClassByIndex(index: number): StandardKeyValue<boolean> {
    return {
      'normal-block': true,
      'left-block': index === 0,
      'right-block': index === this.options.length - 1,
      'middle-block': index > 0 && index < this.options.length - 1,
      active: this.activeIndex === index && !this.disabled,
      disabled: this.disabled
    };
  }
}
