import { Component, OnInit, ViewChild } from '@angular/core';
import { json2String } from 'src/app/shared/tools';
import { Configuration, CfgCardObjects } from './cfg.models';
import { CfgCardsService } from './cfg-cards.service';
import 'src/assets/js/FileSaver.js';
import { Router } from '@angular/router';
import { HttpErrorResponse } from '@angular/common/http';
import { User } from '../account/account.model';

import { OthersComponent } from './others/others.component';
import { ApiserverComponent } from './apiserver/apiserver.component';
import { GogitsComponent } from './gogits/gogits.component';
import { JenkinsComponent } from './jenkins/jenkins.component';
import { KvmComponent } from './kvm/kvm.component';
import { LdapComponent } from './ldap/ldap.component';
import { EmailComponent } from './email/email.component';

declare var saveAs: any;

@Component({
  selector: 'app-cfg-cards',
  templateUrl: './cfg-cards.component.html',
  styleUrls: ['./cfg-cards.component.css']
})
export class CfgCardsComponent implements OnInit {
  config: Configuration;
  cardList: CfgCardObjects;
  applyCfgModal = false;
  user: User;
  loadingFlag = false;
  disableApply = false;

  @ViewChild('others') others: OthersComponent;
  @ViewChild('apiserver') apiserver: ApiserverComponent;
  @ViewChild('gogits') gogits: GogitsComponent;
  @ViewChild('jenkins') jenkins: JenkinsComponent;
  @ViewChild('kvm') kvm: KvmComponent;
  @ViewChild('ldap') ldap: LdapComponent;
  @ViewChild('email') email: EmailComponent;

  constructor(private cfgCardsService: CfgCardsService,
    private router: Router) {
    this.config = new Configuration();
    this.cardList = new CfgCardObjects();
    this.user = new User();
  }

  ngOnInit() {
    this.getCfg();
  }

  getCfg(whichOne?: string) {
    this.cfgCardsService.getConfig(whichOne ? whichOne : '').subscribe(
      (res: Configuration) => {
        this.config = new Configuration(res);
        document.getElementById('container').scrollIntoView();
      },
      (err: HttpErrorResponse) => { this.commonError(err); }
    );
  }

  checkFrom(): string {
    if (!this.others.verifyInputExValid()) {
      return 'others';
    } else if (!this.apiserver.verifyInputExValid()) {
      return 'apiserver';
    } else if (!this.gogits.verifyInputExValid()) {
      return 'gogits';
    } else if (!this.jenkins.verifyInputExValid()) {
      return 'jenkins';
    } else if (!this.kvm.verifyInputExValid()) {
      return 'kvm';
    } else if (this.config.others.authMode == 'ldap_auth' && !this.ldap.verifyInputExValid()) {
      return 'ldap';
    } else if (!this.email.verifyInputExValid()) {
      return 'email';
    } else {
      return '';
    }

  }

  saveCfg() {
    const checkResult = this.checkFrom();
    if (checkResult) {
      document.getElementById(checkResult).scrollIntoView();
    } else {
      this.cfgCardsService.postConfig(this.config).subscribe(
        () => { this.applyCfgModal = true; },
        (err: HttpErrorResponse) => { this.commonError(err); }
      );
    }
  }

  saveAsCfg() {
    let result = [json2String(this.config.PostBody())];
    let file = new File(result, 'board.cfg', { type: 'text/plain;charset=utf-8' });
    saveAs(file);
  }

  applyCfg() {
    this.loadingFlag = true;
    this.disableApply = true;
    this.cfgCardsService.applyCfg(this.user).subscribe(
      () => {
        this.loadingFlag = false;
        this.disableApply = false;
        this.applyCfgModal = false;
        this.router.navigateByUrl('/dashboard');
      },
      (err: HttpErrorResponse) => {
        this.loadingFlag = false;
        this.disableApply = false;
        this.commonError(err);
      }
    );
  }

  cancelApply() {
    this.applyCfgModal = false;
    this.getCfg('tmp');
  }

  commonError(err: HttpErrorResponse) {
    if (err.status === 401) {
      alert('User status error! Please login again!');
      this.router.navigateByUrl('account/login');
    } else {
      alert('Unknown Error!');
    }
  }
}

