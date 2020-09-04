import { AfterViewInit, Component, ElementRef, OnInit, ViewChild } from '@angular/core';
import { Terminal } from 'xterm';
import { JobLogSection, TermColCount, TermOneRowHeight } from '../image.types';

@Component({
  selector: 'app-job-log',
  templateUrl: './job-log.component.html',
  styleUrls: ['./job-log.component.css']
})
export class JobLogComponent implements OnInit, AfterViewInit {
  @ViewChild('terminalContainer') terminalContainer: ElementRef;
  sections: Array<JobLogSection>;
  term: Terminal;
  writing = false;
  curRowsCount = 1;
  curStartNum = 1;
  curStartTop = TermOneRowHeight;

  constructor() {
    this.sections = new Array<JobLogSection>();
  }

  ngOnInit() {
    this.createTerm();
  }

  ngAfterViewInit(): void {
    this.initTerm();
  }

  initTerm() {
    const terminalContainerElement = (this.terminalContainer.nativeElement as HTMLElement);
    while (terminalContainerElement.firstChild) {
      terminalContainerElement.firstChild.remove();
    }
    this.term.open(terminalContainerElement);
    this.term.focus();
  }

  createTerm() {
    this.term = new Terminal({
      cursorBlink: false,
      disableStdin: true,
      cols: TermColCount,
      cursorStyle: 'underline',
      rows: this.curRowsCount
    });
  }

  toggleContent(section: JobLogSection): void {
    section.showContent = !section.showContent;
    this.writeAll();
  }

  writeAll(): void {
    let rows = 0;
    this.term.clear();
    this.term.resize(TermColCount, 1);
    this.sections.forEach(value => rows += value.termRowsCount);
    this.curRowsCount = rows;
    this.curStartTop = TermOneRowHeight;
    this.term.resize(TermColCount, this.curRowsCount);
    this.sections.forEach(section => {
      section.startTop = this.curStartTop;
      if (section.showContent) {
        section.contents.forEach(value => this.term.write(`\r\n${value}`));
      } else {
        this.term.write(`\r\n${section.contents[0]}`);
      }
      this.curStartTop += section.termRowsCount * TermOneRowHeight;
    });
    this.term.scrollLines(this.curRowsCount);
  }

  clear(): void {
    this.curRowsCount = 1;
    this.curStartNum = 1;
    this.curStartTop = TermOneRowHeight;
    this.term.clear();
    this.term.resize(TermColCount, this.curRowsCount);
    this.sections.splice(0, this.sections.length);
  }

  getTermRowCount(str: string): number {
    if (str.length < TermColCount) {
      return 1;
    } else {
      return 1 + this.getTermRowCount(str.substr(TermColCount));
    }
  }

  writeBuffer(arrContent: Array<string>): void {
    arrContent.forEach(content => {
      if (content.indexOf('section_start') > -1) {
        const section = new JobLogSection();
        section.startContent = content.trim();
        section.startNum = this.curStartNum;
        section.startTop = this.curStartTop;
        this.sections.push(section);
      } else if (content.indexOf('section_end') > -1) {
        const lastSection = this.sections[this.sections.length - 1];
        if (lastSection.contents.length === 0) {
          this.sections.splice(this.sections.length - 1, 1);
        } else {
          lastSection.endContent = content.trim();
        }
      } else {
        if (this.sections.length === 0) {
          const section = new JobLogSection();
          const rowCount = this.getTermRowCount(content);
          this.curRowsCount += rowCount;
          section.startNum = this.curStartNum;
          section.startTop = this.curStartTop;
          section.contents.push(content.trim());
          this.term.resize(TermColCount, this.curRowsCount);
          this.term.write(`\r\n${content.trim()}`);
          this.term.scrollLines(this.curStartNum);
          this.curStartNum += 1;
          this.sections.push(section);
        } else {
          const lastSection = this.sections[this.sections.length - 1];
          if (lastSection.isOpenSection) {
            const rowCount = this.getTermRowCount(content);
            this.curRowsCount += rowCount;
            this.curStartTop += rowCount * TermOneRowHeight;
            this.curStartNum += 1;
            lastSection.contents.push(content.trim());
            this.term.resize(TermColCount, this.curRowsCount);
            this.term.write(`\r\n${content.trim()}`);
            this.term.scrollLines(this.curStartNum);
          } else {
            const section = new JobLogSection();
            const rowCount = this.getTermRowCount(content);
            this.curStartTop += rowCount * TermOneRowHeight;
            this.curRowsCount += rowCount;
            section.startTop = this.curStartTop;
            section.startNum = this.curStartNum;
            section.contents.push(content.trim());
            this.sections.push(section);
            this.term.resize(TermColCount, this.curRowsCount);
            this.term.write(`\r\n${content.trim()}`);
            this.term.scrollLines(this.curStartNum);
            this.curStartNum += 1;
          }
        }
      }
    });
  }

  appendContent(content: string): void {
    const arrContent = content.split(/\r\n|\r|\n/);
    this.writeBuffer(arrContent);
  }

  appendContentArray(arrContent: Array<string>): void {
    this.writeBuffer(arrContent);
  }
}
