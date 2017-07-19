import {Component, OnInit, OnDestroy} from "@angular/core";
import {UserService} from "app/profile/user-center/user-service/user-service";
import {user, User} from "app/profile/user-center/user";
import {editModel} from "../user-new-edit/user-new-edit.component"
import {ConfirmationMessage} from "app/shared/service/confirmation-message";
import {MessageService} from "app/shared/service/message.service";
import {Subscription} from "rxjs/Subscription";
import {TranslateService} from "@ngx-translate/core";

@Component({
	selector: "user-list",
	templateUrl: "./user-list.component.html",
	styleUrls: ["./user-list.component.css"],
	providers: [TranslateService]
})

export class UserList implements OnInit, OnDestroy {
	userListData: Array<user> = Array<user>();
	userListErrMsg: string = "";
	curUser: user;
	curEditModel: editModel = editModel.emNew;
	showNewUser: boolean = false;
	_deleteSubscription: Subscription;

	constructor(private userService: UserService,
							private translateService: TranslateService,
							private messageService: MessageService) {
	}

	refreshData(username?: string,
							user_list_page: number = 0,
							user_list_page_size: number = 0): void {
		this.userService.getUserList().then(
			res => (this.userListData = res),
			(reason: string) => {
				this.userListErrMsg = reason;
			}
		)
	}

	addUser() {
		this.curUser = new User();
		this.curEditModel = editModel.emNew;
		this.showNewUser = true;
	}

	editUser(user: user) {
		this.curEditModel = editModel.emEdit;
		this.userService.getUser(user.user_id).then(
			user => {
				this.curUser = user;
				this.showNewUser = true;
			},
			(reason: string) => {
				this.messageService.globalWarning = reason;
			});
	}

	deleteUser(user: user) {
		let m: ConfirmationMessage = new ConfirmationMessage();
		m.title = "USER_CENTER.DELETE_USER";
		m.data = user;
		this.translateService.get("USER_CENTER.CONFIRM_DELETE_USER")
			.subscribe((res: string) => {
				m.message = res.concat(`[${user.user_name}]?`);
				this.messageService.announceMessage(m);
			});
	}

	ngOnInit() {
		this._deleteSubscription = this.messageService.messageConfirmed$.subscribe(next => {
			this.userService.deleteUser(next.data).then(
				(res: User) => {
					this.refreshData();
					this.messageService.globalInfo = "USER_CENTER.DELETE_USER_SUCCESS";
				},
				(reason: string) => {
					this.messageService.globalWarning = reason;
				});
		});
		this.refreshData();
	}

	ngOnDestroy(): void {
		if (this._deleteSubscription) {
			this._deleteSubscription.unsubscribe();
		}
	}
}
