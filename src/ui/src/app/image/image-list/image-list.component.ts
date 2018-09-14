import { Component, ComponentFactoryResolver, OnDestroy, OnInit, ViewContainerRef } from '@angular/core';
import { Subscription } from 'rxjs/Subscription';
import { Image } from "../image";
import { ImageService } from "../image-service/image-service"
import { MessageService } from "../../shared/message-service/message.service";
import { AppInitService } from "../../app.init.service";
import { Project } from "../../project/project";
import { SharedActionService } from "../../shared/shared-action.service";
import { SharedService } from "../../shared/shared.service";
import { CsModalParentBase } from "../../shared/cs-modal-base/cs-modal-parent-base";
import { TranslateService } from "@ngx-translate/core";
import { CreateImageMethod, Message, RETURN_STATUS } from "../../shared/shared.types";
import { CreateImageComponent } from "../image-create/image-create.component";

@Component({
  selector: 'image-list',
  templateUrl: './image-list.component.html',
  styleUrls: ["./image-list.component.css"]
})
export class ImageListComponent extends CsModalParentBase implements OnInit, OnDestroy {
  curImage: Image;
  isShowDetail: boolean = false;
  isBuildImageWIP: boolean = false;
  isOpenNewImage: boolean = false;
  selectedProjectName: string = "";
  selectedProjectId: number = 0;
  imageListErrMsg: string = "";
  imageList: Image[] = Array<Image>();
  imageCountPerPage: number = 10;
  loadingWIP: boolean;
  projectsList: Array<Project>;
  createImageMethod: CreateImageMethod = CreateImageMethod.None;
  dropdownDefaultText: string = "";
  _subscription: Subscription;

  constructor(private imageService: ImageService,
              private sharedActionService: SharedActionService,
              private sharedService: SharedService,
              private translateService: TranslateService,
              private view: ViewContainerRef,
              private resolver: ComponentFactoryResolver,
              private appInitService: AppInitService,
              private messageService: MessageService) {
    super(resolver, view);
    this.projectsList = Array<Project>();
  }

  ngOnInit() {
    this.dropdownDefaultText = "IMAGE.CREATE_IMAGE_SELECT_PROJECT";
    this.imageService.getProjects().then((res: Array<Project>) => {
      let createNewProject: Project = new Project();
      createNewProject.project_name = "IMAGE.CREATE_IMAGE_CREATE_PROJECT";
      createNewProject.project_id = -1;
      createNewProject["isSpecial"] = true;
      createNewProject["OnlyClick"] = true;
      this.projectsList.push(createNewProject);
      if (res && res.length > 0) {
        this.projectsList = this.projectsList.concat(res);
      }
    });
    this.retrieve();
  }

  ngOnDestroy() {
    if (this._subscription) {
      this._subscription.unsubscribe();
    }
  }

  get isSystemAdmin(): boolean {
    if(this.appInitService.currentUser) {
      return this.appInitService.currentUser["user_system_admin"] == 1;
    }
    return false;
  }

  setDropdownDefaultText(): void {
    let selected = this.projectsList.find((project: Project) => project.project_id === this.selectedProjectId);
    this.dropdownDefaultText = selected ? selected.project_name : "IMAGE.CREATE_IMAGE_SELECT_PROJECT";
  }

  clickSelectProject() {
    this.sharedActionService.createProjectComponent(this.selfView).subscribe((projectName: string) => {
      if (projectName) {
        this.sharedService.getOneProject(projectName).then((res: Array<Project>) => {
          this.selectedProjectId = res[0].project_id;
          this.selectedProjectName = res[0].project_name;
          let project = this.projectsList.shift();
          this.projectsList.unshift(res[0]);
          this.projectsList.unshift(project);
          this.setDropdownDefaultText();
        })
      }
    });
  }

  changeSelectProject(project: Project) {
    this.selectedProjectName = project.project_name;
    this.selectedProjectId = project.project_id;
    this.setDropdownDefaultText();
  }

  retrieve() {
    this.loadingWIP = true;
    this.imageService.getImages("", 0, 0)
      .then(res => {
        this.loadingWIP = false;
        this.imageList = res || [];
      })
      .catch(() => this.loadingWIP = false);
  }

  showImageDetail(image: Image) {
    //need add get one Image from server
    this.curImage = image;
    this.isShowDetail = true;
  }

  confirmToDeleteImage(imageName: string) {
    if (this.isSystemAdmin){
      this.translateService.get('IMAGE.CONFIRM_TO_DELETE_IMAGE', [imageName]).subscribe((msg: string) => {
        this.messageService.showDeleteDialog(msg, 'IMAGE.DELETE_IMAGE').subscribe((message: Message) => {
          if (message.returnStatus == RETURN_STATUS.rsConfirm) {
            this.imageService.deleteImages(imageName).then(() => {
              this.messageService.showAlert('IMAGE.SUCCESSFUL_DELETED_IMAGE');
              this.retrieve();
            })
          }
        })
      })
    }
  }

  createImage() {
    this.isBuildImageWIP = true;
    this.selectedProjectName = "";
    this.selectedProjectId = 0;
    this.dropdownDefaultText = "IMAGE.CREATE_IMAGE_SELECT_PROJECT";
    this.createImageMethod = CreateImageMethod.None;
  }

  setCreateImageMethod(method: CreateImageMethod): void {
    this.createImageMethod = method;
  }

  createNewImage() {
    let component = this.createNewModal(CreateImageComponent);
    component.initCustomerNewImage(this.selectedProjectId, this.selectedProjectName);
    component.initBuildMethod(this.createImageMethod);
    component.closeNotification.subscribe((res: any) => {
      this.isBuildImageWIP = false;
      this.createImageMethod = CreateImageMethod.None;
      if (res) {
        this.retrieve();
      }
    })
  }
}