import { Component, OnInit } from "@angular/core";
import { HelmService } from "../helm.service";
import { IChartReleaseDetail } from "../helm.type";
import { CsModalChildBase } from "../../shared/cs-modal-base/cs-modal-child-base";

@Component({
  templateUrl: './chart-release-detail.component.html',
  styleUrls: ['./chart-release-detail.component.html']
})
export class ChartReleaseDetailComponent extends CsModalChildBase {
  detail: IChartReleaseDetail;

  constructor() {
    super();
  }
}