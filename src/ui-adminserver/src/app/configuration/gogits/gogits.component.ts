import { Component, OnInit, Input } from '@angular/core';
import { cardSlide, rotateNega90 } from 'src/app/shared/animations';
import { initTipsArray } from 'src/app/shared/tools';
import { CfgComponent } from '../cfg-component.template';
import { Gogits } from 'src/app/shared.service/configuration.model';

@Component({
  selector: 'app-gogits',
  templateUrl: './gogits.component.html',
  styleUrls: [
    './gogits.component.css',
    '../cfg-cards.component.css'
  ],
  animations: [cardSlide, rotateNega90]
})
export class GogitsComponent extends CfgComponent implements OnInit {
  cfgNum = 3;
  tipsList: Array<boolean>;
  @Input() gogits: Gogits;

  constructor() {
    super();
    this.tipsList = initTipsArray(this.cfgNum, false);
    this.gogits = new Gogits();
  }

  ngOnInit() {
  }

  onEdit(num: number) {
    this.tipsList = initTipsArray(this.cfgNum, false);
    this.tipsList[num] = true;
  }
}
