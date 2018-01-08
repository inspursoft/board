import { Component, ElementRef, EventEmitter, HostListener, Input, OnDestroy, Output, ViewChild } from "@angular/core"

@Component({
  selector: 'cs-guide',
  styleUrls: ['./cs-guide.component.css'],
  templateUrl: './cs-guide.component.html'
})
export class CsGuideComponent implements OnDestroy {
  private _show: boolean = false;
  private _testShow: boolean = false;
  @ViewChild("clrInfoIcon") clrInfoIconRef: ElementRef;
  @Output("onNextStep") nextStepEvent: EventEmitter<any> = new EventEmitter<any>();
  @Output("onClose") closeEvent: EventEmitter<any> = new EventEmitter<any>();
  @Input("description") description: string;
  @Input("position") position: string = 'right-middle';
  @Input("isEndStep") isEndStep: boolean = false;


  @Input("show")
  get show() {
    return this._show;
  }

  set show(value: boolean) {
    this._show = value;
    if (value) {
      setTimeout(() => {
        this.clrInfoIconRef.nativeElement.click();
        let signpostElement: HTMLElement = this.clrInfoIconRef.nativeElement.parentElement;
        let divNodeList: NodeListOf<HTMLDivElement> = signpostElement.getElementsByTagName("div");
        for (let i = 0; i < divNodeList.length; i++) {
          let div = divNodeList.item(i);
          if (div.className == "signpost-flex-wrap") {
            div.style.border = "none";
            div.style.background = "rgba(139, 188, 255, 0.5)";
          }
          if (div.className == "signpost-content-header") {
            let clrIcon: HTMLElement = div.getElementsByTagName("clr-icon").item(0) as HTMLElement;
            let buttonClose: HTMLButtonElement = div.getElementsByTagName("button").item(0) as HTMLButtonElement;
            let btnClassName = buttonClose.className;
            buttonClose.addEventListener("click", (evt: Event) => {
              this.closeEvent.emit(true);
            });
            buttonClose.removeChild(clrIcon);
            buttonClose.innerText = "X";
            buttonClose.className = `${btnClassName} signpost-content-header-btn-close`;
          }
        }
      }, 500);
    }
  }

  ngOnDestroy() {
    console.log("ngOnDestroy");
  }

  @HostListener("window.click", ["$event"]) windowClick() {
    console.log("window.click")
  }

  signpostClickEvent(event: Event) {
    event.stopPropagation();
    return false;
  }

  nextStepClick(event: Event) {
    this.clrInfoIconRef.nativeElement.click();
    this.nextStepEvent.emit();
    event.stopPropagation();
    return false;
  }
}