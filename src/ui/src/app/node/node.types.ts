import { HttpBind, HttpBindArray, HttpBindObject, ResponseBase } from '../shared/ui-model/model-types';

export enum NodeStatusType {
  Schedulable = 1, Unschedulable, Unknown
}

export class NodeStatus extends ResponseBase {
  readonly masterKey = 'node-role.kubernetes.io/master';
  @HttpBind('node_name') nodeName: string;
  @HttpBind('node_ip') nodeIp: string;
  @HttpBind('create_time') createTime: number;
  @HttpBind('status') status: NodeStatusType;
  @HttpBindObject('labels') labels: { [p: string]: string };

  get isMaster(): boolean {
    return Reflect.has(this.labels, this.masterKey);
  }
}

export class NodeGroupStatus extends ResponseBase {
  projectName = '';
  @HttpBind('nodegroup_id') id: number;
  @HttpBind('nodegroup_name') name: string;
  @HttpBind('nodegroup_comment') comment: string;
  @HttpBind('nodegroup_owner_id') ownerId: number;
  @HttpBind('nodegroup_creation_time') creationTime: string;
  @HttpBind('nodegroup_update_time') updateTime: string;
  @HttpBind('nodegroup_deleted') deleted: number;

  postBody(): { [p: string]: string } {
    return {
      nodegroup_project: this.projectName,
      nodegroup_name: this.name,
      nodegroup_comment: this.comment
    };
  }

}

export class ServiceInstance extends ResponseBase {
  @HttpBind('project_name') projectName: string;
  @HttpBind('service_instance_name') serviceInstanceName: string;
}

export class NodeControlStatus extends ResponseBase {
  @HttpBind('node_name') nodeName: string;
  @HttpBind('node_ip') nodeIp: string;
  @HttpBind('node_phase') nodePhase: string;
  @HttpBind('node_deletable') deletable: boolean;
  @HttpBind('node_unschedulable') nodeUnschedulable: boolean;
  @HttpBindArray('service_instances', ServiceInstance) serviceInstances: Array<ServiceInstance>;

  protected prepareInit() {
    super.prepareInit();
    this.serviceInstances = Array<ServiceInstance>();
  }
}

export class NodeDetail extends ResponseBase {
  @HttpBind('node_name') nodeName: string;
  @HttpBind('node_ip') nodeIp: string;
  @HttpBind('create_time') createTime: string;
  @HttpBind('memory_size') memorySize: string;
  @HttpBind('cpu_usage') cpuUsage: number;
  @HttpBind('memory_usage') memoryUsage: number;
  @HttpBind('storage_total') storageTotal: string;
  @HttpBind('storage_use') storageUse: string;
}
