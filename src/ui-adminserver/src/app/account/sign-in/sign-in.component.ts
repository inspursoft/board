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

    this.accountService.signIn(this.user).subscribe(
      (res: MyToken) => {
        if (res) {
          window.sessionStorage.setItem('token', res.token);
          window.sessionStorage.setItem('user', this.user.username);
          this.router.navigateByUrl('dashboard');
        } else {
          this.messageService.showOnlyOkDialog('ERROR.HTTP_UNK', 'ACCOUNT.ERROR');
        }
      },
      (err: HttpErrorResponse) => {
        if (err.status === 403) {
          this.messageService.showOnlyOkDialog('ACCOUNT.FORBIDDEN', 'ACCOUNT.ERROR');
        } else if (err.status === 500) {
          this.messageService.showOnlyOkDialog('ACCOUNT.INCORRECT_USERNAME_OR_PASSWORD', 'ACCOUNT.ERROR');
        } else {
          this.messageService.showOnlyOkDialog('ERROR.HTTP_UNK', 'ACCOUNT.ERROR');
        }
      }
    );
  }

  forgetPassword() {
    this.messageService.showOnlyOkDialog('ACCOUNT.FORGOT_PASSWORD_HELPER', 'ACCOUNT.ERROR');
  }
}
