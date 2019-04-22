import { AfterViewInit, ChangeDetectionStrategy, Component, HostListener, OnDestroy, OnInit } from '@angular/core';
import { Observable, Subject, Subscription } from "rxjs";
import { debounceTime, map, tap } from "rxjs/operators";
import { DashboardComponentParent } from "./dashboard.component.parent"
import { DashboardService, IQuery, IResponse, LineType } from "./dashboard.service";
import { TranslateService } from "@ngx-translate/core";
import { scaleOption } from "./time-range-scale.component/time-range-scale.component";
import { MessageService } from "../shared.service/message.service";
import { AppInitService } from "../shared.service/app-init.service";
import { SharedService } from "../shared.service/shared.service";

const MAX_COUNT_PER_PAGE: number = 200;
const MAX_COUNT_PER_DRAG: number = 100;
const AUTO_REFRESH_CUR_SEED: number = 5;

class ThirdLine {
  data: Array<[Date, number]>;

  constructor() {
    this.data = Array<[Date, number]>();
    this.data.push([new Date(), 0]);
    this.data.push([new Date(), 0]);
  }

  get values(): Array<[Date, number]> {
    return this.data;
  }

  set maxDate(date: Date) {
    this.data[0][0] = date;
  }

  set minDate(date: Date) {
    this.data[1][0] = date;
  }

  get maxDate(): Date {
    return this.data[0][0];
  }

  get minDate(): Date {
    return this.data[1][0];
  }
}

@Component({
  selector: 'dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['dashboard.component.css'],
  changeDetection: ChangeDetectionStrategy.Default
})
export class DashboardComponent extends DashboardComponentParent implements OnInit, AfterViewInit, OnDestroy {
  scaleOptions: Array<scaleOption> = [
    {id: 1, description: "DASHBOARD.MIN", value: "second", valueOfSecond: 5},
    {id: 2, description: "DASHBOARD.HR", value: "minute", valueOfSecond: 60},
    {id: 3, description: "DASHBOARD.DAY", value: "hour", valueOfSecond: 60 * 60},
    {id: 4, description: "DASHBOARD.MTH", value: "day", valueOfSecond: 60 * 60 * 24}];
  _serverTimeStamp: number;
  _autoRefreshCurInterval: number = AUTO_REFRESH_CUR_SEED;
  intervalAutoRefresh: any;
  lineOptions: Map<LineType, Object>;
  lineStateInfo: Map<LineType, { inRefreshWIP: boolean, inDrop: boolean, isDropBack: boolean, isCanAutoRefresh: boolean }>;
  lineResponses: Map<LineType, IResponse>;
  lineThirdLine: Map<LineType, ThirdLine>;
  curValue: Map<LineType, { curFirst: number, curSecond: number }>;
  noData: Map<LineType, boolean>;
  lineTypeSet: Set<LineType>;
  query: Map<LineType, { list_name: string, scale: scaleOption, baseLineTimeStamp: number, time_count: number, timestamp_base: number }>;
  eventDragChange: Subject<{ lineType: LineType, isDragBack: boolean }>;
  eventZoomBarChange: Subject<{ start: number, end: number }>;
  // eventInitChangeDetector: Subject<LineType>;
  eventLangChangeSubscription: Subscription;
  eChartInstance: Map<LineType, Object>;
  autoRefreshInterval: Map<LineType, number>;
  curRealTimeValue: Map<LineType, { curFirst: number, curSecond: number }>;
  noDataErrMsg: Map<LineType, string>;

  constructor(private service: DashboardService,
              private appInitService: AppInitService,
              private messageService: MessageService,
              private translateService: TranslateService,
              private shardService: SharedService) {
    super();
    // this.changeDetectorRef.detach();
    this.eventDragChange = new Subject<{ lineType: LineType, isDragBack: boolean }>();
    this.eventZoomBarChange = new Subject<{ start: number, end: number }>();
    // this.eventInitChangeDetector = new Subject<LineType>();
    this.lineResponses = new Map<LineType, IResponse>();
    this.lineThirdLine = new Map<LineType, ThirdLine>();
    this.query = new Map<LineType, { list_name: string, scale: scaleOption, baseLineTimeStamp: number, time_count: number, timestamp_base: number }>();
    this.lineStateInfo = new Map<LineType, { inRefreshWIP: boolean, inDrop: boolean, isDropBack: boolean, isCanAutoRefresh: boolean }>();
    this.autoRefreshInterval = new Map<LineType, number>();
    this.noData = new Map<LineType, boolean>();
    this.curValue = new Map<LineType, { curFirst: number, curSecond: number }>();
    this.curRealTimeValue = new Map<LineType, { curFirst: number, curSecond: number }>();
    this.noDataErrMsg = new Map<LineType, string>();
    this.lineOptions = new Map<LineType, Object>();
    this.eChartInstance = new Map<LineType, Object>();
    this.lineTypeSet = new Set<LineType>();
  }

  ngOnInit() {
    this.lineTypeSet.add(LineType.ltService);
    this.lineTypeSet.add(LineType.ltNode);
    this.lineTypeSet.add(LineType.ltStorage);
    // this.eventInitChangeDetector.asObservable().bufferCount(this.lineTypeSet.size).subscribe(() => this.changeDetectorRef.reattach());
    this.eventDragChange.asObservable().pipe(debounceTime(300)).subscribe(dragInfo => {
      this.lineTypeSet.forEach((value) => {
        this.refreshLineDataByDrag(value, dragInfo.isDragBack);
      });
    });
    this.eventZoomBarChange.asObservable().pipe(debounceTime(300)).subscribe((zoom: { start: number, end: number }) => {
      this.lineTypeSet.forEach((value) => {
        this.lineOptions.get(value)["dataZoom"][0]["start"] = zoom.start;
        this.lineOptions.get(value)["dataZoom"][0]["end"] = zoom.end;
        this.resetBaseLinePos(value);
      });
    });
    this.eventLangChangeSubscription = this.translateService.onLangChange.subscribe(() => {
      this.lineTypeSet.forEach((lineType: LineType) => {
        this.setLineBaseOption(lineType).subscribe(res => {
          this.lineOptions.set(lineType, res);
          this.detectChartData(lineType);
          this.clearEChart(lineType);
        });
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
    this.initAsyncLines();
    this.intervalAutoRefresh = setInterval(() => {
      this.autoRefreshCurDada();
      // this.lineTypeSet.forEach(value => {
      //   this.autoRefreshDada(value);
      // });
    }, 1000);
  }

  private detectChartData(lineType: LineType) {
    let thirdLine = this.lineThirdLine.get(lineType);
    let query = this.query.get(lineType);
    let minTimeStrap = query.timestamp_base - MAX_COUNT_PER_PAGE * query.scale.valueOfSecond;
    let maxTimeStrap = query.timestamp_base;
    let lineSeries = this.lineOptions.get(lineType);
    let data = this.lineResponses.get(lineType);
    thirdLine.maxDate = new Date(maxTimeStrap * 1000);
    thirdLine.minDate = new Date(minTimeStrap * 1000);
    let newLineOption = Object.create({});
    this.lineOptions.delete(lineType);
    lineSeries["series"][0]["data"] = data.firstLineData;
    lineSeries["series"][1]["data"] = data.secondLineData;
    lineSeries["series"][2]["data"] = this.lineThirdLine.get(lineType).values;
    Object.assign(newLineOption, lineSeries);
    this.lineOptions.set(lineType, newLineOption);
  }

  private initAsyncLines() {
    this.service.getServerTimeStamp().subscribe((res: number) => {
      this._serverTimeStamp = res;
      this.lineTypeSet.forEach((lineType: LineType) => {
        this.initLine(lineType).subscribe(() => {
          this.getOneLineData(lineType).subscribe((res: IResponse) => {
            this.lineResponses.set(lineType, res);
            this.detectChartData(lineType);
            this.curRealTimeValue.set(lineType, {
              curFirst: res.firstLineData.length > 0 ? res.firstLineData[0][1] : 0,
              curSecond: res.secondLineData.length > 0 ? res.secondLineData[0][1] : 0
            });
            // this.eventInitChangeDetector.next(lineType)
          })
        });
      });
    });
  };

  private initThirdLineDate(lineType: LineType) {
    let query = this.query.get(lineType);
    let maxTimeStamp = query.timestamp_base;
    let minTimeStamp = query.timestamp_base - query.time_count * query.scale.valueOfSecond;
    let thirdLine: ThirdLine = new ThirdLine();
    thirdLine.maxDate = new Date(maxTimeStamp * 1000);
    thirdLine.minDate = new Date(minTimeStamp * 1000);
    this.lineThirdLine.set(lineType, thirdLine);
  }

  private initLine(lineType: LineType): Observable<Object> {
    this.curValue.set(lineType, {curFirst: 0, curSecond: 0});
    this.curRealTimeValue.set(lineType, {curFirst: 0, curSecond: 0});
    this.lineStateInfo.set(lineType, {isCanAutoRefresh: true, isDropBack: false, inDrop: false, inRefreshWIP: false});
    this.autoRefreshInterval.set(lineType, this.scaleOptions[0].valueOfSecond);
    this.query.set(lineType, {
      time_count: MAX_COUNT_PER_PAGE,
      list_name: "total",
      baseLineTimeStamp: 0,
      timestamp_base: this._serverTimeStamp,
      scale: this.scaleOptions[0]
    });
    this.initThirdLineDate(lineType);
    return this.setLineBaseOption(lineType)
      .pipe(tap(option => this.lineOptions.set(lineType, option)));
  }

  private getBaseLineTimeStamp(lineType: LineType): number {
    let option = this.lineOptions.get(lineType);
    let thirdLine = this.lineThirdLine.get(lineType);
    let start = option["dataZoom"][0]["start"] / 100;
    let end = option["dataZoom"][0]["end"] / 100;
    let middlePos: number = (start + end) / 2;
    let maxDate: Date = thirdLine.maxDate;
    let minDate: Date = thirdLine.minDate;
    let maxTimeStamp = Math.round(maxDate.getTime() / 1000);
    let minTimeStamp = Math.round(minDate.getTime() / 1000);
    let screenMaxTimeStamp = maxTimeStamp - (maxTimeStamp - minTimeStamp) * (1 - end);
    let screenTimeStamp = (maxTimeStamp - minTimeStamp) * (end - start);
    return Math.round(screenMaxTimeStamp - screenTimeStamp * (1 - middlePos));
  }

  private setLineZoomByCount(lineType: LineType, isDragBack: boolean): void {
    let lineData: IResponse = this.lineResponses.get(lineType);
    let lineOption = this.lineOptions.get(lineType);
    let lineZoomStart = lineOption["dataZoom"][0]["start"];
    let lineZoomEnd = lineOption["dataZoom"][0]["end"];
    let lineZoomHalf: number = (lineZoomEnd - lineZoomStart) / 2;
    if (lineData.firstLineData.length > 0) {
      let countPercent = Math.min((MAX_COUNT_PER_DRAG / MAX_COUNT_PER_PAGE) * 100, 99);
      if (isDragBack) {
        lineOption["dataZoom"][0]["start"] = Math.min(countPercent - lineZoomHalf, 99 - 2 * lineZoomHalf);
        lineOption["dataZoom"][0]["end"] = lineOption["dataZoom"][0]["start"] + 2 * lineZoomHalf;
      } else {
        lineOption["dataZoom"][0]["end"] = Math.min(99 - countPercent + lineZoomHalf, 99);
        lineOption["dataZoom"][0]["start"] = lineOption["dataZoom"][0]["end"] - 2 * lineZoomHalf;
      }
    }
  }

  private setLineZoomByTimeStamp(lineType: LineType, lineTimeStamp: number): void {
    let thirdLine = this.lineThirdLine.get(lineType);
    let lineOption = this.lineOptions.get(lineType);
    let lineZoomStart = lineOption["dataZoom"][0]["start"];
    let lineZoomEnd = lineOption["dataZoom"][0]["end"];
    let lineZoomHalf: number = (lineZoomEnd - lineZoomStart) / 2;
    let maxDate: Date = thirdLine.maxDate;
    let minDate: Date = thirdLine.minDate;
    let maxTimeStrap = Math.round(maxDate.getTime() / 1000);
    let minTimeStrap = Math.round(minDate.getTime() / 1000);
    let percent = ((maxTimeStrap - lineTimeStamp) / (maxTimeStrap - minTimeStrap)) * 100;
    lineOption["dataZoom"][0]["start"] = Math.max(percent - lineZoomHalf, 1);
    lineOption["dataZoom"][0]["end"] = Math.min(lineOption["dataZoom"][0]["start"] + 2 * lineZoomHalf, 99);
  }

  private getOneLineData(lineType: LineType): Observable<IResponse> {
    let query = this.query.get(lineType);
    let httpQuery: IQuery = {
      time_count: query.time_count,
      time_unit: query.scale.value,
      list_name: ["total", "average"].find(value => value == query.list_name) ? "" : query.list_name,
      timestamp_base: query.timestamp_base
    };
    this.lineStateInfo.get(lineType).inRefreshWIP = true;
    return this.service.getLineData(lineType, httpQuery)
      .pipe(tap(() => {
        this.noData.set(lineType, false);
        this.lineStateInfo.get(lineType).inRefreshWIP = false;
      }, () => {
        this.lineStateInfo.get(lineType).inRefreshWIP = false;
        this.noData.set(lineType, true);
      }));
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
    let chartInstance = this.eChartInstance.get(lineType);
    if (option && chartInstance && chartInstance["getWidth"]) {
      let zoomStart = option["dataZoom"][0]["start"] / 100;
      let zoomEnd = option["dataZoom"][0]["end"] / 100;
      let eChartWidth = chartInstance["getWidth"]() - 70;
      let zoomBarWidth = eChartWidth * (zoomEnd - zoomStart);
      option["graphic"][0]["left"] = eChartWidth * (1 - zoomEnd) + zoomBarWidth * (1 - (zoomEnd + zoomStart) / 2) + 38;
      chartInstance["setOption"](option, true, false);
    }
  }


  private autoRefreshCurDada() {
    this._autoRefreshCurInterval--;
    if (this._autoRefreshCurInterval == 0) {
      this._autoRefreshCurInterval = AUTO_REFRESH_CUR_SEED;
      this.service.getServerTimeStamp().subscribe((res: number) => {
        this._serverTimeStamp = res;
        this.lineTypeSet.forEach(lineType => {
          let query: IQuery = {
            time_count: 2,
            time_unit: "second",
            list_name: "",
            timestamp_base: this._serverTimeStamp
          };
          this.service.getLineData(lineType, query).subscribe((res: IResponse) => {
            if (res.firstLineData.length > 0) {
              this.curRealTimeValue.set(lineType, {
                curFirst: res.firstLineData[0][1],
                curSecond: res.secondLineData[0][1]
              });
            }
          });
        });
      });
    }
  }

  private autoRefreshDada(lineType: LineType): void {
    if (this.autoRefreshInterval.get(lineType) > 0) {
      this.autoRefreshInterval.set(lineType, this.autoRefreshInterval.get(lineType) - 1);
      if (this.autoRefreshInterval.get(lineType) == 0) {
        this.autoRefreshInterval.set(lineType, this.query.get(lineType).scale.valueOfSecond);
        if (this.lineStateInfo.get(lineType).isCanAutoRefresh) {
          this.query.get(lineType).time_count = MAX_COUNT_PER_PAGE;
          this.query.get(lineType).timestamp_base = this._serverTimeStamp;
          this.getOneLineData(lineType).subscribe((res: IResponse) => {
            this.lineResponses.set(lineType, res);
            this.detectChartData(lineType);
            this.clearEChart(lineType);
          });
        }
      }
    }
  }

  private updateAfterDragTimeStamp(lineType: LineType, isDropBack: boolean): void {
    let query = this.query.get(lineType);
    let thirdLine = this.lineThirdLine.get(lineType);
    let maxDate: Date = thirdLine.maxDate;
    let minDate: Date = thirdLine.minDate;
    let minTimeStrap = Math.round(minDate.getTime() / 1000);
    let maxTimeStrap = Math.round(maxDate.getTime() / 1000);
    let newMaxTimeStrap: number = 0;
    let newMinTimeStrap: number = 0;
    if (isDropBack) {
      newMaxTimeStrap = maxTimeStrap - MAX_COUNT_PER_DRAG * query.scale.valueOfSecond;
      newMinTimeStrap = newMaxTimeStrap - MAX_COUNT_PER_PAGE * query.scale.valueOfSecond;
      query.timestamp_base = newMaxTimeStrap;
    } else {
      newMaxTimeStrap = Math.min(maxTimeStrap + MAX_COUNT_PER_DRAG * query.scale.valueOfSecond, this._serverTimeStamp);
      newMinTimeStrap = newMaxTimeStrap - MAX_COUNT_PER_PAGE * query.scale.valueOfSecond;
      query.timestamp_base = newMaxTimeStrap;
    }
    thirdLine.maxDate = new Date(newMaxTimeStrap * 1000);
    thirdLine.minDate = new Date(newMinTimeStrap * 1000);
  }

  private resetAfterDragTimeStamp(lineType: LineType): void {
    let query = this.query.get(lineType);
    let thirdLine = this.lineThirdLine.get(lineType);
    let maxDate: Date = thirdLine.maxDate;
    let minDate: Date = thirdLine.minDate;
    let maxTimeStrap = Math.round(maxDate.getTime() / 1000);
    let minTimeStrap = Math.round(minDate.getTime() / 1000);
    let newMaxTimeStrap = maxTimeStrap + MAX_COUNT_PER_DRAG * query.scale.valueOfSecond;
    let newMinTimeStrap = newMaxTimeStrap - MAX_COUNT_PER_PAGE * query.scale.valueOfSecond;
    thirdLine.maxDate = new Date(newMaxTimeStrap * 1000);
    thirdLine.minDate = new Date(newMinTimeStrap * 1000);
  }

  private refreshLineDataByDrag(lineType: LineType, isDragBack): void {
    let lineState = this.lineStateInfo.get(lineType);
    lineState.inDrop = true;
    lineState.isDropBack = isDragBack;
    lineState.isCanAutoRefresh = false;
    this.updateAfterDragTimeStamp(lineType, isDragBack);
    if (isDragBack) {
      this.getOneLineData(lineType).subscribe((res: IResponse) => {
        this.lineStateInfo.get(lineType).inRefreshWIP = false;
        if (!res.limit.isMin) {
          this.lineResponses.set(lineType, res);
          this.setLineZoomByCount(lineType, true);
          this.resetBaseLinePos(lineType);
          this.detectChartData(lineType);
          this.clearEChart(lineType);
        } else {
          this.resetAfterDragTimeStamp(lineType);
        }
      });
    } else {
      this.getOneLineData(lineType).subscribe((res: IResponse) => {
        if (!res.limit.isMax) {
          this.lineResponses.set(lineType, res);
          this.setLineZoomByCount(lineType, false);
          this.resetBaseLinePos(lineType);
          this.detectChartData(lineType);
          this.clearEChart(lineType);
        }
      });
    }
  }

  private afterDragZoomBar(lineType: LineType) {
    let zoomStart = this.lineOptions.get(lineType)["dataZoom"][0]["start"];
    let zoomEnd = this.lineOptions.get(lineType)["dataZoom"][0]["end"];
    let lineState = this.lineStateInfo.get(lineType);
    if (zoomStart == 0 && zoomEnd < 100 && !lineState.inRefreshWIP) {//get backup data
      this.eventDragChange.next({lineType: lineType, isDragBack: true});
    } else if (zoomEnd == 100 && zoomStart > 0 && !lineState.inRefreshWIP && !lineState.isCanAutoRefresh) {//get forward data
      this.eventDragChange.next({lineType: lineType, isDragBack: false});
    }
  }

  private setLineBaseOption(lineType: LineType): Observable<Object> {
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
    return this.translateService.get([firstKey, secondKey])
      .pipe(map(res => {
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
      }));
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
        let query = this.query.get(value);
        query.scale = data;
        query.time_count = MAX_COUNT_PER_PAGE;
        query.timestamp_base = queryTimeStamp;
        query.baseLineTimeStamp = baseLineTimeStamp;
        let maxTimeStamp = queryTimeStamp;
        let minTimeStamp = queryTimeStamp - query.scale.valueOfSecond * MAX_COUNT_PER_PAGE;
        let thirdLine = this.lineThirdLine.get(value);
        thirdLine.minDate = new Date(minTimeStamp * 1000);
        thirdLine.maxDate = new Date(maxTimeStamp * 1000);
        this.getOneLineData(value).subscribe((res: IResponse) => {
          this.lineResponses.set(value, res);
          this.setLineZoomByTimeStamp(value, query.baseLineTimeStamp);
          this.resetBaseLinePos(value);
          this.detectChartData(value);
          this.clearEChart(value);
        });
      });
    }
  }

  dropDownChange(lineType: LineType, listName: string) {
    if (!this.lineStateInfo.get(lineType).inRefreshWIP) {
      this.query.get(lineType).list_name = listName;
      this.query.get(lineType).time_count = MAX_COUNT_PER_PAGE;
      this.getOneLineData(lineType).subscribe((res: IResponse) => {
        this.lineResponses.set(lineType, res);
        this.detectChartData(lineType);
        this.clearEChart(lineType);
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
    this.eventZoomBarChange.next({start: event["start"], end: event["end"]});
  }

  getCurLineName(lineType: LineType): string {
    return this.lineResponses.get(lineType) ? this.lineResponses.get(lineType).curListName : "";
  }

  getLinesName(lineType: LineType): Array<string> {
    return this.lineResponses.get(lineType) ? this.lineResponses.get(lineType).list : null;
  }

  get StorageUnit(): string {
    return this.service.CurStorageUnit;
  };

  get grafanaViewUrl(): string {
    return `http://${this.appInitService.systemInfo['board_host']}/grafana/dashboard/db/kubernetes/`
  }

  get showGrafanaWindow(): boolean {
    return this.appInitService.isSystemAdmin;
  }

  get showMaxGrafanaWindow(): boolean {
    return this.shardService.showMaxGrafanaWindow;
  }

  set showMaxGrafanaWindow(value: boolean) {
    this.shardService.showMaxGrafanaWindow = value;
  }

  get hideMaxGrafanaWindow(): boolean {
    return !this.shardService.showMaxGrafanaWindow;
  }

  refreshLine() {
    this.lineTypeSet.forEach(value => {
      this.query.get(value).timestamp_base = this._serverTimeStamp;
      this.query.get(value).time_count = MAX_COUNT_PER_PAGE;
      this.getOneLineData(value).subscribe((res: IResponse) => {
        this.lineResponses.set(value, res);
        this.detectChartData(value);
        this.clearEChart(value);
      });
    });
  }
}
