import { Injectable } from '@angular/core';
import { HttpClient, HttpResponse } from "@angular/common/http"
import { Observable } from "rxjs/Observable";
import { InvalidServiceName } from "../shared/shared.const";

export enum LineType {ltService, ltNode, ltStorage}

export interface ILineListQueryModel {
  readonly query_name_Key: string;
  readonly data_list_filed_key: string;
  readonly data_list_name_key: string;
  readonly data_list_time_stamp_key: string;
  readonly data_list_cur_name: string;
  readonly data_filed_key: string;
  readonly data_url_key: string;
  readonly data_time_stamp_key: string;
  readonly data_first_line_key: string;
  readonly data_second_line_key: string;
}

export interface IQuery {
  time_count: number,
  time_unit: string,
  list_name: string,
  timestamp_base: number
}

export interface IResponse {
  list: Array<string>,
  firstLineData: Array<[Date,number]>,
  secondLineData: Array<[Date,number]>,
  curListName: string,
  limit: {isMax: boolean, isMin: boolean}
}

class ResponseLineData implements IResponse {
  public list: Array<string>;
  public firstLineData: Array<[Date,number]>;
  public secondLineData: Array<[Date,number]>;
  public curListName: string;
  public limit: {isMax: boolean, isMin: boolean};

  constructor() {
    this.list = Array<string>();
    this.firstLineData = Array<[Date,number]>();
    this.secondLineData = Array<[Date,number]>();
    this.limit = {isMax: false, isMin: false};
  }
}

const BASE_URL = "/api/v1";
const ARR_SIZE_UNIT: Array<string> = ["B", "KB", "MB", "GB", "TB"];

@Injectable()
export class DashboardService {
  private lineNameMap: Map<LineType, ILineListQueryModel>;
  StorageUnit: string;

  constructor(private http: HttpClient) {
    this.lineNameMap = new Map<LineType, ILineListQueryModel>();
    this.lineNameMap.set(LineType.ltService, {
      query_name_Key: "service_name",
      data_list_filed_key: "service_list_data",
      data_list_name_key: "service_name",
      data_list_time_stamp_key: "timestamp",
      data_list_cur_name: "service_name",
      data_filed_key: "service_logs_data",
      data_url_key: "service",
      data_time_stamp_key: "timestamp",
      data_first_line_key: "pod_number",
      data_second_line_key: "container_number"
    });
    this.lineNameMap.set(LineType.ltNode, {
      query_name_Key: "node_name",
      data_list_filed_key: "node_list_data",
      data_list_name_key: "node_name",
      data_list_time_stamp_key: "﻿timestamp",
      data_list_cur_name: "node_name",
      data_filed_key: "node_logs_data",
      data_url_key: "node",
      data_time_stamp_key: "timestamp",
      data_first_line_key: "cpu_usage",
      data_second_line_key: "memory_usage"
    });
    this.lineNameMap.set(LineType.ltStorage, {
      query_name_Key: "node_name",
      data_list_filed_key: "node_list_data",
      data_list_name_key: "node_name",
      data_list_time_stamp_key: "﻿timestamp",
      data_list_cur_name: "node_name",
      data_filed_key: "node_logs_data",
      data_url_key: "node",
      data_time_stamp_key: "timestamp",
      data_first_line_key: "storage_use",
      data_second_line_key: "storage_total"
    });
  };

  private getUnitMultipleValue(sample: number): number {
    let result: number = 1;
    let nameIndex: number = 0;
    while (sample > 1024) {
      sample = sample / 1024;
      nameIndex += 1;
      result *= 1024;
    }
    this.StorageUnit = ARR_SIZE_UNIT[nameIndex];
    return result;
  }

  get CurStorageUnit(): string {
    return this.StorageUnit;
  }

  getQueryBody(lineType: LineType, query: IQuery): Object {
    if (lineType == LineType.ltService) {
      return {
        service_time_unit: query.time_unit,
        service_time_count: query.time_count,
        service_timestamp: query.timestamp_base
      };
    } else {
      return {
        node_time_unit: query.time_unit,
        node_time_count: query.time_count,
        node_timestamp: query.timestamp_base
      };
    }
  }

  getServerTimeStamp(): Observable<number> {
    return this.http.get(`${BASE_URL}/dashboard/time`, {observe: "response"})
      .map((res:HttpResponse<Object>)=>res.body["time_now"]);
  }

  getLineData(lineType: LineType, query: IQuery): Observable<IResponse> {
    let lineListQuery: ILineListQueryModel = this.lineNameMap.get(lineType);
    let requestParams = {};
    let body = this.getQueryBody(lineType, query);
    requestParams[lineListQuery.query_name_Key] = query.list_name;
    return this.http.post(`${BASE_URL}/dashboard/${lineListQuery.data_url_key}`, body, {
      observe: "response",
      params: requestParams
    }).map((obs: HttpResponse<Object>) => {
      let result = new ResponseLineData();
      let time_key = lineListQuery.data_time_stamp_key;
      let first_key = lineListQuery.data_first_line_key;
      let second_key = lineListQuery.data_second_line_key;
      let dataLogs: Array<Object> = obs.body[lineListQuery.data_filed_key];
      if (dataLogs && dataLogs.length > 0) {//for data
        let multiple: number = 1;
        if (lineType == LineType.ltStorage) {
          multiple = this.getUnitMultipleValue(dataLogs[0][first_key]);
        }
        dataLogs.forEach((item: Object) => {
          let timeValue = new Date(item[time_key] * 1000);
          let firstValue = Math.round(item[first_key] / multiple * 100) / 100;
          let secondValue = Math.round(item[second_key] / multiple * 100) / 100;
          result.firstLineData.push([timeValue,firstValue]);
          result.secondLineData.push([timeValue,secondValue]);
        });
      }
      result.curListName = obs.body[lineListQuery.data_list_cur_name];
      result.limit.isMax = obs.body["is_over_max_limit"];
      result.limit.isMin = obs.body["is_over_min_limit"];
      result.list.push(lineType == LineType.ltService ? "total" :"average");
      let listLogs: Array<Object> = obs.body[lineListQuery.data_list_filed_key];
      if (listLogs && listLogs.length > 0) {//for list
        listLogs.sort((a: Object, b: Object) => a[lineListQuery.data_list_name_key] > (b[lineListQuery.data_list_name_key]) ? 1 : -1);
        listLogs.forEach((item: Object) => {
          const serviceName: string = item[lineListQuery.data_list_name_key];
          result.list.push(serviceName);
          // if (InvalidServiceName.indexOf(serviceName.toLocaleLowerCase()) === -1){
          //   result.list.push(serviceName);
          // }
        });
      }
      return result;
    });
  }
}
