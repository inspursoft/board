import { Component } from '@angular/core';
import { AppInitService } from "../shared.service/app-init.service";


@Component({
  selector: 'profile',
  styleUrls:['./profile.component.css'],
  templateUrl: './profile.component.html'
})
export class ProfileComponent {
  version: string = "";
  k8sVersion = '';
  constructor(private appInitService: AppInitService) {
    this.version = this.appInitService.systemInfo.board_version;
    this.k8sVersion = this.appInitService.systemInfo.kubernetes_version;
  }
}
