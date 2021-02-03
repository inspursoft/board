import { ChangeDetectorRef, Component, EventEmitter, OnInit } from '@angular/core';
import { CsModalChildBase } from '../../../shared/cs-modal-base/cs-modal-child-base';
import { ResourceService } from '../../resource.service';
import { MessageService } from '../../../shared.service/message.service';
import { SharedService } from '../../../shared.service/shared.service';
import { SharedConfigMap } from '../../../shared/shared.types';
import { ConfigMapProject } from '../../resource.types';

@Component({
  templateUrl: './create-config-map.component.html',
  styleUrls: ['./create-config-map.component.css']
})
export class CreateConfigMapComponent extends CsModalChildBase implements OnInit {
  isCreateWip = false;
  onAfterCommit: EventEmitter<SharedConfigMap>;
  newConfigMap: SharedConfigMap;
  projectList: Array<ConfigMapProject>;
  isLoadWip = false;
  configMapNamePattern: RegExp = /^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$/;

  constructor(private sharedService: SharedService,
              private change: ChangeDetectorRef,
              private resourceService: ResourceService,
              private messageService: MessageService) {
    super();
    this.onAfterCommit = new EventEmitter<SharedConfigMap>();
    this.newConfigMap = new SharedConfigMap();
    this.projectList = Array<ConfigMapProject>();
  }

  ngOnInit(): void {
    this.isLoadWip = true;
    this.resourceService.getAllProjects().subscribe(
      (res: Array<ConfigMapProject>) => this.projectList = res,
      () => this.isLoadWip = false,
      () => this.isLoadWip = false
    );
  }

  changeSelectProject(project: ConfigMapProject) {
    this.newConfigMap.namespace = project.projectName;
  }

  createConfigMap() {
    if (this.verifyDropdownExValid() && this.verifyInputExValid()) {
      this.isCreateWip = true;
      this.resourceService.createConfigMap(this.newConfigMap).subscribe(
        () => {
          this.messageService.showAlert(`RESOURCE.CREATE_CONFIG_MAP_SUCCESS`);
          this.onAfterCommit.emit(this.newConfigMap);
        },
        () => this.modalOpened = false,
        () => this.modalOpened = false
      );
    }
  }

  addKeyValue() {
    this.newConfigMap.dataList.push({key: '', value: ''});
  }

  removeKeyValue(index: number) {
    this.newConfigMap.dataList.splice(index, 1);
  }
}
