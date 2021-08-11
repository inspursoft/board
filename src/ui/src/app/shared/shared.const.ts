import { ICsMenuItemData } from './shared.types';

export const DISMISS_ALERT_INTERVAL = 4;
export const DISMISS_CHECK_DROPDOWN = 2000;

export enum SERVICE_STATUS {
  PREPARING,
  RUNNING,
  STOPPED,
  WARNING,
  Deploying,
  Completed,
  Failed,
  DELETED,
  UnKnown = 8,
  AutonomousOffline = 9,
  PartAutonomousOffline
}

export enum GUIDE_STEP {
  NONE_STEP,
  PROJECT_LIST,
  CREATE_PROJECT,
  SERVICE_LIST,
  CREATE_SERVICE
}

export const AUDIT_RECORD_HEADER_KEY = 'audit';
export const AUDIT_RECORD_HEADER_VALUE = 'true';
export const RouteInitialize = 'initialize-page';
export const RouteSignIn = 'account/sign-in';
export const RouteSignUp = 'account/sign-up';
export const RouteForgotPassword = 'account/forgot-password';
export const RouteDashboard = 'dashboard';
export const RouteServices = 'services';
export const RouteHelm = 'helm';
export const RouteRepoList = 'repo-list';
export const RouteReleaseList = 'release-list';
export const RouteProjects = 'projects';
export const RouteNodes = 'nodes';
export const RouteImages = 'images';
export const RouteTrainingJob = 'training-job';
export const RouteResource = 'resource';
export const RouteConfigMap = 'config-map';
export const RouteUserCenters = 'user-center';
export const RouteAdmin = 'admin';
export const RouteSystemSetting = 'system-setting';
export const RouteAudit = 'audit';
export const RouteKibana = 'kibana-url';
export const RouteProfile = 'profile';
export const RouteStorage = 'storage';
export const RoutePV = 'pv';
export const RoutePvc = 'pvc';
export const MAIN_MENU_DATA: Array<ICsMenuItemData> = [
  {
    caption: 'SIDE_NAV.DASHBOARD', visible: true, icon: 'dashboard', url: `/${RouteDashboard}`, children: [
      {caption: 'SIDE_NAV.DASHBOARD', visible: true, icon: 'dashboard', url: `/${RouteDashboard}`},
      {caption: 'SIDE_NAV.KIBANA', visible: true, icon: 'curve-chart', url: `/${RouteKibana}`}]
  },
  {caption: 'SIDE_NAV.SERVICES', visible: true, icon: 'applications', url: `/${RouteServices}`},
  {
    caption: 'SIDE_NAV.HELM', visible: true, icon: 'shopping-cart', url: `/${RouteHelm}`, children: [
      {caption: 'SIDE_NAV.MARKET_LIST', visible: true, icon: 'list', url: `/${RouteHelm}/${RouteRepoList}`},
      {caption: 'SIDE_NAV.SERVICE_LIST', visible: true, icon: 'list', url: `/${RouteHelm}/${RouteReleaseList}`}]
  },
  {caption: 'SIDE_NAV.PROJECTS', visible: true, icon: 'vm', url: `/${RouteProjects}`},
  {caption: 'SIDE_NAV.NODES', visible: true, icon: 'layers', url: `/${RouteNodes}`},
  {caption: 'SIDE_NAV.IMAGES', visible: true, icon: 'cluster', url: `/${RouteImages}`},
  {
    caption: 'SIDE_NAV.RESOURCE', visible: true, icon: 'euro', url: `/${RouteResource}`, children: [
      {caption: 'SIDE_NAV.CONFIG_MAP', visible: true, icon: 'map', url: `/${RouteResource}/${RouteConfigMap}`}]
  },
  {caption: 'SIDE_NAV.TRAINING_JOB', visible: true, icon: 'event', url: `/${RouteTrainingJob}`},
  {
    caption: 'SIDE_NAV.STORAGE', visible: true, icon: 'storage', url: `/${RouteStorage}`, children: [
      {caption: 'PV', visible: true, icon: '', url: `/${RouteStorage}/${RoutePV}`},
      {caption: 'PVC', visible: true, icon: '', url: `/${RouteStorage}/${RoutePvc}`}]
  },
  {
    caption: 'SIDE_NAV.ADMIN_OPTIONS', visible: true, icon: 'administrator', url: `/${RouteAdmin}`, children: [
      {caption: 'SIDE_NAV.SYSTEM_SETTING', visible: true, icon: 'cog', url: `/${RouteAdmin}/${RouteSystemSetting}`},
      {caption: 'SIDE_NAV.USER_MANAGEMENT', visible: true, icon: 'users', url: `/${RouteAdmin}/${RouteUserCenters}`},
      {caption: '', visible: true, icon: '', url: '', isAdminServer: true}
    ]
  },
  {caption: 'SIDE_NAV.AUDIT', visible: true, icon: 'library', url: `/${RouteAudit}`},
  {caption: 'SIDE_NAV.PROFILES', visible: true, icon: 'help-info', url: `/${RouteProfile}`}
];

export const UsernameInUseKey: Array<string> = ['explore', 'create', 'assets', 'css', 'img', 'js', 'less', 'plugins', 'debug', 'raw',
  'install', 'api', 'avatar', 'user', 'org', 'help', 'stars', 'issues', 'pulls', 'commits', 'repo', 'template', 'new', '.', '..'];
export const InvalidServiceName: Array<string> = [
  'istio-egressgateway',
  'istio-ingress',
  'istio-ingressgateway',
  'istio-pilot',
  'istio-policy',
  'istio-statsd-prom-bridge',
  'istio-telemetry',
  'prometheus',
  'kube-dns',
  'tiller-deploy',
  'istio-citadel'
];
