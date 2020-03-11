import { DatePipe } from "@angular/common";

export abstract class DashboardComponentParent {
  datePipe: DatePipe = new DatePipe("en-US");

  private static getHoverValue(params: object, firstHint: string, secondHint: string): string {
    return `<div style="display: flex;flex-direction: column">
    <div style="display: flex;align-items: center">
        <div style='width: 16px;height: 16px; background-color: #2f4554;border-radius: 50%'></div>
        <div>${firstHint}:${params[0].value[1]}</div>
    </div>
    <div style="display: flex;align-items: center">
        <div style='width: 16px;height: 16px; background-color: #d48265;border-radius: 50%'></div>
        <div>${secondHint}:${params[1].value[1]}</div>
    </div>
    </div>`
  }

  public static getBaseOptions(): Object {
    return {
      grid: {x: 40, y: 30, x2: 30, y2: 60},
      dataZoom: [{type: "slider", show: true, xAxisIndex: 0}],
      graphic: [{
        type: "line",
        style: {lineWidth: 2, stroke: '#91c7ae'},
        shape: {x1: 0, y1: 30, x2: 0, y2: 240}
      }],
      xAxis: [{
        inverse: "true",
        type: "time",
        splitNumber: 10,
        splitLine: {show: false}
      }],
      yAxis: {type: "value", splitLine: {show: true}},
      color: ['#2f4554', '#d48265', '#749f83', '#ca8622', '#bda29a', '#6e7074', '#546570', '#c4ccd3']
    };
  }

  public static getBaseSeries(): object {
    return {
      type: "line",
      symbol: "circle",
      showSymbol: true,
      smooth: true,
      symbolSize: 5,
      hoverAnimation: false
    }
  }

  public static getBaseSeriesThirdLine(): object {
    return {
      type: "line",
      showSymbol: false,
      hoverAnimation: false,
      lineStyle: {normal: {opacity: 0}}
    }
  }

  public abstract onToolTipEvent(params: Object, lineType: number);

  public getTooltip(hint1: string, hint2: string, lineType: number): Object {
    return {
      trigger: "axis",
      formatter: (params) => {
        if ((params as Array<any>).length > 1) {
          this.onToolTipEvent(params, lineType);
          let xDate: Date = new Date(params[0].value[0]);
          let sDate = this.datePipe.transform(xDate, 'short');
          return sDate + DashboardComponentParent.getHoverValue(params, hint1, hint2);
        }
      },
      axisPointer: {animation: false}
    }
  }


}