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
  _intervalRead: any;
  _onLangChangeSubscription: Subscription;
  _ZoomChange: Subject<Object>;
  scaleOptions: Array<scaleOption> = [
    {"id": 1, "description": "DASHBOARD.MIN", "value": "second", valueOfSecond: 5},
    {"id": 2, "description": "DASHBOARD.HR", "value": "minute", valueOfSecond: 60},
    {"id": 3, "description": "DASHBOARD.DAY", "value": "hour", valueOfSecond: 60 * 60},
    {"id": 4, "description": "DASHBOARD.MTH", "value": "day", valueOfSecond: 60 * 60 * 24}];

  ServiceLineOption: Object;
  NodeLineOption: Object;
  StorageLineOption: Object;


  LineNamesList: Map<LineType, LineListDataModel[]>;
  LineData: Map<LineType, LinesData>;
  DropdownText: Map<LineType, string>;
  OptionsBuffer: Map<LineType, {lastZoomStart: number, lastZoomEnd: number}>;
  Query: Map<LineType, {model: LineListDataModel, scale: scaleOption, time_count: number, timestamp_base: number}>;
  DropInfo: Map<LineType, {isInDrop: boolean, isDropNext: boolean}>;
  IntervalSeed: Map<LineType, number>;
  Already: Map<LineType, boolean>;
  InRefreshIng: Map<LineType, boolean>;
  NoData: Map<LineType, boolean>;
  CurValue: Map<LineType, {curFirst: number, curSecond: number}>;
  NoDataErrMsg: Map<LineType, string>;

  constructor(private service: DashboardService,
              private messageService: MessageService,
              private translateService: TranslateService) {
    this._ZoomChange = new Subject<Object>();
    this._ZoomChange.asObservable()
      .debounceTime(300)
      .subscribe(value => {
        this.NodeLineOption["dataZoom"][0]["start"] = value["start"];
        this.NodeLineOption["dataZoom"][0]["end"] = value["end"];
        let old = Object.create(this.LineData.get(LineType.ltNode));
        this.LineData.set(LineType.ltNode, old);
      });
    this.LineNamesList = new Map<LineType, LineListDataModel[]>();
    this.DropdownText = new Map<LineType, string>();
    this.OptionsBuffer = new Map<LineType, {lastZoomStart: number, lastZoomEnd: number}>();
    this.Query = new Map<LineType, {model: LineListDataModel, scale: scaleOption, time_count: number, timestamp_base: number}>();
    this.DropInfo = new Map<LineType, {isInDrop: boolean, isDropNext: boolean}>();
    this.IntervalSeed = new Map<LineType, number>();
    this.Already = new Map<LineType, boolean>();
    this.InRefreshIng = new Map<LineType, boolean>();
    this.LineData = new Map<LineType, LinesData>();
    this.NoData = new Map<LineType, boolean>();
    this.CurValue = new Map<LineType, {curFirst: number, curSecond: number}>();
    this.NoDataErrMsg = new Map<LineType, string>();
  }

  ngOnInit() {
    this.initLine(LineType.ltService);
    this.initLine(LineType.ltNode);
    this.initLine(LineType.ltStorage);
    this._onLangChangeSubscription = this.translateService.onLangChange.subscribe(() => {
      this.setLineBaseOption(LineType.ltService).then(res => this.ServiceLineOption = res);
      this.setLineBaseOption(LineType.ltNode).then(res => this.NodeLineOption = res);
      this.setLineBaseOption(LineType.ltStorage).then(res => this.StorageLineOption = res);
    });
  }

  ngOnDestroy() {
    clearInterval(this._intervalRead);
    if (this._onLangChangeSubscription) {
      this._onLangChangeSubscription.unsubscribe();
    }
  }

  ngAfterViewInit() {
    this._intervalRead = setInterval(() => {
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
    console.log(date);
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

  private initLine(lineType: LineType) {
    this.CurValue.set(lineType, {curFirst: 0, curSecond: 0});
    this.InRefreshIng.set(lineType, false);
    this.Already.set(lineType, false);
    this.IntervalSeed.set(lineType, REFRESH_SEED_SERVICE);
    this.DropInfo.set(lineType, {isInDrop: false, isDropNext: false});
    this.OptionsBuffer.set(lineType, {lastZoomStart: 100, lastZoomEnd: 80});
    switch (lineType) {
      case LineType.ltNode:
        this.setLineBaseOption(lineType).then(res => this.NodeLineOption = res);
        break;
      case LineType.ltService:
        this.setLineBaseOption(lineType).then(res => this.ServiceLineOption = res);
        break;
      case LineType.ltStorage:
        this.setLineBaseOption(lineType).then(res => this.StorageLineOption = res);
        break;
    }

    this.service.getLineNameList(lineType)
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
        this.refreshData(lineType);
      })
      .catch(err => this.messageService.dispatchError(err));
  }

  private autoRefreshDada(lineType: LineType): void {
    if (this.IntervalSeed.get(lineType) > 0 && this.Already.get(lineType)) {
      this.IntervalSeed.set(lineType, this.IntervalSeed.get(lineType) - 1);
      if (this.IntervalSeed.get(lineType) == 0 && !this.DropInfo.get(lineType).isInDrop) {
        this.refreshData(lineType);
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
        result["legend"] = {data: [firstLineTitle, secondLineTitle], x: "left"};
        return result;
      })
  }

  private resetState(lineType: LineType) {
    this.IntervalSeed.set(lineType, REFRESH_SEED_SERVICE);
    this.Query.get(lineType).time_count = MAX_COUNT_PER_PAGE;
    this.DropInfo.get(lineType).isDropNext = false;
    this.DropInfo.get(lineType).isInDrop = false;
    this.Already.set(lineType, false);
    this.OptionsBuffer.get(lineType).lastZoomEnd = 80;
    this.OptionsBuffer.get(lineType).lastZoomStart = 100;
    this.NoData.set(lineType, false);
  }

  private refreshData(lineType: LineType) {
    let lineQuery = this.Query.get(lineType);
    let lineDropInfo = this.DropInfo.get(lineType);

    let query = {
      time_count: lineQuery.time_count,
      time_unit: lineQuery.scale.value,
      list_name: lineQuery.model.list_name == "total" ? "" : lineQuery.model.list_name,
      timestamp_base: lineDropInfo.isInDrop ? lineQuery.timestamp_base : Math.round(new Date().getTime() / 1000)
    };
    this.InRefreshIng.set(lineType, true);
    this.service.getLineData(lineType, query)
      .then(res => {
        if (this.DropInfo.get(lineType).isInDrop && !this.DropInfo.get(lineType).isDropNext &&
          DashboardComponent.haveLessMaxTimeValue(this.LineData.get(lineType), res)) {
          this.resetState(lineType);
          this.refreshData(lineType);
        } else {
          if (this.DropInfo.get(lineType).isInDrop) {
            this.LineData.set(lineType, DashboardComponent.concatLineData(this.LineData.get(lineType), res, lineDropInfo.isDropNext));
            // lineOption["dataZoom"][0]["start"] = DashboardComponent.calculateZoom(this.LineData.get(lineType), res, lineDropInfo.isDropNext).start;
            // lineOption["dataZoom"][0]["end"] = DashboardComponent.calculateZoom(this.LineData.get(lineType), res, lineDropInfo.isDropNext).end;
          }
          else {
            this.LineData.set(lineType, res);
            // lineOption["dataZoom"][0]["start"] = this.OptionsBuffer.get(lineType).lastZoomStart;
            // lineOption["dataZoom"][0]["end"] = this.OptionsBuffer.get(lineType).lastZoomEnd;
            // this.CurValue.set(lineType, {
            //   curFirst: DashboardComponent.getValue(this.LineData.get(lineType), 0, true),
            //   curSecond: DashboardComponent.getValue(this.LineData.get(lineType), 1, true)
            // });
          }
          this.NoData.set(lineType, false);
          this.Already.set(lineType, true);
          this.InRefreshIng.set(lineType, false);
          this.IntervalSeed.set(lineType, REFRESH_SEED_SERVICE);
        }
      })
      .catch(err => {
        this.CurValue.set(lineType, {curFirst: 0, curSecond: 0});
        this.NoData.set(lineType, true);
        this.Already.set(lineType, true);
        this.InRefreshIng.set(lineType, false);
        this.IntervalSeed.set(lineType, REFRESH_SEED_SERVICE);
        this.messageService.dispatchError(err);
      });
  }

  scaleChange(lineType: LineType, data: scaleOption) {
    if (this.Already.get(lineType)) {
      this.Query.get(lineType).scale = data;
      this.resetState(lineType);
      this.refreshData(lineType);
    }
  }

  dropDownChange(lineType: LineType, lineListData: LineListDataModel) {
    if (this.Already.get(lineType)) {
      this.Query.get(lineType).model = lineListData;
      this.DropdownText.set(lineType, lineListData.list_name);
      this.resetState(lineType);
      this.refreshData(lineType);
    }
  }

  chartDataZoom(lineType: LineType, event: Object) {
    this.OptionsBuffer.get(lineType).lastZoomStart = event["start"];
    this.OptionsBuffer.get(lineType).lastZoomEnd = event["end"];
    this._ZoomChange.next(event);
    // if (event["start"] == 0 && !this.InRefreshIng.get(lineType)) {//get backup data
    //   this.InRefreshIng.set(lineType, true);
    //   this.Already.set(lineType, false);
    //   this.NoData.set(lineType, false);
    //   this.DropInfo.get(lineType).isInDrop = true;
    //   this.DropInfo.get(lineType).isDropNext = true;
    //   this.Query.get(lineType).timestamp_base = DashboardComponent.getTimeStamp(this.LineData.get(lineType), true);
    //   this.Query.get(lineType).time_count = MAX_COUNT_PER_DRAG;
    //   this.refreshData(lineType);
    // }
    // else if (event["end"] == 100 && this.DropInfo.get(lineType).isInDrop && !this.InRefreshIng.get(lineType)) {//get forward data
    //   this.InRefreshIng.set(lineType, true);
    //   this.Already.set(lineType, false);
    //   this.NoData.set(lineType, false);
    //   this.DropInfo.get(lineType).isInDrop = true;
    //   this.DropInfo.get(lineType).isDropNext = false;
    //   this.Query.get(lineType).timestamp_base = DashboardComponent.getTimeStamp(this.LineData.get(lineType), false) +
    //     this.Query.get(lineType).scale.valueOfSecond * MAX_COUNT_PER_DRAG;
    //   this.Query.get(lineType).time_count = MAX_COUNT_PER_DRAG;
    //   this.refreshData(lineType);
    // }
  }
}