import { Component, OnInit } from '@angular/core';
import { AdminService } from '../admin.service';
import { MessageService } from '../../shared.service/message.service';

@Component({
  selector: 'app-system-setting',
  templateUrl: './system-setting.component.html',
  styleUrls: ['./system-setting.component.css']
})
export class SystemSettingComponent implements OnInit {
  k8sProxyStatus: { enable: boolean };

  constructor(private adminService: AdminService,
              private messageService: MessageService) {
    this.k8sProxyStatus = {enable: false};
  }

  ngOnInit() {
    this.reloadData();
  }

  reloadData() {
    this.adminService.getK8sProxyConfig().subscribe(res => this.k8sProxyStatus = res);
  }

  setK8sProxyEnable(enable: boolean) {
    if (enable !== this.k8sProxyStatus.enable) {
      this.adminService.setK8sProxyConfig(enable).subscribe(
        () => this.messageService.showAlert('SystemSetting.SetSuccessfully'),
        () => this.messageService.showAlert('SystemSetting.SetFailed', {alertType: 'danger'}),
        () => this.reloadData()
      );
    }
  }

}
