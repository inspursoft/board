import { HttpBind, ResponseArrayBase, ResponseBase } from '../shared/shared.type';

export enum WsNodeResponseStatus {
  UnKnown = 0,
  Start = 1,
  Normal = 2,
  Error = 3,
  Warning = 4,
  Success = 5,
  Failed = 6,
}

export enum NodeActionsType {Add, Remove}

export class NodeLogResponse extends ResponseBase {
  @HttpBind('status') status: WsNodeResponseStatus;
  @HttpBind('message') message: string;
}

export class NodeListType extends ResponseBase {
  @HttpBind('ip') Ip: string;
  @HttpBind('creation_time') CreationTime: number;
}

export class ResponseArrayNode extends ResponseArrayBase<NodeListType> {
  CreateOneItem(res: object): NodeListType {
    return new NodeListType(res);
  }
}
