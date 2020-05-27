import { ChangeDetectorRef, Component } from '@angular/core';
import { NodeService } from '../node.service';
import { tap } from 'rxjs/operators';
import { zip } from 'rxjs';
import { NodeDetail } from '../node.types';

@Component({
  selector: 'app-node-detail',
  templateUrl: './node-detail.component.html'
})
export class NodeDetailComponent {
  nodeDetailOpened: boolean;
  nodeDetail: NodeDetail;
  nodeGroups: string;

  constructor(private nodeService: NodeService,
              private changeDetectorRef: ChangeDetectorRef) {
    this.changeDetectorRef.detach();
  }

  openNodeDetailModal(nodeName: string): void {
    this.changeDetectorRef.detach();
    this.nodeDetailOpened = true;
    this.nodeGroups = '';
    const obs1 = this.nodeService.getNodeDetailByName(nodeName)
      .pipe(tap((nodeDetail: NodeDetail) => this.nodeDetail = nodeDetail));
    const obs2 = this.nodeService.getNodeGroupsOfOneNode(nodeName)
      .pipe(tap((res: Array<string>) => res.forEach(value => this.nodeGroups = this.nodeGroups.concat(`${value};`))));
    zip(obs1, obs2).subscribe(
      () => this.changeDetectorRef.reattach(),
      () => this.nodeDetailOpened = false);
  }

  toPercentage(num: number) {
    return Math.round(num * 100) / 100 + '%';
  }

  storagePercentage(nodeDetail: NodeDetail): number {
    return Number.parseInt(nodeDetail.storageUse, 10) / Number.parseInt(nodeDetail.storageTotal, 10);
  }

  toGigaBytes(num: string, baseUnit?: string) {
    let denominator = 1024 * 1024 * 1024;
    if (baseUnit === 'KiB') {
      denominator = 1024 * 1024;
    }
    return Math.round(Number.parseInt(num, 10) / denominator) + 'GB';
  }
}
