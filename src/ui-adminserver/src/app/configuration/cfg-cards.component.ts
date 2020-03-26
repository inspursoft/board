import { Component, OnInit, ViewChildren, QueryList } from '@angular/core';
import { json2String } from 'src/app/shared/tools';
import { Configuration, CfgCardObjects } from './cfg.models';
import { CfgCardsService } from './cfg-cards.service';
import 'src/assets/js/FileSaver.js';
import { InputExComponent } from 'board-components-library';
import { Router } from '@angular/router';
import { HttpErrorResponse } from '@angular/common/http';
import { User } from '../account/account.model';

declare var saveAs: any;

@Component({
  selector: 'app-cfg-cards',
  templateUrl: './cfg-cards.component.html',
  styleUrls: ['./cfg-cards.component.css']
})
export class CfgCardsComponent implements OnInit {
  config: Configuration;
  cardList: CfgCardObjects;
  @ViewChildren(InputExComponent) inputExComponents: QueryList<InputExComponent>;
  applyCfgModal = false;
  user: User;

  constructor(private cfgCardsService: CfgCardsService,
    private router: Router) {
    this.config = new Configuration();
    this.cardList = new CfgCardObjects();
    this.user = new User();
  }

  ngOnInit() {
    this.getCfg();
    // this.cfgCardsService.getPubKey().subscribe((res: string) => {
    //   sessionStorage.setItem('pubKey', res);
    // });
    // this.route.data
    //   .subscribe((data: { configuration: Configuration }) => {
    //     this.config = data.configuration;
    //     console.log(data.configuration);
    //   });
  }

  getCfg(whichOne?: string) {
    this.cfgCardsService.getConfig(whichOne ? whichOne : '').subscribe(
      (res: Configuration) => { this.config = new Configuration(res); },
      (err: HttpErrorResponse) => { this.commonError(err); }
    );
  }

  saveCfg() {
    console.log('' + this.verifyInputExValid() + this.inputExComponents.toArray().length);
    this.cfgCardsService.postConfig(this.config).subscribe(
      // if response Status Code is 200: success
      () => {
        // alert('apply success!');
        // location.reload();
        // window.scrollTo({
        //   top: 0
        // });
        this.applyCfgModal = true;
      },
      // if error
      (err: HttpErrorResponse) => { this.commonError(err); }
    );
  }

  saveAsCfg() {
    let result = [json2String(this.config.PostBody())];
    let file = new File(result, 'board.cfg', { type: 'text/plain;charset=utf-8' });
    saveAs(file);
  }

  verifyInputExValid(): boolean {
    return this.inputExComponents.toArray().every((component: InputExComponent) => {
      if (!component.isValid && !component.inputDisabled) {
        component.checkSelf();
      }
      return component.isValid || component.inputDisabled;
    });
  }

  applyCfg() {
    this.cfgCardsService.applyCfg(this.user).subscribe(
      () => {
        this.applyCfgModal = false;
        this.router.navigateByUrl('/dashboard');
      },
      (err: HttpErrorResponse) => { this.commonError(err); }
    );
  }

  cancelApply() {
    this.applyCfgModal = false;
    window.scrollTo({
      top: 0,
      behavior: 'smooth'
    });
    window.location.reload();
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

