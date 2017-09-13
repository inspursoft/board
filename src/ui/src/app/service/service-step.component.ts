export interface ServiceStepComponent {
  data: any;
}

export class ServiceEnvOutput {
  constructor(public key: string = "",
              public value: string = "") {
  }
}

export class ServiceStep1Output {
  constructor(public  project_id: number = 0,
              public project_name: string = "") {
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
    image_copy?: Array<{dockerfile_copyfrom?: string, dockerfile_copyto?: string}>,
    image_run?: Array<string>,
    image_env?: Array<{dockerfile_envname?: string, dockerfile_envvalue?: string}>,
    image_expose?: Array<number>
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
      image_volume: Array<string>(),
      image_run: Array<string>(),
      image_expose: Array<number>(),
      image_copy: Array<{dockerfile_copyfrom?: string, dockerfile_copyto?: string}>(),
      image_env: Array<{dockerfile_envname?: string, dockerfile_envvalue?: string}>()
    }
  }
}
