import { AfterViewInit, ChangeDetectionStrategy, Component, HostListener, OnDestroy, OnInit } from '@angular/core';
import { Observable, Subject, Subscription } from 'rxjs';
import { debounceTime, map, tap } from 'rxjs/operators';
import { DashboardComponentParent } from './dashboard.component.parent';
import { DashboardService, IQuery, IResponse, LineType } from './dashboard.service';
import { TranslateService } from '@ngx-translate/core';
import { scaleOption } from './time-range-scale.component/time-range-scale.component';
import { MessageService } from '../shared.service/message.service';
import { AppInitService } from '../shared.service/app-init.service';
import { SharedService } from '../shared.service/shared.service';

const MAX_COUNT_PER_PAGE = 200;
const MAX_COUNT_PER_DRAG = 100;
const AUTO_REFRESH_CUR_SEED = 5;
const dataZoomKey = 'dataZoom';
const seriesKey = 'series';
const tooltipKey = 'tooltip';
const legendKey = 'legend';
const getWidthKey = 'getWidth';
const graphicKey = 'graphic';
const setOptionKey = 'setOption';
const clearKey = 'clear';

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

  get maxDate(): Date {
    return this.data[0][0];
  }

  set minDate(date: Date) {
    this.data[1][0] = date;
  }

  get minDate(): Date {
    return this.data[1][0];
  }
}

@Component({
  templateUrl: './dashboard.component.html',
  styleUrls: ['dashboard.component.css'],
  changeDetection: ChangeDetectionStrategy.Default
})
export class DashboardComponent extends DashboardComponentParent implements OnInit, AfterViewInit, OnDestroy {
  scaleOptions: Array<scaleOption> = [
    {id: 1, description: 'DASHBOARD.MIN', value: 'second', valueOfSecond: 5},
    {id: 2, description: 'DASHBOARD.HR', value: 'minute', valueOfSecond: 60},
    {id: 3, description: 'DASHBOARD.DAY', value: 'hour', valueOfSecond: 60 * 60},
    {id: 4, description: 'DASHBOARD.MTH', value: 'day', valueOfSecond: 60 * 60 * 24}];
  serverTimeStamp: number;
  autoRefreshCurInterval: number = AUTO_REFRESH_CUR_SEED;
  intervalAutoRefresh: any;
  lineOptions: Map<LineType, object>;
  lineStateInfo: Map<LineType, {
    inRefreshWIP: boolean,
    inDrop: boolean,
    isDropBack: boolean,
    isCanAutoRefresh: boolean
  }>;
  lineResponses: Map<LineType, IResponse>;
  lineThirdLine: Map<LineType, ThirdLine>;
  curValue: Map<LineType, {
    curFirst: number,
    curFirstUnit: string,
    curSecond: number,
    curSecondUnit: string
  }>;
  noData: Map<LineType, boolean>;
  lineTypeSet: Set<LineType>;
  query: Map<LineType, {
    list_name: string,
    scale: scaleOption,
    baseLineTimeStamp: number,
    time_count: number,
    timestamp_base: number
  }>;
  eventDragChange: Subject<{ lineType: LineType, isDragBack: boolean }>;
  eventZoomBarChange: Subject<{ start: number, end: number }>;
  // eventInitChangeDetector: Subject<LineType>;
  eventLangChangeSubscription: Subscription;
  eChartInstance: Map<LineType, object>;
  autoRefreshInterval: Map<LineType, number>;
  curRealTimeValue: Map<LineType, {
    curFirst: number,
    curFirstUnit: string,
    curSecond: number,
    curSecondUnit: string
  }>;
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
    this.query = new Map<LineType, {
      list_name: string,
      scale: scaleOption,
      baseLineTimeStamp: number,
      time_count: number,
      timestamp_base: number
    }>();
    this.lineStateInfo = new Map<LineType, {
      inRefreshWIP: boolean,
      inDrop: boolean,
      isDropBack: boolean,
      isCanAutoRefresh: boolean
    }>();
    this.autoRefreshInterval = new Map<LineType, number>();
    this.noData = new Map<LineType, boolean>();
    this.curValue = new Map<LineType, {
      curFirst: number,
      curFirstUnit: string,
      curSecond: number,
      curSecondUnit: string
    }>();
    this.curRealTimeValue = new Map<LineType, {
      curFirst: number,
      curFirstUnit: string,
      curSecond: number,
      curSecondUnit: string
    }>();
    this.noDataErrMsg = new Map<LineType, string>();
    this.lineOptions = new Map<LineType, object>();
    this.eChartInstance = new Map<LineType, object>();
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
        this.lineOptions.get(value)[dataZoomKey][0].start = zoom.start;
        this.lineOptions.get(value)[dataZoomKey][0].end = zoom.end;
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
    this.lineTypeSet.forEach((value) => {// for update at after destroy
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
    const thirdLine = this.lineThirdLine.get(lineType);
    const query = this.query.get(lineType);
    const minTimeStrap = query.timestamp_base - MAX_COUNT_PER_PAGE * query.scale.valueOfSecond;
    const maxTimeStrap = query.timestamp_base;
    const lineSeries = this.lineOptions.get(lineType);
    const data = this.lineResponses.get(lineType);
    thirdLine.maxDate = new Date(maxTimeStrap * 1000);
    thirdLine.minDate = new Date(minTimeStrap * 1000);
    const newLineOption = Object.create({});
    this.lineOptions.delete(lineType);
    lineSeries[seriesKey][0].data = data.firstLineData;
    lineSeries[seriesKey][1].data = data.secondLineData;
    lineSeries[seriesKey][2].data = this.lineThirdLine.get(lineType).values;
    Object.assign(newLineOption, lineSeries);
    this.lineOptions.set(lineType, newLineOption);
  }

  private initAsyncLines() {
    this.service.getServerTimeStamp().subscribe((res: number) => {
      this.serverTimeStamp = res;
      this.lineTypeSet.forEach((lineType: LineType) => {
        this.initLine(lineType).subscribe(() => {
          this.getOneLineData(lineType).subscribe((response: IResponse) => {
            this.lineResponses.set(lineType, response);
            this.detectChartData(lineType);
            this.curRealTimeValue.set(lineType, {
              curFirst: response.firstLineData.length > 0 ? response.firstLineData[0][1] : 0,
              curFirstUnit: response.firstLineData.length > 0 ? response.firstLineData[0][2] : '',
              curSecond: response.secondLineData.length > 0 ? response.secondLineData[0][1] : 0,
              curSecondUnit: response.secondLineData.length > 0 ? response.secondLineData[0][2] : ''
            });
            // this.eventInitChangeDetector.next(lineType)
          });
        });
      });
    });
  }

  private initThirdLineDate(lineType: LineType) {
    const query = this.query.get(lineType);
    const maxTimeStamp = query.timestamp_base;
    const minTimeStamp = query.timestamp_base - query.time_count * query.scale.valueOfSecond;
    const thirdLine: ThirdLine = new ThirdLine();
    thirdLine.maxDate = new Date(maxTimeStamp * 1000);
    thirdLine.minDate = new Date(minTimeStamp * 1000);
    this.lineThirdLine.set(lineType, thirdLine);
  }

  private initLine(lineType: LineType): Observable<object> {
    this.curValue.set(lineType, {curFirst: 0, curFirstUnit: '', curSecond: 0, curSecondUnit: ''});
    this.curRealTimeValue.set(lineType, {curFirst: 0, curFirstUnit: '', curSecond: 0, curSecondUnit: ''});
    this.lineStateInfo.set(lineType, {isCanAutoRefresh: true, isDropBack: false, inDrop: false, inRefreshWIP: false});
    this.autoRefreshInterval.set(lineType, this.scaleOptions[0].valueOfSecond);
    this.query.set(lineType, {
      time_count: MAX_COUNT_PER_PAGE,
      list_name: 'total',
      baseLineTimeStamp: 0,
      timestamp_base: this.serverTimeStamp,
      scale: this.scaleOptions[0]
    });
    this.initThirdLineDate(lineType);
    return this.setLineBaseOption(lineType)
      .pipe(tap(option => this.lineOptions.set(lineType, option)));
  }

  private getBaseLineTimeStamp(lineType: LineType): number {
    const option = this.lineOptions.get(lineType);
    const thirdLine = this.lineThirdLine.get(lineType);
    const start = option[dataZoomKey][0].start / 100;
    const end = option[dataZoomKey][0].end / 100;
    const middlePos = (start + end) / 2;
    const maxDate = thirdLine.maxDate;
    const minDate = thirdLine.minDate;
    const maxTimeStamp = Math.round(maxDate.getTime() / 1000);
    const minTimeStamp = Math.round(minDate.getTime() / 1000);
    const screenMaxTimeStamp = maxTimeStamp - (maxTimeStamp - minTimeStamp) * (1 - end);
    const screenTimeStamp = (maxTimeStamp - minTimeStamp) * (end - start);
    return Math.round(screenMaxTimeStamp - screenTimeStamp * (1 - middlePos));
  }

  private setLineZoomByCount(lineType: LineType, isDragBack: boolean): void {
    const lineData: IResponse = this.lineResponses.get(lineType);
    const lineOption = this.lineOptions.get(lineType);
    const lineZoomStart = lineOption[dataZoomKey][0].start;
    const lineZoomEnd = lineOption[dataZoomKey][0].end;
    const lineZoomHalf = (lineZoomEnd - lineZoomStart) / 2;
    if (lineData.firstLineData.length > 0) {
      const countPercent = Math.min((MAX_COUNT_PER_DRAG / MAX_COUNT_PER_PAGE) * 100, 99);
      if (isDragBack) {
        lineOption[dataZoomKey][0].start = Math.min(countPercent - lineZoomHalf, 99 - 2 * lineZoomHalf);
        lineOption[dataZoomKey][0].end = lineOption[dataZoomKey][0].start + 2 * lineZoomHalf;
      } else {
        lineOption[dataZoomKey][0].end = Math.min(99 - countPercent + lineZoomHalf, 99);
        lineOption[dataZoomKey][0].start = lineOption[dataZoomKey][0].end - 2 * lineZoomHalf;
      }
    }
  }

  private setLineZoomByTimeStamp(lineType: LineType, lineTimeStamp: number): void {
    const thirdLine = this.lineThirdLine.get(lineType);
    const lineOption = this.lineOptions.get(lineType);
    const lineZoomStart = lineOption[dataZoomKey][0].start;
    const lineZoomEnd = lineOption[dataZoomKey][0].end;
    const lineZoomHalf: number = (lineZoomEnd - lineZoomStart) / 2;
    const maxDate: Date = thirdLine.maxDate;
    const minDate: Date = thirdLine.minDate;
    const maxTimeStrap = Math.round(maxDate.getTime() / 1000);
    const minTimeStrap = Math.round(minDate.getTime() / 1000);
    const percent = ((maxTimeStrap - lineTimeStamp) / (maxTimeStrap - minTimeStrap)) * 100;
    lineOption[dataZoomKey][0].start = Math.max(percent - lineZoomHalf, 1);
    lineOption[dataZoomKey][0].end = Math.min(lineOption[dataZoomKey][0].start + 2 * lineZoomHalf, 99);
  }

  private getOneLineData(lineType: LineType): Observable<IResponse> {
    const query = this.query.get(lineType);
    const httpQuery: IQuery = {
      time_count: query.time_count,
      time_unit: query.scale.value,
      list_name: ['total', 'average'].find(value => value === query.list_name) ? '' : query.list_name,
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
    const iter: IterableIterator<LineType> = this.lineTypeSet.values();
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
    const eChart = this.eChartInstance.get(lineType);
    if (eChart && eChart[clearKey]) {
      eChart[clearKey]();
    }
  }

  private resetBaseLinePos(lineType: LineType) {
    const option = this.lineOptions.get(lineType);
    const chartInstance = this.eChartInstance.get(lineType);
    if (option && chartInstance && chartInstance[getWidthKey]) {
      const zoomStart = option[dataZoomKey][0].start / 100;
      const zoomEnd = option[dataZoomKey][0].end / 100;
      const eChartWidth = chartInstance[getWidthKey]() - 70;
      const zoomBarWidth = eChartWidth * (zoomEnd - zoomStart);
      option[graphicKey][0].left = eChartWidth * (1 - zoomEnd) + zoomBarWidth * (1 - (zoomEnd + zoomStart) / 2) + 38;
      chartInstance[setOptionKey](option, true, false);
    }
  }


  private autoRefreshCurDada() {
    this.autoRefreshCurInterval--;
    if (this.autoRefreshCurInterval === 0) {
      this.autoRefreshCurInterval = AUTO_REFRESH_CUR_SEED;
      this.service.getServerTimeStamp().subscribe((res: number) => {
        this.serverTimeStamp = res;
        this.lineTypeSet.forEach(lineType => {
          const query: IQuery = {
            time_count: 2,
            time_unit: 'second',
            list_name: '',
            timestamp_base: this.serverTimeStamp
          };
          this.service.getLineData(lineType, query).subscribe((response: IResponse) => {
            if (response.firstLineData.length > 0) {
              this.curRealTimeValue.set(lineType, {
                curFirst: response.firstLineData[0][1],
                curSecondUnit: response.firstLineData[0][2],
                curSecond: response.secondLineData[0][1],
                curFirstUnit: response.secondLineData[0][2]
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
      if (this.autoRefreshInterval.get(lineType) === 0) {
        this.autoRefreshInterval.set(lineType, this.query.get(lineType).scale.valueOfSecond);
        if (this.lineStateInfo.get(lineType).isCanAutoRefresh) {
          this.query.get(lineType).time_count = MAX_COUNT_PER_PAGE;
          this.query.get(lineType).timestamp_base = this.serverTimeStamp;
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
    const query = this.query.get(lineType);
    const thirdLine = this.lineThirdLine.get(lineType);
    const maxDate: Date = thirdLine.maxDate;
    const minDate: Date = thirdLine.minDate;
    const minTimeStrap = Math.round(minDate.getTime() / 1000);
    const maxTimeStrap = Math.round(maxDate.getTime() / 1000);
    let newMaxTimeStrap = 0;
    let newMinTimeStrap = 0;
    if (isDropBack) {
      newMaxTimeStrap = maxTimeStrap - MAX_COUNT_PER_DRAG * query.scale.valueOfSecond;
      newMinTimeStrap = newMaxTimeStrap - MAX_COUNT_PER_PAGE * query.scale.valueOfSecond;
      query.timestamp_base = newMaxTimeStrap;
    } else {
      newMaxTimeStrap = Math.min(maxTimeStrap + MAX_COUNT_PER_DRAG * query.scale.valueOfSecond, this.serverTimeStamp);
      newMinTimeStrap = newMaxTimeStrap - MAX_COUNT_PER_PAGE * query.scale.valueOfSecond;
      query.timestamp_base = newMaxTimeStrap;
    }
    thirdLine.maxDate = new Date(newMaxTimeStrap * 1000);
    thirdLine.minDate = new Date(newMinTimeStrap * 1000);
  }

  private resetAfterDragTimeStamp(lineType: LineType): void {
    const query = this.query.get(lineType);
    const thirdLine = this.lineThirdLine.get(lineType);
    const maxDate: Date = thirdLine.maxDate;
    const minDate: Date = thirdLine.minDate;
    const maxTimeStrap = Math.round(maxDate.getTime() / 1000);
    const minTimeStrap = Math.round(minDate.getTime() / 1000);
    const newMaxTimeStrap = maxTimeStrap + MAX_COUNT_PER_DRAG * query.scale.valueOfSecond;
    const newMinTimeStrap = newMaxTimeStrap - MAX_COUNT_PER_PAGE * query.scale.valueOfSecond;
    thirdLine.maxDate = new Date(newMaxTimeStrap * 1000);
    thirdLine.minDate = new Date(newMinTimeStrap * 1000);
  }

  private refreshLineDataByDrag(lineType: LineType, isDragBack): void {
    const lineState = this.lineStateInfo.get(lineType);
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
    const zoomStart = this.lineOptions.get(lineType)[dataZoomKey][0].start;
    const zoomEnd = this.lineOptions.get(lineType)[dataZoomKey][0].end;
    const lineState = this.lineStateInfo.get(lineType);
    if (zoomStart === 0 && zoomEnd < 100 && !lineState.inRefreshWIP) {// get backup data
      this.eventDragChange.next({lineType, isDragBack: true});
    } else if (zoomEnd === 100 && zoomStart > 0 && !lineState.inRefreshWIP &&
      !lineState.isCanAutoRefresh) {// get forward data
      this.eventDragChange.next({lineType, isDragBack: false});
    }
  }

  private setLineBaseOption(lineType: LineType): Observable<object> {
    let firstKey = '';
    let secondKey = '';
    switch (lineType) {
      case LineType.ltService: {
        firstKey = 'DASHBOARD.PODS';
        secondKey = 'DASHBOARD.CONTAINERS';
        break;
      }
      case LineType.ltNode: {
        firstKey = 'DASHBOARD.CPU';
        secondKey = 'DASHBOARD.MEMORY';
        break;
      }
      case LineType.ltStorage: {
        firstKey = 'DASHBOARD.USAGE';
        secondKey = 'DASHBOARD.TOTAL';
        break;
      }
    }
    return this.translateService.get([firstKey, secondKey])
      .pipe(map(res => {
        const firstLineTitle = res[firstKey];
        const secondLineTitle = res[secondKey];
        const result = DashboardComponentParent.getBaseOptions();
        result[tooltipKey] = this.getTooltip(firstLineTitle, secondLineTitle, lineType);
        result[seriesKey] = [
          DashboardComponentParent.getBaseSeries(),
          DashboardComponentParent.getBaseSeries(),
          DashboardComponentParent.getBaseSeriesThirdLine()];
        result[seriesKey][0].name = firstLineTitle;
        result[seriesKey][1].name = secondLineTitle;
        result[dataZoomKey][0].start = 80;
        result[dataZoomKey][0].end = 100;
        result[legendKey] = {data: [firstLineTitle, secondLineTitle], x: 'left'};
        return result;
      }));
  }

  scaleChange(lineType: LineType, data: scaleOption) {
    if (!this.getLineInRefreshWIP()) {
      let baseLineTimeStamp = this.getBaseLineTimeStamp(lineType);
      let queryTimeStamp = 0;
      const maxLineTimeStamp = baseLineTimeStamp + data.valueOfSecond * MAX_COUNT_PER_PAGE / 2;
      if (maxLineTimeStamp > this.serverTimeStamp) {
        queryTimeStamp = this.serverTimeStamp;
        baseLineTimeStamp -= maxLineTimeStamp - this.serverTimeStamp;
      } else {
        queryTimeStamp = maxLineTimeStamp;
      }
      this.lineTypeSet.forEach((value: LineType) => {
        const query = this.query.get(value);
        query.scale = data;
        query.time_count = MAX_COUNT_PER_PAGE;
        query.timestamp_base = queryTimeStamp;
        query.baseLineTimeStamp = baseLineTimeStamp;
        const maxTimeStamp = queryTimeStamp;
        const minTimeStamp = queryTimeStamp - query.scale.valueOfSecond * MAX_COUNT_PER_PAGE;
        const thirdLine = this.lineThirdLine.get(value);
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

  public onToolTipEvent(params: object, lineType: LineType) {
    if ((params as Array<any>).length > 1) {
      this.curValue.set(lineType, {
        curFirst: params[0].value[1],
        curFirstUnit: params[0].value[2],
        curSecond: params[1].value[1],
        curSecondUnit: params[1].value[2]
      });
    }
  }

  public onEChartInit(lineType: LineType, eChart: object) {
    this.eChartInstance.set(lineType, eChart);
    this.resetBaseLinePos(lineType);
  }

  @HostListener('window:resize', ['$event'])
  onEChartWindowResize(event: object) {
    this.lineTypeSet.forEach(value => {
      this.resetBaseLinePos(value);
    });
  }

  chartMouseUp(lineType: LineType, event: object) {
    this.afterDragZoomBar(lineType);
  }

  chartDataZoom(lineType: LineType, data: any) {
    this.eventZoomBarChange.next({start: data.start, end: data.end});
  }

  getCurLineName(lineType: LineType): string {
    return this.lineResponses.get(lineType) ? this.lineResponses.get(lineType).curListName : '';
  }

  getLinesName(lineType: LineType): Array<string> {
    return this.lineResponses.get(lineType) ? this.lineResponses.get(lineType).list : null;
  }

  get grafanaViewUrl(): string {
    return `http://${this.appInitService.systemInfo.board_host}/grafana/dashboard/db/kubernetes/`;
  }

  get showGrafanaWindow(): boolean {
    return this.appInitService.isSystemAdmin &&
      !this.appInitService.isMipsSystem &&
      !this.appInitService.isArmSystem;
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
      this.query.get(value).timestamp_base = this.serverTimeStamp;
      this.query.get(value).time_count = MAX_COUNT_PER_PAGE;
      this.getOneLineData(value).subscribe((res: IResponse) => {
        this.lineResponses.set(value, res);
        this.detectChartData(value);
        this.clearEChart(value);
      });
    });
  }
}
