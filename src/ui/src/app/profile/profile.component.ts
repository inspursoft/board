import { Component } from '@angular/core';
import { AppInitService } from "../app.init.service";

@Component({
  selector: 'profile',
  styleUrls:['./profile.component.css'],
  templateUrl: './profile.component.html'
})
export class ProfileComponent {
  version: string = "";

  constructor(private appInitService: AppInitService) {
    this.version = this.appInitService.systemInfo.board_version;
  }
}