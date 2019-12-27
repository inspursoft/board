import { ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { ValidationErrors } from '@angular/forms';
import { HttpErrorResponse } from '@angular/common/http';
import { Observable, of } from 'rxjs';
import { catchError, map } from 'rxjs/operators';
import { HelmChartVersion, IHelmRepo, Questions, QuestionType } from '../helm.type';
import { CsModalChildBase } from '../../shared/cs-modal-base/cs-modal-child-base';
import { Project } from '../../project/project';
import { HelmService } from '../helm.service';
import { MessageService } from '../../shared.service/message.service';
import { AppInitService } from '../../shared.service/app-init.service';
import { GlobalAlertType } from '../../shared/shared.types';

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
  releaseName = '';
  chartValue = '';
  questions: Questions;

  constructor(private helmService: HelmService,
              private appInitService: AppInitService,
              private messageService: MessageService,
              private changeRef: ChangeDetectorRef) {
    super();
    this.projectsList = Array<Project>();
    this.questions = new Questions({});
    this.changeRef.detach();
  }

  ngOnInit(): void {
    this.helmService.getProjects().subscribe(
      (res: Array<Project>) => this.projectsList = res || Array<Project>()
    );
    this.helmService.getChartRelease(this.repoInfo.id, this.chartVersion.name, this.chartVersion.version).subscribe(
      (res: object) => {
        this.chartValue = Reflect.get(res, 'values');
        this.questions = new Questions(Reflect.get(res, 'questions'));
        this.changeRef.reattach();
        this.updateYamlContainer();
      },
      (error: HttpErrorResponse) => {
        this.messageService.showGlobalMessage(error.message, {
          errorObject: error,
          globalAlertType: GlobalAlertType.gatShowDetail
        });
        this.modalOpened = false;
      }
    );
  }

  updateYamlContainer() {
    setTimeout(() => {
      const collection = document.getElementsByClassName('language-yaml');
      if (collection.length > 0) {
        (collection.item(0) as HTMLPreElement).style.margin = '0';
        (collection.item(0) as HTMLPreElement).style.maxHeight = '100%';
      }
    }, 500);
  }

  get checkChartReleaseNameFun() {
    return this.checkChartReleaseName.bind(this);
  }

  checkChartReleaseName(control: HTMLInputElement): Observable<ValidationErrors | null> {
    return this.helmService.checkChartReleaseName(control.value).pipe(
      map(() => null),
      catchError((err: HttpErrorResponse) => {
        this.messageService.cleanNotification();
        if (err.status === 409) {
          return of({nodeGroupExist: 'HELM.RELEASE_CHART_NAME_EXISTING'});
        } else {
          return of(null);
        }
      })
    );
  }

  changeSelectProject(project: Project) {
    this.selectProject = project;
  }

  setAnswer(variable: string, $event: any) {
    const question = this.questions.getQuestionByVariable(variable);
    if (question.questionType === QuestionType.qtBoolean) {
      question.answer = (($event as Event).target as HTMLInputElement).checked;
    } else if (question.questionType === QuestionType.qtString || question.questionType === QuestionType.qtInteger) {
      question.answer = $event;
    }
  }

  chartRelease() {
    if (!this.selectProject) {
      this.messageService.showAlert('HELM.RELEASE_CHART_SELECT_PROJECT_TIP', {
        view: this.alertView,
        alertType: 'warning'
      });
    } else if (this.verifyDropdownExValid() && this.verifyInputExValid()) {
      this.isReleaseWIP = true;
      this.helmService.releaseChartVersion({
        name: this.releaseName,
        chartversion: this.chartVersion.version,
        repository_id: this.repoInfo.id,
        project_id: this.selectProject.project_id,
        owner_id: this.appInitService.currentUser.user_id,
        chart: this.chartVersion.name,
        Answers: this.questions.postAnswers
      }).subscribe(
        () => this.messageService.showAlert('HELM.RELEASE_CHART_RELEASE_SUCCESS'),
        (error: HttpErrorResponse) => {
          this.messageService.showGlobalMessage(error.message, {
            errorObject: error,
            globalAlertType: GlobalAlertType.gatShowDetail
          });
          this.modalOpened = false;
        },
        () => this.modalOpened = false);
    }
  }
}
