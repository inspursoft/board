import { ChangeDetectorRef, Component, OnInit } from "@angular/core";
import { ValidationErrors } from "@angular/forms";
import { HttpErrorResponse } from "@angular/common/http";
import { Observable, of } from "rxjs";
import { catchError, map } from "rxjs/operators";
import { HelmChartVersion, IHelmRepo } from "../helm.type";
import { CsModalChildBase } from "../../shared/cs-modal-base/cs-modal-child-base";
import { Project } from "../../project/project";
import { HelmService } from "../helm.service";
import { MessageService } from "../../shared.service/message.service";
import { AppInitService } from "../../shared.service/app-init.service";

@Component({
  templateUrl: './chart-release.component.html',
  styleUrls: ['./chart-release.component.css']
})
export class ChartReleaseComponent extends CsModalChildBase implements OnInit {
  repoInfo: IHelmRepo;
  chartVersion: HelmChartVersion;
  projectsList: Array<Project>;
  selectProject: Project = null;
  isReleaseWIP = false;
  isCheckNameWip = false;
  releaseName = '';
  chartValue = '';

  constructor(private helmService: HelmService,
              private appInitService: AppInitService,
              private messageService: MessageService,
              private changeRef: ChangeDetectorRef) {
    super();
    this.changeRef.detach();
    this.projectsList = Array<Project>();
  }

  ngOnInit(): void {
    this.helmService.getProjects().subscribe(
      (res: Array<Project>) => this.projectsList = res || Array<Project>()
    );
    this.helmService.getChartRelease(this.repoInfo.id, this.chartVersion.name, this.chartVersion.version).subscribe(
      (res: Object) => this.chartValue = res['values'], null,
      () => this.changeRef.reattach()
    );
  }

  get checkChartReleaseNameFun() {
    return this.checkChartReleaseName.bind(this);
  }

  checkChartReleaseName(control: HTMLInputElement): Observable<ValidationErrors | null> {
    this.isCheckNameWip = true;
    return this.helmService.checkChartReleaseName(control.value)
      .pipe(
        map(() => {
          setTimeout(() => this.isCheckNameWip = false);
          return null;
        }),
        catchError((err: HttpErrorResponse) => {
          this.messageService.cleanNotification();
          setTimeout(() => this.isCheckNameWip = false);
          if (err.status == 409) {
            return of({nodeGroupExist: "HELM.RELEASE_CHART_NAME_EXISTING"})
          } else {
            return of(null)
          }
        }));
  }

  changeSelectProject(project: Project) {
    this.selectProject = project;
  }

  chartRelease() {
    if (!this.selectProject) {
      this.messageService.showAlert('HELM.RELEASE_CHART_SELECT_PROJECT_TIP', {view: this.alertView, alertType: "warning"})
    } else if (this.verifyInputValid()) {
      this.isReleaseWIP = true;
      this.helmService.releaseChartVersion({
        name: this.releaseName,
        chartVersion: this.chartVersion.version,
        repoId: this.repoInfo.id,
        projectId: this.selectProject.project_id,
        ownerId: this.appInitService.currentUser.user_id,
        chart: this.chartVersion.name
      }).subscribe(() => {
        this.modalOpened = false;
        this.messageService.showAlert('HELM.RELEASE_CHART_RELEASE_SUCCESS')
      }, () => this.modalOpened = false)
    }
  }
}
