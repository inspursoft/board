import { ChangeDetectorRef, Component } from '@angular/core';
import { INodeDetail, NodeService } from "../node.service";
import { tap } from "rxjs/operators";
import { zip } from "rxjs";

@Component({
  selector: 'node-detail',
  templateUrl: './node-detail.component.html'
})
export class NodeDetailComponent {
  nodeDetailOpened: boolean;
  nodeDetail: INodeDetail;
  nodeGroups: string;

  constructor(private nodeService: NodeService,
              private changeDetectorRef: ChangeDetectorRef) {
    this.changeDetectorRef.detach();
  }

  openNodeDetailModal(nodeName: string): void {
    this.changeDetectorRef.detach();
    this.nodeDetailOpened = true;
    this.nodeGroups = "";
    let obs1 = this.nodeService.getNodeByName(nodeName)
      .pipe(tap((nodeDetail: INodeDetail) => this.nodeDetail = nodeDetail));
    let obs2 = this.nodeService.getNodeGroupsOfOneNode(nodeName)
      .pipe(tap((res: Array<string>) => res.forEach(value => this.nodeGroups = this.nodeGroups.concat(`${value};`))));
    zip(obs1, obs2).subscribe(
      () => this.changeDetectorRef.reattach(),
      () => this.nodeDetailOpened = false);
  }

  toPercentage(num: number) {
    return Math.round(num * 100) / 100 + '%';
  }

  storagePercentage(nodeDetail: INodeDetail): number {
    return Number.parseInt(nodeDetail.storage_use) / Number.parseInt(nodeDetail.storage_total);
  }

  toGigaBytes(num: string, baseUnit?: string) {
    let denominator = 1024 * 1024 * 1024;
    if (baseUnit === 'KiB') {
      denominator = 1024 * 1024;
    }
    return Math.round(Number.parseInt(num) / denominator) + 'GB';
  }
}
