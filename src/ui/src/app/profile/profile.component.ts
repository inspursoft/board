import { Component } from '@angular/core';
import { AppInitService } from '../shared.service/app-init.service';


@Component({
  selector: 'profile',
  styleUrls: ['./profile.component.css'],
  templateUrl: './profile.component.html'
})
export class ProfileComponent {
  version = '';
  k8sVersion = '';
  processorType = '';

  constructor(private appInitService: AppInitService) {
    this.version = this.appInitService.systemInfo.board_version;
    this.k8sVersion = this.appInitService.systemInfo.kubernetes_version;
    this.processorType = this.appInitService.systemInfo.processor_type;
  }

  get isShowProcessorType(): boolean {
    return this.processorType !== '' && !this.processorType.startsWith('unknown');
  }
}
