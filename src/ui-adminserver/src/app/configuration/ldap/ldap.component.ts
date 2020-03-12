import { Component, OnInit, Input, OnChanges, SimpleChanges } from '@angular/core';
import { cardSlide, rotateNega90 } from 'src/app/shared/animations';
import { initTipsArray } from 'src/app/shared/tools';
import { CfgComponent } from '../cfg-component.template';
import { Ldap } from '../cfg.models';

@Component({
  selector: 'app-ldap',
  templateUrl: './ldap.component.html',
  styleUrls: [
    './ldap.component.css',
    '../cfg-cards.component.css'
  ],
  animations: [cardSlide, rotateNega90]
})
export class LdapComponent extends CfgComponent implements OnInit {
  cfgNum = 5;
  tipsList: Array<boolean>;
  @Input() ldap: Ldap;
  patternLdapScope: RegExp = /^((LDAP_SCOPE_BASE)|(LDAP_SCOPE_ONELEVEL)|(LDAP_SCOPE_SUBTREE))$/;

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