import { Component, ComponentFactoryResolver, EventEmitter, OnInit, Output, ViewChild, ViewContainerRef } from "@angular/core";
import { JobContainer, JobDeployment } from "../job.type";
import { JobService } from "../job.service";
import { Project } from "../../project/project";
import { SharedActionService } from "../../shared.service/shared-action.service";
import { SharedService } from "../../shared.service/shared.service";
import { CsModalParentBase } from "../../shared/cs-modal-base/cs-modal-parent-base";
import { JobContainerCreateComponent } from "../job-container-create/job-container-create.component";
import { JobContainerConfigComponent } from "../job-container-config/job-container-config.component";
import { MessageService } from "../../shared.service/message.service";
import { IDropdownTag, Message, RETURN_STATUS } from "../../shared/shared.types";
import { CsDropdownComponent } from "../../shared/cs-components-library/cs-dropdown/cs-dropdown.component";
import { Observable, of } from "rxjs";
import { ValidationErrors } from "@angular/forms";
import { catchError, map } from "rxjs/operators";
import { HttpErrorResponse } from "@angular/common/http";
import { JobAffinityComponent } from "../job-affinity/job-affinity.component";

@Component({
  selector: 'job-create',
  styleUrls: ['./job-create.component.css'],
  templateUrl: './job-create.component.html'
})
export class JobCreateComponent extends CsModalParentBase implements OnInit {
  @Output() afterDeployment: EventEmitter<boolean>;
  @ViewChild('selectProject') selectProject: CsDropdownComponent;
  patternServiceName: RegExp = /[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/;
  projectList: Array<Project>;
  newJobDeployment: JobDeployment;
  isActionWip = false;
  nodeSelectorList: Array<{name: string, value: string, tag: IDropdownTag}>;

  constructor(private resolver: ComponentFactoryResolver,
              private view: ViewContainerRef,
              private jobService: JobService,
              private sharedService: SharedService,
              private sharedActionService: SharedActionService,
              private messageService: MessageService) {
    super(resolver, view);
    this.afterDeployment = new EventEmitter<boolean>();
    this.projectList = Array<Project>();
    this.newJobDeployment = new JobDeployment();
    this.nodeSelectorList = Array<{name: string, value: string, tag: IDropdownTag}>();
  }

  ngOnInit(): void {
    this.jobService.getProjectList().subscribe((res: Array<Project>) => {
      const createNewProject: Project = new Project();
      createNewProject.project_name = "JOB.JOB_CREATE_CREATE_PROJECT";
      createNewProject.project_id = 0;
      createNewProject["isSpecial"] = true;
      createNewProject["OnlyClick"] = true;
      this.projectList.push(createNewProject);
      if (res && res.length > 0) {
        this.projectList = this.projectList.concat(res);
      }
    });

    this.nodeSelectorList.push({name: 'JOB.JOB_CREATE_NODE_DEFAULT', value: '', tag: null});
    this.jobService.getNodeSelectors().subscribe((res: Array<{name: string, status: number}>) => {
      res.forEach((value: {name: string, status: number}) => {
        this.nodeSelectorList.push({
          name: value.name, value: value.name, tag: {
            type: value.status == 1 ? 'success' : 'warning',
            description: value.status == 1 ? 'JOB.JOB_CREATE_NODE_STATUS_SCHEDULABLE' : 'JOB.JOB_CREATE_NODE_STATUS_UNSCHEDULABLE'
          }
        })
      });
    });
  }

  get checkJobNameFun() {
    return this.checkJobName.bind(this);
  }

  get curNodeSelector() {
    return this.nodeSelectorList.find(value => value.name === this.newJobDeployment.node_selector);
  }

  get projectDropdownText() {
    return this.newJobDeployment.project_name ? this.newJobDeployment.project_name : 'JOB.JOB_CREATE_SELECT_PROJECT';
  }

  get nodeSelectorDropdownText() {
    return this.newJobDeployment.node_selector ? this.newJobDeployment.node_selector : 'JOB.JOB_CREATE_NODE_DEFAULT';
  }

  get canChangeSelectImageFun() {
    return this.canChangeSelectImage.bind(this);
  }

  checkJobName(control: HTMLInputElement): Observable<ValidationErrors | null> {
    return this.jobService.checkJobNameExists(this.newJobDeployment.project_name, control.value)
      .pipe(map(() => null), catchError((err: HttpErrorResponse) => {
        if (err.status == 409) {
          this.messageService.cleanNotification();
          return of({serviceExist: "JOB.JOB_CREATE_JOB_NAME_EXISTS"});
        } else if (err.status == 404) {
          this.messageService.cleanNotification();
        }
        return of(null);
      }));
  }

  clickSelectProject() {
    this.sharedActionService.createProjectComponent(this.view).subscribe((projectName: string) => {
      if (projectName) {
        this.jobService.getOneProject(projectName).subscribe((res: Project) => {
          if (res) {
            this.newJobDeployment.project_name = res.project_name;
            this.newJobDeployment.project_id = res.project_id;
            this.projectList.push(res);
          }
        })
      }
    });
  }

  canChangeSelectImage(project: Project): boolean {
    if (this.newJobDeployment.container_list.length > 0) {
      this.messageService.showDeleteDialog('JOB.JOB_CREATE_CHANGE_PROJECT').subscribe((msg: Message) => {
        if (msg.returnStatus == RETURN_STATUS.rsConfirm) {
          this.newJobDeployment.container_list.splice(0, this.newJobDeployment.container_list.length);
          this.selectProject.changeSelect(project);
          return true;
        } else {
          return false;
        }
      })
    } else {
      return true;
    }
  }

  changeSelectProject(project: Project) {
    this.newJobDeployment.project_id = project.project_id;
    this.newJobDeployment.project_name = project.project_name;
  }

  editContainer(container: JobContainer, isEditModel: boolean) {
    const component = this.createNewModal(JobContainerConfigComponent);
    component.container = container;
    component.containerList = this.newJobDeployment.container_list;
    component.projectId = this.newJobDeployment.project_id;
    component.projectName = this.newJobDeployment.project_name;
    component.isEditModel = isEditModel;
    component.createSuccess.subscribe((container: JobContainer) => {
      if (!isEditModel) {
        this.newJobDeployment.container_list.push(container);
      }
    })
  }

  deleteContainer(index: number) {
    this.messageService.showDeleteDialog('JOB.JOB_CREATE_DELETE_CONTAINER_COMFIRM').subscribe((msg: Message) => {
      if (msg.returnStatus == RETURN_STATUS.rsConfirm) {
        this.newJobDeployment.container_list.splice(index, 1);
      }
    })
  }

  addNewContainer() {
    if (this.newJobDeployment.project_id > 0) {
      const component = this.createNewModal(JobContainerCreateComponent);
      component.containerList = this.newJobDeployment.container_list;
      component.createSuccess.subscribe((container: JobContainer) => this.editContainer(container, false))
    }
  }

  setAffinity() {
    if (!this.isActionWip && this.newJobDeployment.project_id > 0) {
      let factory = this.factoryResolver.resolveComponentFactory(JobAffinityComponent);
      let componentRef = this.selfView.createComponent(factory);
      componentRef.instance.jobName = this.newJobDeployment.job_name;
      componentRef.instance.projectName = this.newJobDeployment.project_name;
      componentRef.instance.affinityList = this.newJobDeployment.affinity_list;
      componentRef.instance.openModal().subscribe(() => this.selfView.remove(this.selfView.indexOf(componentRef.hostView)));
    }
  }

  cancelDeploymentJob() {
    this.afterDeployment.next(false);
  }

  verifyContainer(): boolean {
    if (this.newJobDeployment.container_list.length === 0) {
      this.messageService.showAlert('JOB.JOB_CREATE_CONTAINER_COUNT', {alertType: "warning"});
      return false;
    } else {
      return true;
    }
  }

  deploymentJob() {
    if (this.verifyContainer() && this.verifyInputValid() && this.verifyInputValid()) {
      this.isActionWip = true;
      this.jobService.deploymentJob(this.newJobDeployment).subscribe(
        () => {
          this.messageService.showAlert('JOB.JOB_CREATE_SUCCESSFULLY');
          this.isActionWip = false;
        },
        () => {
          this.messageService.showAlert('JOB.JOB_CREATE_FAILED', {alertType: "warning"});
          this.isActionWip = false;
        },
        () => this.afterDeployment.next(true)
      )
    }
  }
}
