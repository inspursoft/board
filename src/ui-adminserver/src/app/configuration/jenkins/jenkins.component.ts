import { Component, OnInit, Input } from '@angular/core';
import { cardSlide, rotateNega90 } from 'src/app/shared/animations';
import { initTipsArray } from 'src/app/shared/tools';
import { CfgComponent } from '../cfg-component.template';
import { Jenkins } from 'src/app/shared.service/configuration.model';

@Component({
  selector: 'app-jenkins',
  templateUrl: './jenkins.component.html',
  styleUrls: [
    './jenkins.component.css',
    '../cfg-cards.component.css'
  ],
  animations: [cardSlide, rotateNega90]
})
export class JenkinsComponent extends CfgComponent implements OnInit {
  cfgNum = 8;
  tipsList: Array<boolean>;
  @Input() jenkins: Jenkins;
  patternJenkinsExecutionMode: RegExp = /^((single)|(multi))$/;
  passwordConfirm: string;

  constructor() {
    super();
    this.tipsList = initTipsArray(this.cfgNum, false);
  }

  ngOnInit() {
  }

  onEdit(num: number) {
    this.tipsList = initTipsArray(this.cfgNum, false);
    this.tipsList[num] = true;
  }
}
