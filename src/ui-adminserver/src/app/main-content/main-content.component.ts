import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';

@Component({
  selector: 'app-main-content',
  templateUrl: './main-content.component.html',
  styleUrls: ['./main-content.component.css']
})
export class MainContentComponent implements OnInit {

  constructor(private router: Router) { }

  ngOnInit() {
    if (!window.sessionStorage.getItem('token') ||
      /^\w{8}(-\w{4}){3}-\w{12}$/.test(window.sessionStorage.getItem('token'))) {
      this.router.navigateByUrl('account/login');
    } else {
      this.router.navigateByUrl('dashboard');
    }
  }

}
