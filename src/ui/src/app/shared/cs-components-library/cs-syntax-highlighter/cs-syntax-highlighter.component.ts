import { AfterContentInit, Component, Directive, ElementRef, Input } from '@angular/core';
import * as Prism from "prismjs"
import * as Languages from "prismjs/components"

@Directive({selector: "[highlighter]"})
export class CsSyntaxHighlighterDirective implements AfterContentInit {
  @Input() public language: string;
  @Input() public content: string;

  constructor(private eltRef: ElementRef) {
  }

  ngAfterContentInit() {
    this.eltRef.nativeElement.innerHTML = Prism.highlight(this.content, Languages.languages[this.language]);
  }
}

@Component({
  selector: 'cs-syntax-highlighter',
  templateUrl: './cs-syntax-highlighter.component.html',
  styleUrls: ['./cs-syntax-highlighter.component.css']
})
export class CsSyntaxHighlighterComponent {
  @Input() public language: string;
  @Input() public content: string;
  @Input() public jsonContent: Object;
}
