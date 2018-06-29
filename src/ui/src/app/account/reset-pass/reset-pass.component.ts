import { Component, OnInit } from '@angular/core';
import { AccountService } from "../account.service";
import { MessageService } from "../../shared/message-service/message.service";
import { ActivatedRoute, Router } from "@angular/router";
import { SignUp } from "../sign-up/sign-up";

@Component({
  selector: 'app-reset-pass',
  templateUrl: './reset-pass.component.html',
  styleUrls: ['./reset-pass.component.css']
})
export class ResetPassComponent implements OnInit {
  private resetUuid: string;
  private signUpModel: SignUp = new SignUp();
  constructor(
    private accountService: AccountService,
    private messageService: MessageService,
    private router: Router,
    private route: ActivatedRoute,
  ) { }

  ngOnInit() {
    this.route.queryParamMap.subscribe(params=>{
      this.resetUuid = params.get("reset_uuid")
    });
  }

}
