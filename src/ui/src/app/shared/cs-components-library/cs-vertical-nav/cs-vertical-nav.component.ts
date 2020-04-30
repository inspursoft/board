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
import { AppTokenService } from '../../../shared.service/app-token.service';
import { AppInitService } from '../../../shared.service/app-init.service';

@Directive({
  selector: 'ng-template[appVerticalNavGuide]'
})
export class AppMenuItemUrlDirective {
  @Input() appVerticalNavGuide = '';

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
  collapsed = false;
  @Input() navSource: Array<ICsMenuItemData>;
  @ViewChildren(AppMenuItemUrlDirective) guideContainers: QueryList<AppMenuItemUrlDirective>;
  @ContentChildren(AppMenuItemUrlDirective) guideTemplates: QueryList<AppMenuItemUrlDirective>;

  constructor(private tokenService: AppTokenService,
              private appInitService: AppInitService,
              private changeRef: ChangeDetectorRef) {
    this.navSource = Array<ICsMenuItemData>();
  }

  ngAfterViewInit() {
    this.guideContainers.forEach(container => {
      const guid = this.guideTemplates.find(value =>
        container.appVerticalNavGuide.includes(value.appVerticalNavGuide)
      );
      if (guid) {
        container.viewContainer.createEmbeddedView(guid.templateRef);
        this.changeRef.detectChanges();
      }
    });
  }

  get queryParams(): { token: string } {
    return {token: this.tokenService.token};
  }

  get adminServerUrl(): string {
    return `http://${this.appInitService.systemInfo.board_host}:8082/account/login`;
  }

  get isShowAdminSever(): boolean {
    return !this.appInitService.isArmSystem && !this.appInitService.isMipsSystem;
  }

  isHasChildren(item: ICsMenuItemData): boolean {
    return Reflect.has(item, 'children');
  }
}
