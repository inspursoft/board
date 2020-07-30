import { Component, EventEmitter, Input, Output } from '@angular/core';
import { UserService } from '../user-service/user-service';
import { CsComponentBase } from '../../../shared/cs-components-library/cs-component-base';
import { User } from '../../../shared/shared.types';
import { MessageService } from '../../../shared.service/message.service';

export enum editModel { emNew, emEdit }
@Component({
  selector: 'app-new-user',
  templateUrl: './user-new-edit.component.html',
  styleUrls: ['./user-new-edit.component.css']
})
export class NewEditUserComponent extends CsComponentBase {
  isOpenValue: boolean;
  isWorkWIP = false;

  constructor(private userService: UserService,
              private messageService: MessageService) {
    super();
  }

  @Input() userModel: User;
  @Input() CurEditModel: editModel;

  @Input()
  get isOpen() {
    return this.isOpenValue;
  }

  set isOpen(open: boolean) {
    this.isOpenValue = open;
    this.isOpenChange.emit(this.isOpenValue);
  }

  @Output() isOpenChange: EventEmitter<boolean> = new EventEmitter<boolean>();
  @Output() SubmitSuccessEvent: EventEmitter<any> = new EventEmitter<any>();

  get Title() {
    return this.CurEditModel === editModel.emNew
      ? 'USER_CENTER.ADD_USER'
      : 'USER_CENTER.EDIT_USER';
  }

  get ActionCaption() {
    return this.CurEditModel === editModel.emNew
      ? 'USER_CENTER.ADD'
      : 'USER_CENTER.SAVE';
  }

  submitUser() {
    if (this.verifyInputExValid()) {
      this.isWorkWIP = true;
      this.CurEditModel === editModel.emEdit ? this.updateUser() : this.addNewUser();
    }
  }

  updateUser() {
    this.userService.updateUser(this.userModel).subscribe(() => {
        this.SubmitSuccessEvent.emit(true);
        this.isOpen = false;
        this.messageService.showAlert('USER_CENTER.EDIT_USER_SUCCESS');
      },
      () => this.isOpen = false
    );
  }

  addNewUser() {
    this.userModel.userCreationTime = new Date(Date.now()).toISOString();
    this.userModel.userUpdateTime = new Date(Date.now()).toISOString();
    console.log(this.userModel);
    this.userService.newUser(this.userModel).subscribe(() => {
        this.SubmitSuccessEvent.emit(true);
        this.isOpen = false;
        this.messageService.showAlert('USER_CENTER.ADD_USER_SUCCESS');
      },
      () => this.isOpen = false
    );
  }

}
