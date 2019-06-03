import { Component, Input, OnInit } from '@angular/core';
import { Job, JobPod } from "../job.type";
import { JobService } from "../job.service";
import { CsModalChildBase } from "../../shared/cs-modal-base/cs-modal-child-base";

@Component({
  selector: 'app-job-detail',
  templateUrl: './job-detail.component.html',
  styleUrls: ['./job-detail.component.css']
})
export class JobDetailComponent extends CsModalChildBase implements OnInit {
  @Input() job: Job;
  jobPods: Array<JobPod>;
  jobLogs: Array<string>;
  isLoading = false;
  sinceTime: Date;
  currentPod: JobPod;

  constructor(private jobService: JobService) {
    super();
    this.jobPods = Array<JobPod>();
    this.jobLogs = Array<string>();
    this.sinceTime = new Date(Date.now());
  }

  ngOnInit() {
    this.isLoading = true;
    this.jobService.getJobPods(this.job).subscribe(
      (res: Array<JobPod>) => this.jobPods = res,
      () => this.isLoading = false,
      () => this.isLoading = false
    );
  }

  getLogs(pod: JobPod) {
    const startSecond = 20;
    this.currentPod = pod;
    this.isLoading = true;
    this.jobService.getJobLogs(this.job, pod, {
      timestamps: true,
      sinceTime: this.sinceTime.toISOString()
    }).subscribe(
      (res: string) => this.jobLogs = res.split(/\n/),
      () => this.isLoading = false,
      () => this.isLoading = false
    );
  }

}
