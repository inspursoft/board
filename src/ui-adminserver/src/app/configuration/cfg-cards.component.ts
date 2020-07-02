import { Component, OnInit } from '@angular/core';
import { json2String } from 'src/app/shared/tools';
import 'src/assets/js/FileSaver.js';
import { Router } from '@angular/router';
import { HttpErrorResponse } from '@angular/common/http';
import { User } from '../account/account.model';

import { MessageService } from '../shared/message/message.service';
import { ConfigurationService } from '../shared.service/configuration.service';
import { Configuration } from '../shared.service/cfg.model';

declare var saveAs: any;

@Component({
  selector: 'app-cfg-cards',
  templateUrl: './cfg-cards.component.html',
  styleUrls: ['./cfg-cards.component.css']
})
export class CfgCardsComponent implements OnInit {
  config: Configuration;
  applyCfgModal = false;
  user: User;
  loadingFlag = false;
  disableApply = false;

  showBaselineHelper = false;
  newDate = new Date('2016-01-01 09:00:00');

  constructor(private configurationService: ConfigurationService,
              private messageService: MessageService,
              private router: Router) {
    this.config = new Configuration();
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
      console.error(err.message);
      this.messageService.showOnlyOkDialog('ERROR.HTTP_UNK', 'ACCOUNT.ERROR');
    }
  }

  onFocusBaselineHelper() {
    this.showBaselineHelper = true;
  }

  onBlurBaselineHelper() {
    this.showBaselineHelper = false;
    const year = this.newDate.getFullYear();
    const month = this.newDate.getMonth() + 1;
    const day = this.newDate.getDate();
    this.config.k8s.imageBaselineTime = '' + year + '-' + month + '-' + day + ' 00:00:00';
  }
}

