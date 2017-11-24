export class Image {
  image_name: string = "";
  image_comment: string = "";
  image_deleted: number = 0;

  constructor() {
  }
}

export class ImageDetail {
  image_name: string = "";
  image_tag: string = "";
  image_detail: string = "";
  image_creationtime: string;
  image_size_number: number;
  image_size_unit: string = "MB";

  constructor() {
  }
}

export class BuildImageDataBase {
  image_name: string = "";
  image_tag: string = "";
  project_id: number = 0;
  project_name: string = "";
  image_template: string = "";
}

export class BuildImageDockerfileData {
  image_base: string = "";
  image_author: string = "";
  image_volume?: Array<string>;
  image_copy?: Array<{dockerfile_copyfrom?: string, dockerfile_copyto?: string}>;
  image_run?: Array<string>;
  image_env?: Array<{dockerfile_envname?: string, dockerfile_envvalue?: string}>;
  image_expose?: Array<string>;
  image_entrypoint?: string;
  image_cmd?: string;

  constructor() {
    this.image_volume = Array<string>();
    this.image_run = Array<string>();
    this.image_expose = Array<string>();
    this.image_copy = Array<{dockerfile_copyfrom?: string, dockerfile_copyto?: string}>();
    this.image_env = Array<{dockerfile_envname?: string, dockerfile_envvalue?: string}>();
  }
}

export class BuildImageData extends BuildImageDataBase {
  image_dockerfile: BuildImageDockerfileData;

  constructor() {
    super();
    this.image_dockerfile = new BuildImageDockerfileData();
  }
}

