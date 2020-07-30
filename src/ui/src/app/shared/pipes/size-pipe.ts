/**
 * Created by liyanq on 22/11/2017.
 */

import { Pipe, PipeTransform } from '@angular/core';

const ARR_SIZE_UNIT: Array<string> = ['B', 'KB', 'MB', 'GB', 'TB'];

@Pipe({name: 'size'})
export class SizePipe implements PipeTransform {
  transform(origin: number): string {
    let times = 0;
    let multiple = 1;
    let value: number = origin;
    while (value > 1024) {
      value = value / 1024;
      times += 1;
      multiple *= 1024;
    }
    return `${Math.round(origin / multiple * 100) / 100}${ARR_SIZE_UNIT[times]}`;
  }
}
