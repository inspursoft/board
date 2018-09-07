import { ICsMenuItemData } from "./shared.types";

export const DISMISS_ALERT_INTERVAL: number = 4 * 1000;

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
export const RouteSignIn = 'sign-in';
export const RouteSignUp = 'sign-up';
export const RouteForgotPassword = 'forgot-password';
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