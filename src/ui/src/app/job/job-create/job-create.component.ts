import { Component, ComponentFactoryResolver, EventEmitter, Input, OnInit, Output, ViewContainerRef } from '@angular/core';
import { Observable, of } from 'rxjs';
import { ValidationErrors } from '@angular/forms';
import { catchError, map } from 'rxjs/operators';
import { HttpErrorResponse } from '@angular/common/http';
import { JobContainer, JobDeployment } from '../job.type';
import { JobService } from '../job.service';
import { SharedActionService } from '../../shared.service/shared-action.service';
import { SharedService } from '../../shared.service/shared.service';
import { CsModalParentBase } from '../../shared/cs-modal-base/cs-modal-parent-base';
import { JobContainerCreateComponent } from '../job-container-create/job-container-create.component';
import { JobContainerConfigComponent } from '../job-container-config/job-container-config.component';
import { MessageService } from '../../shared.service/message.service';
import { IDropdownTag, Message, RETURN_STATUS, SharedProject } from '../../shared/shared.types';
import { JobAffinityComponent } from '../job-affinity/job-affinity.component';

@Component({
  selector: 'app-job-create',
  styleUrls: ['./job-create.component.css'],
  templateUrl: './job-create.component.html'
})
export class JobCreateComponent extends CsModalParentBase implements OnInit {
  @Output() afterDeployment: EventEmitter<boolean>;
  @Input() newJobDeployment: JobDeployment;
  patternServiceName: RegExp = /^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/;
  projectList: Array<SharedProject>;
  isActionWip = false;
  projectDefaultIndex = -1;
  nodeSelectorDefaultIndex = -1;
  nodeSelectorList: Array<{ name: string, value: string, tag: IDropdownTag }>;

  constructor(private resolver: ComponentFactoryResolver,
              private view: ViewContainerRef,
              private jobService: JobService,
              private sharedService: SharedService,
              private sharedActionService: SharedActionService,
              private messageService: MessageService) {
    super(resolver, view);
    this.afterDeployment = new EventEmitter<boolean>();
    this.projectList = Array<SharedProject>();
    this.nodeSelectorList = Array<{ name: string, value: string, tag: IDropdownTag }>();
  }

  ngOnInit(): void {
    this.jobService.getProjectList().subscribe((res: Array<SharedProject>) => {
      this.projectList = res;
      if (this.newJobDeployment.projectId > 0) {
        const project = this.projectList.find(value => value.projectId === this.newJobDeployment.projectId);
        this.projectDefaultIndex = this.projectList.indexOf(project);
      }
    });

    this.nodeSelectorList.push({name: 'JOB.JOB_CREATE_NODE_DEFAULT', value: '', tag: null});
    this.jobService.getNodeSelectors().subscribe((res: Array<{ name: string, status: number }>) => {
      res.forEach((value: { name: string, status: number }) => {
        this.nodeSelectorList.push({
          name: value.name, value: value.name, tag: {
            type: value.status === 1 ? 'success' : 'warning',
            description: value.status === 1 ? 'JOB.JOB_CREATE_NODE_STATUS_SCHEDULABLE' : 'JOB.JOB_CREATE_NODE_STATUS_UNSCHEDULABLE'
          }
        });
      });
      if (this.newJobDeployment.nodeSelector !== '') {
        const nodeSelector = this.nodeSelectorList.find(value => value.name === this.newJobDeployment.nodeSelector);
        this.nodeSelectorDefaultIndex = this.nodeSelectorList.indexOf(nodeSelector);
      }
    });
  }

  get checkJobNameFun() {
    return this.checkJobName.bind(this);
  }

  get canChangeSelectImageFun() {
    return this.canChangeSelectImage.bind(this);
  }

  getItemTagClass(dropdownTag: IDropdownTag) {
    return {
      'label-info': dropdownTag.type === 'success',
      'label-warning': dropdownTag.type === 'warning',
      'label-danger': dropdownTag.type === 'danger'
    };
  }

  checkJobName(control: HTMLInputElement): Observable<ValidationErrors | null> {
    return this.jobService.checkJobNameExists(this.newJobDeployment.projectName, control.value)
      .pipe(map(() => null), catchError((err: HttpErrorResponse) => {
        if (err.status === 409) {
          this.messageService.cleanNotification();
          return of({jobNameExists: 'JOB.JOB_CREATE_JOB_NAME_EXISTS'});
        } else if (err.status === 404) {
          this.messageService.cleanNotification();
        }
        return of(null);
      }));
  }

  clickSelectProject() {
    this.sharedActionService.createProjectComponent(this.view).subscribe((projectName: string) => {
      if (projectName) {
        this.jobService.getOneProject(projectName).subscribe((res: SharedProject) => {
          if (res) {
            this.newJobDeployment.projectName = res.projectName;
            this.newJobDeployment.projectId = res.projectId;
            this.projectList.push(res);
          }
        });
      }
    });
  }

  canChangeSelectImage(): Observable<boolean> {
    if (this.newJobDeployment.containerList.length > 0) {
      return this.messageService.showDeleteDialog('JOB.JOB_CREATE_CHANGE_PROJECT')
        .pipe(map((msg: Message) => {
          if (msg.returnStatus === RETURN_STATUS.rsConfirm) {
            this.newJobDeployment.containerList.splice(0, this.newJobDeployment.containerList.length);
            return true;
          } else {
            return false;
          }
        }));
    } else {
      return of(true);
    }
  }

  changeSelectProject(project: SharedProject) {
    this.newJobDeployment.projectId = project.projectId;
    this.newJobDeployment.projectName = project.projectName;
  }

  editContainer(container: JobContainer, isEditModel: boolean) {
    const component = this.createNewModal(JobContainerConfigComponent);
    component.container = container;
    component.containerList = this.newJobDeployment.containerList;
    component.projectId = this.newJobDeployment.projectId;
    component.projectName = this.newJobDeployment.projectName;
    component.isEditModel = isEditModel;
    component.createSuccess.subscribe((jobContainer: JobContainer) => {
      if (!isEditModel) {
        this.newJobDeployment.containerList.push(jobContainer);
      }
    });
  }

  deleteContainer(index: number) {
    this.messageService.showDeleteDialog('JOB.JOB_CREATE_DELETE_CONTAINER_COMFIRM').subscribe(
      (msg: Message) => {
        if (msg.returnStatus === RETURN_STATUS.rsConfirm) {
          this.newJobDeployment.containerList.splice(index, 1);
        }
      });
  }

  addNewContainer() {
    if (this.newJobDeployment.projectId > 0) {
      const component = this.createNewModal(JobContainerCreateComponent);
      component.containerList = this.newJobDeployment.containerList;
      component.createSuccess.subscribe((container: JobContainer) => this.editContainer(container, false));
    }
  }

  setAffinity() {
    if (!this.isActionWip && this.newJobDeployment.projectId > 0) {
      const factory = this.factoryResolver.resolveComponentFactory(JobAffinityComponent);
      const componentRef = this.selfView.createComponent(factory);
      componentRef.instance.jobName = this.newJobDeployment.jobName;
      componentRef.instance.projectName = this.newJobDeployment.projectName;
      componentRef.instance.affinityList = this.newJobDeployment.affinityList;
      componentRef.instance.openModal().subscribe(() => this.selfView.remove(this.selfView.indexOf(componentRef.hostView)));
    }
  }

  cancelDeploymentJob() {
    this.afterDeployment.next(false);
  }

  verifyContainer(): boolean {
    if (this.newJobDeployment.containerList.length === 0) {
      this.messageService.showAlert('JOB.JOB_CREATE_CONTAINER_COUNT', {alertType: 'warning'});
      return false;
    } else {
      return true;
    }
  }

  deploymentJob() {
    if (this.verifyDropdownExValid() && this.verifyInputExValid() && this.verifyContainer()) {
      this.isActionWip = true;
      this.jobService.deploymentJob(this.newJobDeployment).subscribe(
        () => {
          this.messageService.showAlert('JOB.JOB_CREATE_SUCCESSFULLY');
          this.isActionWip = false;
        },
        () => this.isActionWip = false,
        () => this.afterDeployment.next(true)
      );
    }
  }
}
