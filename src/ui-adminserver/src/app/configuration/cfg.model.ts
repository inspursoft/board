import { HttpBind, ResponseBase, RequestBase } from '../shared/shared.type';

export class CardObject {
  cardStatus = true;
  showContent = true;

  toObj(): { [param: string]: boolean } {
    return {
      cardStatus: this.cardStatus,
      showContent: this.showContent
    };
  }
}

export class CfgCardObjects {
  apiserver: CardObject;
  gogits: CardObject;
  jenkins: CardObject;
  kvm: CardObject;
  ldap: CardObject;
  email: CardObject;
  others: CardObject;

  constructor() {
    this.apiserver = new CardObject();
    this.gogits = new CardObject();
    this.jenkins = new CardObject();
    this.kvm = new CardObject();
    this.ldap = new CardObject();
    this.email = new CardObject();
    this.others = new CardObject();
  }

  toObj(): { [param: string]: object } {
    return {
      apiserver: this.apiserver.toObj(),
      gogits: this.gogits.toObj(),
      jenkins: this.jenkins.toObj(),
      kvm: this.kvm.toObj(),
      ldap: this.ldap.toObj(),
      email: this.email.toObj(),
      others: this.others.toObj(),
    };
  }
}

export class VerifyPassword implements RequestBase {
  which = '';
  value = '';

  PostBody(): object {
    return {
      which: this.which,
      value: this.value
    };
  }
}
