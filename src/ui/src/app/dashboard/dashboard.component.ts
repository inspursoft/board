import { Component } from '@angular/core';

@Component({
  selector: 'dashboard',
  templateUrl: 'dashboard.component.html',
  styleUrls: [ 'dashboard.component.css' ]
})
export class DashboardComponent {
  get serviceIcon(): string {
    return '../../images/service_icon.png';
  }

  get nodeIcon(): string {
    return '../../images/node_icon.png';
  }

  get storageIcon(): string {
    return '../../images/storage_icon.png'
  }

}