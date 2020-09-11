import { HttpBase, HttpBind, HttpBindArray, HttpBindObject } from '../shared/ui-model/model-types';

export enum NodeStatusType {
  Schedulable = 1, Unschedulable, Unknown, AutonomousOffline
}

export enum AddNodeType {
  normal, edge
}

export class NodeStatus extends HttpBase {
  readonly masterKey = 'node-role.kubernetes.io/master';
  @HttpBind('node_name') nodeName = '';
  @HttpBind('node_ip') nodeIp = '';
  @HttpBind('node_type') nodeType = '';
  @HttpBind('create_time') createTime = 0;
  @HttpBind('status') status: NodeStatusType;
  @HttpBind('labels') labels: { [p: string]: string };

  get isMaster(): boolean {
    return Reflect.has(this.labels, this.masterKey);
  }

  get isEdge(): boolean {
    return this.nodeType === 'edge';
  }

  get nodeTypeDescribe(): string {
    if (this.nodeType === 'node') {
      return 'NODE.NODE_TYPE_NODE';
    } else if (this.nodeType === 'master') {
      return 'NODE.NODE_TYPE_MASTER';
    } else if (this.nodeType === 'edge') {
      return 'NODE.NODE_TYPE_EDGE';
    } else {
      return 'NODE.NODE_TYPE_UNKNOWN';
    }
  }
}

export class NodeGroupStatus extends HttpBase {
  projectName = '';
  @HttpBind('nodegroup_id') id = 0;
  @HttpBind('nodegroup_name') name = '';
  @HttpBind('nodegroup_comment') comment = '';
  @HttpBind('nodegroup_owner_id') ownerId = 0;
  @HttpBind('nodegroup_creation_time') creationTime = '';
  @HttpBind('nodegroup_update_time') updateTime = '';
  @HttpBind('nodegroup_deleted') deleted = 0;

  postBody(): { [p: string]: any } {
    return {
      nodegroup_id: this.id,
      nodegroup_project: this.projectName,
      nodegroup_name: this.name,
      nodegroup_comment: this.comment
    };
  }
}

export class ServiceInstance extends HttpBase {
  @HttpBind('project_name') projectName: string;
  @HttpBind('service_instance_name') serviceInstanceName: string;
}

export class NodeControlStatus extends HttpBase {
  @HttpBind('node_name') nodeName = '';
  @HttpBind('node_ip') nodeIp = '';
  @HttpBind('node_phase') nodePhase = '';
  @HttpBind('node_deletable') deletable = false;
  @HttpBind('node_unschedulable') nodeUnschedulable = false;
  @HttpBindArray('service_instances', ServiceInstance) serviceInstances: Array<ServiceInstance>;

  prepareInit() {
    super.prepareInit();
    this.serviceInstances = Array<ServiceInstance>();
  }
}

export class NodeDetail extends HttpBase {
  @HttpBind('node_name') nodeName: string;
  @HttpBind('node_ip') nodeIp: string;
  @HttpBind('create_time') createTime: number;
  @HttpBind('memory_size') memorySize: string;
  @HttpBind('cpu_usage') cpuUsage: number;
  @HttpBind('memory_usage') memoryUsage: number;
  @HttpBind('storage_total') storageTotal: string;
  @HttpBind('storage_use') storageUse: string;
}

export class NodeGroupDetail extends HttpBase {
  @HttpBindObject('nodegroup', NodeGroupStatus) nodeGroup: NodeGroupStatus;
  @HttpBind('nodelist') nodeList: Array<string>;

  protected prepareInit() {
    this.nodeGroup = new NodeGroupStatus();
    this.nodeList = new Array<string>();
  }
}

export class EdgeNode extends HttpBase {
  @HttpBind('name') name = '';
  @HttpBind('node_ip') nodeIp = '';
  @HttpBind('node_password') nodePassword = '';
  @HttpBind('cpu_type') cpuType = '';
  @HttpBind('master') master = '';
  @HttpBind('registry_mode') registryMode = '';
}
