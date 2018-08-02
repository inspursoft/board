import { Component, OnDestroy, OnInit, ViewContainerRef } from '@angular/core';
import { Subscription } from 'rxjs/Subscription';
import { Image } from "../image";
import { ImageService } from "../image-service/image-service"
import { MessageService } from "../../shared/message-service/message.service";
import { BUTTON_STYLE, MESSAGE_TARGET } from '../../shared/shared.const';
import { Message } from '../../shared/message-service/message';
import { AppInitService } from "../../app.init.service";
import { Project } from "../../project/project";
import { SharedActionService } from "../../shared/shared-action.service";
import { SharedService } from "../../shared/shared.service";

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
  dropdownDefaultText: string = "";
  _subscription: Subscription;

  constructor(private imageService: ImageService,
              private sharedActionService: SharedActionService,
              private sharedService: SharedService,
              private selfView: ViewContainerRef,
              private appInitService: AppInitService,
              private messageService: MessageService) {
    this.projectsList = Array<Project>();
    this._subscription = this.messageService.messageConfirmed$.subscribe((msg:Message) => {
      if (msg.target == MESSAGE_TARGET.DELETE_IMAGE) {
        let imageName = <string>msg.data;
        this.imageService
          .deleteImages(imageName)
          .then(() => {
            let m: Message = new Message();
            m.message = 'IMAGE.SUCCESSFUL_DELETED_IMAGE';
            this.messageService.inlineAlertMessage(m);
            this.retrieve();
          })
          .catch(err => this.messageService.dispatchError(err));
      }
    });
  }

  ngOnInit() {
    this.dropdownDefaultText = "IMAGE.CREATE_IMAGE_SELECT_PROJECT";
    this.imageService.getProjects()
      .then((res: Array<Project>) => {
        let createNewProject: Project = new Project();
        createNewProject.project_name = "IMAGE.CREATE_IMAGE_CREATE_PROJECT";
        createNewProject.project_id = -1;
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
    this.dropdownDefaultText = "IMAGE.CREATE_IMAGE_SELECT_PROJECT";
    this.createImageMethod = CreateImageMethod.None;
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