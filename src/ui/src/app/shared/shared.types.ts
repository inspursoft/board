import { HttpErrorResponse } from '@angular/common/http';
import { Type } from '@angular/core';
import { TimeoutError } from 'rxjs';
import { HttpBase, HttpBind, HttpBindBoolean, HttpBindObject } from './ui-model/model-types';
import { ConfigMapDetailMetadata } from '../resource/resource.types';

export interface ICsMenuItemData {
  caption: string;
  icon: string;
  url: string;
  visible: boolean;
  children?: Array<ICsMenuItemData>;
  isAdminServer?: boolean;
}

export enum RETURN_STATUS {
  rsNone, rsConfirm, rsCancel
}

export enum EXECUTE_STATUS {
  esNotExe = 'NotExe', esExecuting = 'Executing', esSuccess = 'Success', esFailed = 'Failed'
}

export enum BUTTON_STYLE {
  CONFIRMATION = 1, DELETION, YES_NO, ONLY_CONFIRM
}

export class Message {
  title = '';
  message = '';
  data: any;
  buttonStyle: BUTTON_STYLE = BUTTON_STYLE.CONFIRMATION;
  returnStatus: RETURN_STATUS = RETURN_STATUS.rsNone;
}

export type AlertType = 'success' | 'danger' | 'info' | 'warning';
export type DropdownMenuPosition = 'bottom-left' | 'bottom-right' | 'top-left' | 'top-right';

export interface IDropdownTag {
  type: AlertType;
  description: string;
}

export class AlertMessage {
  message = '';
  alertType: AlertType = 'success';
}

export enum GlobalAlertType {
  gatNormal, gatShowDetail, gatLogin
}

export class GlobalAlertMessage {
  type: GlobalAlertType = GlobalAlertType.gatNormal;
  message = '';
  alertType: AlertType = 'danger';
  errorObject: HttpErrorResponse | Type<Error> | TimeoutError;
  endMessage = '';
}

export class SharedProject extends HttpBase {
  @HttpBind('project_id') projectId = 0;
  @HttpBind('project_name') projectName = '';
  @HttpBind('project_public') projectPublic = 0;
  @HttpBind('project_creation_time') creationTime: Date;
  @HttpBind('project_update_time') updateTime: Date;
  @HttpBind('project_comment') projectComment = '';
  @HttpBind('project_owner_id') projectOwnerId = 0;
  @HttpBind('project_owner_name') projectOwnerName = '';
  @HttpBind('project_current_user_role_id') currentUserRoleId = 0;
  @HttpBind('project_deletable') projectDeletable = false;
  @HttpBind('project_deleted') projectDeleted = 0;
  @HttpBind('project_istio_support') istioSupport = false;
  @HttpBind('project_service_count') serviceCount = 0;
  @HttpBind('project_toggleable') toggleable = false;
}

export class SharedCreateProject extends HttpBase {
  @HttpBind('project_name') projectName = '';
  @HttpBind('project_comment') comment = '';
  @HttpBindBoolean('project_public', 1, 0) publicity = false;
}

export class SharedMember extends HttpBase {
  @HttpBind('project_member_id') id = 0;
  @HttpBind('project_member_user_id') userId = 0;
  @HttpBind('project_member_username') userName = '';
  @HttpBind('project_member_role_id') roleId = 0;
  isMember?: boolean;
}

export class SharedRole extends HttpBase {
  @HttpBind('role_id') roleId = 0;
  @HttpBind('role_name') roleName = '';
}

export class SharedEnvType {
  envName = '';
  envValue = '';
  envConfigMapName = '';
  envConfigMapKey = '';
}

export class SharedConfigMap extends HttpBase {
  @HttpBind('namespace') namespace = '';
  @HttpBind('name') name = '';
  @HttpBind('datalist') data: object;
  dataList: Array<{ key: string, value: string }>;

  protected prepareInit() {
    this.dataList = Array<{ key: string, value: string }>();
  }

  protected afterInit() {
    if (this.data) {
      Reflect.ownKeys(this.data).forEach((key: string) =>
        this.dataList.push({key, value: Reflect.get(this.data, key)})
      );
    }
  }

  getPostBody(): { [p: string]: any } {
    const obj = Object.create({});
    this.dataList.forEach(value =>
      Object.defineProperties(obj, {[value.key]: {enumerable: true, value: value.value}})
    );
    return {
      namespace: this.namespace,
      name: this.name,
      datalist: obj
    };
  }
}

export class SharedConfigMapDetail extends HttpBase {
  @HttpBindObject('metadata', ConfigMapDetailMetadata) metadata: ConfigMapDetailMetadata;
  @HttpBind('data') data: object;
  dataList: Array<{ key: string, value: string }>;

  protected prepareInit() {
    this.metadata = new ConfigMapDetailMetadata();
    this.dataList = Array<{ key: string, value: string }>();
  }

  protected afterInit() {
    if (this.data) {
      Reflect.ownKeys(this.data).forEach((key: string) =>
        this.dataList.push({key, value: Reflect.get(this.data, key)})
      );
    }
  }
}


export interface INode {
  node_name: string;
  node_ip: string;
  status: number;
}

export interface INodeGroup {
  nodegroup_id: number;
  nodegroup_project: string;
  nodegroup_name: string;
  nodegroup_comment: string;
}

export class NodeAvailableResources {
  node_id = 0;
  node_name = '';
  cpu_available = '';
  mem_available = '';
  storage_available = '';
}

export class SystemInfo {
  board_host = '';
  auth_mode = '';
  set_auth_password = '';
  init_project_repo = '';
  sync_k8s = '';
  redirection_url = '';
  board_version = '';
  kubernetes_version = '';
  processor_type = '';
  dns_suffix = '';
  mode = '';

  constructor() {
  }
}

export class User {
  user_id = 0;
  user_name = '';
  user_email = '';
  user_password = '';
  user_confirm_password = '';
  user_realname = '';
  user_comment = '';
  user_deleted = 0;
  user_system_admin = 0;
  user_reset_uuid = '';
  user_salt: string = '';
  user_creation_time: Date;
  user_update_time: Date;

  constructor() {
    this.user_creation_time = new Date();
    this.user_update_time = new Date();
  }
}

export enum DragStatus {
  dsReady = 'ready', dsStart = 'start', dsDragIng = 'drag', dsEnd = 'end'
}

export enum PvAccessMode {
  ReadWriteOnce = 'ReadWriteOnce',
  ReadOnlyMany = 'ReadOnlyMany',
  ReadWriteMany = 'ReadWriteMany'
}

export enum PvcAccessMode {
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
  @HttpBind('pv_type') type = 0;
  @HttpBind('pv_state') state = 0;
  @HttpBind('pv_accessmode') accessMode: PvAccessMode;
  @HttpBind('pv_reclaim') reclaim: PvReclaimMode;
  capacity = 0;

  protected afterInit() {
    this.capacity = Number.parseFloat(Reflect.get(this.res, 'pv_capacity'));
  }

  getPostBody(): object {
    return {
      pv_name: this.name,
      pv_capacity: `${this.capacity}Gi`,
      pv_type: this.type,
      pv_accessmode: this.accessMode,
      pv_reclaim: this.reclaim
    };
  }
}

export class PersistentVolumeClaim extends HttpBase {
  @HttpBind('pvc_id') id = 0;
  @HttpBind('pvc_name') name = '';
  @HttpBind('pvc_projectid') projectId = 0;
  @HttpBind('pvc_projectname') projectName = '';
  @HttpBind('pvc_state') state = 0;
  @HttpBind('pvc_accessmode') accessMode: PvcAccessMode;
  @HttpBind('pvc_class') class = '';
  @HttpBind('pvc_designatedpv') designatedPv = '';
  volume = '';
  events: Array<string>;
  capacity = 0;

  protected afterInit() {
    this.capacity = Number.parseFloat(Reflect.get(this.res, 'pvc_capacity'));
  }

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

  getPostBody(): object {
    return {
      pvc_name: this.name,
      pvc_projectid: this.projectId,
      pvc_capacity: `${this.capacity}Gi`,
      pvc_accessmode: this.accessMode,
      pvc_class: this.class,
      pvc_designatedpv: this.designatedPv
    };
  }
}

export class Tools {
  static isValidString(str: string, reg?: RegExp): boolean {
    if (str === undefined || str == null || str.trim() === '') {
      return false;
    } else if (reg) {
      return reg.test(str);
    }
    return true;
  }

  static isInvalidString(str: string, reg?: RegExp): boolean {
    return !Tools.isValidString(str, reg);
  }

  static isValidObject(obj: any): boolean {
    return obj !== null && obj !== undefined && typeof obj === 'object';
  }

  static isInvalidObject(obj: any): boolean {
    return !Tools.isValidObject(obj);
  }

  static isValidArray(obj: any): boolean {
    return Array.isArray(obj);
  }

  static isInvalidArray(obj: any): boolean {
    return !Tools.isValidArray(obj);
  }
}
