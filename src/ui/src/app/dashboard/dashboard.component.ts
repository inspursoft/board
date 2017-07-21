import { OnInit, AfterViewInit, Component, OnDestroy } from '@angular/core';
import { Assist } from "./dashboard-assist"
import { scaleOption } from "app/dashboard/time-range-scale.component/time-range-scale.component";
import { DashboardService, ServiceListModel, LinesData } from "app/dashboard/dashboard.service";
import { TranslateService } from "@ngx-translate/core";
import { Subscription } from "rxjs/Subscription";

@Component({
  selector: 'dashboard',
  templateUrl: 'dashboard.component.html',
  styleUrls: ['dashboard.component.css']
})

export class DashboardComponent implements OnInit, AfterViewInit, OnDestroy {
  _intervalRead: any;
  _onLangChangeSubscription: Subscription;
  scaleOptions: Array<scaleOption> = [
    {"id": 1, "description": "DASHBOARD.MIN", "value": "second"},
    {"id": 2, "description": "DASHBOARD.HR", "value": "minute"},
    {"id": 3, "description": "DASHBOARD.DAY", "value": "hour"},
    {"id": 4, "description": "DASHBOARD.MTH", "value": "day"}];

  podCount: number = 0;
  containerCount: number = 0;

  memoryPercent: string = '70%';
  cpuPercent: string = '40%';
  usageVolume: string = '3T';
  totalVolume: string = '10T';

  _serviceIntervalSeed: number = 10;
  _serviceQuery: {model: ServiceListModel, scale: scaleOption, count: number};
  _serviceOptionsBuffer: {lastZoomStart: number, lastZoomEnd: number};
  serviceBtnValue: string;
  serviceList: Array<ServiceListModel>;
  serviceOptions: object = {};
  serviceAlready: boolean = false;
  serviceNoData: boolean = false;
  serviceData: LinesData;

  nodeBtnValue: string;
  nodeOptions: object = {};

  storageBtnValue: string;
  storageOptions: object = {};

  constructor(private service: DashboardService,
              private translateService: TranslateService) {
  }

  serviceScaleChange(data: scaleOption) {
    if (this.serviceAlready) {
      this._serviceIntervalSeed = 10;
      this._serviceQuery.scale = data;
      this.serviceAlready = false;
      this.serviceNoData = false;
      this.refreshServiceData();
    }
  }

  serviceDropDownChange(service: ServiceListModel) {
    if (this.serviceAlready) {
      this._serviceIntervalSeed = 10;
      this.serviceBtnValue = service.service_name;
      this.serviceAlready = false;
      this.serviceNoData = false;
      this._serviceQuery.model = service;
      this.refreshServiceData();
    }
  }

  setServiceOption() {
    this.translateService.get(["DASHBOARD.CONTAINERS", "DASHBOARD.PODS"])
      .subscribe(res => {
        let podsTranslate: string = res["DASHBOARD.PODS"];
        let containersTranslate: string = res["DASHBOARD.CONTAINERS"];
        this.serviceOptions = Assist.getServiceOptions();
        this.serviceOptions["tooltip"] = Assist.getTooltip(podsTranslate, containersTranslate);
        this.serviceOptions["series"] = [Assist.getBaseSeries(), Assist.getBaseSeries()];
        this.serviceOptions["series"][0]["name"] = podsTranslate;
        this.serviceOptions["series"][1]["name"] = containersTranslate;
        this.serviceOptions["legend"] = {data: [podsTranslate, containersTranslate], x: "left"};
      });
  }

  serviceChartDataZoom(event: object) {
    this._serviceOptionsBuffer.lastZoomStart = event["start"];
    this._serviceOptionsBuffer.lastZoomEnd = event["end"];
  }

  refreshServiceData() {
    this.service.getServiceData(this._serviceQuery.count, this._serviceQuery.scale.value, this._serviceQuery.model.service_name)
      .then(res => {
        this.serviceOptions["dataZoom"][0]["start"] = this._serviceOptionsBuffer.lastZoomStart;
        this.serviceOptions["dataZoom"][0]["end"] = this._serviceOptionsBuffer.lastZoomEnd;
        this.serviceData = res;
        this._serviceIntervalSeed = 10;
        this.serviceNoData = false;
        this.serviceAlready = true;
        if (this.serviceData[0] && this.serviceData[0].length > 0) {
          this.podCount = this.serviceData[0][this.serviceData[0].length - 1][1] | 0;
          this.containerCount = this.serviceData[1][this.serviceData[1].length - 1][1] | 0;
        }
      })
      .catch(() => {
        this.serviceData = [[], []];
        this.podCount = 0;
        this.containerCount = 0;
        this._serviceIntervalSeed = 10;
        this.serviceNoData = true;
        this.serviceAlready = true;
      });
  }

  ngOnInit() {
    this._onLangChangeSubscription = this.translateService.onLangChange.subscribe(() => {
      this.setServiceOption();
    });
  }

  ngOnDestroy() {
    clearInterval(this._intervalRead);
    if (this._onLangChangeSubscription) {
      this._onLangChangeSubscription.unsubscribe();
    }
  }

  ngAfterViewInit() {
    this.setServiceOption();
    this._serviceOptionsBuffer = Object.create({
      lastZoomStart: 100,
      lastZoomEnd: 80
    });
    this.service.getServiceList()
      .then(res => {
        this.serviceList = res;
        this._intervalRead = setInterval(() => {
          if (this._serviceIntervalSeed > 0 && this.serviceAlready) {
            this._serviceIntervalSeed--;
            if (this._serviceIntervalSeed == 0) {
              this.refreshServiceData();
            }
          }
        }, 1000);
        this.serviceBtnValue = this.serviceList[0].service_name;
        this.nodeBtnValue = this.serviceList[0].service_name;
        this.storageBtnValue = this.serviceList[0].service_name;
        this._serviceQuery = Object.create({
          count: 300,
          model: this.serviceList[0],
          scale: this.scaleOptions[0]
        });
        this.serviceAlready = false;
        this.refreshServiceData();
      })
      .catch(() => null);

    let serviceSimulateData = DashboardService.getBySimulateData(0, 1);

    this.nodeOptions = Assist.getBaseOptions();
    this.nodeOptions["tooltip"] = Assist.getTooltip("CPU", "Memory");
    this.nodeOptions["series"] = [Assist.getBaseSeries(), Assist.getBaseSeries()];
    this.nodeOptions["series"][0]["data"] = serviceSimulateData[0];
    this.nodeOptions["series"][1]["data"] = serviceSimulateData[1];

    this.storageOptions = Assist.getBaseOptions();
    this.storageOptions["tooltip"] = Assist.getTooltip("", "Total");
    this.storageOptions["series"] = [Assist.getBaseSeries(), Assist.getBaseSeries()];
    this.storageOptions["series"][0]["data"] = serviceSimulateData[0];
    this.storageOptions["series"][1]["data"] = serviceSimulateData[1];
  }
}