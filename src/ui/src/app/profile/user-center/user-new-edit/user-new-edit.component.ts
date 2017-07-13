import {
  Component,
  Input,
  Output,
  EventEmitter,
  AfterViewChecked
} from "@angular/core";
import { User } from 'app/profile/user-center/user';
import { UserService } from "../user-service/user-service"

export enum editModel {
  emNew = 0,
  emEdit = 1
}

@Component({
  selector: "new-user",
  templateUrl: "./user-new-edit.component.html",
  styleUrls: ["./user-new-edit.component.css"]
})
export class NewUser implements AfterViewChecked {
  _isOpen: boolean;
  afterCommitErr: string = "";
  constructor(private userService: UserService) { };
  @Input() userModel: User;
  @Input() CurEditModel: editModel;
  @Input()
    
  get isOpen() {
    return this._isOpen;
  }
  
  set isOpen(open: boolean) {
    this._isOpen = open;
    this.isOpenChange.emit(this._isOpen);
  }
  @Output() isOpenChange: EventEmitter<boolean> = new EventEmitter<boolean>();
  @Output() SubmitSuccessEvent: EventEmitter<any> = new EventEmitter<any>();

  ngAfterViewChecked() {
    this.afterCommitErr = "";
  }

  get Title() {
    return this.CurEditModel == editModel.emNew
      ? "USER_CENTER.ADD_USER"
      : "USER_CENTER.EDIT_USER";
  }

  get ActionCaption() {
     return this.CurEditModel == editModel.emNew
      ? "USER_CENTER.ADD"
      : "USER_CENTER.EDIT";
  }

  submitUser() {
    this.CurEditModel == editModel.emEdit ? this.editUser() : this.addNewUser();
  }

  editUser() {
    this.userService.editUser(this.userModel).then(
      res => {
        this.SubmitSuccessEvent.emit(true);
        this.isOpen = false
      },
      (res: Response) => {
        switch (res.status) {
          default: this.afterCommitErr = `${res.status}:服务器错误`;
        }
      })
  }

  addNewUser() {
    this.userService.newUser(this.userModel).then(
      res => {
        this.SubmitSuccessEvent.emit(true);
        this.isOpen = false
      },
      (res: Response) => {
        switch (res.status) {
          case 409: {
            this.afterCommitErr = `${res.status}:用户名或邮箱重复`;
            break;
          }
          default: this.afterCommitErr = `${res.status}:服务器错误`;
        }
      })
  }

}
