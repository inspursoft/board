import { IPagination } from "../shared/shared.types";

export interface IHelmRepo {
  id: number;
  name: string;
  url: string;
  type: number;
}

export interface IChartRelease {
  id: number;
  name: string;
  project_id: number;
  project_name: string;
  repository_id: number;
  repository: string;
  chart: string;
  chartversion: string;
  owner_id: number;
  owner_name: string;
  status: string;
  values: string;
  update_time: string;
  creation_time: string;
}

export interface IChartReleaseDetail {
  id: number;
  name: string;
  project_id: number;
  project_name: string;
  repository_id: number;
  repository: string;
  chart: string;
  chartversion: string;
  owner_id: number;
  owner_name: string;
  status: string;
  values: string;
  update_time: string;
  creation_time: string;
  notes: string;
  workloads: string;
  workloadstatus: string;
}

export class HelmChartVersion {
  name = '';
  version = '';
  description = '';
  urls: Array<string>;
  digest: string;
  icon: string;

  constructor() {
    this.urls = Array<string>();
  }

  static newFromServe(response: Object): HelmChartVersion {
    let version = new HelmChartVersion();
    version.name = response['name'];
    version.version = response['version'];
    version.description = response['description'];
    version.urls = response['urls'];
    version.digest = response['digest'];
    if (Reflect.has(response, 'icon')) {
      version.icon = response['icon'];
    }
    return version;
  }
}

export class HelmChart {
  name = '';
  versions: Array<HelmChartVersion>;

  constructor() {
    this.versions = Array<HelmChartVersion>();
  }

  static newFromServe(response: Object): HelmChart {
    let chart = new HelmChart();
    chart.name = response['name'];
    let resVersions: Array<Object> = response['versions'];
    resVersions.forEach((resVersion: Object) => {
      let version = HelmChartVersion.newFromServe(resVersion);
      chart.versions.push(version);
    });
    return chart;
  }
}

export class HelmRepoDetail {
  baseInfo: IHelmRepo;
  pagination: IPagination;
  charts: Array<HelmChart>;

  constructor() {
    this.charts = Array<HelmChart>();
    this.pagination = {page_count: 1, page_index: 0, page_size: 15, total_count: 0}
  }

  get versionList(): Array<HelmChartVersion> {
    let list = Array<HelmChartVersion>();
    this.charts.forEach((chart: HelmChart) => list.push(...chart.versions));
    return list
  }

  static newFromServe(response: Object): HelmRepoDetail {
    let detail = new HelmRepoDetail();
    detail.baseInfo = {
      id: response['id'],
      name: response['name'],
      url: response['url'],
      type: response['type']
    };
    if (response['pagination']) {
      detail.pagination = response['pagination'];
    }
    if (response['charts']) {
      let resCharts: Array<Object> = response['charts'];
      resCharts.forEach((resChart: Object) => {
        let chart = HelmChart.newFromServe(resChart);
        detail.charts.push(chart);
      });
    }
    return detail;
  }
}

export enum HelmViewType {
  RepoList, ChartList,
}

export class HelmViewData {
  description = '';

  constructor(public type: HelmViewType, public data: any = null) {

  }
}
