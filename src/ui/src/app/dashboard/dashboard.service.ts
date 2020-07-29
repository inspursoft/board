import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

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
  time_count: number;
  time_unit: string;
  list_name: string;
  timestamp_base: number;
}

export interface IResponse {
  list: Array<string>;
  firstLineData: Array<[Date, number, string]>;
  secondLineData: Array<[Date, number, string]>;
  curListName: string;
  limit: { isMax: boolean, isMin: boolean };
}

class ResponseLineData implements IResponse {
  public list: Array<string>;
  public firstLineData: Array<[Date, number, string]>;
  public secondLineData: Array<[Date, number, string]>;
  public curListName: string;
  public limit: { isMax: boolean, isMin: boolean };

  constructor() {
    this.list = Array<string>();
    this.firstLineData = Array<[Date, number, string]>();
    this.secondLineData = Array<[Date, number, string]>();
    this.limit = {isMax: false, isMin: false};
  }
}

const BASE_URL = '/api/v1';


@Injectable()
export class DashboardService {
  private lineNameMap: Map<LineType, ILineListQueryModel>;
  ArrSizeUnit: Array<string> = ['B', 'KB', 'MB', 'GB', 'TB'];

  constructor(private http: HttpClient) {
    this.lineNameMap = new Map<LineType, ILineListQueryModel>();
    this.lineNameMap.set(LineType.ltService, {
      query_name_Key: 'service_name',
      data_list_filed_key: 'service_list_data',
      data_list_name_key: 'service_name',
      data_list_time_stamp_key: 'timestamp',
      data_list_cur_name: 'service_name',
      data_filed_key: 'service_logs_data',
      data_url_key: 'service',
      data_time_stamp_key: 'timestamp',
      data_first_line_key: 'pod_number',
      data_second_line_key: 'container_number'
    });
    this.lineNameMap.set(LineType.ltNode, {
      query_name_Key: 'node_name',
      data_list_filed_key: 'node_list_data',
      data_list_name_key: 'node_name',
      data_list_time_stamp_key: '﻿timestamp',
      data_list_cur_name: 'node_name',
      data_filed_key: 'node_logs_data',
      data_url_key: 'node',
      data_time_stamp_key: 'timestamp',
      data_first_line_key: 'cpu_usage',
      data_second_line_key: 'memory_usage'
    });
    this.lineNameMap.set(LineType.ltStorage, {
      query_name_Key: 'node_name',
      data_list_filed_key: 'node_list_data',
      data_list_name_key: 'node_name',
      data_list_time_stamp_key: '﻿timestamp',
      data_list_cur_name: 'node_name',
      data_filed_key: 'node_logs_data',
      data_url_key: 'node',
      data_time_stamp_key: 'timestamp',
      data_first_line_key: 'storage_use',
      data_second_line_key: 'storage_total'
    });
  }

  testGrafana(grafanaUrl: string): Observable<any> {
    return this.http.get(grafanaUrl, {observe: 'response'});
  }

  private getUnitValue(size: number): { value: number, unit: string } {
    let index = 0;
    let unitValue = 1;
    let originValue = size;
    while (originValue > 1024) {
      originValue = originValue / 1024;
      index += 1;
      unitValue *= 1024;
    }
    return {value: Math.round(size / unitValue * 100) / 100, unit: this.ArrSizeUnit[index]};
  }

  getQueryBody(lineType: LineType, query: IQuery): object {
    if (lineType === LineType.ltService) {
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
    return this.http.get(`${BASE_URL}/dashboard/time`)
      .pipe(map((res: object) => Reflect.get(res, 'time_now')));
  }

  getLineData(lineType: LineType, query: IQuery): Observable<IResponse> {
    const lineListQuery: ILineListQueryModel = this.lineNameMap.get(lineType);
    const requestParams = {};
    const body = this.getQueryBody(lineType, query);
    requestParams[lineListQuery.query_name_Key] = query.list_name;
    return this.http.post(`${BASE_URL}/dashboard/${lineListQuery.data_url_key}`, body, {
      params: requestParams
    }).pipe(map((res: object) => {
      const result = new ResponseLineData();
      const timeKey = lineListQuery.data_time_stamp_key;
      const firstKey = lineListQuery.data_first_line_key;
      const secondKey = lineListQuery.data_second_line_key;
      const dataLogs: Array<object> = Reflect.get(res, lineListQuery.data_filed_key);
      if (dataLogs && dataLogs.length > 0) {
        dataLogs.forEach((item: object) => {
          const timeValue = new Date(item[timeKey] * 1000);
          if (lineType === LineType.ltStorage) {
            const firstValue = this.getUnitValue(Reflect.get(item, firstKey));
            const secondValue = this.getUnitValue(Reflect.get(item, secondKey));
            result.firstLineData.push([timeValue, firstValue.value, firstValue.unit]);
            result.secondLineData.push([timeValue, secondValue.value, firstValue.unit]);
          } else {
            const firstValue = Math.round(Reflect.get(item, firstKey) * 100) / 100;
            const secondValue = Math.round(Reflect.get(item, secondKey) * 100) / 100;
            result.firstLineData.push([timeValue, firstValue, '']);
            result.secondLineData.push([timeValue, secondValue, '']);
          }
        });
      }
      result.curListName = Reflect.get(res, lineListQuery.data_list_cur_name);
      result.limit.isMax = Reflect.get(res, 'is_over_max_limit');
      result.limit.isMin = Reflect.get(res, 'is_over_min_limit');
      result.list.push(lineType === LineType.ltService ? 'total' : 'average');
      const listLogs: Array<object> = Reflect.get(res, lineListQuery.data_list_filed_key);
      if (listLogs && listLogs.length > 0) {
        listLogs.sort((a: object, b: object) =>
          Reflect.get(a, lineListQuery.data_list_name_key) > Reflect.get(b, lineListQuery.data_list_name_key) ? 1 : -1);
        listLogs.forEach((item: object) => {
          const serviceName: string = item[lineListQuery.data_list_name_key];
          result.list.push(serviceName);
        });
      }
      return result;
    }));
  }
}
