import { MESSAGE_TARGET, MESSAGE_TYPE } from '../shared.const';

export class Message {
  title: string;
  message: string;
  target: MESSAGE_TARGET;
  type: MESSAGE_TYPE;
  data: any;
}