import { Directive, HostListener, ElementRef, Renderer2, Input, Inject, ChangeDetectorRef, Output } from '@angular/core';
import { DOCUMENT } from '@angular/common';

@Directive({
  selector: '[appToTop]',
})
export class ToTopDirective {
  private target: HTMLElement | null = null;

  @Input('topTarget')
  set topTarget(el: string | HTMLElement) {
    this.target = typeof el === 'string' ? this.doc.querySelector(el) : el;
    console.log(el);
  }

  constructor(
    // private el: ElementRef,
    @Inject(DOCUMENT) private doc: any,
    private cd: ChangeDetectorRef
  ) { }

  // @HostListener('click', ['$event.target'])
  // onClick() {
  //   // window.scrollTo({
  //   //   top: 0,
  //   //   behavior: 'smooth',
  //   // });
  //   this.target.scrollTo({
  //     top: 0,
  //     behavior: 'smooth',
  //   });
  //   this.cd.markForCheck();
  //   console.log(this.target);
  // }

  onClick() {
    // window.scrollTo({
    //   top: 0,
    //   behavior: 'smooth',
    // });
    this.target.scrollTo({
      top: 0,
      behavior: 'smooth',
    });
    this.cd.markForCheck();
    console.log(this.target);
  }
}
