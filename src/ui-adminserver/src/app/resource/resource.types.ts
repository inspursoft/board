import { HttpBind, ResponseArrayBase, ResponseBase, ResponsePaginationBase } from '../shared/shared.type';

export enum NodeActionsType {Add, Remove, Log}

export enum ActionStatus {Ready, Preparing, Executing, Finished}

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
  @HttpBind('id') id: string;
  @HttpBind('ip') ip: string;
  @HttpBind('log_type') type: NodeActionsType;
  @HttpBind('success') success: boolean;
  @HttpBind('completed') completed: boolean;
  @HttpBind('creation_time') creationTime: number;
}

export class NodeLogs extends ResponsePaginationBase<NodeLog> {
  ListKeyName(): string {
    return 'log_list';
  }

  CreateOneItem(res: object): NodeLog {
    return new NodeLog(res);
  }
}

export class NodePostData {
  nodeIp = '';
  nodePassword = '';
  masterIp = '';
  masterPassword = '';
  hostUsername = 'root';
  hostPassword = '';

  getPostData(): object {
    return {
      node_ip: this.nodeIp,
      node_password: this.nodePassword,
      master_ip: this.masterIp,
      master_password: this.masterPassword,
      host_username: this.hostUsername,
      host_password: this.hostPassword
    };
  }

  getParamsData(): { [param: string]: string } {
    return {
      node_ip: this.nodeIp,
      node_password: this.nodePassword,
      master_password: this.masterPassword,
      host_username: this.hostUsername,
      host_password: this.hostPassword
    };
  }
}

export class NodeDetail extends ResponseBase {
  @HttpBind('status') status: NodeLogStatus;
  @HttpBind('message') message: string;
}

export class NodePreparationData extends ResponseBase {
  @HttpBind('master_ip') masterIp: string;
  @HttpBind('host_ip') hostIp: string;
}

export class NodeDetails extends ResponseArrayBase<NodeDetail> {
  CreateOneItem(res: object): NodeDetail {
    return new NodeDetail(res);
  }
}

export class NodeListType extends ResponseBase {
  @HttpBind('ip') ip: string;
  @HttpBind('node_name') nodeName: string;
  @HttpBind('creation_time') creationTime: number;
  @HttpBind('log_time') logTime: number;
  @HttpBind('origin') origin: number;
  @HttpBind('status') status: number;
  @HttpBind('is_master') isMaster: boolean;
  @HttpBind('is_edge') isEdge: boolean;
}

export class NodeList extends ResponseArrayBase<NodeListType> {
  CreateOneItem(res: object): NodeListType {
    return new NodeListType(res);
  }
}

export class NodeControlStatus extends ResponseBase {
  @HttpBind('node_name') nodeName: string;
  @HttpBind('node_type') nodeType: string;
  @HttpBind('node_ip') nodeIp: string;
  @HttpBind('node_phase') nodePhase: string;
  @HttpBind('node_deletable') nodeDeletable: boolean;
  @HttpBind('node_unschedulable') nodeUnschedulable: boolean;
}
