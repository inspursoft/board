import { Component, OnInit, ViewChild } from '@angular/core';
import { json2String } from 'src/app/shared/tools';
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
import { MessageService } from '../shared/message/message.service';
import { ConfigurationService } from '../shared.service/configuration.service';
import { Configuration } from '../shared.service/configuration.model';
import { CfgCardObjects } from './cfg.model';

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

  constructor(private configurationService: ConfigurationService,
              private messageService: MessageService,
              private router: Router) {
    this.config = new Configuration();
    this.cardList = new CfgCardObjects();
    this.user = new User();
  }

  ngOnInit() {
    this.getCfg();
  }

  getCfg(whichOne?: string) {
    this.configurationService.getConfig(whichOne ? whichOne : '').subscribe(
      (res: Configuration) => {
        this.config = new Configuration(res);
        document.getElementById('container').scrollIntoView();
      },
      (err: HttpErrorResponse) => { this.commonError(err); }
    );
  }

  saveAsCfg() {
    let result = [json2String(this.config.PostBody())];
    let file = new File(result, 'board.cfg', { type: 'text/plain;charset=utf-8' });
    saveAs(file);
  }

  commonError(err: HttpErrorResponse) {
    if (err.status === 401) {
      this.messageService.showOnlyOkDialog('ACCOUNT.TOKEN_ERROR', 'ACCOUNT.ERROR');
      this.router.navigateByUrl('account/login');
    } else {
      this.messageService.showOnlyOkDialog('ERROR.HTTP_UNK', 'ACCOUNT.ERROR');
    }
  }
}

