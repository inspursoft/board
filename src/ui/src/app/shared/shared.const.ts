import { ICsMenuItemData } from "./shared.types";

export enum MESSAGE_TARGET {
  TOGGLE_PROJECT, DELETE_PROJECT, TOGGLE_NODE, DELETE_SERVICE,
  TOGGLE_SERVICE, DELETE_IMAGE, CANCEL_BUILD_IMAGE, CANCEL_BUILD_SERVICE,
  CANCEL_BUILD_SERVICE_GUARD, FORCE_QUIT_BUILD_IMAGE, DELETE_TAG, DELETE_NODE_GROUP,
  DELETE_SERVICE_DEPLOYMENT,DELETE_USER,SIGN_IN_ERROR,SIGN_UP_ERROR,SIGN_UP_SUCCESSFUL,FORGOT_PASSWORD,RESET_PASSWORD
}
export const RouteSignIn = 'sign-in';
export const RouteSignUp = 'sign-up';
export const RouteForgotPassword = 'forgot-password';
export const RouteForgotDashboard = 'dashboard';

export const DISMISS_ALERT_INTERVAL: number = 4 * 1000;
export const DISMISS_GLOBAL_ALERT_INTERVAL: number = 10 * 1000;

export enum MESSAGE_TYPE {
  NONE, COMMON_ERROR = 1, HTTP_401, INTERNAL_ERROR, SHOW_DETAIL
}

export const ROLES: {[key: number]: string} = {
  1: 'PROJECT.PROJECT_ADMIN', 2: 'PROJECT.DEVELOPER', 3: 'PROJECT.VISITOR'
};

export enum BUTTON_STYLE {
  CONFIRMATION = 1, DELETION, YES_NO, ONLY_CONFIRM
}

export enum SERVICE_STATUS{
  PREPARING,
  RUNNING,
  STOPPED,
  WARNING
}

export enum GUIDE_STEP{
  NONE_STEP,
  PROJECT_LIST,
  CREATE_PROJECT,
  SERVICE_LIST,
  CREATE_SERVICE
}

export const AUDIT_RECORD_HEADER_KEY = "audit";
export const AUDIT_RECORD_HEADER_VALUE = "true";

export const RouteDashboard = "dashboard";
export const RouteServices = "services";
export const RouteProjects = "projects";
export const RouteNodes = "nodes";
export const RouteImages = "images";
export const RouteUserCenters = "user-center";
export const RouteAudit = "audit";
export const RouteProfile = "profile";
export const MAIN_MENU_DATA: Array<ICsMenuItemData> = [
  {caption: 'SIDE_NAV.DASHBOARD', visible: true, icon: 'dashboard', url: `/${RouteDashboard}`},
  {caption: 'SIDE_NAV.SERVICES', visible: true, icon: 'applications', url: `/${RouteServices}`},
  {caption: 'SIDE_NAV.PROJECTS', visible: true, icon: 'vm', url: `/${RouteProjects}`},
  {caption: 'SIDE_NAV.NODES', visible: true, icon: 'layers', url: `/${RouteNodes}`},
  {caption: 'SIDE_NAV.IMAGES', visible: true, icon: 'cluster', url: `/${RouteImages}`},
  {caption: 'SIDE_NAV.ADMIN_OPTIONS', visible: true, icon: 'administrator', url: `/${RouteUserCenters}`},
  {caption: 'SIDE_NAV.AUDIT', visible: true, icon: 'library', url: `/${RouteAudit}`},
  {caption: 'SIDE_NAV.PROFILES', visible: true, icon: 'help-info', url: `/${RouteProfile}`}
];