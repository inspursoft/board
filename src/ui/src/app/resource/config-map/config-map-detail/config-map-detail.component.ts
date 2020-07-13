import { Component, OnInit } from '@angular/core';
import { CsModalChildBase } from '../../../shared/cs-modal-base/cs-modal-child-base';
import { ConfigMapDetail } from '../../resource.types';
import { ResourceService } from '../../resource.service';

@Component({
  templateUrl: './config-map-detail.component.html',
  styleUrls: ['./config-map-detail.component.css']
})
export class ConfigMapDetailComponent extends CsModalChildBase implements OnInit {
  configMapDetail: ConfigMapDetail;
  configMapName = '';
  projectName = '';
  isLoadWip = false;

  constructor(private resourceService: ResourceService) {
    super();
    this.configMapDetail = new ConfigMapDetail();
  }

  ngOnInit(): void {
    this.isLoadWip = true;
    this.resourceService.getConfigMapDetail(this.configMapName, this.projectName).subscribe(
      res => this.configMapDetail = res,
      () => this.modalOpened = false,
      () => this.isLoadWip = false
    );
  }
}

