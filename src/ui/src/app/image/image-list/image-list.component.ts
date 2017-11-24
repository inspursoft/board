import { OnInit, Component, OnDestroy } from '@angular/core';
import { Subscription } from 'rxjs/Subscription';
import { Image } from "../image";
import { ImageService } from "../image-service/image-service"
import { MessageService } from "../../shared/message-service/message.service";
import { MESSAGE_TARGET, BUTTON_STYLE, MESSAGE_TYPE } from '../../shared/shared.const';
import { Message } from '../../shared/message-service/message';
import { AppInitService } from "../../app.init.service";
import { Project } from "../../project/project";
import { Router } from "@angular/router";

@Component({
  selector: 'image-list',
  templateUrl: './image-list.component.html',
  styleUrls: ["./image-list.component.css"]
})
export class ImageListComponent implements OnInit, OnDestroy {
  curImage: Image;
  isShowDetail: boolean = false;
  isInBuildingImage: boolean = false;
  isOpenNewImage: boolean = false;
  selectedProjectName: string = "";
  selectedProjectId: number = 0;
  imageListErrMsg: string = "";
  imageList: Image[] = Array<Image>();
  imageCountPerPage: number = 10;
  loadingWIP: boolean;
  projectsList: Array<Project>;
  _subscription: Subscription;

  constructor(private imageService: ImageService,
              private router: Router,
              private appInitService: AppInitService,
              private messageService: MessageService) {
    this.projectsList = Array<Project>();
    this._subscription = this.messageService.messageConfirmed$.subscribe(m => {
      let confirmationMessage = <Message>m;
      if (confirmationMessage) {
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
            m.message = 'IMAGE.FAILED_TO_DELETE_IMAGE';
            m.type = MESSAGE_TYPE.COMMON_ERROR;
            this.messageService.inlineAlertMessage(m);
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

  get createActionIsDisabled(): boolean {
    let user = this.appInitService.currentUser;
    if (user && user.hasOwnProperty("user_project_admin") && user.hasOwnProperty("user_system_admin")) {
      return user["user_project_admin"] == 0 && user["user_system_admin"] == 0;
    }
    return true;
  }

  clickSelectProject(project: Project) {
    this.router.navigate(["/projects"]);
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
    let announceMessage = new Message();
    announceMessage.title = 'IMAGE.DELETE_IMAGE';
    announceMessage.message = 'IMAGE.CONFIRM_TO_DELETE_IMAGE';
    announceMessage.params = [imageName];
    announceMessage.target = MESSAGE_TARGET.DELETE_IMAGE;
    announceMessage.buttons = BUTTON_STYLE.DELETION;
    announceMessage.data = imageName;
    this.messageService.announceMessage(announceMessage);
  }

  createImage() {
    this.isInBuildingImage = true;
    this.selectedProjectName = "";
    this.selectedProjectId = 0;
  }

  onBuildImageCompleted(imageName: string) {
    this.retrieve();
  }

}