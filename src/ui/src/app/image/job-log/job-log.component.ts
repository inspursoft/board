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
  sectionsBuffer: Array<JobLogSection>;
  term: Terminal;
  curRowsCount = 1;
  curStartNum = 1;
  curStartTop = TermOneRowHeight;

  constructor() {
    this.sections = new Array<JobLogSection>();
    this.sectionsBuffer = new Array<JobLogSection>();
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
    this.sectionsBuffer.splice(0, this.sectionsBuffer.length);
  }

  writeBuffer(): void {
    let rowsBuffer = 0;
    this.sectionsBuffer.forEach(value => rowsBuffer += value.termRowsCount);
    this.curRowsCount += rowsBuffer;
    this.term.resize(TermColCount, this.curRowsCount);
    this.sectionsBuffer.forEach(section => {
      section.startTop = this.curStartTop;
      if (section.showContent) {
        section.contents.forEach(value => this.term.write(`\r\n${value}`));
      } else {
        this.term.write(`\r\n${section.contents[0]}`);
      }
      this.curStartTop += section.termRowsCount * TermOneRowHeight;
      this.sections.push(section);
    });
    this.term.scrollLines(this.curRowsCount);
  }

  appendContent(content: string): void {
    const arrContent = content.split(/\r\n|\r|\n/);
    this.parseContent(arrContent);
    this.writeBuffer();
  }

  appendContentArray(arrContent: Array<string>): void {
    this.parseContent(arrContent);
    this.writeBuffer();
  }

  parseContent(arrContent: Array<string>): void {
    this.sectionsBuffer.splice(0, this.sectionsBuffer.length);
    arrContent.forEach(value => {
      if (value.indexOf('section_start') > -1) {
        const lastSectionEndNum = this.sectionsBuffer.length > 0 ?
          this.sectionsBuffer[this.sectionsBuffer.length - 1].endNum : this.curStartNum;
        const section = new JobLogSection();
        section.startContent = value.trim();
        section.startNum = lastSectionEndNum;
        this.sectionsBuffer.push(section);
      } else if (value.indexOf('section_end') > -1) {
        const lastSection = this.sectionsBuffer[this.sectionsBuffer.length - 1];
        if (lastSection.contents.length === 0) {
          this.sectionsBuffer.splice(this.sectionsBuffer.length - 1, 1);
        } else {
          this.curStartNum = lastSection.endNum;
          lastSection.endContent = value.trim();
        }
      } else {
        if (this.sectionsBuffer.length === 0) {
          const section = new JobLogSection();
          section.startNum = this.curStartNum;
          section.contents.push(value.trim());
          this.sectionsBuffer.push(section);
        } else {
          const lastSection = this.sectionsBuffer[this.sectionsBuffer.length - 1];
          if (lastSection.isOpenSection) {
            lastSection.contents.push(value.trim());
          } else {
            const section = new JobLogSection();
            section.startNum = lastSection.endNum;
            section.contents.push(value.trim());
            this.sectionsBuffer.push(section);
          }
        }
      }
    });
  }
}
