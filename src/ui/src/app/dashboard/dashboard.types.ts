import { HttpBase, HttpBind, HttpBindArray } from '../shared/ui-model/model-types';

export const ArrSizeUnit: Array<string> = ['B', 'KB', 'MB', 'GB', 'TB'];
export const BASE_URL = '/api/v1';

export enum LineType {ltService, ltNode, ltStorage}

export class RealtimeData {
  curFirst = 0;
  curFirstUnit = '';
  curSecond = 0;
  curSecondUnit = '';
}

export interface ScaleOption {
  readonly description: string;
  readonly value: string;
  readonly valueOfSecond: number;
  readonly id: number;
}

export class ThirdLine {
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

export class ResponseLineData {
  list: Array<string>;
  firstLineData: Array<[Date, number, string]>;
  secondLineData: Array<[Date, number, string]>;
  curListName = '';
  limit: { isMax: boolean, isMin: boolean };

  constructor() {
    this.list = Array<string>();
    this.firstLineData = Array<[Date, number, string]>();
    this.secondLineData = Array<[Date, number, string]>();
    this.limit = {isMax: false, isMin: false};
  }

  get isHaveData(): boolean {
    return this.firstLineData.length > 0;
  }
}

export class NodeLog extends HttpBase {
  @HttpBind('cpu_usage') cpuUsage = 0;
  @HttpBind('memory_usage') memoryUsage = 0;
  @HttpBind('storage_total') storageTotal = 0;
  @HttpBind('storage_used') storageUsed = 0;
  @HttpBind('timestamp') timestamp = 0;
}

export class ServiceLog extends HttpBase {
  @HttpBind('container_number') containerNumber = 0;
  @HttpBind('pod_number') podNumber = 0;
  @HttpBind('timestamp') timestamp = 0;
}

export class NodeData extends HttpBase {
  @HttpBind('name') name = '';
  @HttpBindArray('node_logs_data', NodeLog) nodeLogsData: Array<NodeLog>;

  protected prepareInit() {
    this.nodeLogsData = new Array<NodeLog>();
  }
}

export class ServiceData extends HttpBase {
  @HttpBind('name') name = '';
  @HttpBindArray('service_logs_data', ServiceLog) serviceLogsData: Array<ServiceLog>;

  protected prepareInit() {
    this.serviceLogsData = new Array<ServiceLog>();
  }
}

export class BodyData extends HttpBase {
  @HttpBind('timestamp') queryTimestamp = 0;
  @HttpBind('time_unit') queryTimeUnit = '';
  @HttpBind('time_count') queryCount = 0;
}

export class QueryData extends HttpBase {
  @HttpBind('service') serviceName = '';
  @HttpBind('node') nodeName = '';
}

export class Prometheus extends HttpBase {
  @HttpBind('is_over_max_limit') isOverMaxLimit = false;
  @HttpBind('is_over_min_limit') isOverMinLimit = false;
  @HttpBind('node_count') nodeCount = 0;
  @HttpBind('service_count') serviceCount = 0;
  @HttpBindArray('node_list_data', NodeData) nodeListData: Array<NodeData>;
  @HttpBindArray('service_list_data', ServiceData) serviceListData: Array<ServiceData>;
  serviceLineData: ResponseLineData;
  nodeLineData: ResponseLineData;
  storageLineData: ResponseLineData;

  static getUnitValue(size: number): { value: number, unit: string } {
    let index = 0;
    let unitValue = 1;
    let originValue = size;
    while (originValue > 1024) {
      originValue = originValue / 1024;
      index += 1;
      unitValue *= 1024;
    }
    return {value: Math.round(size / unitValue * 100) / 100, unit: ArrSizeUnit[index]};
  }

  protected prepareInit() {
    this.nodeListData = new Array<NodeData>();
    this.serviceListData = new Array<ServiceData>();
    this.serviceLineData = new ResponseLineData();
    this.nodeLineData = new ResponseLineData();
    this.storageLineData = new ResponseLineData();
  }

  get serviceRealtimeData(): RealtimeData {
    const isHaveData = this.serviceLineData.isHaveData;
    return {
      curFirst: isHaveData ? this.serviceLineData.firstLineData[0][1] : 0,
      curFirstUnit: isHaveData ? this.serviceLineData.firstLineData[0][2] : '',
      curSecond: isHaveData ? this.serviceLineData.secondLineData[0][1] : 0,
      curSecondUnit: isHaveData ? this.serviceLineData.secondLineData[0][2] : ''
    };
  }

  get nodeRealtimeData(): RealtimeData {
    const isHaveData = this.nodeLineData.isHaveData;
    return {
      curFirst: isHaveData ? this.nodeLineData.firstLineData[0][1] : 0,
      curFirstUnit: isHaveData ? this.nodeLineData.firstLineData[0][2] : '',
      curSecond: isHaveData ? this.nodeLineData.secondLineData[0][1] : 0,
      curSecondUnit: isHaveData ? this.nodeLineData.secondLineData[0][2] : ''
    };
  }

  get storageRealtimeData(): RealtimeData {
    const isHaveData = this.storageLineData.isHaveData;
    return {
      curFirst: isHaveData ? this.storageLineData.firstLineData[0][1] : 0,
      curFirstUnit: isHaveData ? this.storageLineData.firstLineData[0][2] : '',
      curSecond: isHaveData ? this.storageLineData.secondLineData[0][1] : 0,
      curSecondUnit: isHaveData ? this.storageLineData.secondLineData[0][2] : ''
    };
  }

  getResponseLineData(lineType: LineType): ResponseLineData {
    if (lineType === LineType.ltService) {
      return this.serviceLineData;
    } else if (lineType === LineType.ltNode) {
      return this.nodeLineData;
    } else {
      return this.storageLineData;
    }
  }

  analyzeData(serviceName, nodeName: string) {
    this.serviceLineData.curListName = serviceName;
    const serviceList = new Array<string>();
    this.serviceListData.forEach(value => {
      serviceList.push(value.name);
      if (value.name === serviceName) {
        value.serviceLogsData.forEach((log: ServiceLog) => {
          if (log.timestamp > 0) {
            const date = new Date(log.timestamp * 1000);
            const pod = Math.round(log.podNumber * 100) / 100;
            const container = Math.round(log.containerNumber * 100) / 100;
            this.serviceLineData.firstLineData.push([date, pod, '']);
            this.serviceLineData.secondLineData.push([date, container, '']);
          }
        });
      }
    });
    this.serviceLineData.list = serviceList;

    this.nodeLineData.curListName = nodeName;
    this.storageLineData.curListName = nodeName;
    const nodeList = new Array<string>();
    this.nodeListData.forEach(value => {
      nodeList.push(value.name);
      if (value.name === nodeName) {
        value.nodeLogsData.forEach((log: NodeLog) => {
          if (log.timestamp > 0) {
            const date = new Date(log.timestamp * 1000);
            const cpuUsage = Math.round(log.cpuUsage * 100) / 100;
            const memoryUsage = Math.round(log.memoryUsage * 100) / 100;
            this.nodeLineData.firstLineData.push([date, cpuUsage, '']);
            this.nodeLineData.secondLineData.push([date, memoryUsage, '']);
            const used = Prometheus.getUnitValue(log.storageUsed);
            const total = Prometheus.getUnitValue(log.storageTotal);
            this.storageLineData.firstLineData.push([date, used.value, used.unit]);
            this.storageLineData.secondLineData.push([date, total.value, total.unit]);
          }
        });
      }
    });
    this.nodeLineData.list = nodeList;
    this.storageLineData.list = nodeList;
  }

}




