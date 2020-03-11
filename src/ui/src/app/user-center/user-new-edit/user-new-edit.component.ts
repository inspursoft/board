import { Component, EventEmitter, Input, Output } from '@angular/core';
import { UserService } from '../user-service/user-service';
import { CsComponentBase } from '../../shared/cs-components-library/cs-component-base';
import { User } from '../../shared/shared.types';
import { MessageService } from '../../shared.service/message.service';

export enum editModel { emNew, emEdit }
@Component({
  selector: 'new-user',
  templateUrl: './user-new-edit.component.html',
  styleUrls: ['./user-new-edit.component.css']
})
export class NewEditUserComponent extends CsComponentBase {
  _isOpen: boolean;
  isWorkWIP = false;

  constructor(private userService: UserService,
              private messageService: MessageService) {
    super();
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
      ? 'USER_CENTER.ADD_USER'
      : 'USER_CENTER.EDIT_USER';
  }

  get ActionCaption() {
    return this.CurEditModel == editModel.emNew
      ? 'USER_CENTER.ADD'
      : 'USER_CENTER.SAVE';
  }

  submitUser() {
    if (this.verifyInputExValid()) {
      this.isWorkWIP = true;
      this.CurEditModel == editModel.emEdit ? this.updateUser() : this.addNewUser();
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
    this.userService.newUser(this.userModel).subscribe(() => {
        this.SubmitSuccessEvent.emit(true);
        this.isOpen = false;
        this.messageService.showAlert('USER_CENTER.ADD_USER_SUCCESS');
      },
      () => this.isOpen = false
    );
  }

}
