import { Component, OnInit } from '@angular/core';
import { UserVerify } from '../account.model';
import { AccountService } from '../account.service';
import { Router } from '@angular/router';
import { CsComponentBase } from 'src/app/shared/cs-components-library/cs-component-base';

@Component({
  selector: 'app-sign-up',
  templateUrl: './sign-up.component.html',
  styleUrls: ['./sign-up.component.css']
})
export class SignUpComponent extends CsComponentBase implements OnInit {
  isSignUpWIP = false;
  user: UserVerify;

  constructor(private accountService: AccountService,
              private router: Router) {
    super();
    this.user = new UserVerify();
  }

  ngOnInit(): void {
  }

  signUp(): void {
    if (this.verifyInputExValid()) {
      this.isSignUpWIP = true;
      this.accountService.postSignUp(this.user.toUser()).subscribe(
        () => alert('success registry!'),
        () => alert('error registry!'),
        () => this.router.navigateByUrl('/account/login')
      );
    }
  }

  // goBack(): void {
  //   this.router.navigateByUrl('/account/login');
  // }
}

