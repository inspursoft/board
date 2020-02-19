import { Component, OnInit, ViewContainerRef } from '@angular/core';

@Component({
  selector: 'app-resource',
  templateUrl: './resource.component.html',
  styleUrls: ['./resource.component.css']
})
export class ResourceComponent implements OnInit {

  constructor(private view: ViewContainerRef) {
  }

  ngOnInit() {
  }

  get viewContainer(): ViewContainerRef {
    return this.view;
  }
}
