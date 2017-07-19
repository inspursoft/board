import { Component, OnChanges, SimpleChanges, Input } from "@angular/core"

interface buttonsOption {
  readonly pageIndex: number;
  readonly description: string;
}

enum HideMode {
  hmNone,
  hmHeader,
  hmFooter,
  hmBoth
}

@Component({
  selector: "pagination",
  templateUrl: "./pagination.component.html",
  styleUrls: ["./pagination.component.css"]
})

export class Pagination implements OnChanges {
  @Input() recordCount: number;
  _lastVisitPage: number = 0;
  _dynamicButtons: Array<buttonsOption> = new Array<buttonsOption>();

  readonly perPageCountList = Array.from([10, 15, 30, 50, 100]);
  curPage: number = 1;
  recordCountPerPage: number = 10;
  pageCount: number;

  ngOnChanges(changes: SimpleChanges) {
    this.calculatePageCount();
  }

  private calculatePageCount() {
    this.pageCount = Math.ceil(this.recordCount / this.recordCountPerPage) || 1;
    this.pageCount = Math.abs(this.pageCount);
  }

  private getHideMode(): HideMode {
    if (!this.pageCount || this.pageCount <= 3)
      return HideMode.hmNone;
    if (this.curPage < 3) {
      return HideMode.hmFooter;
    }
    else if (this.pageCount - this.curPage > 1) {
      return HideMode.hmBoth;
    }
    else if (this.pageCount - this.curPage <= 1) {
      return HideMode.hmHeader;
    }
  }

  changePerPageCount(perPageCount: number) {
    this._lastVisitPage = 0;
    this.curPage = 1;
    this.recordCountPerPage = perPageCount;
    this.calculatePageCount();
  }

  changePage(page: number) {
    this.curPage = page;
  }

  get onePageMaxRecord() {
    return Math.min(this.curPage * this.recordCountPerPage, this.recordCount);
  }

  get getDynamicButtons() {
    let curHideMode: HideMode = this.getHideMode();
    if (this._lastVisitPage == this.curPage) {
      return this._dynamicButtons
    }
    ;
    this._dynamicButtons = [];
    switch (curHideMode) {
      case HideMode.hmNone: {
        for (let i = 1; i <= this.pageCount; i++) {
          this._dynamicButtons.push({pageIndex: i, description: i.toString()})
        }
        break;
      }
      case HideMode.hmFooter: {
        for (let i = 1; i <= 3; i++) {
          this._dynamicButtons.push({pageIndex: i, description: i.toString()})
        }
        this._dynamicButtons.push({pageIndex: this.pageCount, description: "..."});
        break;
      }
      case HideMode.hmHeader: {
        this._dynamicButtons.push({pageIndex: 1, description: "..."})
        for (let i = this.pageCount - 2; i <= this.pageCount; i++) {
          this._dynamicButtons.push({pageIndex: i, description: i.toString()})
        }
        break;
      }
      case HideMode.hmBoth: {
        this._dynamicButtons.push({pageIndex: 1, description: "..."})
        for (let i = this.curPage - 1; i <= this.curPage + 1; i++) {
          this._dynamicButtons.push({pageIndex: i, description: i.toString()})
        }
        this._dynamicButtons.push({pageIndex: this.pageCount, description: "..."});
        break;
      }
    }
    this._lastVisitPage = this.curPage;
    return this._dynamicButtons;
  }
}