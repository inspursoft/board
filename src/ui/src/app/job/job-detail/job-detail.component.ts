import { Component, Input, OnInit } from '@angular/core';
import { HttpErrorResponse } from '@angular/common/http';
import { Job } from '../job.type';
import { JobService } from '../job.service';
import { CsModalChildMessage } from '../../shared/cs-modal-base/cs-modal-child-base';
import { MessageService } from '../../shared.service/message.service';

@Component({
  selector: 'app-job-detail',
  templateUrl: './job-detail.component.html',
  styleUrls: ['./job-detail.component.css']
})
export class JobDetailComponent extends CsModalChildMessage implements OnInit {
  @Input() job: Job;
  isLoading = false;
  jobDetail: string;

  constructor(private jobService: JobService,
              protected messageService: MessageService) {
    super(messageService);
  }

  ngOnInit() {
    this.isLoading = true;
    this.jobService.getJobStatus(this.job.jobId).subscribe(
      res => this.jobDetail = res,
      (error: HttpErrorResponse) => this.messageService.showAlert(error.message),
      () => this.isLoading = false
    );
  }

}
