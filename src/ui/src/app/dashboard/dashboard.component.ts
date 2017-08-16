import { OnInit, AfterViewInit, Component, OnDestroy } from '@angular/core';
import { Assist } from "./dashboard-assist"
import { scaleOption } from "app/dashboard/time-range-scale.component/time-range-scale.component";
import {
  DashboardService,
  LinesData,
  LineDataModel,
  LineType,
  LineListDataModel
} from "app/dashboard/dashboard.service";
import { TranslateService } from "@ngx-translate/core";
import { Subscription } from "rxjs/Subscription";
import { Subject } from "rxjs/Subject";
import { MessageService } from "../shared/message-service/message.service";

const MAX_COUNT_PER_PAGE: number = 200;
const MAX_COUNT_PER_DRAG: number = 50;
const REFRESH_SEED_SERVICE: number = 10;

@Component({
  selector: 'dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['dashboard.component.css']
})
export class DashboardComponent implements OnInit, AfterViewInit, OnDestroy {
  scaleOptions: Array<scaleOption> = [
    {"id": 1, "description": "DASHBOARD.MIN", "value": "second", valueOfSecond: 5},
    {"id": 2, "description": "DASHBOARD.HR", "value": "minute", valueOfSecond: 60},
    {"id": 3, "description": "DASHBOARD.DAY", "value": "hour", valueOfSecond: 60 * 60},
    {"id": 4, "description": "DASHBOARD.MTH", "value": "day", valueOfSecond: 60 * 60 * 24}];
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
  Query: Map<LineType, {model: LineListDataModel, scale: scaleOption, time_count: number, timestamp_base: number}>;
  IntervalSeed: Map<LineType, number>;
  NoData: Map<LineType, boolean>;
  CurValue: Map<LineType, {curFirst: number, curSecond: number}>;
  NoDataErrMsg: Map<LineType, string>;


  constructor(private service: DashboardService,
              private messageService: MessageService,
              private translateService: TranslateService) {
    this.EventZoomChange = new Subject<Object>();
    this.EventScaleChange = new Subject<Object>();
    this.LineNamesList = new Map<LineType, LineListDataModel[]>();
    this.DropdownText = new Map<LineType, string>();
    this.LineStateInfo = new Map<LineType, {InRefreshIng: boolean, InDrop: boolean, IsDropBack: boolean, IsCanAutoRefresh: boolean}>();
    this.IntervalSeed = new Map<LineType, number>();
    this.LineData = new Map<LineType, LinesData>();
    this.NoData = new Map<LineType, boolean>();
    this.CurValue = new Map<LineType, {curFirst: number, curSecond: number}>();
    this.NoDataErrMsg = new Map<LineType, string>();
    this.LineOptions = new Map<LineType, Object>();
    this.EChartInstance = new Map<LineType, Object>();
    this.LineTypeSet = new Set<LineType>();
    this.Query = new Map<LineType, {model: LineListDataModel, scale: scaleOption, time_count: number, timestamp_base: number}>();
  }

  ngOnInit() {
    this.LineTypeSet.add(LineType.ltService);
    this.LineTypeSet.add(LineType.ltNode);
    this.LineTypeSet.add(LineType.ltStorage);
    this.LineTypeSet.forEach((lineType: LineType) => {
      this.initLine(lineType)
        .then(() => this.refreshLineList(lineType))
        .then(() => this.refreshOneLineData(lineType));
    });
    this.EventZoomChange.asObservable().debounceTime(300).subscribe(ZoomInfo => {
      this.LineTypeSet.forEach((value) => {
        if (ZoomInfo["lineType"] != value) {
          this.DragZoomBar(value, {start: ZoomInfo["value"]["start"], end: ZoomInfo["value"]["end"]});
          let newLineData = Object.create(this.LineData.get(value));
          this.LineData.set(value, newLineData);
        }
      });
    });
    this.EventScaleChange.asObservable().debounceTime(300).subscribe(ScaleInfo => {
      this.LineTypeSet.forEach((value: LineType) => {
        if (ScaleInfo["lineType"] != value) {
          this.refreshOneLineData(value);
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
      this.autoRefreshDada(LineType.ltService);
      this.autoRefreshDada(LineType.ltNode);
      this.autoRefreshDada(LineType.ltStorage);
    }, 1000);
  }

  private static concatLineData(source: LinesData, res: LinesData, isAppend: boolean): LinesData {
    let result: LinesData = [Array<[Date, number]>(0), Array<[Date, number]>(0)];
    if (isAppend) {
      let bufArr1: LineDataModel[] = source[0].slice(0, source[0].length - res[0].length);
      let bufArr2: LineDataModel[] = source[1].slice(0, source[0].length - res[1].length);
      result[0] = res[0].concat(bufArr1);
      result[1] = res[1].concat(bufArr2);
    } else {
      let bufArr1: LineDataModel[] = source[0].slice(res[0].length, source[0].length);
      let bufArr2: LineDataModel[] = source[1].slice(res[1].length, source[1].length);
      result[0] = bufArr1.concat(res[0]);
      result[1] = bufArr2.concat(res[1]);
    }
    return result;
  }

  private static getTimeStamp(source: LinesData, isMinTimeStamp: boolean): number {
    let timeArr: LineDataModel[] = source[0];
    if (timeArr.length == 0) {
      return Math.round(new Date().getTime() / 1000);
    }
    let date: Date = isMinTimeStamp ? timeArr[0][0] : timeArr[timeArr.length - 1][0];
    return Math.round(date.getTime() / 1000);
  }

  private static getTimeStampByZoomBar(source: LinesData, zoomBar: {start: number, end: number}): number {
    if (source[0].length == 0) {
      return Math.round(new Date().getTime() / 1000);
    }
    let middlePos = Math.round((zoomBar.start + zoomBar.end) / 2);
    let index = source[0].length - Math.round(source[0].length / 100) * middlePos;
    let date: Date = source[0][index][0];
    return Math.round(date.getTime() / 1000);
  }

  private static getValue(source: LinesData, lineNum: 0 | 1, isMin: boolean): number {
    let valueArr: LineDataModel[] = source[lineNum];
    if (valueArr.length == 0) {
      return 0;
    }
    return isMin ? valueArr[0][1] : valueArr[valueArr.length - 1][1]
  }

  private static haveLessMaxTimeValue(source: LinesData, res: LinesData): boolean {
    let curMaxTimeStamp = DashboardComponent.getTimeStamp(source, false);
    for (let item of res[0]) {
      if (Math.round(item[0].getTime() / 1000) < curMaxTimeStamp) {
        return true;
      }
    }
    return false;
  }

  private static calculateZoom(source: LinesData, resCount: number, isDropBack: boolean): {start: number, end: number} {
    let result = {start: 100, end: 80};
    if (resCount > 0) {
      if (isDropBack) {
        result.start = Math.max((resCount / source[0].length) * 100 + 10, 20);
        result.start = result.start > 100 ? 100 : result.start;
        result.end = result.start - 20;
      } else {
        result.start = 100 - Math.max((resCount / source[0].length) * 100 - 10, 20);
        result.start = result.start < 0 ? 0 : result.start;
        result.end = result.start - 20;
      }
    }
    return result;
  }

  private refreshLineList(lineType: LineType): Promise<boolean> {
    return this.service.getLineNameList(lineType)
      .then(res => {
        let firstLineList: LineListDataModel = res[0];//default total
        this.LineNamesList.set(lineType, res);
        this.DropdownText.set(lineType, firstLineList.list_name);
        this.Query.set(lineType, {
          time_count: MAX_COUNT_PER_PAGE,
          model: firstLineList,
          timestamp_base: Math.round(new Date().getTime() / 1000),
          scale: this.scaleOptions[0]
        });
        return true;
      })
      .catch(err => this.messageService.dispatchError(err));
  }

  private DragZoomBar(lineType: LineType, ZoomInfo: {start: number, end: number}) {
    this.LineOptions.get(lineType)["dataZoom"][0]["start"] = ZoomInfo.start;
    this.LineOptions.get(lineType)["dataZoom"][0]["end"] = ZoomInfo.end;
    let lineState = this.LineStateInfo.get(lineType);
    if (ZoomInfo.start == 0 && !lineState.InRefreshIng) {//get backup data
      lineState.InDrop = true;
      lineState.IsDropBack = true;
      lineState.IsCanAutoRefresh = false;
      this.Query.get(lineType).timestamp_base = DashboardComponent.getTimeStamp(this.LineData.get(lineType), true);
      this.Query.get(lineType).time_count = MAX_COUNT_PER_DRAG;
      this.refreshOneLineData(lineType).then(resCount => {
        let lineOption = this.LineOptions.get(lineType);
        let posCount = this.LineStateInfo.get(lineType).IsCanAutoRefresh ? 0 : resCount;
        let newZoomPos = DashboardComponent.calculateZoom(this.LineData.get(lineType), posCount, true);
        lineOption["dataZoom"][0]["start"] = newZoomPos.start;
        lineOption["dataZoom"][0]["end"] = newZoomPos.end;
      });
    }
    else if (ZoomInfo.end == 100 && !lineState.InRefreshIng && !lineState.IsCanAutoRefresh) {//get forward data
      lineState.InDrop = true;
      lineState.IsDropBack = false;
      lineState.IsCanAutoRefresh = false;
      this.Query.get(lineType).timestamp_base = DashboardComponent.getTimeStamp(this.LineData.get(lineType), false) +
        this.Query.get(lineType).scale.valueOfSecond * MAX_COUNT_PER_DRAG;
      this.Query.get(lineType).time_count = MAX_COUNT_PER_DRAG;
      this.refreshOneLineData(lineType).then(resCount => {
        let lineOption = this.LineOptions.get(lineType);
        let posCount = this.LineStateInfo.get(lineType).IsCanAutoRefresh ? 0 : resCount;
        let newZoomPos = DashboardComponent.calculateZoom(this.LineData.get(lineType), posCount, false);
        lineOption["dataZoom"][0]["start"] = newZoomPos.start;
        lineOption["dataZoom"][0]["end"] = newZoomPos.end;
      });
    }
  }

  private initLine(lineType: LineType): Promise<boolean> {
    this.CurValue.set(lineType, {curFirst: 0, curSecond: 0});
    this.LineStateInfo.set(lineType, {IsCanAutoRefresh: true, IsDropBack: false, InDrop: false, InRefreshIng: false});
    this.IntervalSeed.set(lineType, REFRESH_SEED_SERVICE);
    return this.setLineBaseOption(lineType).then(res => {
      this.LineOptions.set(lineType, res);
      return true;
    });
  }

  private autoRefreshDada(lineType: LineType): void {
    if (this.IntervalSeed.get(lineType) > 0 && this.LineStateInfo.get(lineType).IsCanAutoRefresh) {
      this.IntervalSeed.set(lineType, this.IntervalSeed.get(lineType) - 1);
      if (this.IntervalSeed.get(lineType) == 0) {
        this.IntervalSeed.set(lineType, REFRESH_SEED_SERVICE);
        this.Query.get(lineType).time_count = MAX_COUNT_PER_PAGE;
        this.Query.get(lineType).timestamp_base = Math.round(new Date().getTime() / 1000);
        this.refreshOneLineData(lineType);
      }
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
        let result = Assist.getBaseOptions();
        result["tooltip"] = Assist.getTooltip(firstLineTitle, secondLineTitle);
        result["series"] = [Assist.getBaseSeries(), Assist.getBaseSeries()];
        result["series"][0]["name"] = firstLineTitle;
        result["series"][1]["name"] = secondLineTitle;
        result["dataZoom"][0]["start"] = 100;
        result["dataZoom"][0]["end"] = 80;
        result["legend"] = {data: [firstLineTitle, secondLineTitle], x: "left"};
        return result;
      })
  }

  private refreshOneLineData(lineType: LineType): Promise<number> {
    let lineState = this.LineStateInfo.get(lineType);
    let query = {
      time_count: this.Query.get(lineType).time_count,
      time_unit: this.Query.get(lineType).scale.value,
      list_name: this.Query.get(lineType).model.list_name == "total" ? "" : this.Query.get(lineType).model.list_name,
      timestamp_base: this.Query.get(lineType).timestamp_base
    };
    lineState.InRefreshIng = true;
    return this.service.getLineData(lineType, query)
      .then(res => {
        let eChart = this.EChartInstance.get(lineType);
        if (eChart && eChart["clear"]) {
          eChart["clear"]();
        }
        if (this.LineStateInfo.get(lineType).InDrop) {
          if (!lineState.IsDropBack && DashboardComponent.haveLessMaxTimeValue(this.LineData.get(lineType), res)) {
            lineState.IsDropBack = false;
            lineState.InDrop = false;
            lineState.IsCanAutoRefresh = true;
          }
          let resData = DashboardComponent.concatLineData(this.LineData.get(lineType), res, lineState.IsDropBack);
          this.LineData.set(lineType, resData);
        } else {
          this.LineData.set(lineType, res);
          this.CurValue.set(lineType, {
            curFirst: DashboardComponent.getValue(this.LineData.get(lineType), 0, false),
            curSecond: DashboardComponent.getValue(this.LineData.get(lineType), 1, false)
          });
        }
        this.NoData.set(lineType, false);
        lineState.InRefreshIng = false;
        return res[0].length;
      })
      .catch(err => {
        this.CurValue.set(lineType, {curFirst: 0, curSecond: 0});
        lineState.InRefreshIng = false;
        this.NoData.set(lineType, true);
        this.messageService.dispatchError(err);
        return 1;
      });
  }

  scaleChange(lineType: LineType, data: scaleOption) {
    let lineStateInfo = this.LineStateInfo.get(lineType);
    if (!lineStateInfo.InRefreshIng) {
      let zoomBar = {
        start: this.LineOptions.get(lineType)["dataZoom"][0]["start"],
        end: this.LineOptions.get(lineType)["dataZoom"][0]["end"]
      };
      let baseTimeStamp: number = 0;
      if (lineStateInfo.IsCanAutoRefresh) {
        baseTimeStamp = DashboardComponent.getTimeStampByZoomBar(this.LineData.get(lineType), zoomBar);
        baseTimeStamp += data.valueOfSecond * MAX_COUNT_PER_DRAG / 2;
      } else {
        baseTimeStamp = this.Query.get(lineType).timestamp_base;
      }
      this.LineTypeSet.forEach((value: LineType) => {
        this.LineStateInfo.get(value).IsCanAutoRefresh = false;
        this.Query.get(value).scale = data;
        this.Query.get(value).timestamp_base = baseTimeStamp;
      });
      this.refreshOneLineData(lineType).then(() => {
        this.EventScaleChange.next({lineType: lineType, value: data});//refresh others lines
      })
    }
  }

  dropDownChange(lineType: LineType, lineListData: LineListDataModel) {
    if (!this.LineStateInfo.get(lineType).InRefreshIng) {
      this.Query.get(lineType).model = lineListData;
      this.DropdownText.set(lineType, lineListData.list_name);
      this.refreshOneLineData(lineType);
    }
  }

  eChartInit(lineType: LineType, eChart: Object) {
    this.EChartInstance.set(lineType, eChart);
  }

  chartDataZoom(lineType: LineType, event: Object) {
    this.DragZoomBar(lineType, {start: event["start"], end: event["end"]});
    this.EventZoomChange.next({lineType: lineType, value: event});
  }

  get StorageUnit(): string {
    return this.service.CurStorageUnit;
  };
}