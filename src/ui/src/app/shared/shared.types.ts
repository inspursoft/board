import { HttpErrorResponse } from "@angular/common/http";
import { Type } from "@angular/core";
import { TimeoutError } from "rxjs/src/util/TimeoutError";
import * as isObject from "isobject";

export interface ICsMenuItemData {
  caption: string,
  icon: string,
  url: string,
  visible: boolean,
  children?: Array<ICsMenuItemData>
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
  title: string = '';
  message: string = '';
  data: any;
  buttonStyle: BUTTON_STYLE = BUTTON_STYLE.CONFIRMATION;
  returnStatus: RETURN_STATUS = RETURN_STATUS.rsNone;
}

export type AlertType = 'alert-success' | 'alert-danger' | 'alert-info' | 'alert-warning';
export type DropdownMenuPosition = 'bottom-left' | 'bottom-right' | 'top-left' | 'top-right';

export interface IDropdownTag {
  type: AlertType,
  description: string
}

export class AlertMessage {
  message: string = '';
  alertType: AlertType = 'alert-success';
}

export enum GlobalAlertType {
  gatNormal, gatShowDetail, gatLogin
}

export class GlobalAlertMessage {
  type: GlobalAlertType = GlobalAlertType.gatNormal;
  message: string = '';
  alertType: AlertType = 'alert-danger';
  errorObject: HttpErrorResponse | Type<Error> | TimeoutError;
  endMessage: string = '';
}

export class SignUp {
  username: string;
  email: string;
  password: string;
  confirmPassword: string;
  realname: string;
  comment: string;
}

export interface INode {
  node_name: string;
  node_ip: string;
  status: number;
}

export interface INodeGroup {
  nodegroup_id: number,
  nodegroup_project: string,
  nodegroup_name: string,
  nodegroup_comment: string;
}

export class NodeAvailableResources {
  node_id: number = 0;
  node_name: string = '';
  cpu_available: string = '';
  mem_available: string = '';
  storage_available: string = '';
}

export class ServiceHPA {
  hpa_id: number;     //The hpa ID of this autoscale. ,
  hpa_name: string = '';   //The hpa name of this autoscale. ,
  hpa_status: number = 0;
  service_id: number; //The service ID of this hpa to control. ,
  min_pod: number = 1;     //The minimum pod number. ,
  max_pod: number = 1;     //The maximum pod number. ,
  cpu_percent: number = 0;//The target CPU percentage.
  isEdit: boolean = false;
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
  dns_suffix = '';
  constructor() {
  }
}

export class User {
  public user_id = 0;
  public user_name = '';
  public user_email = '';
  public user_password = '';
  public user_confirm_password = '';
  public user_realname = '';
  public user_comment = '';
  public user_deleted = 0;
  public user_system_admin = 0;
  public user_reset_uuid = '';
  public user_salt: string = '';
  public user_creation_time: Date;
  public user_update_time: Date;

  constructor() {
    this.user_creation_time = new Date();
    this.user_update_time = new Date();
  }
}

export enum DragStatus {
  dsReady = 'ready', dsStart = 'start', dsDragIng = 'drag', dsEnd = 'end'
}

export enum CreateImageMethod {None, Template, DockerFile, ImagePackage}

export class PersistentVolumeOptions {
  public path = '';
  public server = '';
}

export class PersistentVolumeOptionsRBD {
  public user = '';
  public keyring = '';
  public pool = 'rbd';
  public image = '';
  public fstype = '';
  public secretname = '';
  public secretnamespace = '';
  public monitors = '';
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

export class PersistentVolume {
  public id = 0;
  public name = '';
  public type = 0;
  public state = 0;
  public capacity = 0;
  public accessMode: PvAccessMode;
  public reclaim: PvReclaimMode;

  initFromRes(res: Object) {
    if (res) {
      this.id = Reflect.get(res, 'pv_id');
      this.name = Reflect.get(res, 'pv_name');
      this.type = Reflect.get(res, 'pv_type');
      this.state = Reflect.get(res, 'pv_state');
      this.capacity = Number.parseFloat(Reflect.get(res, 'pv_capacity'));
      this.accessMode = Reflect.get(res, 'pv_accessmode');
      this.reclaim = Reflect.get(res, 'pv_reclaim');
    }
  }

  get typeDescription(): string {
    return ['Unknown', 'NFS', 'RBD'][this.type];
  }

  get statusDescription(): string{
    return [
      'STORAGE.STATE_UNKNOWN',
      'STORAGE.STATE_PENDING',
      'STORAGE.STATE_AVAILABLE',
      'STORAGE.STATE_BOUND',
      'STORAGE.STATE_RELEASED',
      'STORAGE.STATE_FAILED',
      'STORAGE.STATE_INVALID'][this.state];
  }

  postObject(): Object {
    return {
      pv_name: this.name,
      pv_capacity: `${this.capacity}Gi`,
      pv_type: this.type,
      pv_accessmode: this.accessMode,
      pv_reclaim: this.reclaim
    }
  }
}

export class NFSPersistentVolume extends PersistentVolume {
  public options: PersistentVolumeOptions;

  constructor() {
    super();
    this.type = 1;
    this.options = new PersistentVolumeOptions();
  }

  postObject(): Object {
    let result = super.postObject();
    Reflect.set(result, 'pv_options', this.options);
    return result;
  }
}

export class RBDPersistentVolume extends PersistentVolume {
  public options: PersistentVolumeOptionsRBD;

  constructor() {
    super();
    this.type = 2;
    this.options = new PersistentVolumeOptionsRBD();
  }

  postObject(): Object {
    let result = super.postObject();
    Reflect.set(result, 'pv_options', this.options);
    return result;
  }
}

export class PersistentVolumeClaim {
  public id = 0;
  public name = '';
  public projectId = 0;
  public projectName = '';
  public capacity = 0;
  public state = 0;
  public accessMode: PvcAccessMode;
  public class = '';
  public designatedPv = '';
  public volume = '';
  public events: Array<string>;

  constructor() {
    this.events = Array<string>();
  }

  initFromRes(res: Object) {
    if (res) {
      this.id = Reflect.get(res, 'pvc_id');
      this.name = Reflect.get(res, 'pvc_name');
      this.projectId = Reflect.get(res, 'pvc_projectid');
      this.projectName = Reflect.get(res, 'pvc_projectname');
      this.capacity = Number.parseFloat(Reflect.get(res, 'pvc_capacity'));
      this.state = Reflect.get(res, 'pvc_state');
      this.accessMode = Reflect.get(res, 'pvc_accessmode');
      this.class = Reflect.get(res, 'pvc_class');
      this.designatedPv = Reflect.get(res, 'pvc_designatedpv');
    }
  }

  get statusDescription(): string {
    return [
      'STORAGE.STATE_UNKNOWN',
      'STORAGE.STATE_PENDING',
      'STORAGE.STATE_BOUND',
      'STORAGE.STATE_LOST',
      'STORAGE.STATE_INVALID'][this.state];
  }

  postObject(): object {
    return {
      pvc_name: this.name,
      pvc_projectid: this.projectId,
      pvc_capacity: `${this.capacity}Gi`,
      pvc_accessmode: this.accessMode,
      pvc_class: this.class,
      pvc_designatedpv: this.designatedPv
    }
  }
}

export class Tools {
  static isValidString(str: string, reg?: RegExp): boolean {
    if (str == undefined || str == null || str.trim() == '') {
      return false;
    } else if (reg) {
      return reg.test(str)
    }
    return true;
  }

  static isInvalidString(str: string, reg?: RegExp): boolean {
    return !Tools.isValidString(str, reg);
  }

  static isValidObject(obj: any): boolean {
    return obj != null && obj != undefined && typeof obj == 'object';
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