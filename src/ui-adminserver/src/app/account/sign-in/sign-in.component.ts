import { Component, OnInit } from '@angular/core';
import { User, MyToken } from '../account.model';
import { AccountService } from '../account.service';
import { Router } from '@angular/router';
import { HttpErrorResponse } from '@angular/common/http';
import { MessageService } from 'src/app/shared/message/message.service';

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

    // TODO
    this.accountService.postSignIn(this.user).subscribe(
      (res: MyToken) => {
        console.log(res);
        if (res) {
          window.sessionStorage.setItem('token', res.token);
          window.sessionStorage.setItem('user', this.user.username);
          this.router.navigateByUrl('dashboard');
        } else {
          alert('Unknown Error!');
        }
      },
      (err: HttpErrorResponse) => {
        const currentLang = (window.localStorage.getItem('currentLang') === 'zh-cn' || window.localStorage.getItem('currentLang') === 'zh');
        const FORBIDDEN = currentLang ? '禁止访问！' : 'Forbidden!';
        const errorUser = currentLang ? '账号或密码错误！' : 'Account or password error!';
        if (err.status === 403) {
          alert(FORBIDDEN);
        } else {
          alert(errorUser);
        }
      }
    );
  }

  forgetPassword() {
    const currentLang = (window.localStorage.getItem('currentLang') === 'zh-cn' || window.localStorage.getItem('currentLang') === 'zh');
    const forgetPwd = currentLang ? '请在board中修改密码!' : 'Please change the password in Board!';
    alert(forgetPwd);
  }
}
