import { Component, OnDestroy, OnInit } from '@angular/core';
import { Subscription } from 'rxjs/Subscription';
import { Image } from "../image";
import { ImageService } from "../image-service/image-service"
import { MessageService } from "../../shared/message-service/message.service";
import { BUTTON_STYLE, MESSAGE_TARGET } from '../../shared/shared.const';
import { Message } from '../../shared/message-service/message';
import { AppInitService } from "../../app.init.service";
import { Project } from "../../project/project";
import { Router } from "@angular/router";

enum CreateImageMethod{None, Template, DockerFile, DevOps}
@Component({
  selector: 'image-list',
  templateUrl: './image-list.component.html',
  styleUrls: ["./image-list.component.css"]
})
export class ImageListComponent implements OnInit, OnDestroy {
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
  _subscription: Subscription;

  constructor(private imageService: ImageService,
              private router: Router,
              private appInitService: AppInitService,
              private messageService: MessageService) {
    this.projectsList = Array<Project>();
    this._subscription = this.messageService.messageConfirmed$.subscribe(m => {
      let confirmationMessage = <Message>m;
      if (confirmationMessage && confirmationMessage.target == MESSAGE_TARGET.DELETE_IMAGE) {
        let imageName = <string>confirmationMessage.data;
        let m: Message = new Message();
        this.imageService
          .deleteImages(imageName)
          .then(res => {
            m.message = 'IMAGE.SUCCESSFUL_DELETED_IMAGE';
            this.messageService.inlineAlertMessage(m);
            this.retrieve();
          })
          .catch(err => {
             this.messageService.dispatchError(err);
          });
      }
    });
  }

  ngOnInit() {
    this.imageService.getProjects()
      .then(res => {
        let createNewProject: Project = new Project();
        createNewProject.project_name = "IMAGE.CREATE_IMAGE_CREATE_PROJECT";
        createNewProject["isSpecial"] = true;
        createNewProject["OnlyClick"] = true;
        this.projectsList.push(createNewProject);
        if (res && res.length > 0) {
          this.projectsList = this.projectsList.concat(res);
        }
      })
      .catch(err => this.messageService.dispatchError(err));
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

  clickSelectProject(project: Project) {
    this.router.navigate(["/projects"], {queryParams: {token: this.appInitService.token}, fragment: "create"});
  }

  changeSelectProject(project: Project) {
    this.selectedProjectName = project.project_name;
    this.selectedProjectId = project.project_id;
  }

  retrieve() {
    this.loadingWIP = true;
    this.imageService.getImages("", 0, 0)
      .then(res => {
        this.loadingWIP = false;
        this.imageList = res || [];
      })
      .catch(err => {
        this.loadingWIP = false;
        this.messageService.dispatchError(err)
      });
  }

  showImageDetail(image: Image) {
    //need add get one Image from server
    this.curImage = image;
    this.isShowDetail = true;
  }

  confirmToDeleteImage(imageName: string) {
    if (this.isSystemAdmin){
      let announceMessage = new Message();
      announceMessage.title = 'IMAGE.DELETE_IMAGE';
      announceMessage.message = 'IMAGE.CONFIRM_TO_DELETE_IMAGE';
      announceMessage.params = [imageName];
      announceMessage.target = MESSAGE_TARGET.DELETE_IMAGE;
      announceMessage.buttons = BUTTON_STYLE.DELETION;
      announceMessage.data = imageName;
      this.messageService.announceMessage(announceMessage);
    }
  }

  createImage() {
    this.isBuildImageWIP = true;
    this.selectedProjectName = "";
    this.selectedProjectId = 0;
  }

  onBuildImageCompleted(imageName: string) {
    this.isBuildImageWIP = false;
    this.createImageMethod = CreateImageMethod.None;
    this.retrieve();
  }

  setCreateImageMethod(method: CreateImageMethod): void {
    this.createImageMethod = method;
  }
}