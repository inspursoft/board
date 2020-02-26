import { Component, HostBinding, OnInit, ViewChild } from '@angular/core';
import { RouterOutlet } from '@angular/router';

@Component({
  selector: 'app-resource',
  templateUrl: './resource.component.html',
  styleUrls: ['./resource.component.css']
})
export class ResourceComponent implements OnInit {
  @HostBinding('style.flex-basis') flexBasis = '100%';
  @HostBinding('style.display') display = 'flex';
  @ViewChild(RouterOutlet) outLet: RouterOutlet;
  computeNavGroupExpanded = false;
  storageNavGroupExpanded = false;

  constructor() {
  }

  ngOnInit() {
    this.setNavGroupsExpandStatus();
    this.outLet.activateEvents.subscribe(() => this.setNavGroupsExpandStatus());
  }

  setNavGroupsExpandStatus() {
    const data = this.outLet.activatedRouteData;
    this.computeNavGroupExpanded = Reflect.get(data, 'group') === 1;
    this.storageNavGroupExpanded = Reflect.get(data, 'group') === 2;
  }
}
