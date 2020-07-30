import { Component, ElementRef, EventEmitter, Input, Output, ViewChild } from '@angular/core';

@Component({
  selector: 'app-cs-guide',
  styleUrls: ['./cs-guide.component.css'],
  templateUrl: './cs-guide.component.html'
})
export class CsGuideComponent {
  showValue = false;
  @ViewChild('clrInfoIcon') clrInfoIconRef: ElementRef;
  @Output() nextStep: EventEmitter<any> = new EventEmitter<any>();
  @Output() closeGuide: EventEmitter<any> = new EventEmitter<any>();
  @Input() description: string;
  @Input() position = 'right-middle';
  @Input() isEndStep = false;
  @Input() isShowIcon = false;

  @Input('show')
  get show() {
    return this.showValue;
  }

  set show(value: boolean) {
    this.showValue = value;
    if (value) {
      this.resetView();
    }
  }

  signpostClickEvent(event: Event) {
    event.stopPropagation();
    return false;
  }

  nextStepClick(event: Event) {
    this.nextStep.emit();
    event.stopPropagation();
    return false;
  }

  resetView() {
    setTimeout(() => {
      const el = this.clrInfoIconRef.nativeElement as HTMLElement;
      el.addEventListener('click', () => {
        const signpostElement: HTMLElement = this.clrInfoIconRef.nativeElement.parentElement;
        const divNodeList = signpostElement.getElementsByTagName('div') as HTMLCollectionOf<HTMLDivElement>;
        for (let i = 0; i < divNodeList.length; i++) {
          const div = divNodeList.item(i);
          if (div.className === 'signpost-flex-wrap' || div.className === 'signpost-wrap') {
            div.style.border = 'none';
            div.style.background = 'rgba(0, 64, 96, 0.8)';
          }
          if (div.className === 'signpost-content-header') {
            const clrIcon: HTMLElement = div.getElementsByTagName('clr-icon').item(0) as HTMLElement;
            const buttonClose: HTMLButtonElement = div.getElementsByTagName('button').item(0) as HTMLButtonElement;
            const btnClassName = buttonClose.className;
            buttonClose.addEventListener('click', (evt: MouseEvent) => {
              this.closeGuide.emit(true);
            });
            buttonClose.removeChild(clrIcon);
            buttonClose.innerText = 'X';
            buttonClose.className = `${btnClassName} signpost-content-header-btn-close`;
          }
        }
      });
      this.clrInfoIconRef.nativeElement.click();
    }, 500);
  }
}
