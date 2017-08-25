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
import { promise } from "selenium-webdriver";
import checkedNodeCall = promise.checkedNodeCall;

const MAX_COUNT_PER_PAGE: number = 200;
const MAX_COUNT_PER_DRAG: number = 50;
const AUTO_REFRESH_SEED: number = 10;
const AUTO_REFRESH_CUR_SEED: number = 2;
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
  _ServerTimeStamp: number;
  _AutoRefreshCurInterval: number = AUTO_REFRESH_CUR_SEED;
  IntervalAutoRefresh: any;
  LineOptions: Map<LineType, Object>;
  LineStateInfo: Map<LineType, {InRefreshIng: boolean, InDrop: boolean, IsDropBack: boolean, IsCanAutoRefresh: boolean}>;
  LineNamesList: Map<LineType, LineListDataModel[]>;
  LineTypeSet: Set<LineType>;
  LineData: Map<LineType, LinesData>;
  EventZoomChange: Subject<Object>;
  EventScaleChange: Subject<Object>;
  EventLangChangeSubscription: Subscription;
  EChartInstance: Map<LineType, Object>;
  DropdownText: Map<LineType, string>;
  Query: Map<LineType, {model: LineListDataModel, scale: scaleOption, baseLineTimeStamp: number, time_count: number, timestamp_base: number}>;
  AutoRefreshInterval: Map<LineType, number>;
  NoData: Map<LineType, boolean>;
  CurValue: Map<LineType, {curFirst: number, curSecond: number}>;
  CurRealTimeValue: Map<LineType, {curFirst: number, curSecond: number}>;
  NoDataErrMsg: Map<LineType, string>;

  constructor(private service: DashboardService,
              private messageService: MessageService,
              private translateService: TranslateService) {
    super();
    this.EventZoomChange = new Subject<Object>();
    this.EventScaleChange = new Subject<Object>();
    this.LineNamesList = new Map<LineType, LineListDataModel[]>();
    this.DropdownText = new Map<LineType, string>();
    this.LineStateInfo = new Map<LineType, {InRefreshIng: boolean, InDrop: boolean, IsDropBack: boolean, IsCanAutoRefresh: boolean}>();
    this.AutoRefreshInterval = new Map<LineType, number>();
    this.LineData = new Map<LineType, LinesData>();
    this.NoData = new Map<LineType, boolean>();
    this.CurValue = new Map<LineType, {curFirst: number, curSecond: number}>();
    this.CurRealTimeValue = new Map<LineType, {curFirst: number, curSecond: number}>();
    this.NoDataErrMsg = new Map<LineType, string>();
    this.LineOptions = new Map<LineType, Object>();
    this.EChartInstance = new Map<LineType, Object>();
    this.LineTypeSet = new Set<LineType>();
    this.Query = new Map<LineType, {model: LineListDataModel, scale: scaleOption, baseLineTimeStamp: number, time_count: number, timestamp_base: number}>();
  }

  ngOnInit() {
    this.LineTypeSet.add(LineType.ltService);
    this.LineTypeSet.add(LineType.ltNode);
    this.LineTypeSet.add(LineType.ltStorage);
    this.service.getServerTimeStamp().then(serverTime => {
      this._ServerTimeStamp = serverTime;
      this.LineTypeSet.forEach((lineType: LineType) => {
        this.initAsyncLine(lineType);
      });
    }).catch(err => this.messageService.dispatchError(err));
    this.EventZoomChange.asObservable().debounceTime(300).subscribe(ZoomInfo => {
      this.LineTypeSet.forEach((value) => {
        let zoomStart = ZoomInfo["value"]["start"];
        let zoomEnd = ZoomInfo["value"]["end"];
        if (ZoomInfo["lineType"] != value) {
          this.DragZoomBar(value, {start: zoomStart, end: zoomEnd});
        }
        this.resetBaseLinePos(value);
      });
    });
    this.EventScaleChange.asObservable().debounceTime(300).subscribe(ScaleInfo => {
      this.LineTypeSet.forEach((value: LineType) => {
        if (ScaleInfo["lineType"] != value) {
          this.getOneLineData(value).then(res => {
            this.clearEChart(value);
            this.LineData.set(value, res);
            let lineOption = this.LineOptions.get(value);
            let baseLineTimeStamp = this.Query.get(value).baseLineTimeStamp;
            let newZoomPos = DashboardComponent.calculateZoomByTimeStamp(this.LineData.get(value), baseLineTimeStamp);
            lineOption["dataZoom"][0]["start"] = newZoomPos.start;
            lineOption["dataZoom"][0]["end"] = newZoomPos.end;
            this.resetBaseLinePos(value);
          }).catch(() => {
          })
        }
      });
    });
    this.EventLangChangeSubscription = this.translateService.onLangChange.subscribe(() => {
      this.LineTypeSet.forEach((lineType: LineType) => {
        this.setLineBaseOption(lineType).then(res => this.LineOptions.set(lineType, res));
      });
    });
  }

  ngOnDestroy() {
    this.LineTypeSet.forEach((value) => {//for update at after destroy
      this.EChartInstance.set(value, null);
    });
    clearInterval(this.IntervalAutoRefresh);
    if (this.EventLangChangeSubscription) {
      this.EventLangChangeSubscription.unsubscribe();
    }
  }

  ngAfterViewInit() {
    this.IntervalAutoRefresh = setInterval(() => {
      this.autoRefreshCurDada();
      this.LineTypeSet.forEach(value => {
        this.autoRefreshDada(value);
      });
    }, 1000);
  }

  private static getTimeStamp(source: LinesData, isMinTimeStamp: boolean): number {
    let timeArr: LineDataModel[] = source[0];
    if (timeArr.length == 0) {
      return Math.round(new Date().getTime() / 1000);
    }
    let date: Date = isMinTimeStamp ? timeArr[0][0] : timeArr[timeArr.length - 1][0];
    return Math.round(date.getTime() / 1000);
  }

  private static  concatLineData(source: LinesData, res: LinesData): LinesData {
    if (source[0].length == 0) {
      return res;
    }
    let result: LinesData = [Array<[Date, number]>(0), Array<[Date, number]>(0)];
    let maxResTimeStamp = DashboardComponent.getTimeStamp(res, false);
    let minResTimeStamp = DashboardComponent.getTimeStamp(res, true);
    let maxSourceTimeStamp = DashboardComponent.getTimeStamp(source, false);
    let minSourceTimeStamp = DashboardComponent.getTimeStamp(source, true);
    if (maxResTimeStamp < minSourceTimeStamp) {
      let bufArr1: LineDataModel[] = source[0].slice(0, source[0].length - res[0].length);
      let bufArr2: LineDataModel[] = source[1].slice(0, source[1].length - res[1].length);
      result[0] = res[0].concat(bufArr1);
      result[1] = res[1].concat(bufArr2);
    } else if (minResTimeStamp >= maxSourceTimeStamp) {
      let bufArr1: LineDataModel[] = source[0].slice(res[0].length, source[0].length);
      let bufArr2: LineDataModel[] = source[1].slice(res[0].length, source[1].length);
      result[0] = bufArr1.concat(res[0]);
      result[1] = bufArr2.concat(res[1]);
    } else if (
      (maxResTimeStamp == maxSourceTimeStamp && minResTimeStamp > minSourceTimeStamp) ||
      (minResTimeStamp == minSourceTimeStamp && maxResTimeStamp < maxSourceTimeStamp)) {
      result = source;
    }
    else if (minResTimeStamp < minSourceTimeStamp &&
      maxResTimeStamp > minSourceTimeStamp &&
      maxResTimeStamp < maxSourceTimeStamp) {//almost impossible
      let maxResIndex: number = 0;
      source[0].forEach((value, index, arr) => {
        if (index < arr.length - 1 && arr[index][1] <= minResTimeStamp && arr[index + 1][1] > minResTimeStamp) {
          maxResIndex = index;
        }
      });
      if (maxResIndex > 0) {
        let bufArr1: LineDataModel[] = source[0].slice(maxResIndex + 1, source[0].length);
        let bufArr2: LineDataModel[] = source[1].slice(maxResIndex + 1, source[0].length);
        result[0] = res[0].concat(bufArr1);
        result[1] = res[1].concat(bufArr2);
      }
    } else if (minResTimeStamp < maxSourceTimeStamp &&
      maxResTimeStamp > maxSourceTimeStamp &&
      minResTimeStamp > minSourceTimeStamp) {//almost impossible
      let minResIndex: number = 0;
      for (let i = 1; i < source[0].length; i++) {
        let curTimeStamp = Math.round(source[0][i][0].getTime() / 1000);
        let preTimeStamp = Math.round(source[0][i - 1][0].getTime() / 1000);
        if (preTimeStamp <= minResTimeStamp && curTimeStamp > minResTimeStamp) {
          minResIndex = i;
        }
      }
      if (minResIndex > 0) {
        let bufArr1: LineDataModel[] = source[0].slice(0, minResIndex - 1);
        let bufArr2: LineDataModel[] = source[1].slice(0, minResIndex - 1);
        result[0] = bufArr1.concat(res[0]);
        result[1] = bufArr2.concat(res[1]);
      }
    }
    return result;
  }

  private static getBaseLineTimeStamp(source: LinesData, zoomBar: {start: number, end: number}): number {
    let start = zoomBar.start / 100;
    let end = zoomBar.end / 100;
    let middlePos: number = (start + end) / 2;
    let maxTimeStamp = DashboardComponent.getTimeStamp(source, false);
    let minTimeStamp = DashboardComponent.getTimeStamp(source, true);
    let screenMaxTimeStamp = maxTimeStamp - (maxTimeStamp - minTimeStamp) * (1 - end);
    let screenTimeStamp = (maxTimeStamp - minTimeStamp) * (end - start);
    return Math.round(screenMaxTimeStamp - screenTimeStamp * (1 - middlePos));
  }

  private setLineZoomByCount(lineType: LineType, resCount: number, isDropBack: boolean): void {
    if (resCount > 0) {
      let lineData = this.LineData.get(lineType);
      let lineOption = this.LineOptions.get(lineType);
      let lineZoomStart = lineOption["dataZoom"][0]["start"];
      let lineZoomEnd = lineOption["dataZoom"][0]["end"];
      let lineZoomDif: number = lineZoomEnd - lineZoomStart;
      let countPercent = (resCount / lineData[0].length) * 100;
      if (isDropBack) {
        if (lineZoomEnd + countPercent > 100) {
          lineZoomEnd = 100;
          lineZoomStart = lineZoomEnd - lineZoomDif;
        } else {
          lineZoomEnd += countPercent;
          lineZoomStart += countPercent;
        }
      } else {
        if (lineZoomStart - countPercent < 0) {
          lineZoomStart = 0;
          lineZoomEnd = lineZoomStart + lineZoomDif;
        } else {
          lineZoomStart -= countPercent;
          lineZoomEnd -= countPercent;
        }
      }
      lineOption["dataZoom"][0]["start"] = lineZoomStart;
      lineOption["dataZoom"][0]["end"] = lineZoomEnd;
    }
  }

  private static calculateZoomByTimeStamp(source: LinesData, baseLineTimeStrap: number): {start: number, end: number} {
    let result = {start: 80, end: 100};
    if (source[0].length > 0) {
      let maxTimeStamp = DashboardComponent.getTimeStamp(source, false);
      let minTimeStamp = DashboardComponent.getTimeStamp(source, true);
      if (baseLineTimeStrap > minTimeStamp && baseLineTimeStrap < maxTimeStamp) {
        let middlePos = (maxTimeStamp - baseLineTimeStrap) / (maxTimeStamp - minTimeStamp) * 100;
        result.start = middlePos - 10;
        result.end = middlePos + 10;
      }
    }
    if (result.end >= 100 || result.start <= 0) {
      result.start = 80;
      result.end = 100;
    }
    return result;
  }

  private async initAsyncLine(lineType: LineType) {
    await this.initLine(lineType);
    await this.getOneLineData(lineType)
      .then(res => this.LineData.set(lineType, res))
      .catch(() => {
      });
  };

  private getOneLineData(lineType: LineType): Promise<LinesData> {
    let query = {
      time_count: this.Query.get(lineType).time_count,
      time_unit: this.Query.get(lineType).scale.value,
      list_name: this.Query.get(lineType).model.list_name == "total" ? "" : this.Query.get(lineType).model.list_name,
      timestamp_base: this.Query.get(lineType).timestamp_base
    };
    this.LineStateInfo.get(lineType).InRefreshIng = true;
    return this.service.getLineData(lineType, query)
      .then((res: {List: Array<LineListDataModel>, Data: LinesData, CurListName: string, Limit: {isMax: boolean, isMin: boolean}}) => {
        this.NoData.set(lineType, false);
        this.LineNamesList.set(lineType, res.List);
        this.DropdownText.set(lineType, res.CurListName);
        this.LineStateInfo.get(lineType).InRefreshIng = false;
        console.log(res.Limit);
        return res.Data;
      })
      .catch(err => {
        this.LineStateInfo.get(lineType).InRefreshIng = false;
        this.NoData.set(lineType, true);
        this.messageService.dispatchError(err);
      });
  }

  private getLineInRefreshIng(): boolean {
    let iter: IterableIterator<LineType> = this.LineTypeSet.values();
    let iterResult: IteratorResult<LineType> = iter.next();
    while (!iterResult.done) {
      if (this.LineStateInfo.get(iterResult.value).InRefreshIng) {
        return true;
      }
      iterResult = iter.next();
    }
    return false;
  }

  private clearEChart(lineType: LineType): void {
    let eChart = this.EChartInstance.get(lineType);
    if (eChart && eChart["clear"]) {
      eChart["clear"]();
    }
  }

  private resetBaseLinePos(lineType: LineType) {
    let option = this.LineOptions.get(lineType);
    let zoomStart = option["dataZoom"][0]["start"] / 100;
    let zoomEnd = option["dataZoom"][0]["end"] / 100;
    let eChartWidth = this.EChartInstance.get(lineType)["getWidth"]() - 70;
    let zoomBarWidth = eChartWidth * (zoomEnd - zoomStart);
    option["graphic"][0]["left"] = eChartWidth * (1 - zoomEnd) + zoomBarWidth * (1 - (zoomEnd + zoomStart) / 2) + 38;
    this.EChartInstance.get(lineType)["setOption"](option, true, false);
    if (this.LineData.get(lineType)) {
      this.clearEChart(lineType);
      this.LineData.set(lineType, Object.create(this.LineData.get(lineType)));
    }
  }

  private  delayNormal(lineType: LineType): Promise<boolean> {
    return new Promise<boolean>((resolve, reject) => {
      this.LineStateInfo.get(lineType).InRefreshIng = true;
      setTimeout(() => {
        this.LineStateInfo.get(lineType).InRefreshIng = false;
        resolve(true);
      }, 200)
    });
  }

  private initLine(lineType: LineType): Promise<boolean> {
    this.CurValue.set(lineType, {curFirst: 0, curSecond: 0});
    this.CurRealTimeValue.set(lineType, {curFirst: 0, curSecond: 0});
    this.LineStateInfo.set(lineType, {IsCanAutoRefresh: true, IsDropBack: false, InDrop: false, InRefreshIng: false});
    this.AutoRefreshInterval.set(lineType, AUTO_REFRESH_SEED);
    this.Query.set(lineType, {
      time_count: MAX_COUNT_PER_PAGE,
      model: {list_name: "", time_stamp: 0},
      baseLineTimeStamp: 0,
      timestamp_base: this._ServerTimeStamp,
      scale: this.scaleOptions[0]
    });
    return this.setLineBaseOption(lineType).then(res => {
      this.LineOptions.set(lineType, res);
      return true;
    });
  }

  private async autoRefreshCurDada() {
    this._AutoRefreshCurInterval--;
    if (this._AutoRefreshCurInterval == 0) {
      this._AutoRefreshCurInterval = AUTO_REFRESH_CUR_SEED;
      await this.service.getServerTimeStamp().then(serverTime => this._ServerTimeStamp = serverTime);
      this.LineTypeSet.forEach(lineType => {
        let query = {
          time_count: 1,
          time_unit: this.Query.get(lineType).scale.value,
          list_name: this.Query.get(lineType).model.list_name == "total" ? "" : this.Query.get(lineType).model.list_name,
          timestamp_base: this._ServerTimeStamp
        };
        this.service.getLineData(lineType, query)
          .then((res: {
            List: Array<LineListDataModel>,
            Data: LinesData,
            CurListName: string,
            Limit: {isMax: boolean, isMin: boolean}
          }) => {
            if (res.Data[0].length > 0) {
              this.CurRealTimeValue.set(lineType, {
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
    if (this.AutoRefreshInterval.get(lineType) > 0) {
      this.AutoRefreshInterval.set(lineType, this.AutoRefreshInterval.get(lineType) - 1);
      if (this.AutoRefreshInterval.get(lineType) == 0) {
        this.AutoRefreshInterval.set(lineType, AUTO_REFRESH_SEED);
        if (this.LineStateInfo.get(lineType).IsCanAutoRefresh) {
          this.Query.get(lineType).time_count = MAX_COUNT_PER_PAGE;
          this.Query.get(lineType).timestamp_base = this._ServerTimeStamp;
          this.getOneLineData(lineType).then(res => {
            this.clearEChart(lineType);
            this.LineData.set(lineType, res);
          }).catch(() => {
          });
        }
      }
    }
  }

  private addForNoneData(lineType: LineType, isDragBack: boolean): LinesData {
    let query = this.Query.get(lineType);
    let result: LinesData = [Array<[Date, number]>(0), Array<[Date, number]>(0)];
    let maxTimeStamp: number = 0;
    let minTimeStamp: number = 0;
    if (isDragBack) {
      maxTimeStamp = query.timestamp_base;
      minTimeStamp = query.timestamp_base - query.scale.valueOfSecond * MAX_COUNT_PER_DRAG;
    } else {
      maxTimeStamp = query.timestamp_base + query.scale.valueOfSecond * MAX_COUNT_PER_DRAG;
      minTimeStamp = query.timestamp_base;
    }
    result[0].push([new Date(minTimeStamp * 1000), 0]);
    result[0].push([new Date(maxTimeStamp * 1000), 0]);
    result[1].push([new Date(minTimeStamp * 1000), 0]);
    result[1].push([new Date(maxTimeStamp * 1000), 0]);
    return result;
  }

  private DragZoomBar(lineType: LineType, ZoomInfo: {start: number, end: number}) {
    this.LineOptions.get(lineType)["dataZoom"][0]["start"] = ZoomInfo.start;
    this.LineOptions.get(lineType)["dataZoom"][0]["end"] = ZoomInfo.end;
    this.resetBaseLinePos(lineType);
    let lineState = this.LineStateInfo.get(lineType);
    if (ZoomInfo.start == 0 && !lineState.InRefreshIng) {//get backup data
      lineState.InDrop = true;
      lineState.IsDropBack = true;
      lineState.IsCanAutoRefresh = false;
      let query = this.Query.get(lineType);
      let lineData = this.LineData.get(lineType);
      query.timestamp_base = DashboardComponent.getTimeStamp(lineData, true) - query.scale.valueOfSecond;
      query.time_count = MAX_COUNT_PER_DRAG;
      this.getOneLineData(lineType).then(res => {
        this.delayNormal(lineType).then(() => {
          let resData: LinesData = [Array<[Date, number]>(0), Array<[Date, number]>(0)];
          if (res[0].length == 0) {
            resData = this.addForNoneData(lineType, true);
          } else {
            resData = DashboardComponent.concatLineData(this.LineData.get(lineType), res);
          }
          this.clearEChart(lineType);
          this.LineData.set(lineType, resData);
          this.setLineZoomByCount(lineType, res[0].length, true);
          this.resetBaseLinePos(lineType);
        });
      }).catch(() => {
      });
    }
    else if (ZoomInfo.end == 100 && lineState.InDrop && !lineState.InRefreshIng && !lineState.IsCanAutoRefresh) {//get forward data
      lineState.InDrop = true;
      lineState.IsDropBack = false;
      lineState.IsCanAutoRefresh = false;
      let query = this.Query.get(lineType);
      let lineData = this.LineData.get(lineType);
      query.time_count = MAX_COUNT_PER_DRAG;
      query.timestamp_base = DashboardComponent.getTimeStamp(lineData, false) + query.scale.valueOfSecond * MAX_COUNT_PER_DRAG;
      console.log(new Date(query.timestamp_base * 1000));
      this.getOneLineData(lineType).then(res => {
        this.delayNormal(lineType).then(() => {//add delay for drag
          // if (DashboardComponent.haveLessMaxTimeValue(lineData, res)) {//need api support
          //   lineState.IsDropBack = false;
          //   lineState.InDrop = false;
          //   lineState.IsCanAutoRefresh = true;
          //   this.IntervalSeed.set(lineType, 1);
          // }
          // else {
          let resData: LinesData = [Array<[Date, number]>(0), Array<[Date, number]>(0)];
          if (res[0].length == 0) {
            resData = this.addForNoneData(lineType, false);
          } else {
            resData = DashboardComponent.concatLineData(this.LineData.get(lineType), res);
          }
          this.clearEChart(lineType);
          this.LineData.set(lineType, resData);
          this.setLineZoomByCount(lineType, res[0].length, false);
          this.resetBaseLinePos(lineType);
          // }
        });
      }).catch(() => {
      });
    }
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
        result["series"] = [DashboardComponentParent.getBaseSeries(), DashboardComponentParent.getBaseSeries()];
        result["series"][0]["name"] = firstLineTitle;
        result["series"][1]["name"] = secondLineTitle;
        result["dataZoom"][0]["start"] = 80;
        result["dataZoom"][0]["end"] = 100;
        result["legend"] = {data: [firstLineTitle, secondLineTitle], x: "left"};
        return result;
      })
  }


  scaleChange(lineType: LineType, data: scaleOption) {
    if (!this.getLineInRefreshIng()) {
      let option = this.LineOptions.get(lineType);
      let zoomBar = {start: option["dataZoom"][0]["start"], end: option["dataZoom"][0]["end"]};
      let baseLineTimeStamp = DashboardComponent.getBaseLineTimeStamp(this.LineData.get(lineType), zoomBar);
      let queryTimeStamp = baseLineTimeStamp + data.valueOfSecond * MAX_COUNT_PER_PAGE / 2;
      this.LineTypeSet.forEach((value: LineType) => {
        this.LineStateInfo.get(value).IsCanAutoRefresh = false;
        this.Query.get(value).scale = data;
        this.Query.get(value).time_count = MAX_COUNT_PER_PAGE;
        this.Query.get(value).timestamp_base = queryTimeStamp;
        this.Query.get(value).baseLineTimeStamp = baseLineTimeStamp;
      });
      this.getOneLineData(lineType).then(res => {
        let query = this.Query.get(lineType);
        if (res[0].length == 0) {
          let maxTimeStamp = query.timestamp_base;
          let minTimeStamp = query.timestamp_base - query.scale.valueOfSecond * MAX_COUNT_PER_PAGE;
          res[0].push([new Date(minTimeStamp * 1000), 0]);
          res[0].push([new Date(maxTimeStamp * 1000), 0]);
          res[1].push([new Date(minTimeStamp * 1000), 0]);
          res[1].push([new Date(maxTimeStamp * 1000), 0]);
        }
        this.clearEChart(lineType);
        this.LineData.set(lineType, res);
        let lineOption = this.LineOptions.get(lineType);
        let baseLineTimeStamp = query.baseLineTimeStamp;
        let newZoomPos = DashboardComponent.calculateZoomByTimeStamp(this.LineData.get(lineType), baseLineTimeStamp);
        lineOption["dataZoom"][0]["start"] = newZoomPos.start;
        lineOption["dataZoom"][0]["end"] = newZoomPos.end;
        this.resetBaseLinePos(lineType);
        this.EventScaleChange.next({lineType: lineType, value: data});//refresh others lines
      }).catch(() => {
      });
    }
  }

  dropDownChange(lineType: LineType, lineListData: LineListDataModel) {
    if (!this.LineStateInfo.get(lineType).InRefreshIng) {
      this.Query.get(lineType).model = lineListData;
      this.Query.get(lineType).time_count = MAX_COUNT_PER_PAGE;
      this.getOneLineData(lineType).then(res => {
        this.clearEChart(lineType);
        this.LineData.set(lineType, res);
      }).catch(() => {
      });
    }
  }

  public onToolTipEvent(params: Object, lineType: LineType) {
    this.CurValue.set(lineType, {curFirst: params[0].value[1], curSecond: params[1].value[1]});
  }

  public onEChartInit(lineType: LineType, eChart: Object) {
    this.EChartInstance.set(lineType, eChart);
    this.resetBaseLinePos(lineType);
  }

  @HostListener("window:resize", ["$event"])
  onEChartWindowResize(event: Object) {
    this.LineTypeSet.forEach(value => {
      this.delayNormal(value).then(() => {
        this.resetBaseLinePos(value);
      })
    });
  }

  chartDataZoom(lineType: LineType, event: Object) {
    this.DragZoomBar(lineType, {start: event["start"], end: event["end"]});
    this.EventZoomChange.next({lineType: lineType, value: event});
  }

  get StorageUnit(): string {
    return this.service.CurStorageUnit;
  };
}