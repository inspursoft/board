import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { CreateMethod, Job, PaginationJob } from '../job.type';

@Component({
  selector: 'app-job-create-guide',
  templateUrl: './job-create-guide.component.html',
  styleUrls: ['./job-create-guide.component.css']
})
export class JobCreateGuideComponent implements OnInit {
  @Input() paginationJobs: PaginationJob;
  @Output() confirmEvent: EventEmitter<{ method: CreateMethod, jobId: number }>;
  @Output() cancelEvent: EventEmitter<any>;
  isActionWip = false;
  createMethod: CreateMethod = CreateMethod.byDefault;
  selectJobId = 0;

  constructor() {
    this.confirmEvent = new EventEmitter();
    this.cancelEvent = new EventEmitter();
  }

  ngOnInit() {
  }

  setCreateMethod(method: CreateMethod) {
    this.createMethod = method;
    if (method === CreateMethod.byExistsJob) {
      this.selectJobId = this.paginationJobs.list[0].jobId;
    }
  }

  setSelectedJob(job: Job) {
    if (this.createMethod === CreateMethod.byExistsJob) {
      this.selectJobId = job.jobId;
    }
  }

  cancel() {
    this.cancelEvent.emit();
  }

  confirm() {
    this.isActionWip = true;
    this.confirmEvent.emit({method: this.createMethod, jobId: this.selectJobId});
  }
}
