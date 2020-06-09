import { Component, Input, OnInit } from '@angular/core';
import { CreateMethod, Job, PaginationJob } from '../job.type';

@Component({
  selector: 'app-job-create-guide',
  templateUrl: './job-create-guide.component.html',
  styleUrls: ['./job-create-guide.component.css']
})
export class JobCreateGuideComponent implements OnInit {
  @Input() paginationJobs: PaginationJob;
  createMethod: CreateMethod = CreateMethod.byDefault;
  selectJobId = 0;

  constructor() {
  }

  ngOnInit() {
  }

  setCreateMethod(method: CreateMethod) {
    this.createMethod = method;
  }

  setSelectedJob(job: Job) {
    if (this.createMethod === CreateMethod.byExistsJob) {
      this.selectJobId = job.jobId;
    }
  }
}
