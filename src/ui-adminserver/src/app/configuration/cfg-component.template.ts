import { Input, ViewChildren, QueryList, ViewChild } from '@angular/core';
import { CardObject } from './cfg.models';
import { CsComponentBase } from 'src/app/shared/cs-components-library/cs-component-base';
import { InputExComponent } from 'board-components-library';

export class CfgComponent extends CsComponentBase {
  @Input() card: CardObject;
  @Input() isInit: boolean;
  @ViewChildren(InputExComponent) inputExComponents: QueryList<InputExComponent>;

  toggle() {
    this.card.cardStatus = !this.card.cardStatus;
    if (!this.card.cardStatus) {
      setTimeout(() => {
        this.card.showContent = !this.card.showContent;
      }, 250);
    } else {
      this.card.showContent = !this.card.showContent;
    }
  }

  verifyInputExValid(): boolean {
    return this.inputExComponents.toArray().every((component: InputExComponent) => {
      if (!component.isValid && !component.inputDisabled) {
        component.checkSelf();
      }
      return component.isValid || component.inputDisabled;
    });
  }
}
