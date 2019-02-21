import { Component } from "@angular/core";
import { IHelmRepoList } from "../helm.type";
import { CsModalChildBase } from "../../shared/cs-modal-base/cs-modal-child-base";
import { HelmService } from "../helm.service";
import { HttpEvent, HttpEventType, HttpProgressEvent } from "@angular/common/http";
import { MessageService } from "../../shared/message-service/message.service";
import { Observable, Subject } from "rxjs";

@Component({
  templateUrl: './upload-chart.component.html',
  styleUrls: ['./upload-chart.component.css']
})
export class UploadChartComponent extends CsModalChildBase {
  repoInfo: IHelmRepoList;
  isUploadChartWIP = false;
  selectedFile: File = null;
  uploadProgressValue: HttpProgressEvent;
  uploadSuccessSubject: Subject<any>;

  constructor(private helmService: HelmService,
              private messageService: MessageService) {
    super();
    this.uploadSuccessSubject = new Subject();
  }

  get uploadSuccessObservable(): Observable<any> {
    return this.uploadSuccessSubject.asObservable();
  }

  selectChartPackage(event: Event) {
    let fileList: FileList = (event.target as HTMLInputElement).files;
    if (fileList.length > 0) {
      this.selectedFile = fileList[0];
    } else {
      this.selectedFile = null;
      (event.target as HTMLInputElement).value = "";
    }
  }

  uploadChart() {
    if (this.selectedFile) {
      this.isUploadChartWIP = true;
      let formData: FormData = new FormData();
      formData.append('upload_file', this.selectedFile, this.selectedFile.name);
      this.helmService.uploadChart(this.repoInfo.id, formData).subscribe((res: HttpEvent<Object>) => {
          if (res.type == HttpEventType.UploadProgress) {
            this.uploadProgressValue = res;
          } else if (res.type == HttpEventType.Response) {
            this.messageService.showAlert('HELM.UPLOAD_CHART_SUCCESS');
            this.isUploadChartWIP = false;
            this.uploadSuccessSubject.next();
            this.modalOpened = false;
          }
        }, () => this.modalOpened = false
      );
    } else {
      this.messageService.showAlert('HELM.UPLOAD_CHART_PACKAGE_TIP', {view: this.alertView, alertType: "alert-warning"});
    }
  }
}