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

const MAX_COUNT_PER_PAGE: number = 100;
const MAX_COUNT_PER_DRAG: number = 40;
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
  InRefreshIng: Map<LineType, boolean>;
  EventZoomChange: Subject<Object>;
  EventScaleChange: Subject<Object>;
  EventLangChangeSubscription: Subscription;
  LineOptions: Map<LineType, Object>;
  EChartInstance: Map<LineType, Object>;
  LineNamesList: Map<LineType, LineListDataModel[]>;
  LineData: Map<LineType, LinesData>;
  DropdownText: Map<LineType, string>;
  Query: Map<LineType, {model: LineListDataModel, scale: scaleOption, time_count: number, timestamp_base: number}>;
  DropInfo: Map<LineType, {isInDrop: boolean, isDropNext: boolean}>;
  IntervalSeed: Map<LineType, number>;
  NoData: Map<LineType, boolean>;
  CurValue: Map<LineType, {curFirst: number, curSecond: number}>;
  NoDataErrMsg: Map<LineType, string>;
  LineTypeSet: Set<LineType>;

  constructor(private service: DashboardService,
              private messageService: MessageService,
              private translateService: TranslateService) {
    this.EventZoomChange = new Subject<Object>();
    this.EventScaleChange = new Subject<Object>();
    this.LineNamesList = new Map<LineType, LineListDataModel[]>();
    this.DropdownText = new Map<LineType, string>();
    this.DropInfo = new Map<LineType, {isInDrop: boolean, isDropNext: boolean}>();
    this.IntervalSeed = new Map<LineType, number>();
    this.LineData = new Map<LineType, LinesData>();
    this.InRefreshIng = new Map<LineType, boolean>();
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
        this.LineOptions.get(value)["dataZoom"][0]["start"] = ZoomInfo["value"]["start"];
        this.LineOptions.get(value)["dataZoom"][0]["end"] = ZoomInfo["value"]["end"];
        if (ZoomInfo["value"]["start"] == 0 && !this.InRefreshIng.get(value)) {//get backup data
          this.NoData.set(value, false);
          this.DropInfo.get(value).isInDrop = true;
          this.DropInfo.get(value).isDropNext = true;
          this.Query.get(value).timestamp_base = DashboardComponent.getTimeStamp(this.LineData.get(value), true);
          this.Query.get(value).time_count = MAX_COUNT_PER_DRAG;
          this.refreshOneLineData(value);
        }
        else if (ZoomInfo["value"]["end"] == 100 && !this.InRefreshIng.get(value)) {//get forward data
          this.NoData.set(value, false);
          this.DropInfo.get(value).isInDrop = true;
          this.DropInfo.get(value).isDropNext = false;
          this.Query.get(value).timestamp_base = DashboardComponent.getTimeStamp(this.LineData.get(value), false) +
            this.Query.get(value).scale.valueOfSecond * MAX_COUNT_PER_DRAG;
          this.Query.get(value).time_count = MAX_COUNT_PER_DRAG;
          this.refreshOneLineData(value);
        } else {
          let newOption = Object.create(this.LineData.get(value));
          this.LineData.set(value, newOption);
        }
      });
    });
    this.EventScaleChange.asObservable().debounceTime(300).subscribe(ScaleInfo => {
      this.LineTypeSet.forEach((value) => {
        if (ScaleInfo["lineType"] != value) {
          let newQuery = Object.create(this.Query.get(value));
          newQuery.scale = ScaleInfo["value"];
          this.Query.set(value, newQuery);
          this.resetOneLineState(value);
          this.refreshOneLineData(value);
        }
      });
    });
    this.EventLangChangeSubscription = this.translateService.onLangChange.subscribe(() => {
      this.setLineBaseOption(LineType.ltService).then(res => this.LineOptions.set(LineType.ltService, res));
      this.setLineBaseOption(LineType.ltNode).then(res => this.LineOptions.set(LineType.ltNode, res));
      this.setLineBaseOption(LineType.ltStorage).then(res => this.LineOptions.set(LineType.ltStorage, res));
    });
  }

  ngOnDestroy() {
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

  private static getValue(source: LinesData, lineNum: 0 | 1, isMin: boolean): number {
    let valueArr: LineDataModel[] = source[lineNum];
    if (valueArr.length == 0) {
      return 0;
    }
    return isMin ? valueArr[0][1] : valueArr[valueArr.length - 1][1]
  }

  private static haveLessMaxTimeValue(source: LinesData, res: LinesData): boolean {
    let curMaxTimeStamp = DashboardComponent.getTimeStamp(source, false);
    let result: boolean = false;
    for (let item of res[0]) {
      if (Math.round(item[0].getTime() / 1000) <= curMaxTimeStamp) {
        result = true;
        break;
      }
    }
    return result;
  }

  private static calculateZoom(source: LinesData, res: LinesData, isDropNext: boolean): {start: number, end: number} {
    let result = {start: 100, end: 80};
    if (res[0].length > 0) {
      if (isDropNext) {
        result.start = Math.max((res[0].length / source[0].length) * 100 + 10, 20);
        result.start = result.start > 100 ? 100 : result.start;
        result.end = result.start - 20;
      } else {
        result.start = 100 - Math.max((res[0].length / source[0].length) * 100 - 10, 20);
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

  private initLine(lineType: LineType): Promise<boolean> {
    this.CurValue.set(lineType, {curFirst: 0, curSecond: 0});
    this.InRefreshIng.set(lineType, false);
    this.IntervalSeed.set(lineType, REFRESH_SEED_SERVICE);
    this.DropInfo.set(lineType, {isInDrop: false, isDropNext: false});
    return this.setLineBaseOption(lineType).then(res => {
      this.LineOptions.set(lineType, res);
      return true;
    });
  }

  private get getAllLineCanRefresh(): boolean {
    let isCan: boolean = true;
    for (let refreshing in this.InRefreshIng.values()) {
      if (refreshing) {
        isCan = false;
      }
    }
    return isCan;
  }

  private autoRefreshDada(lineType: LineType): void {
    if (this.IntervalSeed.get(lineType) > 0 && this.getAllLineCanRefresh && !this.DropInfo.get(lineType).isInDrop) {
      this.IntervalSeed.set(lineType, this.IntervalSeed.get(lineType) - 1);
      if (this.IntervalSeed.get(lineType) == 0) {
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

  private resetOneLineState(lineType: LineType) {
    this.IntervalSeed.set(lineType, REFRESH_SEED_SERVICE);
    this.Query.get(lineType).time_count = MAX_COUNT_PER_PAGE;
    this.DropInfo.get(lineType).isDropNext = false;
    this.DropInfo.get(lineType).isInDrop = false;
    this.NoData.set(lineType, false);
  }

  private refreshOneLineData(lineType: LineType): Promise<boolean> {
    let lineDropInfo = this.DropInfo.get(lineType);
    let lineOption = this.LineOptions.get(lineType);
    let query = {
      time_count: this.Query.get(lineType).time_count,
      time_unit: this.Query.get(lineType).scale.value,
      list_name: this.Query.get(lineType).model.list_name == "total" ? "" : this.Query.get(lineType).model.list_name,
      timestamp_base: this.Query.get(lineType).timestamp_base
    };
    this.InRefreshIng.set(lineType, true);
    return this.service.getLineData(lineType, query)
      .then(res => {
        if (this.DropInfo.get(lineType).isInDrop && !this.DropInfo.get(lineType).isDropNext &&
          DashboardComponent.haveLessMaxTimeValue(this.LineData.get(lineType), res)) {
          this.resetOneLineState(lineType);
          this.refreshOneLineData(lineType);//here is very danger!!
        } else {
          let eChart = this.EChartInstance.get(lineType);
          if (eChart && eChart["clear"]) {
            eChart["clear"]();
          }
          if (this.DropInfo.get(lineType).isInDrop) {
            this.LineData.set(lineType, DashboardComponent.concatLineData(this.LineData.get(lineType), res, lineDropInfo.isDropNext));
            lineOption["dataZoom"][0]["start"] = DashboardComponent.calculateZoom(this.LineData.get(lineType), res, lineDropInfo.isDropNext).start;
            lineOption["dataZoom"][0]["end"] = DashboardComponent.calculateZoom(this.LineData.get(lineType), res, lineDropInfo.isDropNext).end;
            if (res.length < MAX_COUNT_PER_DRAG){
              this.DropInfo.get(lineType).isInDrop = false;
            }
          }
          else {
            this.LineData.set(lineType, res);
            this.CurValue.set(lineType, {
              curFirst: DashboardComponent.getValue(this.LineData.get(lineType), 0, false),
              curSecond: DashboardComponent.getValue(this.LineData.get(lineType), 1, false)
            });
          }
          this.NoData.set(lineType, false);
          this.IntervalSeed.set(lineType, REFRESH_SEED_SERVICE);
        }
        this.InRefreshIng.set(lineType, false);
        return true;
      })
      .catch(err => {
        this.CurValue.set(lineType, {curFirst: 0, curSecond: 0});
        this.NoData.set(lineType, true);
        this.InRefreshIng.set(lineType, false);
        this.IntervalSeed.set(lineType, REFRESH_SEED_SERVICE);
        this.messageService.dispatchError(err);
        return false;
      });
  }

  scaleChange(lineType: LineType, data: scaleOption) {
    if (!this.InRefreshIng.get(lineType)) {
      this.Query.get(lineType).scale = data;
      this.Query.get(lineType).timestamp_base = DashboardComponent.getTimeStamp(this.LineData.get(lineType), false);
      this.resetOneLineState(lineType);
      this.refreshOneLineData(lineType).then(() => {
        if (this.getAllLineCanRefresh) {//refresh others lines
          this.EventScaleChange.next({lineType: lineType, value: data});
        }
      })
    }
  }

  dropDownChange(lineType: LineType, lineListData: LineListDataModel) {
    if (!this.InRefreshIng.get(lineType)) {
      this.Query.get(lineType).model = lineListData;
      this.DropdownText.set(lineType, lineListData.list_name);
      this.resetOneLineState(lineType);
      this.refreshOneLineData(lineType);
    }
  }

  eChartInit(lineType: LineType, eChart: Object) {
    this.EChartInstance.set(lineType, eChart);
  }

  chartDataZoom(lineType: LineType, event: Object) {
    this.EventZoomChange.next({lineType: lineType, value: event});
  }
}