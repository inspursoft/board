import { Component, OnInit, Input, OnChanges, SimpleChanges } from '@angular/core';
import { cardSlide, rotateNega90 } from 'src/app/shared/animations';
import { initTipsArray } from 'src/app/shared/tools';
import { CfgComponent } from '../cfg-component.template';
import { Others, VerifyPassword } from '../cfg.models';

import * as JsEncryptModule from 'jsencrypt';
import { CfgCardsService } from '../cfg-cards.service';
// import 'src/assets/js/jsencrypt.min.js';
// declare var JSEncrypt: any;

@Component({
  selector: 'app-others',
  templateUrl: './others.component.html',
  styleUrls: [
    './others.component.css',
    '../cfg-cards.component.css'
  ],
  animations: [cardSlide, rotateNega90]
})
export class OthersComponent extends CfgComponent implements OnInit {
  cfgNum = 12;
  tipsList: Array<boolean>;
  @Input() others: Others;
  patternArchType: RegExp = /^((x86_64)|(mips))$/;
  patternAuditDebug: RegExp = /^((true)|(false))$/;
  private publicKey: string;
  @Input() isInit: boolean;
  passwordConfirm: string;
  passwordOld: VerifyPassword;
  showVerify = false;
  passwordVerify = false;
  verifySpinner = false;

  constructor(private cfgCardsService: CfgCardsService) {
    super();
    this.tipsList = initTipsArray(this.cfgNum, false);
    this.passwordOld = new VerifyPassword();
  }

  ngOnInit() { }

  onEdit(num: number) {
    this.tipsList = initTipsArray(this.cfgNum, false);
    this.tipsList[num] = true;
  }

  encryptStr(password: string) {
    this.publicKey = sessionStorage.getItem('pubKey');
    const encrypt = new JsEncryptModule.JSEncrypt();
    encrypt.setPublicKey(this.publicKey);
    const encrypted = encrypt.encrypt(password.trim());
    this.others.boardAdminPassword = encrypted;
  }

  verifyPassword(key:string, oldPwd: string) {
    this.showVerify = true;
    this.verifySpinner = true;

    this.passwordOld.which = key;
    this.publicKey = sessionStorage.getItem('pubKey');
    const encrypt = new JsEncryptModule.JSEncrypt();
    encrypt.setPublicKey(this.publicKey);
    this.passwordOld.value = encrypt.encrypt(oldPwd);
    this.cfgCardsService.verifyPassword(this.passwordOld).subscribe(
      (res) => {
        this.passwordVerify = false;
        if (res === 'success') {
          this.passwordVerify = true;
        }
        this.verifySpinner = false;
      },
      () => {
        this.verifySpinner = false;
      }
    );
  }
}
