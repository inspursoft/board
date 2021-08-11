import { Component, ComponentFactoryResolver, OnInit, ViewContainerRef } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { ImageService } from '../image.service';
import { MessageService } from '../../shared.service/message.service';
import { AppInitService } from '../../shared.service/app-init.service';
import { SharedActionService } from '../../shared.service/shared-action.service';
import { CsModalParentBase } from '../../shared/cs-modal-base/cs-modal-parent-base';
import { TranslateService } from '@ngx-translate/core';
import { Message, RETURN_STATUS } from '../../shared/shared.types';
import { CreateImageComponent } from '../image-create/image-create.component';
import { CreateImageMethod, Image, ImageProject } from '../image.types';

@Component({
  templateUrl: './image-list.component.html',
  styleUrls: ['./image-list.component.css']
})
export class ImageListComponent extends CsModalParentBase implements OnInit {
  curImage: Image;
  isShowDetail = false;
  isBuildImageWIP = false;
  selectedProjectName = '';
  selectedProjectId = 0;
  imageListErrMsg = '';
  imageList: Array<Image> = Array<Image>();
  imageCountPerPage = 10;
  loadingWIP: boolean;
  projectsList: Array<ImageProject>;
  createImageMethod: CreateImageMethod = CreateImageMethod.None;

  constructor(private activatedRoute: ActivatedRoute,
              private imageService: ImageService,
              private sharedActionService: SharedActionService,
              private translateService: TranslateService,
              private view: ViewContainerRef,
              private resolver: ComponentFactoryResolver,
              private appInitService: AppInitService,
              private messageService: MessageService) {
    super(resolver, view);
    this.projectsList = Array<ImageProject>();
  }

  ngOnInit() {
    this.imageService.getProjects().subscribe((res: Array<ImageProject>) => this.projectsList = res);
    this.retrieve();
    this.activatedRoute.fragment.subscribe((fragment: string) => {
      if (fragment === 'createImage') {
        this.createImage();
      }
    });
  }

  get isSystemAdmin(): boolean {
    return this.appInitService.currentUser.userSystemAdmin === 1;
  }

  clickSelectProject() {
    this.sharedActionService.createProjectComponent(this.selfView).subscribe((projectName: string) => {
      if (projectName) {
        this.imageService.getProjects(projectName).subscribe((res: Array<ImageProject>) => {
          this.selectedProjectId = res[0].projectId;
          this.selectedProjectName = res[0].projectName;
          this.projectsList.push(res[0]);
        });
      }
    });
  }

  changeSelectProject(project: ImageProject) {
    this.selectedProjectName = project.projectName;
    this.selectedProjectId = project.projectId;
  }

  retrieve() {
    this.loadingWIP = true;
    this.imageService.getImages('', 0, 0).subscribe((res: Array<Image>) => {
        this.loadingWIP = false;
        this.imageList = res || [];
      }, () => this.loadingWIP = false
    );
  }

  showImageDetail(image: Image) {
    this.curImage = image;
    this.isShowDetail = true;
  }

  confirmToDeleteImage(imageName: string) {
    if (this.isSystemAdmin) {
      this.translateService.get('IMAGE.CONFIRM_TO_DELETE_IMAGE', [imageName]).subscribe((msg: string) => {
        this.messageService.showDeleteDialog(msg, 'IMAGE.DELETE_IMAGE').subscribe((message: Message) => {
          if (message.returnStatus === RETURN_STATUS.rsConfirm) {
            this.imageService.deleteImages(imageName).subscribe(() => {
              this.messageService.showAlert('IMAGE.SUCCESSFUL_DELETED_IMAGE');
              this.retrieve();
            });
          }
        });
      });
    }
  }

  createImage() {
    this.isBuildImageWIP = true;
    this.selectedProjectName = '';
    this.selectedProjectId = 0;
    this.createImageMethod = CreateImageMethod.None;
  }

  setCreateImageMethod(method: CreateImageMethod): void {
    this.createImageMethod = method;
  }

  createNewImage() {
    const component = this.createNewModal(CreateImageComponent);
    component.initCustomerNewImage(this.selectedProjectId, this.selectedProjectName);
    component.initBuildMethod(this.createImageMethod);
    component.closeNotification.subscribe((res: any) => {
        this.createImageMethod = CreateImageMethod.None;
        this.isBuildImageWIP = false;
    });
    component.refreshNotification.subscribe(() => this.retrieve());
  }
}
