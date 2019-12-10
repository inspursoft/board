import { IPagination } from '../shared/shared.types';
import { HttpBind, ResponseArrayBase, ResponseBase } from '../shared.service/shared-model-types';
import YAML from 'yaml';

export enum QuestionType {
  qtUnknown, qtBoolean, qtString, qtInteger
}

export enum HelmViewType {
  RepoList, ChartList,
}

export enum ViewMethod {List = 'list', Card = 'card'}

export class ChartFiles extends ResponseArrayBase<ChartFile> {
  CreateOneItem(res: object): ChartFile {
    let chartFile = new ChartFile(res);
    if (chartFile.isQuestionFile) {
      chartFile = new QuestionsChartFile(res);
      chartFile.parseYamlFile();
    }
    return chartFile;
  }

  constructor(protected res: object) {
    super(res);
  }

  get questionsChartFile(): QuestionsChartFile {
    return this.data.find(value => value.isQuestionFile) as QuestionsChartFile;
  }
}

export class ChartFile extends ResponseBase {
  @HttpBind('name') fileName: string;
  @HttpBind('contents') fileContents: string;

  parseYamlFile() {

  }

  get isQuestionFile(): boolean {
    return this.fileName.includes('question');
  }
}

export class QuestionsChartFile extends ChartFile {
  questions: Array<Question>;

  constructor(res: object) {
    super(res);
    this.questions = Array<Question>();
  }

  parseYamlFile() {
    const yml = YAML.parse(this.fileContents);
    (Reflect.get(yml, 'questions') as Array<object>).forEach(value => {
      const question = new Question(value);
      this.questions.push(question);
    });
  }

  get postAnwsers(): { [key: string]: string } {
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

export class Question extends ResponseBase {
  @HttpBind('default') default: string;
  @HttpBind('label') label: string;
  @HttpBind('show_subquestion_if') showSubQuestion: string;
  @HttpBind('type') type: string;
  @HttpBind('variable') variable: string;
  @HttpBind('description') description: string;

  subQuestions: Array<Question>;
  answerValue = '';

  constructor(res: object) {
    super(res);
    this.subQuestions = Array<Question>();
    if (Reflect.has(res, 'subquestions')) {
      (Reflect.get(res, 'subquestions') as Array<object>).forEach(value => {
        const question = new Question(value);
        this.subQuestions.push(question);
      });
    }
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


export interface IHelmRepo {
  id: number;
  name: string;
  url: string;
  type: number;
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

export class HelmChartVersion {
  name = '';
  version = '';
  description = '';
  urls: Array<string>;
  digest: string;
  icon: string;

  constructor() {
    this.urls = Array<string>();
  }

  static newFromServe(response: object): HelmChartVersion {
    const version = new HelmChartVersion();
    version.name = Reflect.get(response, 'name');
    version.version = Reflect.get(response, 'version');
    version.description = Reflect.get(response, 'description');
    version.urls = Reflect.get(response, 'urls');
    version.digest = Reflect.get(response, 'digest');
    if (Reflect.has(response, 'icon')) {
      version.icon = Reflect.get(response, 'icon');
    }
    return version;
  }
}

export class HelmChart {
  name = '';
  versions: Array<HelmChartVersion>;

  constructor() {
    this.versions = Array<HelmChartVersion>();
  }

  static newFromServe(response: object): HelmChart {
    const chart = new HelmChart();
    chart.name = Reflect.get(response, 'name');
    const resVersions: Array<object> = Reflect.get(response, 'versions');
    resVersions.forEach((resVersion: object) => {
      const version = HelmChartVersion.newFromServe(resVersion);
      chart.versions.push(version);
    });
    return chart;
  }
}

export class HelmRepoDetail {
  baseInfo: IHelmRepo;
  pagination: IPagination;
  charts: Array<HelmChart>;

  constructor() {
    this.charts = Array<HelmChart>();
    this.pagination = {page_count: 1, page_index: 0, page_size: 15, total_count: 0};
  }

  get versionList(): Array<HelmChartVersion> {
    const list = Array<HelmChartVersion>();
    this.charts.forEach((chart: HelmChart) => list.push(...chart.versions));
    return list;
  }

  static newFromServe(response: object): HelmRepoDetail {
    const detail = new HelmRepoDetail();
    detail.baseInfo = {
      id: Reflect.get(response, 'id'),
      name: Reflect.get(response, 'name'),
      url: Reflect.get(response, 'url'),
      type: Reflect.get(response, 'type')
    };
    if (Reflect.has(response, 'pagination')) {
      detail.pagination = Reflect.get(response, 'pagination');
    }
    if (Reflect.has(response, 'charts')) {
      const resCharts: Array<object> = Reflect.get(response, 'charts');
      resCharts.forEach((resChart: object) => {
        const chart = HelmChart.newFromServe(resChart);
        detail.charts.push(chart);
      });
    }
    return detail;
  }
}

export class HelmViewData {
  description = '';

  constructor(public type: HelmViewType, public data: any = null) {

  }
}
