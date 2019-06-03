import { ChangeDetectionStrategy, ChangeDetectorRef, Component, ComponentFactoryResolver, OnInit, ViewContainerRef } from "@angular/core";
import { Job, PaginationJob } from "../job.type";
import { JobService } from "../job.service";
import { MessageService } from "../../shared.service/message.service";
import { Message, RETURN_STATUS } from "../../shared/shared.types";
import { JobDetailComponent } from "../job-detail/job-detail.component";
import { CsModalParentBase } from "../../shared/cs-modal-base/cs-modal-parent-base";

@Component({
  templateUrl: './job-list.component.html',
  styleUrls: ['./job-list.component.css'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class JobListComponent extends CsModalParentBase implements OnInit {
  loadingWIP = false;
  createNewJob = false;
  paginationJobs: PaginationJob;

  constructor(private resolver: ComponentFactoryResolver,
              private view: ViewContainerRef,
              private jobService: JobService,
              private messageService: MessageService,
              private changeRef: ChangeDetectorRef) {
    super(resolver, view);
    this.paginationJobs = new PaginationJob();
  }

  /*preparing = iota
    running
    stopped
    uncompleted
    warning
    deploying
    completed
    failed*/
  getJobStatus(status: number) {
    switch (status) {
      case 0:
        return 'JOB.STATUS_PREPARING';
      case 1:
        return 'JOB.STATUS_RUNNING';
      case 2:
        return 'JOB.STATUS_STOPPED';
      case 3:
        return 'JOB.STATUS_UNCOMPLETED';
      case 4:
        return 'JOB.STATUS_WARNING';
      case 5:
        return 'JOB.STATUS_DEPLOYING';
      case 6:
        return 'JOB.STATUS_COMPLETED';
      case 7:
        return 'JOB.STATUS_FAILED';
    }
  }

  ngOnInit(): void {
    this.retrieve();
  }

  retrieve() {
    this.loadingWIP = true;
    this.changeRef.detectChanges();
    this.jobService.getJobList(this.paginationJobs.pagination.page_index, this.paginationJobs.pagination.page_size).subscribe(
      (res: PaginationJob) => this.paginationJobs = res,
      () => {
        this.loadingWIP = false;
        this.changeRef.detectChanges();
      },
      () => {
        this.loadingWIP = false;
        this.changeRef.detectChanges();
      }
    );
  }

  deleteJob(job: Job) {
    this.messageService.showDeleteDialog('JOB.JOB_LIST_DELETE_CONFIRM').subscribe((msg: Message) => {
      if (msg.returnStatus == RETURN_STATUS.rsConfirm) {
        this.jobService.deleteJob(job).subscribe(
          () => this.messageService.showAlert('JOB.JOB_LIST_DELETE_SUCCESSFULLY'),
          () => this.messageService.showAlert('JOB.JOB_LIST_DELETE_FAILED', {alertType: "warning"}),
          () => this.retrieve())
      }
    })
  }

  showJobDetail(job: Job) {
    const component = this.createNewModal(JobDetailComponent);
    component.job = job;
  }

  createJob() {
    this.createNewJob = true;
  }

  afterDeployment(isSuccess: boolean) {
    if (isSuccess) {
      this.retrieve();
    }
    this.createNewJob = false;
  }
}
