import { Component } from '@angular/core';
import { AppTokenService } from "../shared.service/app-token.service";
import { AppInitService } from "../shared.service/app-init.service";
import { MessageService } from "../shared.service/message.service";

@Component({
  templateUrl: './node.component.html'
})
export class NodeComponent {
  constructor(private appTokenService: AppTokenService,
              private appInitService: AppInitService,
              private messageService: MessageService) {
    console.log('appTokenService', this.appTokenService)
    console.log('appInitService', this.appInitService)
    console.log('messageService', this.messageService)
  }
}
