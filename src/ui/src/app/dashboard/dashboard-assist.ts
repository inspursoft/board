import { DatePipe } from "@angular/common";

export class Assist {
  static datePipe: DatePipe = new DatePipe("lt");

  static getHoverValue(params: object, firstHint: string, secondHint: string): string {
    return `<div style="display: flex;flex-direction: column">
    <div style="display: flex;align-items: center">
        <div style='width: 16px;height: 16px; background-color: #61a0a8;border-radius: 50%'></div>
        <div>${firstHint}:${params[0].value[1]}</div>
    </div>
    <div style="display: flex;align-items: center">
        <div style='width: 16px;height: 16px; background-color: #d48265;border-radius: 50%'></div>
        <div>${secondHint}:${params[1].value[1]}</div>
    </div>
</div>`
  }

  static getTooltip(hint1: string, hint2: string): object {
    return {
      trigger: "axis",
      formatter: (params) => {
        let xDate: Date = new Date(params[0].value[0]);
        let sDate = Assist.datePipe.transform(xDate, 'yyyy/MM/dd HH:mm:ss');
        return sDate + this.getHoverValue(params, hint1, hint2);
      },
      axisPointer: {animation: false}
    }
  }

  static getBaseOptions(): Object {
    return {
      dataZoom: [
        {
          show: true,
          xAxisIndex: 0
        }
      ],
      xAxis: [{
        inverse: "true",
        type: "time",
        splitNumber: 10,
        splitLine: {show: false}
      }],
      yAxis: {
        type: "value",
        splitLine: {show: true}
      },
      color:['#91c7ae','#749f83',  '#ca8622', '#bda29a','#6e7074', '#546570', '#c4ccd3']
    };
  }

  static getBaseSeries(): object {
    return {
      type: "line",
      showSymbol: false,
      smooth: true,
      symbolSize: 5,
      hoverAnimation: false
    }
  }
}