import { Component, OnInit } from "@angular/core";
import { HelmService } from "../helm.service";
import { HelmRepoDetail, HelmViewData, HelmViewType, IHelmRepo } from "../helm.type";

@Component({
  selector: 'helm-repo-list',
  templateUrl: './repo-list.component.html',
  styleUrls: ["./repo-list.component.css"]
})
export class RepoListComponent implements OnInit {
  loadingWIP = false;
  repoList: Array<IHelmRepo>;

  constructor(private helmService: HelmService) {
    this.repoList = Array<IHelmRepo>();
  }

  ngOnInit() {
    this.retrieve();
  }

  retrieve() {
    this.loadingWIP = true;
    this.helmService.getRepoList().subscribe((res: Array<IHelmRepo>) => {
        this.loadingWIP = false;
        this.repoList = res;
      }, () => this.loadingWIP = false
    );
  }

  showRepoDetail(repo: IHelmRepo) {
    let viewData = new HelmViewData(HelmViewType.ChartList, repo);
    viewData.description = `HELM.CHART_LIST_TITTLE`;
    this.helmService.pushNewView(viewData);
  }
}
