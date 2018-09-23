/**
 * Dropdown Component
 * v1.0
 * Created by liyanq on 9/4/17.
 */

import { Component, EventEmitter, Input, OnChanges, OnInit, Output, SimpleChanges, ViewChild } from "@angular/core"
import { DropdownMenuPositon } from "../../shared.types";
import { Subject } from "rxjs/Subject";
import { animate, state, style, transition, trigger } from "@angular/animations";
import { DISMISS_CHECK_DROPDOWN } from "../../shared.const";

export const ONLY_FOR_CLICK = "OnlyClick";
const DROP_DOWN_SHOW_COUNT = 20;
export type EnableSelectCallBack = (item: any) => boolean;

@Component({
  selector: "cs-dropdown",
  templateUrl: "./cs-dropdown.component.html",
  styleUrls: ["./cs-dropdown.component.css"],
  animations: [
    trigger('check', [
      state('begin', style({backgroundColor: '#ebafa6'})),
      state('end', style({backgroundColor: 'transparent'})),
      transition('begin => end', animate(500))
    ])
  ]
})
export class CsDropdownComponent implements OnChanges, OnInit {
  @ViewChild("csDropdown") csDropdown: Object;
  @Input() dropdownPosition: DropdownMenuPositon = 'bottom-left';
  @Input() dropdownDisabled = false;
  @Input() dropdownHideSearch = false;
  @Input() dropdownCanSelect: EnableSelectCallBack;
  @Input() dropdownDefaultText = '';
  @Input() dropdownWidth = 100;
  @Input() dropdownList: Array<any>;
  @Input() dropdownListTextKey = '';
  @Input() dropdownTitleFontSize = 14;
  @Input() dropdownMustBeSelect = true;
  @Output("onChange") dropdownChange: EventEmitter<any>;
  @Output("onOnlyClickItem") dropdownClick: EventEmitter<any>;
  isShowDefaultText = true;
  dropdownSearchText = '';
  dropdownShowTimes = 1;
  dropdownText = '';
  subFilterDropdownList: Subject<string>;
  shownDropdownList: Array<any>;
  filterDropdownList: Array<any>;
  animation: string;

  constructor() {
    this.dropdownChange = new EventEmitter<any>();
    this.dropdownClick = new EventEmitter<any>();
    this.subFilterDropdownList = new Subject<string>();
    this.shownDropdownList = Array<any>();
    this.filterDropdownList = Array<any>();
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (changes["dropdownList"] && changes["dropdownList"].currentValue) {
      this.isShowDefaultText = true;
      this.filterDropdownList = this.dropdownList;
      if (this.dropdownHideSearch) {
        this.shownDropdownList = this.dropdownList;
      } else {
        this.shownDropdownList = this.dropdownList.filter((value, index) => index < this.dropdownShowTimes * DROP_DOWN_SHOW_COUNT);
      }
    }
  }

  ngOnInit() {
    this.subFilterDropdownList.asObservable().debounceTime(300).subscribe((searchText: string) => {
      this.filterDropdownList = this.dropdownList.filter(value => {
        const text = this.getItemDescription(value);
        return searchText !== '' ? text.indexOf(searchText) > -1 : true;
      });
      this.shownDropdownList = this.filterDropdownList.filter((value, index) => index < this.dropdownShowTimes * DROP_DOWN_SHOW_COUNT)
    })
  }

  filterExecute($event: KeyboardEvent) {
    this.dropdownSearchText = ($event.target as HTMLInputElement).value;
    this.subFilterDropdownList.next(this.dropdownSearchText);
  }

  getItemClass(item: any) {
    return {
      'special': (typeof item == "object") && item['isSpecial'],
      'active': this.dropdownText === this.getItemDescription(item)
    }
  }

  get dropdownEnabled(): boolean {
    return !this.dropdownDisabled;
  }

  get hasMoreItems(): boolean {
    return this.dropdownShowSearch && this.shownDropdownList.length < this.filterDropdownList.length;
  }

  get dropdownShowSearch(): boolean {
    return !this.dropdownHideSearch;
  }

  get active(): boolean {
    /*Todo:this is bad method, but no way better than it at present.2018/1/3*/
    return this.csDropdown && this.csDropdown["ifOpenService"]["open"];
  }

  public get valid(): boolean {
    return this.dropdownDisabled || (this.dropdownMustBeSelect ? !this.isShowDefaultText : true)
  }

  getItemDescription(item: any): string {
    if (typeof item == "object") {
      return item[this.dropdownListTextKey].toString();
    }
    return item.toString();
  }

  changeSelect(item: any) {
    if (typeof item == "object" && item[ONLY_FOR_CLICK]) {
      this.dropdownClick.emit(item);
    } else {
      if (this.dropdownCanSelect && !this.dropdownCanSelect(item)) {
        return;
      }
      if (this.dropdownText != this.getItemDescription(item)) {
        this.isShowDefaultText = false;
        this.dropdownText = this.getItemDescription(item);
        this.dropdownChange.emit(item);
      }
    }
  }

  incShowTimes(event: MouseEvent): void {
    this.dropdownShowTimes += 1;
    this.subFilterDropdownList.next(this.dropdownSearchText);
    event.stopImmediatePropagation();
  }

  public checkInputSelf() {
    if (this.dropdownEnabled && this.isShowDefaultText && this.dropdownMustBeSelect) {
      this.animation = 'begin';
      setTimeout(() => this.animation = 'end', DISMISS_CHECK_DROPDOWN);
    }
  }
}