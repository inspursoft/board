import { ChangeDetectorRef, Component, EventEmitter, OnInit } from "@angular/core";
import { CsModalChildBase } from "../../../shared/cs-modal-base/cs-modal-child-base";
import { ConfigMap } from "../../resource.types";
import { ResourceService } from "../../resource.service";
import { MessageService } from "../../../shared.service/message.service";
import { SharedService } from "../../../shared.service/shared.service";
import { Project } from "../../../project/project";

@Component({
  templateUrl: './create-config-map.component.html',
  styleUrls: ['./create-config-map.component.css']
})
export class CreateConfigMapComponent extends CsModalChildBase implements OnInit {
  isCreateWip = false;
  onAfterCommit: EventEmitter<ConfigMap>;
  newConfigMap: ConfigMap;
  projectList: Array<Project>;
  isLoadWip = false;

  constructor(private sharedService: SharedService,
              private change: ChangeDetectorRef,
              private resourceService: ResourceService,
              private messageService: MessageService) {
    super();
    this.onAfterCommit = new EventEmitter<ConfigMap>();
    this.newConfigMap = new ConfigMap();
    this.projectList = Array<Project>();
  }

  ngOnInit(): void {
    this.isLoadWip = true;
    this.sharedService.getAllProjects().subscribe(
      (res: Array<Project>) => this.projectList = res,
      () => this.isLoadWip = false,
      () => this.isLoadWip = false
    );
  }

  changeSelectProject(project: Project) {
    this.newConfigMap.namespace = project.project_name;
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
      )
    }
  }

  addKeyValue() {
    this.newConfigMap.dataList.push({key: '', value: ''});
  }

  removeKeyValue(index: number) {
    this.newConfigMap.dataList.splice(index, 1);
  }
}
