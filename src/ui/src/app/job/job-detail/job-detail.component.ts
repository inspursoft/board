import { Component, Input, OnInit } from '@angular/core';
import { Job, JobPod } from "../job.type";
import { JobService } from "../job.service";
import { CsModalChildBase, CsModalChildMessage } from "../../shared/cs-modal-base/cs-modal-child-base";
import { MessageService } from "../../shared.service/message.service";

@Component({
  selector: 'app-job-detail',
  templateUrl: './job-detail.component.html',
  styleUrls: ['./job-detail.component.css']
})
export class JobDetailComponent extends CsModalChildMessage implements OnInit {
  @Input() job: Job;
  jobPods: Array<JobPod>;
  jobLogs: Array<{datetime: string, content: string}>;
  isLoading = false;
  sinceTime: string;
  sinceDate: Date;
  currentPod: JobPod;

  constructor(private jobService: JobService,
              protected messageService: MessageService) {
    super(messageService);
    const now: Date = new Date();
    this.jobPods = Array<JobPod>();
    this.jobLogs = Array<{datetime: string, content: string}>();
    this.sinceDate = new Date(now.getFullYear(), now.getMonth(), now.getDate(), 0, 0, 0, 0);
    const num = Number(now.getHours());
    console.log();
    this.sinceTime = `${this.getFormatNumber(now.getHours())}:${this.getFormatNumber(now.getMinutes())}`;
  }

  ngOnInit() {
    this.isLoading = true;
    this.jobService.getJobPods(this.job).subscribe(
      (res: Array<JobPod>) => this.jobPods = res,
      () => this.isLoading = false,
      () => this.isLoading = false
    );
  }

  getFormatNumber(num: number): string {
    const r = Number(num).toString();
    return r.length <= 1 ? `0${r}` : r;
  }

  getSearchDateTime(): Date {
    const h = Number(this.sinceTime.split(':')[0]);
    const m = Number(this.sinceTime.split(':')[1]);
    return new Date(
      this.sinceDate.getFullYear(),
      this.sinceDate.getMonth(),
      this.sinceDate.getDate(),
      h, m, 0);
  }

  setSinceDate(date: Date) {
    this.sinceDate = date;
    this.getLogs(this.currentPod);
  }

  setSinceTime(event: Event) {
    const time = (event.target as HTMLInputElement).value;
    if (time && time !== '') {
      this.sinceTime = time;
    }
  }

  refreshLogs() {
    if (this.currentPod && !this.isLoading) {
      this.getLogs(this.currentPod);
    }
  }

  getLogs(pod: JobPod) {
    this.currentPod = pod;
    this.isLoading = true;
    this.jobLogs.splice(0, this.jobLogs.length);
    this.jobService.getJobLogs(this.job, this.currentPod, {
      timestamps: true,
      limitBytes: 1024,
      sinceTime: this.getSearchDateTime().toISOString()
    }).subscribe((res: string) => {
        res.split(/\n/).forEach((log: string) => {
          const arrLog = log.split(" ");
          if (arrLog.length == 2){
            this.jobLogs.push({datetime: arrLog[0], content: `   ${arrLog[1]}`});
          }
        })
      },
      () => this.isLoading = false,
      () => this.isLoading = false
    );
  }

}
