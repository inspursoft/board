import { Injectable } from '@angular/core';
import { Http, RequestOptions, Headers, Response } from "@angular/http"
import { MessageService } from "app/shared/message-service/message.service"
import "rxjs/add/operator/map";
import 'rxjs/add/operator/toPromise';
import 'rxjs/add/operator/timeout';

export type LineDataModel = [Date, number];
export type LinesData = [LineDataModel[], LineDataModel[]];

export interface ServiceListModel {
  readonly service_name: string;
}

export interface ServiceDataModel {
  readonly podcontainer_timestamp: number;
  readonly pods_number: number;
  readonly container_number: number;
}

export interface NodeDataModel {
  readonly date: Date;
  readonly value: number;
}

export interface StorageDataModel {
  readonly date: Date;
  readonly value: number;
}

const BASE_URL = "/api/v1";
@Injectable()
export class DashboardService {

  static baseDate: Date = new Date();

  static getOneStepTime(dateScaleId: number): number {
    switch (dateScaleId) {
      case 1:
        return 5 * 1000;
      case 2:
        return 60 * 1000;
      case 3:
        return 24 * 60 * 1000;
      case 4:
        return 12 * 24 * 60 * 1000;

      default:
        return 1;
    }
  }

  private static getSimulateData(serviceID: number): number {
    switch (serviceID) {
      case 0:
        return 130 + Math.round(Math.random() * 50);
      case 1:
        return 30 + Math.round(Math.random() * 10);
      case 2:
        return 20 + Math.round(Math.random() * 10);
      case 3:
        return 50 + Math.round(Math.random() * 10);
      case 4:
        return 30 + Math.round(Math.random() * 20);
    }
  }

  private static  getSimulateDate(dateScaleId: number): Date {
    DashboardService.baseDate.setTime(DashboardService.baseDate.getTime()
      + DashboardService.getOneStepTime(dateScaleId));
    return new Date(DashboardService.baseDate.getTime());
  }

  /**
   *getServiceData
   * @param serviceID
   * @param dateScaleId
   *node:dateScaleId=>1:1min;2:1hr;3:1day;4:1mth
   */
  static getBySimulateData(serviceID: number, dateScaleId: number): Map<number, LineDataModel[]> {
    if (!dateScaleId || dateScaleId < 1 || dateScaleId > 4) return null;
    let r: Map<number, LineDataModel[]> = new Map<number, LineDataModel[]>();
    r[0] = new Array<LineDataModel>(0);
    r[1] = new Array<LineDataModel>(0);
    for (let i = 0; i < 11; i++) {
      let date: Date = DashboardService.getSimulateDate(dateScaleId);
      let arrBuf1 = [date, DashboardService.getSimulateData(serviceID)];
      let arrBuf2 = [date, DashboardService.getSimulateData(serviceID)];
      r[0].push(arrBuf1);
      r[1].push(arrBuf2);
    }
    return r;
  }

  constructor(private http: Http,
              private messageService: MessageService) {
  };

  readonly defaultHeaders: Headers = new Headers({
    contentType: "application/json"
  });

  static getErrorMsg(reason: Response | Error, statusArr: Array<number>, errorKey: string): string {
    if (reason instanceof Response) {
      return statusArr.indexOf(reason.status) > -1 ?
        `DASHBOARD.${errorKey}_ERR_${reason.status}` :
        `${reason.status}:${reason.statusText}`;
    }
    else {
      return `${reason.name}:${reason.message}`;
    }
  }

  getServiceList(): Promise<ServiceListModel[]> {
    let options = new RequestOptions({headers: this.defaultHeaders});
    return this.http.get(`${BASE_URL}/dashboard/service/list`, options)
      .toPromise()
      .then(res => {
        let arr = Array.from(res.json()).sort((a, b) => {
          return a["service_name"] == b["service_name"] ? 0 :
            a["service_name"] > b["service_name"] ? 1 : -1;
        });
        arr.unshift({service_name: "total"});//add total service
        return arr;
      })
      .catch(reason => {
        let errMsg: string = DashboardService.getErrorMsg(reason, Array.from([404, 409, 504]), "GET");
        return Promise.reject(errMsg);
      });
  };

  /**data origin
   * {"service_name": "mysql-read",
   * "service_timeunit": "second",
   * "service_count": "11",
   * "service_statuslogs": [{
   * "pods_number": 2,
   * "container_number": 4,
   * "time_stamp":1499842237
   * }]}
   */
  getServiceData(query: {time_count: number, time_unit: string, service_name: string, timestamp_base: number}): Promise<LinesData> {
    let params: Map<string, string> = new Map<string, string>();
    params["service_name"] = query.service_name;
    let options = new RequestOptions({
      headers: this.defaultHeaders,
      search: params
    });
    return this.http.post(`${BASE_URL}/dashboard/service`, {
      time_count: query.time_count.toString(),
      time_unit: query.time_unit,
      timestamp_base: query.timestamp_base.toString()
    }, options)
      .toPromise()
      .then((res: Response) => {
        let resJson: object = res.json();
        let logs: ServiceDataModel[] = resJson["service_statuslogs"];
        let r: LinesData = [[[new Date(), 1]], [[new Date(), 1]]];
        if (logs && logs.length > 0) {
          r[0] = r[0].slice(0,0);
          r[1] = r[1].slice(0,0);
          logs.forEach((item: ServiceDataModel) => {
            r[0].push([new Date(item.podcontainer_timestamp * 1000), item.pods_number]);
            r[1].push([new Date(item.podcontainer_timestamp * 1000), item.container_number]);
          });
        }
        return r;
      })
      .catch(reason => {
        let errMsg: string = DashboardService.getErrorMsg(reason, Array.from([404, 409, 504]), "GET");
        return Promise.reject(errMsg);
      });
  }
}
