import { Injectable } from '@angular/core';
import { Http, RequestOptions, Headers, Response } from "@angular/http"
import { AppInitService } from "../app.init.service";
import "rxjs/add/operator/map";
import 'rxjs/add/operator/toPromise';

export enum LineType {ltService, ltNode, ltStorage}
export type LineDataModel = [Date, number];
export type LinesData = [LineDataModel[], LineDataModel[], LineDataModel[]];

export interface LineListQueryModel {
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

export interface LineListDataModel {
  list_name: string;
  time_stamp?: number;
}
const BASE_URL = "/api/v1";
const ARR_SIZE_UNIT: Array<string> = ["B", "KB", "MB", "GB", "TB"];
@Injectable()
export class DashboardService {
  LineNameMap: Map<LineType, LineListQueryModel>;
  StorageUnit: string;

  constructor(private http: Http,
              private appInitService: AppInitService) {
    this.LineNameMap = new Map<LineType, LineListQueryModel>();
    this.LineNameMap.set(LineType.ltService, {
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
    this.LineNameMap.set(LineType.ltNode, {
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
    this.LineNameMap.set(LineType.ltStorage, {
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

  get defaultHeader(): Headers {
    let headers = new Headers();
    headers.append('Content-Type', 'application/json');
    headers.append('token', this.appInitService.token);
    return headers;
  }

  getServerTimeStamp(): Promise<number> {
    let options = new RequestOptions({
      headers: this.defaultHeader
    });
    return this.http.get(`${BASE_URL}/dashboard/time`, options).toPromise()
      .then((res: Response) => {
        this.appInitService.chainResponse(res);
        return res.json()["time_now"];
      }).catch(err => Promise.reject(err));
  }

  getLineData(lineType: LineType, query: {
    time_count: number,
    time_unit: string,
    list_name: string,
    timestamp_base: number,
    service_duration_time?: number
  }): Promise<{
    List: Array<LineListDataModel>,
    Data: LinesData,
    CurListName: string,
    Limit: {isMax: boolean, isMin: boolean}
  }> {
    let lineKey = this.LineNameMap.get(lineType).query_name_Key;
    let requestParams = {};
    requestParams[lineKey] = query.list_name;
    let options = new RequestOptions({
      headers: this.defaultHeader,
      params: requestParams
    });
    let body: Object;
    switch (lineType) {
      case LineType.ltService: {
        body = {
          service_time_unit: query.time_unit,
          service_time_count: query.time_count,
          service_timestamp: query.timestamp_base
        };
        break;
      }
      case LineType.ltStorage:
      case LineType.ltNode: {
        body = {
          node_time_unit: query.time_unit,
          node_time_count: query.time_count,
          node_timestamp: query.timestamp_base
        };
        break;
      }
    }
    return this.http.post(`${BASE_URL}/dashboard/${this.LineNameMap.get(lineType).data_url_key}`, body, options)
      .toPromise()
      .then((res: Response) => {
        this.appInitService.chainResponse(res);
        let lineNameMap = this.LineNameMap.get(lineType);
        let time_key = lineNameMap.data_time_stamp_key;
        let first_key = lineNameMap.data_first_line_key;
        let second_key = lineNameMap.data_second_line_key;
        let resJson: Object = res.json();
        let dataLogs: Array<Object> = resJson[lineNameMap.data_filed_key];
        let dataListLogs: Array<Object> = resJson[lineNameMap.data_list_filed_key];
        let resultData: LinesData = [Array<[Date, number]>(0), Array<[Date, number]>(0),Array<[Date, number]>(0)];
        let resultList: Array<LineListDataModel> = new Array<LineListDataModel>(0);
        if (dataLogs && dataLogs.length > 0) {//for data
          let multiple: number = 1;
          if (lineType == LineType.ltStorage) {
            multiple = this.getUnitMultipleValue(dataLogs[0][first_key]);
          }
          dataLogs.forEach((item: Object) => {
            resultData[0].push([new Date(item[time_key] * 1000), Math.round(item[first_key] / multiple * 100) / 100]);
            resultData[1].push([new Date(item[time_key] * 1000), Math.round(item[second_key] / multiple * 100) / 100]);
          });
        }
        if (dataListLogs && dataListLogs.length > 0) {//for list
          resultList.push({list_name: "total", time_stamp: 0});
          dataListLogs.forEach((item: Object) => {
            resultList.push({
              list_name: item[lineNameMap.data_list_name_key],
              time_stamp: item[lineNameMap.data_list_time_stamp_key]
            });
          });
        }
        resultData[2].push([new Date(query.timestamp_base * 1000), 0.5]);
        resultData[2].push([new Date(query.service_duration_time * 1000), 0.5]);
        return {
          List: resultList, Data: resultData,
          CurListName: resJson[lineNameMap.data_list_cur_name],
          Limit: {isMax: resJson["is_over_max_limit"], isMin: resJson["is_over_min_limit"]}
        };
      })
      .catch(err => Promise.reject(err));
  }
}
