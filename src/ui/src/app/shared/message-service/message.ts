import { MESSAGE_TARGET, MESSAGE_TYPE, BUTTON_STYLE } from '../shared.const';

export class Message {
  title: string;
  message: string;
  target?: MESSAGE_TARGET;
  data?: any;

  _params: object;
  _buttons: BUTTON_STYLE;
  _type: MESSAGE_TYPE;

  set params(currentParams: object) {
    currentParams ? this._params = currentParams : this._params = [];
  }

  get params(): object {
    return this._params;
  }

  set buttons(currentStyle: BUTTON_STYLE) {
    currentStyle ? this._buttons = currentStyle : this._buttons = BUTTON_STYLE.CONFIRMATION;
  }

  get buttons(): BUTTON_STYLE {
    return this._buttons;
  }

  set type(currentType: MESSAGE_TYPE) {
    currentType ? this._type = currentType : this._type = MESSAGE_TYPE.INTERNAL_ERROR;
  }

  get type(): MESSAGE_TYPE {
    return this._type;
  }
}