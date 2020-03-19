import { Input } from '@angular/core';
import { CardObject } from './cfg.models';
import { CsComponentBase } from 'src/app/shared/cs-components-library/cs-component-base';

export class CfgComponent extends CsComponentBase {
  @Input() card: CardObject;
  @Input() isInit: boolean;

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
}
