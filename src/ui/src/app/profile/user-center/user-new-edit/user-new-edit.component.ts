import { Component, Input, Output, EventEmitter, AfterViewChecked } from "@angular/core";
import { User } from 'app/profile/user-center/user';
import { UserService } from "../user-service/user-service"

export enum editModel { emNew, emEdit }

@Component({
  selector: "new-user",
  templateUrl: "./user-new-edit.component.html",
  styleUrls: ["./user-new-edit.component.css"]
})
export class NewUser implements AfterViewChecked {
  _isOpen: boolean;
  isAlertOpen: boolean = false;
  afterCommitErr: string = "";

  constructor(private userService: UserService) {
  };

  ngAfterViewChecked() {
    this.isAlertOpen = false;
  }

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
    this.CurEditModel == editModel.emEdit ? this.updateUser() : this.addNewUser();
  }

  updateUser() {
    this.userService.updateUser(this.userModel)
      .then(() => {
        this.SubmitSuccessEvent.emit(true);
        this.isOpen = false
      })
      .catch((reason: string) => {
        this.afterCommitErr = reason;
        this.isAlertOpen = true;
      });
  }

  addNewUser() {
    this.userService.newUser(this.userModel)
      .then(() => {
        this.SubmitSuccessEvent.emit(true);
        this.isOpen = false;
      })
      .catch((reason: string) => {
        this.afterCommitErr = reason;
        this.isAlertOpen = true;
      })
  }

}
