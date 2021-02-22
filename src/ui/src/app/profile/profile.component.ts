import { Component } from '@angular/core';
import { AppInitService } from '../shared.service/app-init.service';


@Component({
  styleUrls: ['./profile.component.css'],
  templateUrl: './profile.component.html'
})
export class ProfileComponent {
  version = '';
  k8sVersion = '';
  processorType = '';

  constructor(private appInitService: AppInitService) {
    this.version = this.appInitService.systemInfo.boardVersion;
    this.k8sVersion = this.appInitService.systemInfo.kubernetesVersion;
    this.processorType = this.appInitService.systemInfo.processorType;
  }

  get isShowProcessorType(): boolean {
    return this.processorType !== '' && !this.processorType.startsWith('unknown');
  }

  get logoPath(): string {
    return this.appInitService.isOpenBoard ? '../../images/board-logo.png' : '../../images/iboard-logo.png';
  }
}
