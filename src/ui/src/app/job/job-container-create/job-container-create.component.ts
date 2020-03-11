import { Component, Input, OnInit } from '@angular/core';
import { JobContainer, JobImage } from "../job.type";
import { CsModalChildBase } from "../../shared/cs-modal-base/cs-modal-child-base";
import { JobService } from "../job.service";
import { Image, ImageDetail } from "../../image/image";
import { MessageService } from "../../shared.service/message.service";
import { Subject } from "rxjs";

@Component({
  selector: 'app-job-container-create',
  templateUrl: './job-container-create.component.html',
  styleUrls: ['./job-container-create.component.css']
})
export class JobContainerCreateComponent extends CsModalChildBase implements OnInit {
  @Input() containerList: Array<JobContainer>;
  imageList: Array<Image>;
  imageDetailList: Array<ImageDetail>;
  isLoading = false;
  selectedImageName: string;
  selectedImageTag: string;
  createSuccess: Subject<JobContainer>;

  constructor(private jobService: JobService,
              private messageService: MessageService) {
    super();
    this.imageList = Array<Image>();
    this.imageDetailList = Array<ImageDetail>();
    this.createSuccess = new Subject();
  }

  ngOnInit() {
    this.selectedImageName = "";
    this.selectedImageTag = "";
    this.isLoading = true;
    this.jobService.getImageList().subscribe((res: Array<Image>) => {
      if (res && res.length > 0) {
        this.imageList = res;
        this.getImageDetail(this.imageList[0]);
      }
    }, () => this.isLoading = false, () => this.isLoading = false)
  }

  getImageDetail(image: Image) {
    this.isLoading = true;
    this.selectedImageName = image.image_name;
    this.selectedImageTag = "";
    this.jobService.getImageDetailList(image.image_name)
      .subscribe((res: Array<ImageDetail>) => {
        if (res && res.length > 0) {
          this.imageDetailList = res;
        }
      }, () => this.isLoading = false, () => this.isLoading = false)
  }

  setImageTag(detail: ImageDetail) {
    this.selectedImageTag = detail.image_tag;
  }

  checkImageAndTag(): boolean {
    return this.containerList.find((jobContainer: JobContainer) =>
      jobContainer.image.image_tag === this.selectedImageTag &&
      jobContainer.image.image_name === this.selectedImageName) !== undefined;
  }

  createNewContainer() {
    if (this.selectedImageName === "" || this.selectedImageTag === "") {
      this.messageService.showAlert("JOB.JOB_CREATE_SELECT_IMAGE_TIP", {alertType: "warning", view: this.alertView})
    } else if (this.checkImageAndTag()) {
      this.messageService.showAlert("JOB.JOB_CREATE_SELECT_IMAGE_EXISTS", {alertType: "warning", view: this.alertView})
    } else {
      const container = new JobContainer();
      const jobImage = new JobImage();
      jobImage.image_name = this.selectedImageName;
      jobImage.image_tag = this.selectedImageTag;
      container.image = jobImage;
      this.createSuccess.next(container);
      this.modalOpened = false;
    }
  }
}
