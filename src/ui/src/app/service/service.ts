export enum ServiceType {
  ServiceTypeUnknown,
  ServiceTypeNormalNodePort,
  ServiceTypeHelm,
  ServiceTypeDeploymentOnly,
  ServiceTypeClusterIP,
  ServiceTypeStatefulSet,
  ServiceTypeJob
}

export enum ServiceSource {
  ServiceSourceBoard,
  ServiceSourceK8s,
  ServiceSourceHelm
}

export class Service {
  service_id: number;
  service_name: string;
  service_project_id: number;
  service_project_name: string;
  service_owner_id: number;
  service_owner_name: string;
  service_creation_time: string;
  service_create_time: Date;
  service_public: number;
  service_status: number;
  service_source: ServiceSource;
  service_is_member: number;
  service_type: ServiceType;
}
