import { Component, OnInit, ViewChild } from '@angular/core';
import { User, UserVerify } from '../account.model';
import { AccountService } from '../account.service';
import { Router } from '@angular/router';
import { ClrWizard } from '@clr/angular';

@Component({
  selector: 'app-sign-in',
  templateUrl: './sign-in.component.html',
  styleUrls: ['./sign-in.component.css']
})
export class SignInComponent implements OnInit {
  @ViewChild("wizard") wizard: ClrWizard;
  _open: boolean = true;
  @ViewChild("myForm") formData: any;
  loadingFlag: boolean = false;
  errorFlag: boolean = false;
  errorVerifyFlag: boolean = false;

  showInitDB = true;
  showInitSSH = true;

  uuid = '';
  dbPwd = '';
  dbPwdConfirm = '';
  disableDbPwdFrom = true;
  sshAccount: User;
  disableAccountFrom = true;
  user: User;
  account: UserVerify;

  constructor(private accountService: AccountService,
    private router: Router) {
    this.account = new UserVerify();
    this.account.username = 'admin';
    this.sshAccount = new User();
    this.user = new User();
  }

  ngOnInit() {
    // this.accountService.checkInit().subscribe(
    //   (res: string) => {
    //     let step = res.toLowerCase();
    //     if (step === 'no') {
    //       this._open = false;
    //     } else if (step === "step2") {
    //       this.showInitDB = false;
    //     } else if (step === "step3") {
    //       this.showInitDB = false;
    //       this.showInitSSH = false;
    //     }
    //   }
    // );

    // for test
    let step = 'step';
    if (step === 'no') {
      this._open = false;
    } else if (step === "step2") {
      this.showInitDB = false;
    } else if (step === "step3") {
      this.showInitDB = false;
      this.showInitSSH = false;
    }

  }

  onWelcome(): void {
    this.loadingFlag = true;
    this.errorFlag = false;
    // this.accountService.createUUID().subscribe(
    //   () => { this.wizard.forceNext(); this.loadingFlag = false; },
    //   () => { this.errorFlag = true; this.loadingFlag = false; },
    // );

    // for test
    setTimeout(() => {
      this.wizard.forceNext();
      this.loadingFlag = false;
    }, 1000);
  }

  onVerify(): void {
    this.loadingFlag = true;
    this.errorFlag = false;
    // this.accountService.validateUUID(this.uuid).subscribe(
    //   () => { this.wizard.forceNext(); this.loadingFlag = false; },
    //   () => { this.errorFlag = true; this.loadingFlag = false; },
    // );

    // for test
    setTimeout(() => {
      if (this.uuid === "42") {
        this.wizard.forceNext();
      } else {
        this.errorFlag = true;
      }
      this.loadingFlag = false;
    }, 1000);
  }

  verifyDBPwd() {
    this.errorVerifyFlag = false;
    if (this.dbPwd == this.dbPwdConfirm) {
      this.disableDbPwdFrom = false;
    } else {
      this.disableDbPwdFrom = true;
      if (this.dbPwd && this.dbPwdConfirm) {
        this.errorVerifyFlag = true;
      }
    }
  }

  onInitDB(): void {
    this.loadingFlag = true;
    this.errorFlag = false;
    // this.accountService.initDB(this.dbPwd).subscribe(
    //   () => { this.wizard.forceNext(); this.loadingFlag = false; },
    //   () => { this.errorFlag = true; this.loadingFlag = false; },
    // );

    // for test
    setTimeout(() => {
      if (this.dbPwd == this.dbPwdConfirm) {
        this.wizard.forceNext();
      } else {
        this.errorFlag = true;
      }
      this.loadingFlag = false;
    }, 1000);
  }

  onInitSSH(): void {
    this.loadingFlag = true;
    this.errorFlag = false;
    // this.accountService.initSSH(this.sshAccount).subscribe(
    //   () => { this.wizard.forceNext(); this.loadingFlag = false; },
    //   () => { this.errorFlag = true; this.loadingFlag = false; },
    // );

    // for test
    setTimeout(() => {
      if (this.sshAccount.username === "1") {
        this.wizard.forceNext();
      } else {
        this.errorFlag = true;
      }
      this.loadingFlag = false;
    }, 1000);
  }

  verifyAccountPwd() {
    this.errorVerifyFlag = false;
    if (this.account.password == this.account.passwordConfirm) {
      this.disableAccountFrom = false;
    } else {
      this.disableAccountFrom = true;
      if (this.account.password && this.account.passwordConfirm) {
        this.errorVerifyFlag = true;
      }
    }
  }

  onInitAccount(): void {
    this.loadingFlag = true;
    this.errorFlag = false;
    // this.accountService.postSignUp(this.account.toUser()).subscribe(
    //   () => { this.wizard.forceFinish(); this.loadingFlag = false; },
    //   () => { this.errorFlag = true; this.loadingFlag = false; },
    // );

    // for test
    setTimeout(() => {
      if (this.account.password === "11111111") {
        this.wizard.forceFinish();
      } else {
        this.errorFlag = true;
      }
      this.loadingFlag = false;
    }, 1000);
  }


  signIn() {
    console.log(this.user.username + this.user.password);
    // test
    // window.sessionStorage.setItem('token', `username=${this.user.username}&password=${this.user.password}`);
    // this.router.navigateByUrl('dashboard');

    // TODO
    this.accountService.postSignIn(this.user).subscribe(
      (res: string) => {
        if (res == 'login success') {
          window.sessionStorage.setItem('token', `It_is_only_a_test`);
          this.router.navigateByUrl('dashboard');
        } else {
          alert('Error account or password!');
        }

      },
      () => {
        alert('Network Error!');
      }
    );
  }
}
