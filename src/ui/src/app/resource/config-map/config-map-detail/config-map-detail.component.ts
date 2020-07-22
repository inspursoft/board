import { Component, OnInit } from '@angular/core';
import { CsModalChildBase } from '../../../shared/cs-modal-base/cs-modal-child-base';
import { ResourceService } from '../../resource.service';
import { SharedConfigMapDetail } from '../../../shared/shared.types';

@Component({
  templateUrl: './config-map-detail.component.html',
  styleUrls: ['./config-map-detail.component.css']
})
export class ConfigMapDetailComponent extends CsModalChildBase implements OnInit {
  configMapDetail: SharedConfigMapDetail;
  configMapName = '';
  projectName = '';
  isLoadWip = false;

  constructor(private resourceService: ResourceService) {
    super();
    this.configMapDetail = new SharedConfigMapDetail();
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

