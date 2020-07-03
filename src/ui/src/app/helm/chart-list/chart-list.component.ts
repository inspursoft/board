import { Component, ComponentFactoryResolver, ViewContainerRef } from '@angular/core';
import { CsModalParentBase } from '../../shared/cs-modal-base/cs-modal-parent-base';
import { HelmChartVersion, HelmRepoDetail, IHelmRepo, ViewMethod } from '../helm.type';
import { HelmService } from '../helm.service';
import { Message, RETURN_STATUS } from '../../shared/shared.types';
import { MessageService } from '../../shared.service/message.service';
import { UploadChartComponent } from '../upload-chart/upload-chart.component';
import { ChartReleaseComponent } from '../chart-release/chart-release.component';
import { AppInitService } from '../../shared.service/app-init.service';

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
  viewMethod: ViewMethod = ViewMethod.List;

  constructor(private helmService: HelmService,
              private resolver: ComponentFactoryResolver,
              private view: ViewContainerRef,
              private messageService: MessageService,
              private appInitService: AppInitService) {
    super(resolver, view);
    this.versionList = Array<HelmChartVersion>();
  }

  get isSystemAdmin(): boolean {
    return this.appInitService.currentUser.user_system_admin === 1;
  }

  setViewMethod(method: string) {
    this.viewMethod = ViewMethod[method];
  }

  retrieve(): void {
    setTimeout(() => {
      this.loadingWIP = true;
      this.helmService.getRepoDetail(this.repoInfo.id, this.curPageIndex, this.curPageSize).subscribe(
        (res: HelmRepoDetail) => {
          this.versionList.splice(0, this.versionList.length);
          this.versionList = res.versionList;
          this.recordTotalCount = res.pagination.TotalCount;
        },
        () => this.loadingWIP = false,
        () => this.loadingWIP = false);
    });
  }

  showUploadChart() {
    const component = this.createNewModal(UploadChartComponent);
    component.repoInfo = this.repoInfo;
    component.uploadSuccessObservable.subscribe(() => this.retrieve());
  }

  showReleaseChartVersion(version: HelmChartVersion) {
    const component = this.createNewModal(ChartReleaseComponent);
    component.repoInfo = this.repoInfo;
    component.chartVersion = version;
  }

  deleteChartVersion(version: HelmChartVersion) {
    if (this.isSystemAdmin) {
      this.messageService.showDeleteDialog('HELM.CHART_LIST_DELETE_MSG', 'HELM.CHART_LIST_DELETE').subscribe(
        (message: Message) => {
          if (message.returnStatus === RETURN_STATUS.rsConfirm) {
            this.helmService.deleteChartVersion(this.repoInfo.id, version.name, version.version).subscribe(() => {
              this.messageService.showAlert('HELM.CHART_LIST_SUCCESS_DELETE_MSG');
              this.retrieve();
            });
          }
        });
    }
  }
}
