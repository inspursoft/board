import { HttpBase, HttpBind, HttpBindObject } from '../shared/ui-model/model-types';

export enum PvcAccessMode {
  ReadWriteOnce = 'ReadWriteOnce',
  ReadOnlyMany = 'ReadOnlyMany',
  ReadWriteMany = 'ReadWriteMany'
}

export enum PvAccessMode {
  ReadWriteOnce = 'ReadWriteOnce',
  ReadOnlyMany = 'ReadOnlyMany',
  ReadWriteMany = 'ReadWriteMany'
}

export enum PvReclaimMode {
  Retain = 'Retain',
  Recycle = 'Recycle',
  Delete = 'Delete'
}

export class PersistentVolume extends HttpBase {
  @HttpBind('pv_id') id = 0;
  @HttpBind('pv_name') name = '';
  @HttpBind('pv_type') type = 1;
  @HttpBind('pv_state') state = 0;
  @HttpBind('pv_capacity') capacity = '';
  @HttpBind('pv_accessmode') accessMode: PvAccessMode;
  @HttpBind('pv_reclaim') reclaim: PvReclaimMode;

  get capacityValue(): number {
    return Number.parseFloat(this.capacity);
  }

  set capacityValue(value: number) {
    this.capacity = `${value}Gi`;
  }

  get typeDescription(): string {
    return ['Unknown', 'NFS', 'RBD'][this.type];
  }

  get statusDescription(): string {
    return [
      'STORAGE.STATE_UNKNOWN',
      'STORAGE.STATE_PENDING',
      'STORAGE.STATE_AVAILABLE',
      'STORAGE.STATE_BOUND',
      'STORAGE.STATE_RELEASED',
      'STORAGE.STATE_FAILED',
      'STORAGE.STATE_INVALID'][this.state];
  }
}

export class PVOptions extends HttpBase {
  @HttpBind('path') path = '';
  @HttpBind('server') server = '';
}

export class RBDOptions extends HttpBase {
  @HttpBind('user') user = '';
  @HttpBind('keyring') keyring = '';
  @HttpBind('pool') pool = 'rbd';
  @HttpBind('image') image = '';
  @HttpBind('fstype') fsType = '';
  @HttpBind('secretname') secretName = '';
  @HttpBind('secretnamespace') secretNamespace = '';
  @HttpBind('monitors') monitors = '';
}

export class NFSPersistentVolume extends PersistentVolume {
  @HttpBindObject('pv_options', PVOptions) options: PVOptions;

  protected prepareInit() {
    this.type = 1;
    this.options = new PVOptions();
  }
}

export class RBDPersistentVolume extends PersistentVolume {
  @HttpBindObject('pv_options', RBDOptions) options: RBDOptions;

  protected prepareInit() {
    this.type = 2;
    this.options = new RBDOptions();
  }
}

export class PersistentVolumeClaim extends HttpBase {
  @HttpBind('pvc_id') id = 0;
  @HttpBind('pvc_name') name = '';
  @HttpBind('pvc_projectid') projectId = 0;
  @HttpBind('pvc_projectname') projectName = '';
  @HttpBind('pvc_capacity') capacity = '';
  @HttpBind('pvc_state') state = 0;
  @HttpBind('pvc_accessmode') accessMode: PvcAccessMode;
  @HttpBind('pvc_class') class = '';
  @HttpBind('pvc_designatedpv') designatedPv = '';
  @HttpBind('pvc_volume') volume = '';
  events: Array<string>;

  protected prepareInit() {
    this.events = Array<string>();
  }

  get statusDescription(): string {
    return [
      'STORAGE.STATE_UNKNOWN',
      'STORAGE.STATE_PENDING',
      'STORAGE.STATE_BOUND',
      'STORAGE.STATE_LOST',
      'STORAGE.STATE_INVALID'][this.state];
  }
}

export class PersistentVolumeClaimDetail extends HttpBase {
  @HttpBindObject('pvclaim', PersistentVolumeClaim) claim: PersistentVolumeClaim;
  @HttpBind('pvc_state') state = 0;
  @HttpBind('pvc_volume') volume = '';
  @HttpBind('pvc_events') events: Array<string>;

  protected prepareInit() {
    this.events = new Array<string>();
  }
}


