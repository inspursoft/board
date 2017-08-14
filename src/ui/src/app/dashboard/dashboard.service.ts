import { Injectable } from '@angular/core';
import { Http, RequestOptions, Headers, Response } from "@angular/http"
import { AppInitService } from "../app.init.service";
import "rxjs/add/operator/map";
import 'rxjs/add/operator/toPromise';
import 'rxjs/add/operator/timeout';

export enum LineType {ltService, ltNode, ltStorage}
export type LineDataModel = [Date, number];
export type LinesData = [LineDataModel[], LineDataModel[]];

export interface LineListQueryModel {
  readonly list_name_Key: string;
  readonly list_filed_Key: string;
  readonly list_url_key: string;
  readonly data_filed_key: string;
  readonly data_url_key: string;
  readonly data_time_stamp_key: string;
  readonly data_first_line_key: string;
  readonly data_second_line_key: string;
}

export interface LineListDataModel {
  list_name: string;
  time_list_name?: number;
}
const BASE_URL = "/api/v1";
@Injectable()
export class DashboardService {
  LineNameMap: Map<LineType, LineListQueryModel>;
  constructor(private http: Http,
              private appInitService: AppInitService) {
    this.LineNameMap = new Map<LineType, LineListQueryModel>();
    this.LineNameMap.set(LineType.ltService, {
      list_name_Key: "service_name",
      list_filed_Key: "service_list",
      list_url_key: "service/list",
      data_filed_key: "service_statuslogs",
      data_url_key: "service",
      data_time_stamp_key: "podcontainer_timestamp",
      data_first_line_key: "pods_number",
      data_second_line_key: "container_number"
    });
    this.LineNameMap.set(LineType.ltNode, {
      list_name_Key: "node_name",
      list_filed_Key: "list",
      list_url_key: "node/list",
      data_filed_key: "node_logs",
      data_url_key: "node",
      data_time_stamp_key: "time_stamp",
      data_first_line_key: "cpu_usage",
      data_second_line_key: "mem_usage"
    });
    this.LineNameMap.set(LineType.ltStorage, {
      list_name_Key: "node_name",
      list_filed_Key: "list",
      list_url_key: "node/list",
      data_filed_key: "node_logs",
      data_url_key: "node",
      data_time_stamp_key: "time_stamp",
      data_first_line_key: "storage_use",
      data_second_line_key: "storage_total"
    });
  };

  readonly defaultHeaders: Headers = new Headers({
    contentType: "application/json"
  });

  getLineNameList(lineType: LineType): Promise<LineListDataModel[]> {
    let options = new RequestOptions({
        headers: this.defaultHeaders,
        params: {'token': this.appInitService.token}
      }
    );
    return this.http.get(`${BASE_URL}/dashboard/${this.LineNameMap.get(lineType).list_url_key}`, options)
      .toPromise()
      .then(res => {
        let resJson = res.json();
        let resultArr = Array<LineListDataModel>();
        let nameKey = this.LineNameMap.get(lineType).list_name_Key;
        Array.from(resJson[this.LineNameMap.get(lineType).list_filed_Key]).forEach(value => {
          resultArr.push({list_name: value[nameKey], time_list_name: value["time_list_name"]});
        });
        resultArr.sort((left, right) => {
          return left.list_name == right.list_name ? 0 : left.list_name > right.list_name ? 1 : -1;
        });
        if (lineType == LineType.ltService) {//only service have total api
          resultArr.unshift({list_name: "total"});
        }
        return resultArr;
      })
      .catch(err => Promise.reject(err));
  }

  getLineData(lineType: LineType, query: {time_count: number, time_unit: string, list_name: string, timestamp_base: number}): Promise<LinesData> {
    let lineKey = this.LineNameMap.get(lineType).list_name_Key;
    let requestParams = {'token': this.appInitService.token};
    requestParams[lineKey] = query.list_name;
    let options = new RequestOptions({
      headers: this.defaultHeaders,
      params: requestParams
    });
    return this.http.post(`${BASE_URL}/dashboard/${this.LineNameMap.get(lineType).data_url_key}`, {
      time_count: query.time_count.toString(),
      time_unit: query.time_unit,
      timestamp_base: (query.timestamp_base).toString()
    }, options)
      .toPromise()
      .then((res: Response) => {
        let resJson: Object = res.json();
        let result: LinesData = [Array<[Date, number]>(0), Array<[Date, number]>(0)];
        let logs: Array<Object> = resJson[this.LineNameMap.get(lineType).data_filed_key];
        let time_key = this.LineNameMap.get(lineType).data_time_stamp_key;
        let first_key = this.LineNameMap.get(lineType).data_first_line_key;
        let second_key = this.LineNameMap.get(lineType).data_second_line_key;
        if (logs && logs.length > 0) {
          logs.forEach((item: Object) => {
            result[0].push([new Date(item[time_key] * 1000), Math.round(item[first_key] * 100) / 100]);
            result[1].push([new Date(item[time_key] * 1000), Math.round(item[second_key] * 100) / 100]);
          });
          result[0].sort((left,right)=>{
            return left[0] == right[0] ? 0 : left[0] > right[0] ? 1 : -1;
          });
          result[1].sort((left,right)=>{
            return left[0] == right[0] ? 0 : left[0] > right[0] ? 1 : -1;
          });
        }
        return result;
      })
      .catch(err => Promise.reject(err));
  }
}
