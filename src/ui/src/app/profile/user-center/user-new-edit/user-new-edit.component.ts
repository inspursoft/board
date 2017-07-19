import {Component, Input, Output, EventEmitter, OnDestroy} from "@angular/core";
import {User, user} from 'app/profile/user-center/user';
import {UserService} from "../user-service/user-service"
import {MessageService} from "../../../shared/service/message.service";

export enum editModel {
	emNew = 0,
	emEdit = 1
}

@Component({
	selector: "new-user",
	templateUrl: "./user-new-edit.component.html",
	styleUrls: ["./user-new-edit.component.css"]
})
export class NewUser implements OnDestroy {
	_isOpen: boolean;
	_afterCommitErrInterval: any;
	_afterCommitErrSeed: number = 0;
	afterCommitErr: string = "";

	constructor(private userService: UserService,
							private messageService: MessageService) {
		this._afterCommitErrInterval = setInterval(() => {
			if (this._afterCommitErrSeed > 0 && this.afterCommitErr != '') {
				this._afterCommitErrSeed--;
				if (this._afterCommitErrSeed == 0) {
					this.afterCommitErr = "";
				}
			}
		}, 1000)
	};

	ngOnDestroy() {
		clearInterval(this._afterCommitErrInterval);
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
		this.userService.updateUser(this.userModel).then(
			res => {
				this.SubmitSuccessEvent.emit(true);
				this.isOpen = false
			},
			(reason: string) => {
				this.afterCommitErr = reason;
				this._afterCommitErrSeed = 3;
			});
	}

	addNewUser() {
		this.userService.newUser(this.userModel)
			.then(
				res => {
					this.SubmitSuccessEvent.emit(true);
					this.isOpen = false;
				}, (reason: string) => {
					this.afterCommitErr = reason;
					this._afterCommitErrSeed = 3;
				})
			.catch(err => {
			})
	}

}
