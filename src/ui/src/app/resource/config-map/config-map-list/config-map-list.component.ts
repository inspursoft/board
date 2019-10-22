import { Component, ComponentFactoryResolver, ViewContainerRef } from "@angular/core";
import { ResourceService } from "../../resource.service";
import { ConfigMap } from "../../resource.types";
import { CsModalParentBase } from "../../../shared/cs-modal-base/cs-modal-parent-base";
import { CreateConfigMapComponent } from "../create-config-map/create-config-map.component";
import { MessageService } from "../../../shared.service/message.service";
import { Message, RETURN_STATUS } from "../../../shared/shared.types";
import { ConfigMapDetailComponent } from "../config-map-detail/config-map-detail.component";
import { ConfigMapUpdateComponent } from "../config-map-update/config-map-update.component";

@Component({
  templateUrl: './config-map-list.component.html',
  styleUrls: ['./config-map-list.component.css']
})
export class ConfigMapListComponent extends CsModalParentBase {
  configMapList: Array<ConfigMap>;
  isInLoading = true;

  constructor(private view: ViewContainerRef,
              private resolver: ComponentFactoryResolver,
              private resourceService: ResourceService,
              private messageService: MessageService) {
    super(resolver, view);
    this.configMapList = new Array<ConfigMap>();
  }

  retrieve() {
    this.isInLoading = true;
    //TODO: Add the pageIndex and pageSize params when backend is ok.
    this.resourceService.getConfigMapList('', 1, 15).subscribe(
      res => this.configMapList = res,
      () => this.isInLoading = false,
      () => this.isInLoading = false
    )
  }

  confirmToDeleteConfigMap(configMap: ConfigMap) {
    this.messageService.showDeleteDialog(`RESOURCE.CONFIG_MAP_LIST_DELETE`).subscribe((message: Message) => {
      if (message.returnStatus == RETURN_STATUS.rsConfirm) {
        this.resourceService.deleteConfigMap(configMap.name, configMap.namespace).subscribe(
          () => this.messageService.showAlert(`RESOURCE.CONFIG_MAP_LIST_DELETE_SUCCESS`), null,
          () => this.retrieve())
      }
    })
  }

  createConfigMap() {
    let component = this.createNewModal(CreateConfigMapComponent);
    component.onAfterCommit.subscribe(() => this.retrieve());
  }

  showConfigMapDetail(configMap: ConfigMap) {
    let component = this.createNewModal(ConfigMapDetailComponent);
    component.projectName = configMap.namespace;
    component.configMapName = configMap.name;
  }

  updateConfigMap(configMap: ConfigMap){
    let component = this.createNewModal(ConfigMapUpdateComponent);
    component.projectName = configMap.namespace;
    component.configMapName = configMap.name;
  }
}
