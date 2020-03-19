import { Component, OnInit, Input, OnChanges, SimpleChanges } from '@angular/core';
import { cardSlide, rotateNega90 } from 'src/app/shared/animations';
import { initTipsArray } from 'src/app/shared/tools';
import { CfgComponent } from '../cfg-component.template';
import { Email } from '../cfg.models';

@Component({
  selector: 'app-email',
  templateUrl: './email.component.html',
  styleUrls: [
    './email.component.css',
    '../cfg-cards.component.css'
  ],
  animations: [cardSlide, rotateNega90]
})
export class EmailComponent extends CfgComponent implements OnInit {
  cfgNum = 7;
  tipsList: Array<boolean>;
  @Input() email: Email;
  patternEmailSsl: RegExp = /^((true)|(false))$/;

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
