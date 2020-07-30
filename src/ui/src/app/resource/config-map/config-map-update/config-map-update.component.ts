import { Component, OnInit } from '@angular/core';
import { CsModalChildBase } from '../../../shared/cs-modal-base/cs-modal-child-base';
import { ResourceService } from '../../resource.service';
import { MessageService } from '../../../shared.service/message.service';
import { SharedConfigMap, SharedConfigMapDetail } from '../../../shared/shared.types';

@Component({
  templateUrl: './config-map-update.component.html',
  styleUrls: ['./config-map-update.component.css']
})

export class ConfigMapUpdateComponent extends CsModalChildBase implements OnInit {
  configMapDetail: SharedConfigMapDetail;
  configMapName = '';
  projectName = '';
  isLoadWip = false;
  isUpdateWip = false;

  constructor(private resourceService: ResourceService,
              private messageService: MessageService) {
    super();
    this.configMapDetail = new SharedConfigMapDetail();
  }

  ngOnInit(): void {
    this.isLoadWip = true;
    this.resourceService.getConfigMapDetail(this.configMapName, this.projectName).subscribe(
      res => this.configMapDetail = res,
      () => this.modalOpened = false,
      () => this.isLoadWip = false
    );
  }

  updateConfigMap() {
    if (this.verifyInputExValid()) {
      const configMap = new SharedConfigMap();
      configMap.name = this.configMapName;
      configMap.namespace = this.projectName;
      this.configMapDetail.dataList.forEach(value => configMap.dataList.push(value));
      this.isUpdateWip = true;
      this.resourceService.updateConfigMap(configMap).subscribe(
        () => this.messageService.showAlert('RESOURCE.CONFIG_MAP_EDIT_UPDATE_SUCCESS'),
        () => this.modalOpened = false,
        () => this.modalOpened = false
      );
    }
  }
}
