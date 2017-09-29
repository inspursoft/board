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


