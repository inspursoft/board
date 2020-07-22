import { Component, ComponentFactoryResolver, ViewContainerRef } from '@angular/core';
import { ResourceService } from '../../resource.service';
import { CsModalParentBase } from '../../../shared/cs-modal-base/cs-modal-parent-base';
import { CreateConfigMapComponent } from '../create-config-map/create-config-map.component';
import { MessageService } from '../../../shared.service/message.service';
import { Message, RETURN_STATUS, SharedConfigMap } from '../../../shared/shared.types';
import { ConfigMapDetailComponent } from '../config-map-detail/config-map-detail.component';
import { ConfigMapUpdateComponent } from '../config-map-update/config-map-update.component';

@Component({
  templateUrl: './config-map-list.component.html',
  styleUrls: ['./config-map-list.component.css']
})
export class ConfigMapListComponent extends CsModalParentBase {
  configMapList: Array<SharedConfigMap>;
  isInLoading = true;

  constructor(private view: ViewContainerRef,
              private resolver: ComponentFactoryResolver,
              private resourceService: ResourceService,
              private messageService: MessageService) {
    super(resolver, view);
    this.configMapList = new Array<SharedConfigMap>();
  }

  retrieve() {
    this.isInLoading = true;
    this.resourceService.getConfigMapList('', 1, 15).subscribe(
      res => this.configMapList = res,
      () => this.isInLoading = false,
      () => this.isInLoading = false
    );
  }

  confirmToDeleteConfigMap(configMap: SharedConfigMap) {
    this.messageService.showDeleteDialog(`RESOURCE.CONFIG_MAP_LIST_DELETE`).subscribe((message: Message) => {
      if (message.returnStatus === RETURN_STATUS.rsConfirm) {
        this.resourceService.deleteConfigMap(configMap.name, configMap.namespace).subscribe(
          () => this.messageService.showAlert(`RESOURCE.CONFIG_MAP_LIST_DELETE_SUCCESS`), null,
          () => this.retrieve());
      }
    });
  }

  createConfigMap() {
    const component = this.createNewModal(CreateConfigMapComponent);
    component.onAfterCommit.subscribe(() => this.retrieve());
  }

  showConfigMapDetail(configMap: SharedConfigMap) {
    const component = this.createNewModal(ConfigMapDetailComponent);
    component.projectName = configMap.namespace;
    component.configMapName = configMap.name;
  }

  updateConfigMap(configMap: SharedConfigMap) {
    const component = this.createNewModal(ConfigMapUpdateComponent);
    component.projectName = configMap.namespace;
    component.configMapName = configMap.name;
  }
}
