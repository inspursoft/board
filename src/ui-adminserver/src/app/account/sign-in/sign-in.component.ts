import { Component, OnInit, ViewChild } from '@angular/core';
import { User, UserVerify, DBInfo } from '../account.model';
import { AccountService } from '../account.service';
import { Router } from '@angular/router';
import { ClrWizard, ClrModal } from '@clr/angular';
import { HeaderComponent } from 'src/app/shared/header/header.component';

@Component({
  selector: 'app-sign-in',
  templateUrl: './sign-in.component.html',
  styleUrls: ['./sign-in.component.css']
})
export class SignInComponent implements OnInit {
  @ViewChild('wizard') wizard: ClrWizard;
  @ViewChild('modal') modal: ClrModal;
  openWizard = true;
  loadingFlag = false;
  errorFlag = false;
  errorVerifyFlag = false;
  errorDBMaxFlag = false;

  showInitDB = true;
  showInitSSH = true;

  disableDbPwdForm = true;
  disableAccountFrom = true;

  uuid = '';
  dbInfo: DBInfo;
  sshAccount: UserVerify;
  account: UserVerify;

  disableInput = false;

  current = 0;
  isCurrent = true;

  user: User;

  constructor(private accountService: AccountService,
    private header: HeaderComponent,
    private router: Router) {
    this.account = new UserVerify();
    this.account.username = 'admin';
    this.sshAccount = new UserVerify();
    this.user = new User();
    this.dbInfo = new DBInfo();
  }

  ngOnInit() {
    this.accountService.checkInit().subscribe(
      (res: string) => {
        const step = res.toString().toLowerCase();
        if (step === 'no') {
          this.openWizard = false;
          this.checkDBAlive();
        } else if (step === 'step2') {
          this.showInitDB = false;
        } else if (step === 'step3') {
          this.showInitDB = false;
          this.showInitSSH = false;
        }
      }
    );

    // for test
    // let step = 'step3';
    // if (step === 'no') {
    //   this.openWizard = false;
    //   // this.checkDBAlive();
    // } else if (step === 'step2') {
    //   this.showInitDB = false;
    // } else if (step === 'step3') {
    //   this.showInitDB = false;
    //   this.showInitSSH = false;
    // }
  }

  public onTranslate(): void {
    this.disableInput = true;
    this.loadingFlag = true;
    const currentLang = window.localStorage.getItem('currentLang');
    const trans = (currentLang === 'en' || currentLang === 'en-us') ? 'zh-cn' : 'en';
    window.localStorage.setItem('currentLang', trans);
    this.header.changLanguage(trans);
  }

  onWelcome(): void {
    this.waitingFlag(true);
    this.accountService.createUUID().subscribe(
      () => {
        this.wizard.forceNext();
        this.successFlag(this.wizard.currentPage._id);
        this.checkBtn();
      },
      () => { this.waitingFlag(false); },
    );

    // for test
    // this.waitingFlag(true);
    // setTimeout(() => {
    //   this.wizard.forceNext();
    //   this.successFlag(this.wizard.currentPage._id);
    //   this.checkBtn();
    // }, 1000);
  }

  onVerify(): void {
    this.waitingFlag(true);
    this.accountService.validateUUID(this.uuid).subscribe(
      () => {
        this.wizard.forceNext();
        this.successFlag(this.wizard.currentPage._id);
        this.checkBtn();
      },
      () => { this.waitingFlag(false); },
    );

    // for test
    // this.waitingFlag(true);
    // setTimeout(() => {
    //   if (this.uuid === '42') {
    //     this.wizard.forceNext();
    //     this.successFlag(this.wizard.currentPage._id);
    //     this.checkBtn();
    //   } else {
    //     this.waitingFlag(false);
    //   }
    // }, 1000);
  }

  verifyDBPwd() {
    this.errorVerifyFlag = false;
    if (this.dbInfo.verify()) {
      this.disableDbPwdForm = false;
    } else {
      this.disableDbPwdForm = true;
      if (this.dbInfo.password && this.dbInfo.passwordConfirm) {
        this.errorVerifyFlag = true;
      }
    }
  }

  onInitDB(): void {
    this.errorDBMaxFlag = false;
    this.waitingFlag(true);
    if (this.dbInfo.maxConnection < 10 || this.dbInfo.maxConnection > 16384) {
      this.disableInput = false;
      this.loadingFlag = false;
      this.errorDBMaxFlag = true;
    } else {
      this.accountService.initDB(this.dbInfo).subscribe(
        () => {
          this.wizard.forceNext();
          this.successFlag(this.wizard.currentPage._id);
          this.checkBtn();
        },
        () => { this.waitingFlag(false); },
      );
    }

    // for test
    // if (this.dbInfo.maxConnection < 10 || this.dbInfo.maxConnection > 16384) {
    //   this.disableInput = false;
    //   this.loadingFlag = false;
    //   this.errorDBMaxFlag = true;
    // } else {
    //   setTimeout(() => {
    //     if (this.dbInfo.verify()) {
    //       this.wizard.forceNext();
    //       this.successFlag(this.wizard.currentPage._id);
    //       this.checkBtn();
    //     } else {
    //       this.waitingFlag(false);
    //     }
    //   }, 1000);
    // }
    // console.log(this.dbInfo);
  }

  onInitSSH(): void {
    this.waitingFlag(true);
    this.accountService.initSSH(this.sshAccount).subscribe(
      () => {
        this.wizard.forceNext();
        this.modal.closable = true;
        this.modal.close();
        this.successFlag(this.wizard.currentPage._id);
        this.checkBtn();
      },
      () => { this.waitingFlag(false); },
    );

    // for test
    // setTimeout(() => {
    //   if (this.sshAccount.username === '1') {
    //     this.wizard.forceNext();
    //     this.modal.closable = true;
    //     this.modal.close();
    //     this.successFlag(this.wizard.currentPage._id);
    //     this.checkBtn();
    //   } else {
    //     this.waitingFlag(false);
    //   }
    // }, 1000);
  }

  verifyAccountPwd() {
    this.errorVerifyFlag = false;
    if (this.account.password === this.account.passwordConfirm) {
      this.disableAccountFrom = false;
    } else {
      this.disableAccountFrom = true;
      if (this.account.password && this.account.passwordConfirm) {
        this.errorVerifyFlag = true;
      }
    }
  }

  onInitAccount(): void {
    this.waitingFlag(true);
    this.accountService.postSignUp(this.account).subscribe(
      () => {
        this.wizard.forceFinish();
        this.loadingFlag = false;
        this.disableInput = false;
      },
      () => { this.waitingFlag(false); },
    );

    // for test
    // setTimeout(() => {
    //   if (this.account.password === '11111111') {
    //     this.wizard.forceFinish();
    //     this.loadingFlag = false;
    //     this.disableInput = false;
    //   } else {
    //     this.waitingFlag(false);
    //   }
    // }, 1000);
  }

  checkDBAlive() {
    this.accountService.checkDB().subscribe(
      () => { },
      () => {
        this.modal.closable = false;
        this.modal.open();
      }
    );

    // for test
    // this.modal.closable = false;
    // this.modal.open();
  }


  signIn() {
    console.log(this.user.username + this.user.password);
    // test
    // window.sessionStorage.setItem('token', `username=${this.user.username}&password=${this.user.password}`);
    // this.router.navigateByUrl('dashboard');

    // TODO
    this.accountService.postSignIn(this.user).subscribe(
      (res: string) => {
        if (res) {
          window.sessionStorage.setItem('token', res);
          this.router.navigateByUrl('dashboard');
        } else {
          alert('Unknown Error!');
        }
      },
      () => {
        alert('账号或密码错误！# Account or password error!');
      }
    );
  }

  checkBtn() {
    if (this.wizard.currentPage._id == this.current) {
      this.isCurrent = true;
    } else {
      this.isCurrent = false;
    }
  }

  waitingFlag(flag: boolean) {
    if (flag) {
      this.errorFlag = false;
      this.loadingFlag = true;
      this.disableInput = true;
    } else {
      this.errorFlag = true;
      this.loadingFlag = false;
      this.disableInput = false;
    }
  }

  successFlag(id: number) {
    this.loadingFlag = false;
    this.disableInput = false;
    this.current = id;
  }
}
