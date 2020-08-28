import { Component, Input } from '@angular/core';
import { Observable } from 'rxjs';
import { HttpErrorResponse } from '@angular/common/http';
import { CsModalChildBase } from '../../shared/cs-modal-base/cs-modal-child-base';
import { Job, JobAffinity, JobAffinityCardData } from '../job.type';
import { JobService } from '../job.service';
import { MessageService } from '../../shared.service/message.service';
import { DragStatus } from '../../shared/shared.types';

@Component({
  templateUrl: './job-affinity.component.html',
  styleUrls: ['./job-affinity.component.css']
})
export class JobAffinityComponent extends CsModalChildBase {
  @Input() jobName: string;
  @Input() projectName: string;
  @Input() affinityList: Array<JobAffinity>;
  isActionWip = false;
  sourceList: Array<JobAffinityCardData>;
  selectedList: Array<{ antiFlag: boolean, list: Array<JobAffinityCardData> }>;

  constructor(private jobService: JobService,
              private messageService: MessageService) {
    super();
    this.sourceList = Array<JobAffinityCardData>();
    this.selectedList = Array<{ antiFlag: boolean, list: Array<JobAffinityCardData> }>();
  }

  addNewAffinity() {
    this.selectedList.push({antiFlag: false, list: Array<JobAffinityCardData>()});
  }

  deleteAffinity(index: number) {
    if (!this.isActionWip) {
      this.selectedList[index].list.forEach((jobCardData: JobAffinityCardData) => {
        jobCardData.status = DragStatus.dsReady;
        this.sourceList.push(jobCardData);
      });
      this.selectedList.splice(index, 1);
    }
  }

  initAffinity() {
    this.isActionWip = true;
    this.affinityList.forEach((jobAffinity: JobAffinity) => {
      const list = Array<JobAffinityCardData>();
      jobAffinity.jobNames.forEach((jobName: string) => {
        const jobCard = new JobAffinityCardData();
        jobCard.status = DragStatus.dsEnd;
        jobCard.jobName = jobName;
        list.push(jobCard);
      });
      this.selectedList.push({antiFlag: jobAffinity.antiFlag === 1, list});
    });
    this.jobService.getCollaborativeJobs(this.projectName).subscribe((res: Array<Job>) => {
      this.isActionWip = false;
      res.forEach((job: Job) => {
        const jobInUsed = this.affinityList.find(
          (jobAffinity: JobAffinity) => jobAffinity.jobNames.find(
            (jobName: string) => jobName === job.jobName) !== undefined);
        if (!jobInUsed) {
          const jobCard = new JobAffinityCardData();
          jobCard.jobName = job.jobName;
          jobCard.status = DragStatus.dsReady;
          this.sourceList.push(jobCard);
        }
      });
    }, (err: HttpErrorResponse) => {
      this.isActionWip = false;
      if (err.status === 404) {
        this.messageService.cleanNotification();
      }
    });
  }

  setAffinity() {
    this.affinityList.splice(0, this.affinityList.length);
    this.selectedList.forEach((selected: { antiFlag: boolean, list: Array<JobAffinityCardData> }) => {
      if (selected.list.length > 0) {
        const affinity = new JobAffinity();
        affinity.antiFlag = selected.antiFlag ? 1 : 0;
        selected.list.forEach((jobCard: JobAffinityCardData) => affinity.jobNames.push(jobCard.jobName));
        this.affinityList.push(affinity);
      }
    });
    this.modalOpened = false;
  }

  openModal(): Observable<any> {
    if (this.affinityList.length === 0) {
      this.addNewAffinity();
    }
    this.initAffinity();
    return super.openModal();
  }

  onDropEvent(jobCardKey: string, targetList: Array<JobAffinityCardData>) {
    const jobCard = this.sourceList.find((value: JobAffinityCardData) => value.key === jobCardKey);
    if (jobCard) {
      jobCard.status = DragStatus.dsEnd;
      const index = this.sourceList.indexOf(jobCard);
      this.sourceList.splice(index, 1);
      targetList.push(jobCard);
    }
  }

  onRemoveEvent(jobCard: JobAffinityCardData, targetList: Array<JobAffinityCardData>) {
    jobCard.status = DragStatus.dsReady;
    const index = targetList.indexOf(jobCard);
    targetList.splice(index, 1);
    this.sourceList.push(jobCard);
  }
}
