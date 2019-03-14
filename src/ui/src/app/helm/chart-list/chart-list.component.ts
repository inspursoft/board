import { Component, ComponentFactoryResolver, ViewContainerRef } from "@angular/core";
import { CsModalParentBase } from "../../shared/cs-modal-base/cs-modal-parent-base";
import { HelmChart, HelmChartVersion, HelmRepoDetail, IHelmRepo } from "../helm.type";
import { HelmService } from "../helm.service";
import { Message, RETURN_STATUS } from "../../shared/shared.types";
import { MessageService } from "../../shared/message-service/message.service";
import { UploadChartComponent } from "../upload-chart/upload-chart.component";
import { ChartReleaseComponent } from "../chart-release/chart-release.component";

enum ViewMethod {vmList = 'list', vmCard = 'card'}

@Component({
  templateUrl: './chart-list.component.html',
  styleUrls: ['./chart-list.component.css']
})
export class ChartListComponent extends CsModalParentBase {
  repoInfo: IHelmRepo;
  loadingWIP = false;
  versionList: Array<HelmChartVersion>;
  curPageSize = 15;
  curPageIndex = 1;
  recordTotalCount = 1;
  ViewMethod = ViewMethod;
  viewMethod: ViewMethod = ViewMethod.vmList;

  constructor(private helmService: HelmService,
              private resolver: ComponentFactoryResolver,
              private view: ViewContainerRef,
              private messageService: MessageService) {
    super(resolver, view);
    this.versionList = Array<HelmChartVersion>();
  }

  retrieve(): void {
    setTimeout(() => {
      this.loadingWIP = true;
      this.helmService.getRepoDetail(this.repoInfo.id, this.curPageIndex, this.curPageSize).subscribe(
        (res: Object) => {
          this.versionList.splice(0, this.versionList.length);
          let repoDetail = HelmRepoDetail.newFromServe(res);
          repoDetail.charts.forEach((chart: HelmChart) => this.versionList.push(...chart.versions));
          this.recordTotalCount = repoDetail.pagination.total_count;
        },
        () => this.loadingWIP = false,
        () => this.loadingWIP = false);
    })
  }

  showUploadChart() {
    let component = this.createNewModal(UploadChartComponent);
    component.repoInfo = this.repoInfo;
    component.uploadSuccessObservable.subscribe(() => this.retrieve());
  }

  showReleaseChartVersion(version: HelmChartVersion) {
    let component = this.createNewModal(ChartReleaseComponent);
    component.repoInfo = this.repoInfo;
    component.chartVersion = version;
  }

  deleteChartVersion(version: HelmChartVersion) {
    this.messageService.showDeleteDialog('HELM.CHART_LIST_DELETE_MSG', 'HELM.CHART_LIST_DELETE').subscribe(
      (message: Message) => {
        if (message.returnStatus == RETURN_STATUS.rsConfirm) {
          this.helmService.deleteChartVersion(this.repoInfo.id, version.name, version.version).subscribe(() => {
            this.messageService.showAlert('HELM.CHART_LIST_SUCCESS_DELETE_MSG');
            this.retrieve();
          })
        }
      })
  }
}
