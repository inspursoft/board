import { OnInit, AfterViewInit, Component, OnDestroy, HostListener } from '@angular/core';
import { DashboardComponentParent } from "./dashboard.component.parent"
import { scaleOption } from "app/dashboard/time-range-scale.component/time-range-scale.component";
import {
  DashboardService, LinesData, LineDataModel, LineType, LineListDataModel
} from "app/dashboard/dashboard.service";
import { TranslateService } from "@ngx-translate/core";
import { Subscription } from "rxjs/Subscription";
import { Subject } from "rxjs/Subject";
import { MessageService } from "../shared/message-service/message.service";

const MAX_COUNT_PER_PAGE: number = 200;
const MAX_COUNT_PER_DRAG: number = 100;
const AUTO_REFRESH_SEED: number = 10;
const AUTO_REFRESH_CUR_SEED: number = 5;
@Component({
  selector: 'dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['dashboard.component.css']
})
export class DashboardComponent extends DashboardComponentParent implements OnInit, AfterViewInit, OnDestroy {
  scaleOptions: Array<scaleOption> = [
    {"id": 1, "description": "DASHBOARD.MIN", "value": "second", valueOfSecond: 5},
    {"id": 2, "description": "DASHBOARD.HR", "value": "minute", valueOfSecond: 60},
    {"id": 3, "description": "DASHBOARD.DAY", "value": "hour", valueOfSecond: 60 * 60},
    {"id": 4, "description": "DASHBOARD.MTH", "value": "day", valueOfSecond: 60 * 60 * 24}];
  _serverTimeStamp: number;
  _autoRefreshCurInterval: number = 1;
  intervalAutoRefresh: any;
  lineOptions: Map<LineType, Object>;
  lineStateInfo: Map<LineType, {inRefreshWIP: boolean, inDrop: boolean, isDropBack: boolean, isCanAutoRefresh: boolean}>;
  lineNamesList: Map<LineType, LineListDataModel[]>;
  lineTypeSet: Set<LineType>;
  lineData: Map<LineType, LinesData>;
  eventDragChange: Subject<{lineType: LineType, isDragBack: boolean}>;
  eventScaleChange: Subject<Object>;
  eventZoomBarChange: Subject<LineType>;
  eventLangChangeSubscription: Subscription;
  eChartInstance: Map<LineType, Object>;
  dropdownText: Map<LineType, string>;
  query: Map<LineType, {model: LineListDataModel, scale: scaleOption, baseLineTimeStamp: number, time_count: number, timestamp_base: number}>;
  autoRefreshInterval: Map<LineType, number>;
  noData: Map<LineType, boolean>;
  curValue: Map<LineType, {curFirst: number, curSecond: number}>;
  curRealTimeValue: Map<LineType, {curFirst: number, curSecond: number}>;
  noDataErrMsg: Map<LineType, string>;

  constructor(private service: DashboardService,
              private messageService: MessageService,
              private translateService: TranslateService) {
    super();
    this.eventDragChange = new Subject<{lineType: LineType, isDragBack: boolean}>();
    this.eventScaleChange = new Subject<Object>();
    this.eventZoomBarChange = new Subject<LineType>();
    this.lineNamesList = new Map<LineType, LineListDataModel[]>();
    this.dropdownText = new Map<LineType, string>();
    this.lineStateInfo = new Map<LineType, {inRefreshWIP: boolean, inDrop: boolean, isDropBack: boolean, isCanAutoRefresh: boolean}>();
    this.autoRefreshInterval = new Map<LineType, number>();
    this.lineData = new Map<LineType, LinesData>();
    this.noData = new Map<LineType, boolean>();
    this.curValue = new Map<LineType, {curFirst: number, curSecond: number}>();
    this.curRealTimeValue = new Map<LineType, {curFirst: number, curSecond: number}>();
    this.noDataErrMsg = new Map<LineType, string>();
    this.lineOptions = new Map<LineType, Object>();
    this.eChartInstance = new Map<LineType, Object>();
    this.lineTypeSet = new Set<LineType>();
    this.query = new Map<LineType, {model: LineListDataModel, scale: scaleOption, baseLineTimeStamp: number, time_count: number, timestamp_base: number}>();
  }

  ngOnInit() {
    this.lineTypeSet.add(LineType.ltService);
    this.lineTypeSet.add(LineType.ltNode);
    this.lineTypeSet.add(LineType.ltStorage);
    this.service.getServerTimeStamp().then(serverTime => {
      this._serverTimeStamp = serverTime;
      this.lineTypeSet.forEach((lineType: LineType) => {
        this.initAsyncLine(lineType);
      });
    }).catch(err => this.messageService.dispatchError(err));
    this.eventDragChange.asObservable().debounceTime(300).subscribe(dragInfo => {
      this.lineTypeSet.forEach((value) => {
        if (dragInfo.lineType != value) {
          this.refreshLineDataByDrag(value, dragInfo.isDragBack);
        }
        this.resetBaseLinePos(value);
      });
    });
    this.eventZoomBarChange.asObservable().debounceTime(300).subscribe(lineType => {
      this.lineTypeSet.forEach((value) => {
        this.resetBaseLinePos(value);
      });
    });
    this.eventScaleChange.asObservable().debounceTime(300).subscribe(ScaleInfo => {
      this.lineTypeSet.forEach((value: LineType) => {
        if (ScaleInfo["lineType"] != value) {
          this.getOneLineData(value).then(res => {
            this.clearEChart(value);
            this.lineData.set(value, res.Data);
            this.setLineZoomByTimeStamp(value, this.query.get(value).baseLineTimeStamp);
            this.resetBaseLinePos(value);
          }).catch(() => {
          })
        }
      });
    });
    this.eventLangChangeSubscription = this.translateService.onLangChange.subscribe(() => {
      this.lineTypeSet.forEach((lineType: LineType) => {
        this.setLineBaseOption(lineType).then(res => this.lineOptions.set(lineType, res));
      });
    });
  }

  ngOnDestroy() {
    this.lineTypeSet.forEach((value) => {//for update at after destroy
      this.eChartInstance.set(value, null);
    });
    clearInterval(this.intervalAutoRefresh);
    if (this.eventLangChangeSubscription) {
      this.eventLangChangeSubscription.unsubscribe();
    }
  }

  ngAfterViewInit() {
    this.intervalAutoRefresh = setInterval(() => {
      this.autoRefreshCurDada();
      this.lineTypeSet.forEach(value => {
        this.autoRefreshDada(value);
      });
    }, 1000);
  }

  private getBaseLineTimeStamp(lineType: LineType): number {
    let option = this.lineOptions.get(lineType);
    let lineData = this.lineData.get(lineType);
    let start = option["dataZoom"][0]["start"] / 100;
    let end = option["dataZoom"][0]["end"] / 100;
    let middlePos: number = (start + end) / 2;
    let maxDate: Date = lineData[2][0][0];
    let minDate: Date = lineData[2][1][0];
    let maxTimeStamp = Math.round(maxDate.getTime() / 1000);
    let minTimeStamp = Math.round(minDate.getTime() / 1000);
    let screenMaxTimeStamp = maxTimeStamp - (maxTimeStamp - minTimeStamp) * (1 - end);
    let screenTimeStamp = (maxTimeStamp - minTimeStamp) * (end - start);
    return Math.round(screenMaxTimeStamp - screenTimeStamp * (1 - middlePos));
  }

  private setLineZoomByCount(lineType: LineType, resCount: number, isDragBack: boolean): void {
    if (resCount > 0) {
      let lineData = this.lineData.get(lineType);
      let lineOption = this.lineOptions.get(lineType);
      let lineZoomStart = lineOption["dataZoom"][0]["start"];
      let lineZoomEnd = lineOption["dataZoom"][0]["end"];
      let lineZoomHalf: number = (lineZoomEnd - lineZoomStart) / 2;
      if (lineData[0].length > 0) {
        let countPercent = Math.min((resCount / lineData[0].length) * 100, 99);
        if (isDragBack) {
          lineOption["dataZoom"][0]["start"] = Math.min(countPercent - lineZoomHalf, 99 - 2 * lineZoomHalf);
          lineOption["dataZoom"][0]["end"] = lineOption["dataZoom"][0]["start"] + 2 * lineZoomHalf;
        } else {
          lineOption["dataZoom"][0]["end"] = Math.min(99 - countPercent + lineZoomHalf, 99);
          lineOption["dataZoom"][0]["start"] = lineOption["dataZoom"][0]["end"] - 2 * lineZoomHalf;
        }
      }
    }
  }

  private  setLineZoomByTimeStamp(lineType: LineType, lineTimeStamp: number): void {
    let lineData = this.lineData.get(lineType);
    let lineOption = this.lineOptions.get(lineType);
    let lineZoomStart = lineOption["dataZoom"][0]["start"];
    let lineZoomEnd = lineOption["dataZoom"][0]["end"];
    let lineZoomHalf: number = (lineZoomEnd - lineZoomStart) / 2;
    let maxDate: Date = lineData[2][0][0];
    let minDate: Date = lineData[2][1][0];
    let maxTimeStrap = Math.round(maxDate.getTime() / 1000);
    let minTimeStrap = Math.round(minDate.getTime() / 1000);
    let percent = ((maxTimeStrap - lineTimeStamp) / (maxTimeStrap - minTimeStrap)) * 100;
    lineOption["dataZoom"][0]["start"] = Math.max(percent - lineZoomHalf, 1);
    lineOption["dataZoom"][0]["end"] = Math.min(lineOption["dataZoom"][0]["start"] + 2 * lineZoomHalf, 99);
  }

  private async initAsyncLine(lineType: LineType) {
    await this.initLine(lineType);
    await this.getOneLineData(lineType)
      .then(res => this.lineData.set(lineType, res.Data))
      .catch(() => {
      });
  };

  private getOneLineData(lineType: LineType): Promise<{Data: LinesData, Limit: {isMax: boolean, isMin: boolean}}> {
    let query = this.query.get(lineType);
    let httpQuery = {
      time_count: query.time_count,
      time_unit: query.scale.value,
      list_name: query.model.list_name == "total" ? "" : query.model.list_name,
      timestamp_base: query.timestamp_base,
      service_duration_time: query.timestamp_base - query.time_count * query.scale.valueOfSecond
    };
    this.lineStateInfo.get(lineType).inRefreshWIP = true;
    return this.service.getlineData(lineType, httpQuery)
      .then((res: {List: Array<LineListDataModel>, Data: LinesData, CurListName: string, Limit: {isMax: boolean, isMin: boolean}}) => {
        this.noData.set(lineType, false);
        this.lineNamesList.set(lineType, res.List);
        this.dropdownText.set(lineType, res.CurListName);
        this.lineStateInfo.get(lineType).inRefreshWIP = false;
        return {Data: res.Data, Limit: res.Limit};
      })
      .catch(err => {
        this.lineStateInfo.get(lineType).inRefreshWIP = false;
        this.noData.set(lineType, true);
        this.messageService.dispatchError(err);
      });
  }

  private getLineInRefreshWIP(): boolean {
    let iter: IterableIterator<LineType> = this.lineTypeSet.values();
    let iterResult: IteratorResult<LineType> = iter.next();
    while (!iterResult.done) {
      if (this.lineStateInfo.get(iterResult.value).inRefreshWIP) {
        return true;
      }
      iterResult = iter.next();
    }
    return false;
  }

  private clearEChart(lineType: LineType): void {
    let eChart = this.eChartInstance.get(lineType);
    if (eChart && eChart["clear"]) {
      eChart["clear"]();
    }
  }

  private resetBaseLinePos(lineType: LineType) {
    let option = this.lineOptions.get(lineType);
    if (option["dataZoom"]) {
      let zoomStart = option["dataZoom"][0]["start"] / 100;
      let zoomEnd = option["dataZoom"][0]["end"] / 100;
      let eChartWidth = this.eChartInstance.get(lineType)["getWidth"]() - 70;
      let zoomBarWidth = eChartWidth * (zoomEnd - zoomStart);
      option["graphic"][0]["left"] = eChartWidth * (1 - zoomEnd) + zoomBarWidth * (1 - (zoomEnd + zoomStart) / 2) + 38;
      this.eChartInstance.get(lineType)["setOption"](option, true, false);
      if (this.lineData.get(lineType)) {
        this.clearEChart(lineType);
        this.lineData.set(lineType, Object.create(this.lineData.get(lineType)));
      }
    }
  }

  private  delayNormal(lineType: LineType): Promise<boolean> {
    return new Promise<boolean>((resolve, reject) => {
      this.lineStateInfo.get(lineType).inRefreshWIP = true;
      setTimeout(() => {
        this.lineStateInfo.get(lineType).inRefreshWIP = false;
        resolve(true);
      }, 200)
    });
  }

  private initLine(lineType: LineType): Promise<boolean> {
    this.curValue.set(lineType, {curFirst: 0, curSecond: 0});
    this.curRealTimeValue.set(lineType, {curFirst: 0, curSecond: 0});
    this.lineStateInfo.set(lineType, {isCanAutoRefresh: true, isDropBack: false, inDrop: false, inRefreshWIP: false});
    this.autoRefreshInterval.set(lineType, AUTO_REFRESH_SEED);
    this.query.set(lineType, {
      time_count: MAX_COUNT_PER_PAGE,
      model: {list_name: "", time_stamp: 0},
      baseLineTimeStamp: 0,
      timestamp_base: this._serverTimeStamp,
      scale: this.scaleOptions[0]
    });
    return this.setLineBaseOption(lineType).then(res => {
      this.lineOptions.set(lineType, res);
      return true;
    });
  }

  private async autoRefreshCurDada() {
    this._autoRefreshCurInterval--;
    if (this._autoRefreshCurInterval == 0) {
      this._autoRefreshCurInterval = AUTO_REFRESH_CUR_SEED;
      await this.service.getServerTimeStamp().then(serverTime => this._serverTimeStamp = serverTime);
      this.lineTypeSet.forEach(lineType => {
        let query = {
          time_count: 1,
          time_unit: "second",
          list_name: "",
          timestamp_base: this._serverTimeStamp
        };
        this.service.getlineData(lineType, query)
          .then((res: {
            List: Array<LineListDataModel>,
            Data: LinesData,
            CurListName: string,
            Limit: {isMax: boolean, isMin: boolean}
          }) => {
            if (res.Data[0].length > 0) {
              this.curRealTimeValue.set(lineType, {
                curFirst: res.Data[0][0][1],
                curSecond: res.Data[1][0][1]
              })
            }
          })
          .catch(err => this.messageService.dispatchError(err));
      });
    }
  }

  private autoRefreshDada(lineType: LineType): void {
    if (this.autoRefreshInterval.get(lineType) > 0) {
      this.autoRefreshInterval.set(lineType, this.autoRefreshInterval.get(lineType) - 1);
      if (this.autoRefreshInterval.get(lineType) == 0) {
        this.autoRefreshInterval.set(lineType, AUTO_REFRESH_SEED);
        if (this.lineStateInfo.get(lineType).isCanAutoRefresh) {
          this.query.get(lineType).time_count = MAX_COUNT_PER_PAGE;
          this.query.get(lineType).timestamp_base = this._serverTimeStamp;
          this.getOneLineData(lineType).then(res => {
            this.clearEChart(lineType);
            this.lineData.set(lineType, res.Data);
          }).catch(() => {
          });
        }
      }
    }
  }

  private  updateAfterDragTimeStamp(lineType: LineType, isDropBack: boolean): void {
    let query = this.query.get(lineType);
    let lineData = this.lineData.get(lineType);
    let maxDate: Date = lineData[2][0][0];
    let minDate: Date = lineData[2][1][0];
    let minTimeStrap = Math.round(minDate.getTime() / 1000);
    let maxTimeStrap = Math.round(maxDate.getTime() / 1000);
    query.time_count = MAX_COUNT_PER_DRAG;
    let newMaxTimeStrap: number = 0;
    let newMinTimeStrap: number = 0;
    if (isDropBack) {
      newMaxTimeStrap = maxTimeStrap - MAX_COUNT_PER_DRAG * query.scale.valueOfSecond;
      newMinTimeStrap = newMaxTimeStrap - MAX_COUNT_PER_PAGE * query.scale.valueOfSecond;
      query.timestamp_base = minTimeStrap;
    } else {
      newMaxTimeStrap = Math.min(maxTimeStrap + MAX_COUNT_PER_DRAG * query.scale.valueOfSecond, this._serverTimeStamp);
      newMinTimeStrap = newMaxTimeStrap - MAX_COUNT_PER_PAGE * query.scale.valueOfSecond;
      query.timestamp_base = newMaxTimeStrap;
    }
    lineData[2][0][0] = new Date(newMaxTimeStrap * 1000);
    lineData[2][1][0] = new Date(newMinTimeStrap * 1000);
  }

  private resetAfterDraglineData(lineType: LineType): void {
    let lineData = this.lineData.get(lineType);
    let maxDate: Date = lineData[2][0][0];
    let minDate: Date = lineData[2][1][0];
    let minTimeStrap = Math.round(minDate.getTime() / 1000);
    let maxTimeStrap = Math.round(maxDate.getTime() / 1000);
    lineData[0] = lineData[0].filter((value) => {
      let date = value[0];
      let timeStrap = Math.round(date.getTime() / 1000);
      return timeStrap > minTimeStrap && timeStrap < maxTimeStrap;
    });
    lineData[1] = lineData[1].filter((value) => {
      let date = value[0];
      let timeStrap = Math.round(date.getTime() / 1000);
      return timeStrap > minTimeStrap && timeStrap < maxTimeStrap;
    });
  }

  private resetAfterDragTimeStamp(lineType: LineType): void {
    let query = this.query.get(lineType);
    let lineData = this.lineData.get(lineType);
    let maxDate: Date = lineData[2][0][0];
    let minDate: Date = lineData[2][1][0];
    let maxTimeStrap = Math.round(maxDate.getTime() / 1000);
    let minTimeStrap = Math.round(minDate.getTime() / 1000);
    let newMaxTimeStrap = maxTimeStrap + MAX_COUNT_PER_DRAG * query.scale.valueOfSecond;
    let newMinTimeStrap = newMaxTimeStrap - MAX_COUNT_PER_PAGE * query.scale.valueOfSecond;
    lineData[2][0][0] = new Date(newMaxTimeStrap * 1000);
    lineData[2][1][0] = new Date(newMinTimeStrap * 1000);
  }

  private filterMaxLineData(this: LineDataModel[], value: [Date, number]): boolean {
    let date = value[0];
    let timeStrap = Math.round(date.getTime() / 1000);
    if (this.length > 0) {
      let maxAlreadyDate = this[this.length - 1][0];
      return timeStrap > Math.round(maxAlreadyDate.getTime() / 1000);
    }
    return true;
  }

  private filterMinlineData(this: LineDataModel[], value: [Date, number]): boolean {
    let date = value[0];
    let timeStrap = Math.round(date.getTime() / 1000);
    if (this.length > 0) {
      let minAlreadyDate = this[0][0];
      return timeStrap < Math.round(minAlreadyDate.getTime() / 1000);
    }
    return true;
  }

  private concatLineData(lineType: LineType, res: LinesData, isDropBack: boolean): number {
    let lineData = this.lineData.get(lineType);
    if (!isDropBack) {
      let newData1 = res[0].filter(this.filterMaxLineData, lineData[0]);
      let newData2 = res[1].filter(this.filterMaxLineData, lineData[1]);
      lineData[0] = lineData[0].concat(newData1);
      lineData[1] = lineData[1].concat(newData2);
      return newData2.length;
    } else {
      let newData1 = res[0].filter(this.filterMinlineData, lineData[0]);
      let newData2 = res[1].filter(this.filterMinlineData, lineData[1]);
      lineData[0] = newData1.concat(lineData[0]);
      lineData[1] = newData2.concat(lineData[1]);
      return newData2.length;
    }
  }

  private refreshLineDataByDrag(lineType: LineType, isDragBack) {
    let lineState = this.lineStateInfo.get(lineType);
    if (isDragBack) {
      lineState.inDrop = true;
      lineState.isDropBack = true;
      lineState.isCanAutoRefresh = false;
      this.updateAfterDragTimeStamp(lineType, true);
      this.getOneLineData(lineType).then(res => {
        this.delayNormal(lineType).then(() => {
          if (!res.Limit.isMin) {
            this.clearEChart(lineType);
            this.resetAfterDraglineData(lineType);
            let newCount = this.concatLineData(lineType, res.Data, true);
            this.setLineZoomByCount(lineType, newCount, true);
            this.resetBaseLinePos(lineType);
          } else {
            this.resetAfterDragTimeStamp(lineType);
          }
        });
      }).catch(() => {
      });
    } else {
      lineState.inDrop = true;
      lineState.isDropBack = false;
      lineState.isCanAutoRefresh = false;
      this.updateAfterDragTimeStamp(lineType, false);
      this.getOneLineData(lineType).then(res => {
        this.delayNormal(lineType).then(() => {//add delay for drag
          if (!res.Limit.isMax) {
            this.clearEChart(lineType);
            this.resetAfterDraglineData(lineType);
            let newCount = this.concatLineData(lineType, res.Data, false);
            this.setLineZoomByCount(lineType, newCount, false);
            this.resetBaseLinePos(lineType);
          }
        });
      }).catch(() => {
      });
    }
  }

  private afterDragZoomBar(lineType: LineType) {
    let zoomStart = this.lineOptions.get(lineType)["dataZoom"][0]["start"];
    let zoomEnd = this.lineOptions.get(lineType)["dataZoom"][0]["end"];
    let lineState = this.lineStateInfo.get(lineType);
    if (zoomStart == 0 && zoomEnd < 100 && !lineState.inRefreshWIP) {//get backup data
      this.refreshLineDataByDrag(lineType, true);
      this.eventDragChange.next({lineType: lineType, isDragBack: true});
    }
    else if (zoomEnd == 100 && zoomStart > 0 && !lineState.inRefreshWIP && !lineState.isCanAutoRefresh) {//get forward data
      this.refreshLineDataByDrag(lineType, false);
      this.eventDragChange.next({lineType: lineType, isDragBack: false});
    }
  }

  private ZoomBarChange(lineType: LineType, ZoomInfo: {start: number, end: number}) {
    this.lineOptions.get(lineType)["dataZoom"][0]["start"] = ZoomInfo.start;
    this.lineOptions.get(lineType)["dataZoom"][0]["end"] = ZoomInfo.end;
    this.eventZoomBarChange.next();
  }

  private setLineBaseOption(lineType: LineType): Promise<Object> {
    let firstKey, secondKey: string;
    switch (lineType) {
      case LineType.ltService: {
        firstKey = "DASHBOARD.PODS";
        secondKey = "DASHBOARD.CONTAINERS";
        break;
      }
      case LineType.ltNode: {
        firstKey = "DASHBOARD.CPU";
        secondKey = "DASHBOARD.MEMORY";
        break;
      }
      case LineType.ltStorage: {
        firstKey = "DASHBOARD.USAGE";
        secondKey = "DASHBOARD.TOTAL";
        break;
      }
    }
    return this.translateService.get([firstKey, secondKey]).toPromise()
      .then(res => {
        let firstLineTitle: string = res[firstKey];
        let secondLineTitle: string = res[secondKey];
        let result = DashboardComponentParent.getBaseOptions();
        result["tooltip"] = this.getTooltip(firstLineTitle, secondLineTitle, lineType);
        result["series"] = [
          DashboardComponentParent.getBaseSeries(),
          DashboardComponentParent.getBaseSeries(),
          DashboardComponentParent.getBaseSeriesThirdLine()];
        result["series"][0]["name"] = firstLineTitle;
        result["series"][1]["name"] = secondLineTitle;
        result["dataZoom"][0]["start"] = 80;
        result["dataZoom"][0]["end"] = 100;
        result["legend"] = {data: [firstLineTitle, secondLineTitle], x: "left"};
        return result;
      })
  }

  scaleChange(lineType: LineType, data: scaleOption) {
    if (!this.getLineInRefreshWIP()) {
      let baseLineTimeStamp = this.getBaseLineTimeStamp(lineType);
      let queryTimeStamp = 0;
      let maxLineTimeStamp = baseLineTimeStamp + data.valueOfSecond * MAX_COUNT_PER_PAGE / 2;
      if (maxLineTimeStamp > this._serverTimeStamp) {
        queryTimeStamp = this._serverTimeStamp;
        baseLineTimeStamp -= maxLineTimeStamp - this._serverTimeStamp
      } else {
        queryTimeStamp = maxLineTimeStamp;
      }
      this.lineTypeSet.forEach((value: LineType) => {
        this.query.get(value).scale = data;
        this.query.get(value).time_count = MAX_COUNT_PER_PAGE;
        this.query.get(value).timestamp_base = queryTimeStamp;
        this.query.get(value).baseLineTimeStamp = baseLineTimeStamp;
      });
      this.getOneLineData(lineType).then(res => {
        let query = this.query.get(lineType);
        if (res.Data[0].length == 0) {
          let maxTimeStamp = query.timestamp_base;
          let minTimeStamp = query.timestamp_base - query.scale.valueOfSecond * MAX_COUNT_PER_PAGE;
          res.Data[0].push([new Date(minTimeStamp * 1000), 0]);
          res.Data[0].push([new Date(maxTimeStamp * 1000), 0]);
          res.Data[1].push([new Date(minTimeStamp * 1000), 0]);
          res.Data[1].push([new Date(maxTimeStamp * 1000), 0]);
        }
        this.clearEChart(lineType);
        this.lineData.set(lineType, res.Data);
        this.setLineZoomByTimeStamp(lineType, query.baseLineTimeStamp);
        this.resetBaseLinePos(lineType);
        this.eventScaleChange.next({lineType: lineType, value: data});//refresh others lines
      }).catch(() => {
      });
    }
  }

  dropDownChange(lineType: LineType, lineListData: LineListDataModel) {
    if (!this.lineStateInfo.get(lineType).inRefreshWIP) {
      this.query.get(lineType).model = lineListData;
      this.query.get(lineType).time_count = MAX_COUNT_PER_PAGE;
      this.getOneLineData(lineType).then(res => {
        this.clearEChart(lineType);
        this.lineData.set(lineType, res.Data);
      }).catch(() => {
      });
    }
  }

  public onToolTipEvent(params: Object, lineType: LineType) {
    if ((params as Array<any>).length > 1) {
      this.curValue.set(lineType, {curFirst: params[0].value[1], curSecond: params[1].value[1]});
    }
  }

  public onEChartInit(lineType: LineType, eChart: Object) {
    this.eChartInstance.set(lineType, eChart);
    this.resetBaseLinePos(lineType);
  }

  @HostListener("window:resize", ["$event"])
  onEChartWindowResize(event: Object) {
    this.lineTypeSet.forEach(value => {
      this.resetBaseLinePos(value);
    });
  }

  chartMouseUp(lineType: LineType, event: Object) {
    this.afterDragZoomBar(lineType);
  }

  chartDataZoom(lineType: LineType, event: Object) {
    this.lineTypeSet.forEach((value) => {
      this.ZoomBarChange(value, {start: event["start"], end: event["end"]});
    });
  }

  get StorageUnit(): string {
    return this.service.CurStorageUnit;
  };
}