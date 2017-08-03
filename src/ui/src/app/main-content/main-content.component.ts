import { Component } from '@angular/core';
import { AppInitService } from '../app.init.service';

@Component({
  selector: 'main-content',
  templateUrl: 'main-content.component.html'
})  
export class MainContentComponent {
  token: string;
  constructor(private appInitService: AppInitService) {
    this.token = this.appInitService.token;
  }
}