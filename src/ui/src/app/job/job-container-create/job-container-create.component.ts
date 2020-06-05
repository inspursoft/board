import { Component, Input, OnInit } from '@angular/core';
import { Subject } from 'rxjs';
import { JobContainer, JobImage, JobImageDetailInfo, JobImageInfo } from '../job.type';
import { CsModalChildBase } from '../../shared/cs-modal-base/cs-modal-child-base';
import { JobService } from '../job.service';
import { MessageService } from '../../shared.service/message.service';

@Component({
  selector: 'app-job-container-create',
  templateUrl: './job-container-create.component.html',
  styleUrls: ['./job-container-create.component.css']
})
export class JobContainerCreateComponent extends CsModalChildBase implements OnInit {
  @Input() containerList: Array<JobContainer>;
  imageList: Array<JobImageInfo>;
  imageDetailList: Array<JobImageDetailInfo>;
  isLoading = false;
  selectedImageName: string;
  selectedImageTag: string;
  createSuccess: Subject<JobContainer>;

  constructor(private jobService: JobService,
              private messageService: MessageService) {
    super();
    this.imageList = Array<JobImageInfo>();
    this.imageDetailList = Array<JobImageDetailInfo>();
    this.createSuccess = new Subject();
  }

  ngOnInit() {
    this.selectedImageName = '';
    this.selectedImageTag = '';
    this.isLoading = true;
    this.jobService.getImageList().subscribe((res: Array<JobImageInfo>) => {
      if (res && res.length > 0) {
        this.imageList = res;
        this.getImageDetail(this.imageList[0]);
      }
    }, () => this.isLoading = false, () => this.isLoading = false);
  }

  getImageDetail(image: JobImageInfo) {
    this.isLoading = true;
    this.selectedImageName = image.imageName;
    this.selectedImageTag = '';
    this.jobService.getImageDetailList(image.imageName).subscribe(
      (res: Array<JobImageDetailInfo>) => {
        if (res && res.length > 0) {
          this.imageDetailList = res;
        }
      }, () => this.isLoading = false, () => this.isLoading = false);
  }

  setImageTag(detail: JobImageDetailInfo) {
    this.selectedImageTag = detail.imageTag;
  }

  checkImageAndTag(): boolean {
    return this.containerList.find((jobContainer: JobContainer) =>
      jobContainer.image.imageTag === this.selectedImageTag &&
      jobContainer.image.imageName === this.selectedImageName) !== undefined;
  }

  createNewContainer() {
    if (this.selectedImageName === '' || this.selectedImageTag === '') {
      this.messageService.showAlert('JOB.JOB_CREATE_SELECT_IMAGE_TIP', {alertType: 'warning', view: this.alertView});
    } else if (this.checkImageAndTag()) {
      this.messageService.showAlert('JOB.JOB_CREATE_SELECT_IMAGE_EXISTS', {alertType: 'warning', view: this.alertView});
    } else {
      const container = new JobContainer();
      const jobImage = new JobImage();
      jobImage.imageName = this.selectedImageName;
      jobImage.imageTag = this.selectedImageTag;
      container.image = jobImage;
      this.createSuccess.next(container);
      this.modalOpened = false;
    }
  }
}
