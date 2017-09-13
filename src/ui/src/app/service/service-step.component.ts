export interface ServiceStepComponent {
  data: any;
}

export class ServiceStep1Output {
  public project_name: string;
  public project_id: number;

  constructor() {
    this.project_id = 0;
    this.project_name = "";
  }
}

export class ServiceStep2Output {
  image_name: string;
  image_tag: string;
  project_id: number;
  project_name: string;
  image_template: string;
  image_dockerfile: {
    image_base: string,
    image_author: string,
    image_volume?: Array<string>,
    image_copy?: [{dockerfile_copyfrom?: string, dockerfile_copyto?: string}],
    image_run?: Array<string>,
    image_entrypoint?: string,
    image_cmd?: string
  };

  constructor() {
    this.image_name = "";
    this.image_tag = "";
    this.project_name = "";
    this.image_template = "";
    this.image_dockerfile = {
      image_base: "",
      image_author: "",
      image_volume: [],
      image_run: [],
      image_copy: [{}]
    }
  }
}
