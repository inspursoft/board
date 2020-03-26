import { Component, OnInit, ViewChild } from '@angular/core';
import { User, UserVerify, DBInfo } from '../account.model';
import { AccountService } from '../account.service';
import { Router } from '@angular/router';
import { ClrWizard, ClrModal } from '@clr/angular';

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

  showInitDB = true;
  showInitSSH = true;

  uuid = '';
  dbInfo: DBInfo;
  errorDBMax = false;
  disableDbPwdForm = true;
  sshAccount: UserVerify;
  disableAccountFrom = true;
  user: User;
  account: UserVerify;

  dbAlive = false;
  disableInput = false;

  current = 0;
  isCurrent = true;

  constructor(private accountService: AccountService,
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
        const step = res.toLowerCase();
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
    // let step = 'no';
    // if (step === 'no') {
    //   this.openWizard = false;
    //   this.checkDBAlive();
    // } else if (step === 'step2') {
    //   this.showInitDB = false;
    // } else if (step === 'step3') {
    //   this.showInitDB = false;
    //   this.showInitSSH = false;
    // }
  }

  onWelcome(): void {
    this.accountService.createUUID().subscribe(
      () => { this.successFlag(this.wizard.currentPage._id); },
      () => { this.waitingFlag(false); },
    );

    // for test
    // setTimeout(() => {
    //   this.wizard.forceNext();
    //   this.loadingFlag = false;
    //   this.disableInput = false;
    //   this.current = this.wizard.currentPage._id;
    //   this.checkBtn();
    // }, 1000);
  }

  onVerify(): void {
    this.waitingFlag(true);
    this.accountService.validateUUID(this.uuid).subscribe(
      () => { this.successFlag(this.wizard.currentPage._id); },
      () => { this.waitingFlag(false); },
    );

    // for test
    // setTimeout(() => {
    //   if (this.uuid === '42') {
    //     this.wizard.forceNext();
    //     this.current = this.wizard.currentPage._id;
    //     this.checkBtn();
    //   } else {
    //     this.errorFlag = true;
    //   }
    //   this.loadingFlag = false;
    //   this.disableInput = false;
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
    this.waitingFlag(true);
    this.accountService.initDB(this.dbInfo).subscribe(
      () => { this.successFlag(this.wizard.currentPage._id); },
      () => { this.waitingFlag(false); },
    );

    // for test
    // setTimeout(() => {
    //   if (this.dbInfo.verify()) {
    //     this.wizard.forceNext();
    //     this.current = this.wizard.currentPage._id;
    //     this.checkBtn();
    //   } else {
    //     this.errorFlag = true;
    //   }
    //   this.loadingFlag = false;
    //   this.disableInput = false;
    // }, 1000);
    // console.log(this.dbInfo);
  }

  onInitSSH(): void {
    this.waitingFlag(true);
    this.accountService.initSSH(this.sshAccount).subscribe(
      () => {
        this.dbAlive = true;
        this.modal.closable = true;
        this.modal.close();
        this.successFlag(this.wizard.currentPage._id);
      },
      () => { this.waitingFlag(false); },
    );

    // for test
    // setTimeout(() => {
    //   if (this.sshAccount.username === '1') {
    //     this.dbAlive = true;
    //     this.wizard.forceNext();
    //     this.modal.closable = true;
    //     this.modal.close();
    //     this.current = this.wizard.currentPage._id;
    //     this.checkBtn();
    //   } else {
    //     this.errorFlag = true;
    //   }
    //   this.loadingFlag = false;
    //   this.disableInput = false;
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
    //   } else {
    //     this.errorFlag = true;
    //   }
    //   this.loadingFlag = false;
    //   this.disableInput = false;
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
    this.wizard.forceNext();
    this.loadingFlag = false;
    this.disableInput = false;
    this.current = id;
    this.checkBtn();
  }
}
