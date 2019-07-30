import {
  AfterViewInit,
  ChangeDetectorRef,
  Component,
  ContentChildren,
  Directive,
  Input,
  QueryList,
  TemplateRef, ViewChild,
  ViewChildren,
  ViewContainerRef
} from '@angular/core';
import { ICsMenuItemData } from '../../shared.types';
import { AppTokenService } from "../../../shared.service/app-token.service";
import { ClrVerticalNavGroup } from "@clr/angular";

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
  collapsed = false;
  private _navSource: Array<ICsMenuItemData>;
  @ContentChildren(CsMenuItemUrlDirective) guideTemplates: QueryList<CsMenuItemUrlDirective>;
  @ViewChildren(CsMenuItemUrlDirective) guideContainers: QueryList<CsMenuItemUrlDirective>;

  @Input()
  set navSource(value: Array<ICsMenuItemData>) {
    this._navSource = value;
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
    });
    const triggers: HTMLCollectionOf<Element> = document.getElementsByClassName('nav-group-trigger');
    for (let i = 0; i < triggers.length; i++) {
      this.replaceButtonTagToDiv(triggers.item(i))
    }
  }

  replaceButtonTagToDiv(buttonEelement: Element) {
    let div = document.createElement('div') as HTMLDivElement;
    div.className = 'nav-group-trigger';
    for (let i = 0; i < buttonEelement.childNodes.length; i++) {
      div.appendChild(buttonEelement.childNodes.item(i).cloneNode(true));
    }
    div.addEventListener('click', this.toggleNavStatus);
    buttonEelement.parentNode.appendChild(div);
    buttonEelement.remove();
  }

  toggleNavStatus(event: Event) {
    const navIcon = (event.target as HTMLDivElement).parentElement.getElementsByClassName('nav-group-trigger-icon');
    const parentElement = (event.target as HTMLDivElement).parentElement.parentElement;
    const childNavParentElement = parentElement.nextElementSibling as HTMLDivElement;
    const open = navIcon.item(0).getAttribute('dir') === 'down';
    if (open) {
      navIcon.item(0).setAttribute('dir', 'right');
      childNavParentElement.style.height = '0px';
      childNavParentElement.style.visibility = 'hidden';
      childNavParentElement.style.overflowY = 'hidden';
    } else {
      navIcon.item(0).setAttribute('dir', 'down');
      childNavParentElement.style.height = 'inherit';
      childNavParentElement.style.visibility = 'inherit';
      childNavParentElement.style.overflowY = 'inherit';
    }
  }

  get queryParams(): {token: string} {
    return {token: this.tokenService.token};
  }

  isHasChildren(item: ICsMenuItemData): boolean {
    return Reflect.has(item, 'children');
  }
}
