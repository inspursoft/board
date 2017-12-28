/**
 * Created by liyanq on 15/11/2017.
 */
import { Component, EventEmitter, Output, OnInit, ViewChild, ElementRef } from "@angular/core"
import { Subject } from "rxjs/Subject";

@Component({
  selector: "cs-search-input",
  templateUrl: "./cs-search-input.component.html",
  styleUrls: ["./cs-search-input.component.css"]
})
export class CsSearchInput implements OnInit {
  searchText: string = "";
  searchInputChange: Subject<string>;
  @ViewChild("inputControl") inputControl: ElementRef;
  @Output() onSearchEvent: EventEmitter<string>;

  constructor() {
    this.onSearchEvent = new EventEmitter<string>();
    this.searchInputChange = new Subject<string>();
  }

  ngOnInit() {
    this.searchInputChange.asObservable().debounceTime(300)
      .subscribe(searchText => {
        this.onSearchEvent.emit(searchText);
      })
  }

  onInputChange(event: KeyboardEvent) {
    setTimeout(() => {
      this.searchInputChange.next(this.searchText);
    }, 10)
  }

  onIconClick() {
    this.searchInputChange.next(this.searchText);
  }
}