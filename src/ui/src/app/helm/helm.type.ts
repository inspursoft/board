import { HttpBase, HttpBind, HttpBindArray, HttpBindObject, Pagination } from '../shared/ui-model/model-types';

export enum QuestionType {
  qtUnknown, qtBoolean, qtString, qtInteger
}

export enum HelmViewType {
  RepoList, ChartList,
}

export enum ViewMethod {List = 'list', Card = 'card'}

export class ReleaseFile extends HttpBase {
  @HttpBind('name') name = '';
  @HttpBind('contents') contents = '';
}

export class ReleaseTemplate extends HttpBase {
  @HttpBind('name') name = '';
  @HttpBind('contents') contents = '';
}

export class ReleaseMetadata extends HttpBase {
  @HttpBind('description') description = '';
  @HttpBind('icon') icon = '';
  @HttpBind('version') version = '';
  @HttpBind('name') name = '';
}

export class Question extends HttpBase {
  @HttpBind('default') default: string;
  @HttpBind('label') label: string;
  @HttpBind('show_subquestion_if') showSubQuestion: string;
  @HttpBind('type') type: string;
  @HttpBind('variable') variable: string;
  @HttpBind('description') description: string;
  @HttpBindArray('subquestions', Question) subQuestions: Array<Question>;
  answerValue = '';

  protected prepareInit() {
    this.subQuestions = Array<Question>();
  }

  set answer(value: any) {
    if (this.questionType === QuestionType.qtBoolean) {
      this.answerValue = value ? 'true' : 'false';
    } else if (this.questionType === QuestionType.qtInteger) {
      this.answerValue = Number(value).toString();
    } else if (this.questionType === QuestionType.qtString) {
      this.answerValue = value;
    }
  }

  get defaultValue(): any {
    if (this.questionType === QuestionType.qtBoolean) {
      return this.default === 'true';
    } else if (this.questionType === QuestionType.qtInteger) {
      return this.default;
    } else if (this.questionType === QuestionType.qtString) {
      return this.default;
    }
  }

  get isShowSubQuestion(): boolean {
    return this.answerValue === this.showSubQuestion;
  }

  get isHasSubQuestion(): boolean {
    return this.subQuestions.length > 0;
  }

  get questionType(): QuestionType {
    if (this.type === 'boolean') {
      return QuestionType.qtBoolean;
    } else if (this.type === 'string') {
      return QuestionType.qtString;
    } else if (this.type === 'int') {
      return QuestionType.qtInteger;
    } else {
      return QuestionType.qtUnknown;
    }
  }
}

export class ChartRelease extends HttpBase {
  @HttpBindArray('files', ReleaseFile) files: Array<ReleaseFile>;
  @HttpBindArray('questions', Question) questions: Array<Question>;
  @HttpBindArray('templates', ReleaseTemplate) templates: Array<ReleaseTemplate>;
  @HttpBindObject('metadata', ReleaseMetadata) metadata: ReleaseMetadata;
  @HttpBind('values') values = '';

  protected prepareInit() {
    this.files = new Array<ReleaseFile>();
    this.metadata = new ReleaseMetadata();
    this.questions = new Array<Question>();
    this.templates = new Array<ReleaseTemplate>();
  }

  get postAnswers(): { [key: string]: string } {
    const result: { [key: string]: string } = {};
    this.questions.forEach(question => {
      if (question.answerValue !== '' && question.answerValue !== question.default) {
        Reflect.set(result, question.variable, question.answerValue);
        if (question.isHasSubQuestion) {
          question.subQuestions.forEach(subQuestion => {
            if (subQuestion.answerValue !== '' && subQuestion.answerValue !== subQuestion.default) {
              Reflect.set(result, subQuestion.variable, subQuestion.answerValue);
            }
          });
        }
      }
    });
    return result;
  }

  getQuestionByVariable(variable: string): Question {
    let result: Question;
    result = this.questions.find(value => value.variable === variable);
    if (result === undefined) {
      this.questions.forEach(value => {
        if (value.isHasSubQuestion && result === undefined) {
          result = value.subQuestions.find(subValue => subValue.variable === variable);
        }
      });
    }
    return result;
  }
}


export interface IHelmRepo {
  id: number;
  name: string;
  url: string;
  type: number;
}

export interface IChartReleasePost {
  name: string;
  project_id: number;
  repository_id: number;
  chart: string;
  chartversion: string;
  owner_id: number;
  Answers: { [index: string]: string };
  values: string;
}

export interface IChartRelease {
  id: number;
  name: string;
  project_id: number;
  project_name: string;
  repository_id: number;
  repository: string;
  chart: string;
  chartversion: string;
  owner_id: number;
  owner_name: string;
  status: string;
  values: string;
  update_time: string;
  creation_time: string;
}

export interface IChartReleaseDetail {
  id: number;
  name: string;
  project_id: number;
  project_name: string;
  repository_id: number;
  repository: string;
  chart: string;
  chartversion: string;
  owner_id: number;
  owner_name: string;
  status: string;
  values: string;
  update_time: string;
  creation_time: string;
  notes: string;
  workloads: string;
  workloadstatus: string;
}

export class HelmChartVersion extends HttpBase {
  @HttpBind('name') name = '';
  @HttpBind('version') version = '';
  @HttpBind('description') description = '';
  @HttpBind('urls') urls: Array<string>;
  @HttpBind('digest') digest = '';
  @HttpBind('icon') icon = '';

  protected prepareInit() {
    this.urls = Array<string>();
  }
}

export class HelmChart extends HttpBase {
  @HttpBind('name') name = '';
  @HttpBindArray('versions', HelmChartVersion) versions: Array<HelmChartVersion>;

  protected prepareInit() {
    this.versions = Array<HelmChartVersion>();
  }
}

export class HelmRepoDetail extends HttpBase {
  @HttpBind('id') id = 0;
  @HttpBind('name') name = '';
  @HttpBind('url') url = '';
  @HttpBind('type') type = 0;
  @HttpBindObject('pagination', Pagination) pagination: Pagination;
  @HttpBindArray('charts', HelmChart) charts: Array<HelmChart>;

  protected prepareInit() {
    this.pagination = new Pagination();
    this.pagination.PageCount = 1;
    this.pagination.PageSize = 15;
    this.charts = Array<HelmChart>();
  }

  get versionList(): Array<HelmChartVersion> {
    const list = Array<HelmChartVersion>();
    this.charts.forEach((chart: HelmChart) => list.push(...chart.versions));
    return list;
  }
}

export class HelmViewData {
  description = '';

  constructor(public type: HelmViewType, public data: any = null) {

  }
}
