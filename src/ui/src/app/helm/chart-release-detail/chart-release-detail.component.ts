import { Component, OnInit } from "@angular/core";
import { HelmService } from "../helm.service";
import { IChartReleaseDetail } from "../helm.type";
import { CsModalChildBase } from "../../shared/cs-modal-base/cs-modal-child-base";

@Component({
  templateUrl: './chart-release-detail.component.html',
  styleUrls: ['./chart-release-detail.component.css']
})
export class ChartReleaseDetailComponent extends CsModalChildBase implements OnInit {
  detail: IChartReleaseDetail;
  releaseId = 0;
  isGotDetail = false;

  constructor(private helmService: HelmService) {
    super();
  }

  ngOnInit(): void {
    this.helmService.getChartReleaseDetail(this.releaseId)
      .subscribe((res: IChartReleaseDetail) => this.detail = res, null, () => this.isGotDetail = true);
  }
}
