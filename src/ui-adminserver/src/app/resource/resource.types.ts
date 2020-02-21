import { HttpBind, ResponseArrayBase, ResponseBase } from '../shared/shared.type';

export enum NodeActionsType {Add, Remove, Log}

export enum ActionStatus {Ready, Executing, Finished}

export enum NodeLogStatus {
  UnKnown = 0,
  Start = 1,
  Normal = 2,
  Error = 3,
  Warning = 4,
  Success = 5,
  Failed = 6,
}

export class NodeLog extends ResponseBase {
  @HttpBind('ip') ip: string;
  @HttpBind('type') type: NodeActionsType;
  @HttpBind('success') success: boolean;
  @HttpBind('pid') pid: number;
  @HttpBind('creation_time') creationTime: number;
  @HttpBind('completed') completed: boolean;
}

export class NodeLogs extends ResponseArrayBase<NodeLog> {
  CreateOneItem(res: object): NodeLog {
    return new NodeLog(res);
  }
}

export class NodeDetail extends ResponseBase {
  @HttpBind('status') status: NodeLogStatus;
  @HttpBind('message') message: string;
}

export class NodeDetails extends ResponseArrayBase<NodeDetail> {
  CreateOneItem(res: object): NodeDetail {
    return new NodeDetail(res);
  }
}

export class NodeListType extends ResponseBase {
  @HttpBind('ip') Ip: string;
  @HttpBind('creation_time') CreationTime: number;
}

export class NodeList extends ResponseArrayBase<NodeListType> {
  CreateOneItem(res: object): NodeListType {
    return new NodeListType(res);
  }
}
