import { Component, OnInit, Input } from '@angular/core';
import { cardSlide, rotateNega90 } from 'src/app/shared/animations';
import { initTipsArray } from 'src/app/shared/tools';
import { CfgComponent } from '../cfg-component.template';
import { ApiServer } from '../cfg.models';

@Component({
  selector: 'app-apiserver',
  templateUrl: './apiserver.component.html',
  styleUrls: [
    './apiserver.component.css',
    '../cfg-cards.component.css'
  ],
  animations: [cardSlide, rotateNega90]
})
export class ApiserverComponent extends CfgComponent implements OnInit {
  cfgNum = 8;
  tipsList: Array<boolean>;
  @Input() apiserver: ApiServer;
  patternHttpScheme: RegExp = /^http(s?)$/;

  constructor() {
    super();
    this.tipsList = initTipsArray(this.cfgNum, false);
  }

  ngOnInit() {
  }

  onEdit(num: number) {
    this.tipsList = initTipsArray(this.cfgNum, false);
    this.tipsList[num] = true;

    // console.log(this.verifyInputExValid());
    if (this.verifyInputExValid()) {
      alert('error');
    }
  }
}
