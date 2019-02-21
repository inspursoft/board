import { Component, OnInit } from "@angular/core";
import { HelmService } from "../helm.service";
import { HelmRepoDetail, HelmViewData, HelmViewType, IHelmRepoList } from "../helm.type";

@Component({
  selector: 'helm-repo-list',
  templateUrl: './repo-list.component.html',
  styleUrls: ["./repo-list.component.css"]
})
export class RepoListComponent implements OnInit {
  loadingWIP = false;
  repoList: Array<IHelmRepoList>;

  constructor(private helmService: HelmService) {

  }

  ngOnInit() {
    this.retrieve();
  }

  retrieve() {
    this.loadingWIP = true;
    this.helmService.getRepoList().subscribe((res: Array<IHelmRepoList>) => {
        this.loadingWIP = false;
        this.repoList = res;
      }, () => this.loadingWIP = false
    );
  }

  showRepoDetail(repo: IHelmRepoList) {
    // this.loadingWIP = true;
    // this.helmService.getRepoDetail(repo.id).subscribe((res: Object) => {
    //   let repoDetail = HelmRepoDetail.newFromServe(res);
    //   let viewData = new HelmViewData(repoDetail, HelmViewType.hvtChartList);
    //   viewData.description = `HELM.CHART_LIST_TITTLE`;
    //   this.helmService.pushNewView(viewData);
    // });


    let viewData = new HelmViewData(repo, HelmViewType.hvtChartList);
    viewData.description = `HELM.CHART_LIST_TITTLE`;
    this.helmService.pushNewView(viewData);
  }
}