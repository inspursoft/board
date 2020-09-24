import { AfterViewInit, ChangeDetectionStrategy, Component, HostListener, OnDestroy, OnInit } from '@angular/core';
import { Observable, Subject } from 'rxjs';
import { debounceTime, map, tap } from 'rxjs/operators';
import { DashboardComponentParent } from './dashboard.component.parent';
import { DashboardService } from './dashboard.service';
import { TranslateService } from '@ngx-translate/core';
import { MessageService } from '../shared.service/message.service';
import { AppInitService } from '../shared.service/app-init.service';
import { SharedService } from '../shared.service/shared.service';
import { BodyData, LineType, Prometheus, QueryData, RealtimeData, ScaleOption, ThirdLine } from './dashboard.types';

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

@Component({
  templateUrl: './dashboard.component.html',
  styleUrls: ['dashboard.component.css'],
  changeDetection: ChangeDetectionStrategy.Default
})
export class DashboardComponent extends DashboardComponentParent implements OnInit, AfterViewInit, OnDestroy {
  scaleOptions: Array<ScaleOption> = [
    {id: 1, description: 'DASHBOARD.MIN', value: 'second', valueOfSecond: 5},
    {id: 2, description: 'DASHBOARD.HR', value: 'minute', valueOfSecond: 60},
    {id: 3, description: 'DASHBOARD.DAY', value: 'hour', valueOfSecond: 60 * 60},
    {id: 4, description: 'DASHBOARD.MTH', value: 'day', valueOfSecond: 60 * 60 * 24}];
  serverTimeStamp = 0;
  autoRefreshCurInterval = AUTO_REFRESH_CUR_SEED;
  intervalAutoRefresh: any;
  lineOptions: Map<LineType, object>;
  prometheus: Prometheus;
  queryData: QueryData;
  bodyData: BodyData;
  thirdLineData: ThirdLine;
  lineStateInfo: {
    inRefreshWIP: boolean,
    inDrop: boolean,
    isDropBack: boolean,
    isCanAutoRefresh: boolean
  };
  curScaleOption: ScaleOption;
  curValue: Map<LineType, {
    curFirst: number,
    curFirstUnit: string,
    curSecond: number,
    curSecondUnit: string
  }>;
  noData: boolean;
  lineTypeSet: Set<LineType>;
  eventDragChange: Subject<{ lineType: LineType, isDragBack: boolean }>;
  eventZoomBarChange: Subject<{ start: number, end: number }>;
  eChartInstance: Map<LineType, object>;
  autoRefreshInterval: Map<LineType, number>;
  curRealTimeValue: Map<LineType, RealtimeData>;
  noDataErrMsg: Map<LineType, string>;

  constructor(private service: DashboardService,
              private appInitService: AppInitService,
              private messageService: MessageService,
              private translateService: TranslateService,
              private shardService: SharedService) {
    super();
    this.prometheus = new Prometheus();
    this.eventDragChange = new Subject<{ lineType: LineType, isDragBack: boolean }>();
    this.eventZoomBarChange = new Subject<{ start: number, end: number }>();
    this.thirdLineData = new ThirdLine();
    this.queryData = new QueryData();
    this.bodyData = new BodyData();
    this.autoRefreshInterval = new Map<LineType, number>();

    this.curValue = new Map<LineType, {
      curFirst: number,
      curFirstUnit: string,
      curSecond: number,
      curSecondUnit: string
    }>();
    this.curRealTimeValue = new Map<LineType, RealtimeData>();
    this.noDataErrMsg = new Map<LineType, string>();
    this.lineOptions = new Map<LineType, object>();
    this.eChartInstance = new Map<LineType, object>();
    this.lineTypeSet = new Set<LineType>();
  }

  ngOnInit() {
    this.lineTypeSet.add(LineType.ltService);
    this.lineTypeSet.add(LineType.ltNode);
    this.lineTypeSet.add(LineType.ltStorage);
    this.curScaleOption = this.scaleOptions[0];
    this.lineStateInfo = {
      isCanAutoRefresh: true,
      isDropBack: false,
      inDrop: false,
      inRefreshWIP: false
    };
    this.eventDragChange.asObservable().pipe(debounceTime(300)).subscribe(dragInfo => {
      this.refreshLineDataByDrag(dragInfo.isDragBack);
    });
    this.eventZoomBarChange.asObservable().pipe(debounceTime(300)).subscribe((zoom: { start: number, end: number }) => {
      this.lineTypeSet.forEach((value) => {
        this.lineOptions.get(value)[dataZoomKey][0].start = zoom.start;
        this.lineOptions.get(value)[dataZoomKey][0].end = zoom.end;
        this.resetBaseLinePos(value);
      });
    });
  }

  ngOnDestroy() {
    this.lineTypeSet.forEach((value) => {// for update at after destroy
      this.eChartInstance.set(value, null);
    });
    clearInterval(this.intervalAutoRefresh);
  }

  ngAfterViewInit() {
    this.initAsyncLines();
    this.intervalAutoRefresh = setInterval(() => {
      this.autoRefreshCurDada();
    }, 5000);
  }

  setLineZoomByTimeStamp(lineType: LineType, lineTimeStamp: number): void {
    const lineOption = this.lineOptions.get(lineType);
    const lineZoomStart = lineOption[dataZoomKey][0].start;
    const lineZoomEnd = lineOption[dataZoomKey][0].end;
    const lineZoomHalf: number = (lineZoomEnd - lineZoomStart) / 2;
    const maxDate: Date = this.thirdLineData.maxDate;
    const minDate: Date = this.thirdLineData.minDate;
    const maxTimeStrap = Math.round(maxDate.getTime() / 1000);
    const minTimeStrap = Math.round(minDate.getTime() / 1000);
    const percent = ((maxTimeStrap - lineTimeStamp) / (maxTimeStrap - minTimeStrap)) * 100;
    lineOption[dataZoomKey][0].start = Math.max(percent - lineZoomHalf, 1);
    lineOption[dataZoomKey][0].end = Math.min(lineOption[dataZoomKey][0].start + 2 * lineZoomHalf, 99);
  }

  getLineData(): Observable<Prometheus> {
    this.lineStateInfo.inRefreshWIP = true;
    return this.service.getLineData(this.queryData, this.bodyData)
      .pipe(tap(() => {
          this.noData = false;
          this.lineStateInfo.inRefreshWIP = false;
        }, () => {
          this.lineStateInfo.inRefreshWIP = false;
          this.noData = true;
        })
      );
  }

  detectChartData() {
    const minTimeStrap = this.bodyData.queryTimestamp - MAX_COUNT_PER_PAGE * this.curScaleOption.valueOfSecond;
    const maxTimeStrap = this.bodyData.queryTimestamp;
    this.thirdLineData.maxDate = new Date(maxTimeStrap * 1000);
    this.thirdLineData.minDate = new Date(minTimeStrap * 1000);
    this.lineTypeSet.forEach(lineType => {
      const lineSeries = this.lineOptions.get(lineType);
      const newLineOption = Object.create({});
      this.lineOptions.delete(lineType);
      switch (lineType) {
        case LineType.ltService: {
          lineSeries[seriesKey][0].data = this.prometheus.serviceLineData.firstLineData;
          lineSeries[seriesKey][1].data = this.prometheus.serviceLineData.secondLineData;
          lineSeries[seriesKey][2].data = this.thirdLineData.values;
          break;
        }
        case LineType.ltNode: {
          lineSeries[seriesKey][0].data = this.prometheus.nodeLineData.firstLineData;
          lineSeries[seriesKey][1].data = this.prometheus.nodeLineData.secondLineData;
          lineSeries[seriesKey][2].data = this.thirdLineData.values;
          break;
        }
        case LineType.ltStorage: {
          lineSeries[seriesKey][0].data = this.prometheus.storageLineData.firstLineData;
          lineSeries[seriesKey][1].data = this.prometheus.storageLineData.secondLineData;
          lineSeries[seriesKey][2].data = this.thirdLineData.values;
          break;
        }
      }
      Object.assign(newLineOption, lineSeries);
      this.lineOptions.set(lineType, newLineOption);
    });
  }

  initThirdLineDate() {
    const maxTimeStamp = this.bodyData.queryTimestamp;
    const minTimeStamp = maxTimeStamp - this.bodyData.queryCount * this.curScaleOption.valueOfSecond;
    this.thirdLineData = new ThirdLine();
    this.thirdLineData.maxDate = new Date(maxTimeStamp * 1000);
    this.thirdLineData.minDate = new Date(minTimeStamp * 1000);
  }

  initAsyncLines() {
    this.service.getServerTimeStamp().subscribe((res: number) => {
      this.serverTimeStamp = res;
      this.queryData.nodeName = 'average';
      this.queryData.serviceName = 'total';
      this.bodyData.queryCount = MAX_COUNT_PER_PAGE;
      this.bodyData.queryTimeUnit = this.curScaleOption.value;
      this.bodyData.queryTimestamp = this.serverTimeStamp;
      this.initThirdLineDate();
      this.lineTypeSet.forEach((lineType: LineType) => {
        this.curValue.set(lineType, {curFirst: 0, curFirstUnit: '', curSecond: 0, curSecondUnit: ''});
        this.curRealTimeValue.set(lineType, {curFirst: 0, curFirstUnit: '', curSecond: 0, curSecondUnit: ''});
        this.autoRefreshInterval.set(lineType, this.curScaleOption.valueOfSecond);
        this.setLineBaseOption(lineType).pipe(tap(option => this.lineOptions.set(lineType, option))).subscribe(() => {
          if (lineType === LineType.ltStorage) {
            this.getLineData().subscribe((prometheus1: Prometheus) => {
              this.prometheus = prometheus1;
              this.detectChartData();
              this.curRealTimeValue.set(LineType.ltService, this.prometheus.serviceRealtimeData);
              this.curRealTimeValue.set(LineType.ltNode, this.prometheus.nodeRealtimeData);
              this.curRealTimeValue.set(LineType.ltStorage, this.prometheus.storageRealtimeData);
            });
          }
        });
      });
    });
  }


  getBaseLineTimeStamp(lineType: LineType): number {
    const option = this.lineOptions.get(lineType);
    const start = option[dataZoomKey][0].start / 100;
    const end = option[dataZoomKey][0].end / 100;
    const middlePos = (start + end) / 2;
    const maxDate = this.thirdLineData.maxDate;
    const minDate = this.thirdLineData.minDate;
    const maxTimeStamp = Math.round(maxDate.getTime() / 1000);
    const minTimeStamp = Math.round(minDate.getTime() / 1000);
    const screenMaxTimeStamp = maxTimeStamp - (maxTimeStamp - minTimeStamp) * (1 - end);
    const screenTimeStamp = (maxTimeStamp - minTimeStamp) * (end - start);
    return Math.round(screenMaxTimeStamp - screenTimeStamp * (1 - middlePos));
  }

  setLineZoomByCount(isDragBack: boolean): void {
    this.lineTypeSet.forEach(lineType => {
      const lineData = this.prometheus.getResponseLineData(lineType);
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
    });
  }

  clearEChart(): void {
    this.lineTypeSet.forEach(lineType => {
      const eChart = this.eChartInstance.get(lineType);
      if (eChart && eChart[clearKey]) {
        eChart[clearKey]();
      }
    });
  }

  resetBaseLinePos(lineType: LineType) {
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

  autoRefreshCurDada() {
    this.autoRefreshCurInterval--;
    if (this.autoRefreshCurInterval === 0) {
      this.autoRefreshCurInterval = AUTO_REFRESH_CUR_SEED;
      this.service.getServerTimeStamp().subscribe((res: number) => {
        this.serverTimeStamp = res;
        const bodyData = new BodyData();
        bodyData.queryTimestamp = this.serverTimeStamp;
        bodyData.queryCount = 2;
        bodyData.queryTimeUnit = 'second';
        this.service.getLineData(this.queryData, bodyData).subscribe((prometheus1: Prometheus) => {
          this.curRealTimeValue.set(LineType.ltService, prometheus1.serviceRealtimeData);
          this.curRealTimeValue.set(LineType.ltNode, prometheus1.nodeRealtimeData);
          this.curRealTimeValue.set(LineType.ltStorage, prometheus1.storageRealtimeData);
        });
      });
    }
  }

  updateAfterDragTimeStamp(isDropBack: boolean): void {
    const maxDate: Date = this.thirdLineData.maxDate;
    const minDate: Date = this.thirdLineData.minDate;
    const minTimeStrap = Math.round(minDate.getTime() / 1000);
    const maxTimeStrap = Math.round(maxDate.getTime() / 1000);
    let newMaxTimeStrap = 0;
    let newMinTimeStrap = 0;
    if (isDropBack) {
      newMaxTimeStrap = maxTimeStrap - MAX_COUNT_PER_DRAG * this.curScaleOption.valueOfSecond;
      newMinTimeStrap = newMaxTimeStrap - MAX_COUNT_PER_PAGE * this.curScaleOption.valueOfSecond;
      this.bodyData.queryTimestamp = newMaxTimeStrap;
    } else {
      newMaxTimeStrap = Math.min(maxTimeStrap + MAX_COUNT_PER_DRAG * this.curScaleOption.valueOfSecond, this.serverTimeStamp);
      newMinTimeStrap = newMaxTimeStrap - MAX_COUNT_PER_PAGE * this.curScaleOption.valueOfSecond;
      this.bodyData.queryTimestamp = newMaxTimeStrap;
    }
    this.thirdLineData.maxDate = new Date(newMaxTimeStrap * 1000);
    this.thirdLineData.minDate = new Date(newMinTimeStrap * 1000);
  }

  private resetAfterDragTimeStamp(): void {
    const maxDate: Date = this.thirdLineData.maxDate;
    const minDate: Date = this.thirdLineData.minDate;
    const maxTimeStrap = Math.round(maxDate.getTime() / 1000);
    const minTimeStrap = Math.round(minDate.getTime() / 1000);
    const newMaxTimeStrap = maxTimeStrap + MAX_COUNT_PER_DRAG * this.curScaleOption.valueOfSecond;
    const newMinTimeStrap = newMaxTimeStrap - MAX_COUNT_PER_PAGE * this.curScaleOption.valueOfSecond;
    this.thirdLineData.maxDate = new Date(newMaxTimeStrap * 1000);
    this.thirdLineData.minDate = new Date(newMinTimeStrap * 1000);
  }

  refreshLineDataByDrag(isDragBack: boolean): void {
    this.lineStateInfo.inDrop = true;
    this.lineStateInfo.isDropBack = isDragBack;
    this.lineStateInfo.isCanAutoRefresh = false;
    this.lineStateInfo.inRefreshWIP = true;
    this.updateAfterDragTimeStamp(isDragBack);
    if (isDragBack) {
      this.getLineData().subscribe((prometheus1: Prometheus) => {
        this.lineStateInfo.inRefreshWIP = false;
        if (!this.prometheus.isOverMinLimit) {
          this.prometheus = prometheus1;
          this.setLineZoomByCount(true);
          this.lineTypeSet.forEach(lineType => this.resetBaseLinePos(lineType));
          this.detectChartData();
          this.clearEChart();
        } else {
          this.resetAfterDragTimeStamp();
        }
      });
    } else {
      this.getLineData().subscribe((prometheus1: Prometheus) => {
        this.lineStateInfo.inRefreshWIP = false;
        if (!this.prometheus.isOverMaxLimit) {
          this.prometheus = prometheus1;
          this.setLineZoomByCount(false);
          this.lineTypeSet.forEach(lineType => this.resetBaseLinePos(lineType));
          this.detectChartData();
          this.clearEChart();
        }
      });
    }
  }

  private afterDragZoomBar(lineType: LineType) {
    const zoomStart = this.lineOptions.get(lineType)[dataZoomKey][0].start;
    const zoomEnd = this.lineOptions.get(lineType)[dataZoomKey][0].end;
    if (zoomStart === 0 && zoomEnd < 100 && !this.lineStateInfo.inRefreshWIP) {// get backup data
      this.eventDragChange.next({lineType, isDragBack: true});
    } else if (zoomEnd === 100 && zoomStart > 0 && !this.lineStateInfo.inRefreshWIP &&
      !this.lineStateInfo.isCanAutoRefresh) {// get forward data
      this.eventDragChange.next({lineType, isDragBack: false});
    }
  }

  setLineBaseOption(lineType: LineType): Observable<object> {
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
        })
      );
  }

  scaleChange(lineType: LineType, data: ScaleOption) {
    if (!this.lineStateInfo.inRefreshWIP) {
      this.curScaleOption = data;
      let baseLineTimeStamp = this.getBaseLineTimeStamp(lineType);
      let queryTimeStamp = 0;
      const maxLineTimeStamp = baseLineTimeStamp + this.curScaleOption.valueOfSecond * MAX_COUNT_PER_PAGE / 2;
      if (maxLineTimeStamp > this.serverTimeStamp) {
        queryTimeStamp = this.serverTimeStamp;
        baseLineTimeStamp -= maxLineTimeStamp - this.serverTimeStamp;
      } else {
        queryTimeStamp = maxLineTimeStamp;
      }
      this.bodyData.queryCount = MAX_COUNT_PER_PAGE;
      this.bodyData.queryTimestamp = queryTimeStamp;
      this.bodyData.queryTimeUnit = data.value;
      const maxTimeStamp = queryTimeStamp;
      const minTimeStamp = queryTimeStamp - data.valueOfSecond * MAX_COUNT_PER_PAGE;
      this.thirdLineData.minDate = new Date(minTimeStamp * 1000);
      this.thirdLineData.maxDate = new Date(maxTimeStamp * 1000);
      this.getLineData().subscribe((res: Prometheus) => {
        this.prometheus = res;
        this.lineTypeSet.forEach(lineType1 => {
          this.setLineZoomByTimeStamp(lineType1, baseLineTimeStamp);
          this.resetBaseLinePos(lineType1);
        });
        this.detectChartData();
        this.clearEChart();
      });
    }
  }

  serviceDropDownChange(serviceName: string) {
    if (!this.lineStateInfo.inRefreshWIP) {
      this.queryData.serviceName = serviceName;
      this.getLineData().subscribe((prometheus1: Prometheus) => {
        this.prometheus = prometheus1;
        this.detectChartData();
        this.clearEChart();
      });
    }
  }

  nodeDropDownChange(nodeName: string) {
    if (!this.lineStateInfo.inRefreshWIP) {
      this.queryData.nodeName = nodeName;
      this.getLineData().subscribe((prometheus1: Prometheus) => {
        this.prometheus = prometheus1;
        this.detectChartData();
        this.clearEChart();
      });
    }
  }

  onToolTipEvent(params: object, lineType: LineType) {
    if ((params as Array<any>).length > 1) {
      this.curValue.set(lineType, {
        curFirst: params[0].value[1],
        curFirstUnit: params[0].value[2],
        curSecond: params[1].value[1],
        curSecondUnit: params[1].value[2]
      });
    }
  }

  onEChartInit(lineType: LineType, eChart: object) {
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

  get showGrafanaWindow(): boolean {
    return this.appInitService.isSystemAdmin &&
      !this.appInitService.isMipsSystem &&
      !this.appInitService.isArmSystem &&
      this.appInitService.isNormalMode;
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
    this.bodyData.queryTimestamp = this.serverTimeStamp;
    this.bodyData.queryCount = MAX_COUNT_PER_PAGE;
    this.getLineData().subscribe((prometheus1: Prometheus) => {
      this.prometheus = prometheus1;
      this.detectChartData();
      this.clearEChart();
    });
  }
}
