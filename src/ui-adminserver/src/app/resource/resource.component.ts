import { Component, HostBinding, OnInit, ViewContainerRef } from '@angular/core';

@Component({
  selector: 'app-resource',
  templateUrl: './resource.component.html',
  styleUrls: ['./resource.component.css']
})
export class ResourceComponent implements OnInit {
  @HostBinding('style.flex-basis') flexBasis = '100%';
  @HostBinding('style.display') display = 'flex';


  constructor(private view: ViewContainerRef) {
  }

  ngOnInit() {
  }

  get viewContainer(): ViewContainerRef {
    return this.view;
  }
}
