import { HttpBase, HttpBind, HttpBindArray, HttpBindObject } from '../shared/ui-model/model-types';

export enum CreateImageMethod {None, Template, DockerFile, ImagePackage}

export class Image extends HttpBase {
  @HttpBind('image_name') imageName = '';
  @HttpBind('image_comment') imageComment = '';
  @HttpBind('image_deleted') imageDeleted = 0;
  @HttpBind('image_update_time') imageUpdateTime = '';
  @HttpBind('image_creation_time') imageCreationTime = '';
}

export class ImageDetail extends HttpBase {
  @HttpBind('image_name') imageName = '';
  @HttpBind('image_id') imageId = '';
  @HttpBind('image_author') imageAuthor = '';
  @HttpBind('image_tag') imageTag = '';
  @HttpBind('image_creationtime') imageCreationTime = '';
  @HttpBind('image_size_number') imageSizeNumber = 0;
  @HttpBind('image_size_unit') imageSizeUnit = 'MB';
}

export class BuildImageDataBase extends HttpBase {
  @HttpBind('image_name') imageName = '';
  @HttpBind('image_tag') imageTag = '';
  @HttpBind('project_id') projectId = 0;
  @HttpBind('project_name') projectName = '';
  @HttpBind('image_template') imageTemplate = '';
}

export class ImageCopy extends HttpBase {
  @HttpBind('dockerfile_copyfrom') copyFrom = '';
  @HttpBind('dockerfile_copyto') copyTo = '';
}

export class ImageEnv extends HttpBase {
  @HttpBind('dockerfile_envname') envName = '';
  @HttpBind('dockerfile_envvalue') envValue = '';
}

export class BuildImageDockerfileData extends HttpBase {
  @HttpBind('image_base') imageBase = '';
  @HttpBind('image_author') imageAuthor = '';
  @HttpBind('image_entrypoint') imageEntryPoint = '';
  @HttpBind('image_cmd') imageCmd = '';
  @HttpBind('image_volume') imageVolume: Array<string>;
  @HttpBind('image_run') imageRun: Array<string>;
  @HttpBind('image_expose') imageExpose: Array<string>;
  @HttpBindArray('image_copy', ImageCopy) imageCopy: Array<ImageCopy>;
  @HttpBindArray('image_env', ImageEnv) imageEnv: Array<ImageEnv>;

  protected prepareInit() {
    this.imageVolume = Array<string>();
    this.imageRun = Array<string>();
    this.imageExpose = Array<string>();
    this.imageCopy = Array<ImageCopy>();
    this.imageEnv = Array<ImageEnv>();
  }
}

export class BuildImageData extends BuildImageDataBase {
  @HttpBindObject('image_dockerfile', BuildImageDockerfileData) imageDockerFile: BuildImageDockerfileData;

  protected prepareInit() {
    this.imageDockerFile = new BuildImageDockerfileData();
  }
}


export class ImageProject extends HttpBase {
  @HttpBind('project_id') projectId = 0;
  @HttpBind('project_name') projectName: string;
  @HttpBind('publicity') publicity: boolean;
  @HttpBind('project_public') public: number;
  @HttpBind('project_creation_time') creationTime: string;
  @HttpBind('project_comment') comment: string;
  @HttpBind('project_owner_id') ownerId: number;
  @HttpBind('project_owner_name') ownerName: string;
}


