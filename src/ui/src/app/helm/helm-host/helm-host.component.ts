import { Component, ComponentFactoryResolver, OnDestroy, OnInit, ViewChild, ViewContainerRef } from '@angular/core';
import { HelmService } from '../helm.service';
import { HelmViewData, HelmViewType } from '../helm.type';
import { ChartListComponent } from '../chart-list/chart-list.component';
import { RepoListComponent } from '../repo-list/repo-list.component';

@Component({
  templateUrl: './helm-host.component.html'
})
export class HelmHostComponent implements OnInit, OnDestroy {
  @ViewChild('host', {read: ViewContainerRef}) hostView: ViewContainerRef;

  constructor(private resolver: ComponentFactoryResolver,
              private helmService: HelmService) {
  }

  ngOnInit(): void {
    this.helmService.viewSubject.asObservable().subscribe((helmViewData: HelmViewData) => {
      switch (helmViewData.type) {
        case HelmViewType.RepoList:
          this.createRepoList();
          return;
        case HelmViewType.ChartList:
          this.createChartList(helmViewData);
          return;
      }
    });
    const repoView = new HelmViewData(HelmViewType.RepoList);
    repoView.description = `HELM.REPO_LIST_TITLE`;
    this.helmService.pushNewView(repoView);
  }

  ngOnDestroy(): void {
    this.helmService.cleanViewData();
  }

  get viewDataList(): Array<HelmViewData> {
    return this.helmService.viewDataList;
  }

  popToView(helmViewData: HelmViewData) {
    this.helmService.popToView(helmViewData);
  }

  createRepoList() {
    this.hostView.clear();
    const factory = this.resolver.resolveComponentFactory(RepoListComponent);
    this.hostView.createComponent(factory);
  }

  createChartList(helmViewData: HelmViewData) {
    this.hostView.clear();
    const factory = this.resolver.resolveComponentFactory(ChartListComponent);
    const component = this.hostView.createComponent(factory);
    component.instance.repoInfo = helmViewData.data;
  }
}
