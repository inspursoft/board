import { Component, OnInit } from '@angular/core';
import { MessageService } from "../../shared/message-service/message.service";
import { Message } from "../../shared/message-service/message";
import { BUTTON_STYLE, MESSAGE_TARGET } from "../../shared/shared.const";

@Component({
  selector: 'app-list-audit',
  templateUrl: './list-audit.component.html',
  styleUrls: ['./list-audit.component.css']
})
export class ListAuditComponent implements OnInit {
  beginDate: Date;
  endDate: Date;

  constructor(private messageService:MessageService) {
  }

  ngOnInit() {
  }

  dateTest() {
    if (!this.beginDate || !this.endDate) {
      return
    } else if (this.beginDate > this.endDate) {
      let msg: Message = new Message();
      msg.title = "AUDIT.ILLEGAL_DATE_TITLE";
      msg.message = "AUDIT.ILLEGAL_DATE_MSG";
      msg.buttons = BUTTON_STYLE.ONLY_CONFIRM;
      this.messageService.announceMessage(msg);
      return false;
    } else {
      return
    }
  }
}
