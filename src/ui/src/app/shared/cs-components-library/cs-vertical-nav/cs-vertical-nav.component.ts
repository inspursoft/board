import {
  AfterViewInit,
  ChangeDetectorRef,
  Component,
  ContentChildren,
  Directive,
  Input,
  QueryList,
  TemplateRef,
  ViewChildren,
  ViewContainerRef
} from '@angular/core';
import { ICsMenuItemData } from '../../shared.types';
import { AppTokenService } from "../../../app.init.service";

@Directive({
  selector: 'ng-template[csVerticalNavGuide]'
})
export class CsMenuItemUrlDirective {
  @Input() csVerticalNavGuide = "";

  constructor(public templateRef: TemplateRef<any>,
              public viewContainer: ViewContainerRef) {
  }
}

@Component({
  selector: 'cs-vertical-nav',
  templateUrl: './cs-vertical-nav.component.html',
  styleUrls: ['./cs-vertical-nav.component.css']
})
export class CsVerticalNavComponent implements AfterViewInit {
  collapsed: boolean = false;
  activeIndex: number = 0;
  private _navSource: Array<ICsMenuItemData>;
  @ContentChildren(CsMenuItemUrlDirective) guideTemplates: QueryList<CsMenuItemUrlDirective>;
  @ViewChildren(CsMenuItemUrlDirective) guideContainers: QueryList<CsMenuItemUrlDirective>;

  @Input()
  set navSource(value: Array<ICsMenuItemData>) {
    this._navSource = value;
    this.changeMenuClick(0);
  }

  get navSource(): Array<ICsMenuItemData> {
    return this._navSource;
  }

  constructor(private tokenService: AppTokenService,
              private changeRef: ChangeDetectorRef) {
    this.navSource = Array<ICsMenuItemData>();
  }

  ngAfterViewInit() {
    this.guideContainers.forEach(container => {
      let guid = this.guideTemplates.find(guid => container.csVerticalNavGuide.includes(guid.csVerticalNavGuide));
      if (guid) {
        container.viewContainer.createEmbeddedView(guid.templateRef);
        this.changeRef.detectChanges();
      }
    })
  }

  changeMenuClick(index: number) {
    this.activeIndex = index;
  }

  get queryParams(): {token: string} {
    return {token: this.tokenService.token};
  }
}
