/**
 * Created by liyanq on 15/11/2017.
 */
import { Component, EventEmitter, Output, OnInit, ViewChild, ElementRef } from '@angular/core';
import { Subject } from 'rxjs';
import { debounceTime } from 'rxjs/operators';

@Component({
  templateUrl: './cs-search-input.component.html',
  styleUrls: ['./cs-search-input.component.css']
})
export class CsSearchInputComponent implements OnInit {
  searchText = '';
  searchInputChange: Subject<string>;
  @ViewChild('inputControl') inputControl: ElementRef;
  @Output() searchEvent: EventEmitter<string>;

  constructor() {
    this.searchEvent = new EventEmitter<string>();
    this.searchInputChange = new Subject<string>();
  }

  ngOnInit() {
    this.searchInputChange.asObservable().pipe(debounceTime(300))
      .subscribe(searchText => {
        this.searchEvent.emit(searchText);
      });
  }

  onInputChange(event: KeyboardEvent) {
    setTimeout(() => {
      this.searchInputChange.next(this.searchText);
    }, 10);
  }

  onIconClick() {
    this.searchInputChange.next(this.searchText);
  }
}
