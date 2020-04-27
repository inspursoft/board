import { Component, OnInit } from '@angular/core';
import { User, MyToken } from '../account.model';
import { AccountService } from '../account.service';
import { Router } from '@angular/router';
<<<<<<< HEAD
import { ClrWizard, ClrModal } from '@clr/angular';
import { HeaderComponent } from 'src/app/shared/header/header.component';
=======
import { HttpErrorResponse } from '@angular/common/http';
import { MessageService } from 'src/app/shared/message/message.service';
>>>>>>> dev_new2

@Component({
  selector: 'app-sign-in',
  templateUrl: './sign-in.component.html',
  styleUrls: ['./sign-in.component.css']
})
export class SignInComponent implements OnInit {
  user: User;

  constructor(private accountService: AccountService,
              private messageService: MessageService,
              private router: Router) {
    this.user = new User();
  }

  ngOnInit() { }

  signIn() {
    // test
    // window.sessionStorage.setItem('token', `username=${this.user.username}&password=${this.user.password}`);
    // this.router.navigateByUrl('dashboard');

    this.accountService.postSignIn(this.user).subscribe(
      (res: MyToken) => {
        if (res) {
<<<<<<< HEAD
          window.sessionStorage.setItem('token', res);
=======
          window.sessionStorage.setItem('token', res.token);
          window.sessionStorage.setItem('user', this.user.username);
>>>>>>> dev_new2
          this.router.navigateByUrl('dashboard');
        } else {
          this.messageService.showOnlyOkDialog('ERROR.HTTP_UNK', 'ACCOUNT.ERROR');
        }
      },
<<<<<<< HEAD
      () => {
        alert('账号或密码错误！# Account or password error!');
=======
      (err: HttpErrorResponse) => {
        if (err.status === 403) {
          this.messageService.showOnlyOkDialog('ACCOUNT.FORBIDDEN', 'ACCOUNT.ERROR');
        } else {
          this.messageService.showOnlyOkDialog('ACCOUNT.INCORRECT_USERNAME_OR_PASSWORD', 'ACCOUNT.ERROR');
        }
>>>>>>> dev_new2
      }
    );
  }

<<<<<<< HEAD
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
=======
  forgetPassword() {
    this.messageService.showOnlyOkDialog('ACCOUNT.FORGOT_PASSWORD_HELPER', 'ACCOUNT.ERROR');
>>>>>>> dev_new2
  }
}
