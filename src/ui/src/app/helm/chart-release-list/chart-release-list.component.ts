import { Component, ComponentFactoryResolver, OnInit, ViewContainerRef } from "@angular/core";
import { IChartReleaseDetail, IChartReleaseList } from "../helm.type";
import { HelmService } from "../helm.service";
import { Message, RETURN_STATUS } from "../../shared/shared.types";
import { MessageService } from "../../shared/message-service/message.service";
import { CsModalParentBase } from "../../shared/cs-modal-base/cs-modal-parent-base";
import { ChartReleaseDetailComponent } from "../chart-release-detail/chart-release-detail.component";

@Component({
  templateUrl: './chart-release-list.component.html',
  styleUrls: ['./chart-release-list.component.html']
})
export class ChartReleaseListComponent extends CsModalParentBase implements OnInit {
  loadingWIP = false;
  chartReleaseList: Array<IChartReleaseList>;

  constructor(private helmService: HelmService,
              private resolver: ComponentFactoryResolver,
              private view: ViewContainerRef,
              private messageService: MessageService) {
    super(resolver, view);
    this.chartReleaseList = Array<IChartReleaseList>();
  }

  ngOnInit() {
    this.retrieve();
  }

  retrieve() {
    this.loadingWIP = true;
    this.helmService.getChartReleaseList().subscribe(
      (res: Array<IChartReleaseList>) => this.chartReleaseList = res,
      () => this.loadingWIP = false,
      () => this.loadingWIP = false
    );
  }

  showReleaseDetail(release: IChartReleaseList) {
    this.helmService.getChartReleaseDetail(release.id).subscribe((res: IChartReleaseDetail) => {
      this.view.clear();
      let component = this.createNewModal(ChartReleaseDetailComponent);
      component.detail = res;
    });
  }

  deleteChartRelease(release: IChartReleaseList) {
    this.messageService.showDeleteDialog('HELM.RELEASE_CHART_LIST_DELETE_MSG', 'HELM.RELEASE_CHART_LIST_DELETE').subscribe(
      (message: Message) => {
        if (message.returnStatus == RETURN_STATUS.rsConfirm) {
          this.helmService.deleteChartRelease(release.id).subscribe(() => {
            this.messageService.showAlert('HELM.RELEASE_CHART_LIST_DELETE_SUCCESS_MSG');
            this.retrieve();
          })
        }
      })
  }
}