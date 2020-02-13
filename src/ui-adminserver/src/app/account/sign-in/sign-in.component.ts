import { Component, OnInit } from '@angular/core';
import { User } from '../account.model';
import { AccountService } from '../account.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-sign-in',
  templateUrl: './sign-in.component.html',
  styleUrls: ['./sign-in.component.css']
})
export class SignInComponent implements OnInit {
  user: User;

  constructor(private accountService: AccountService,
    private router: Router) {
    this.user = new User();
  }

  ngOnInit() {
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
