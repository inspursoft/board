import { Component, OnInit, OnDestroy } from "@angular/core";
import { UserService } from "app/profile/user-center/user-service/user-service";
import { user, User } from "app/profile/user-center/user";
import { editModel } from "../user-new-edit/user-new-edit.component"
import { Message } from "app/shared/message-service/message";
import { MessageService } from "app/shared/message-service/message.service";
import { Subscription } from "rxjs/Subscription";
import { TranslateService } from "@ngx-translate/core";

@Component( {
  selector: "user-list",
  templateUrl: "./user-list.component.html",
  styleUrls: [ "./user-list.component.css" ],
  providers: [ TranslateService ]
} )

export class UserList implements OnInit, OnDestroy {
  _deleteSubscription: Subscription;
  userListData: Array<user> = Array<user>();
  userListErrMsg: string = "";
  userCountPerPage: number = 2;
  curUser: user;
  curEditModel: editModel = editModel.emNew;
  curPage: number = 1;
  showNewUser: boolean = false;

  constructor( private userService: UserService,
               private translateService: TranslateService,
               private messageService: MessageService ) {
  }

  refreshData( username?: string,
               user_list_page: number = 0,
               user_list_page_size: number = 0 ): void {
    this.userService.getUserList()
      .then( res => this.userListData = res )
      .catch( ( reason: Response ) => {
        this.userListErrMsg = `${reason.status}:${reason.statusText}`;
      } );
  }

  addUser() {
    this.curUser = new User();
    this.curEditModel = editModel.emNew;
    this.showNewUser = true;
  }

  editUser( user: user ) {
    this.curEditModel = editModel.emEdit;
    this.userService.getUser( user.user_id )
      .then( user => {
        this.curUser = user;
        this.showNewUser = true;
      } )
      .catch( () => {
      } );
  }

  pageChange( page: number ) {
    this.curPage = page;
  }

  deleteUser( user: user ) {
    this.translateService.get( "USER_CENTER.CONFIRM_DELETE_USER", [ user.user_name ] )
      .subscribe( ( res: string ) => {
        let m: Message = new Message();
        m.title = "USER_CENTER.DELETE_USER";
        m.data = user;
        m.message = res;
        this.messageService.announceMessage( m );
      } );
  }

  ngOnInit() {
    this._deleteSubscription = this.messageService.messageConfirmed$.subscribe( next => {
      this.userService.deleteUser( next.data )
        .then( ( res: User ) => {
          this.refreshData();
          let m: Message = new Message();
          m.message = "USER_CENTER.DELETE_USER_SUCCESS";
          this.messageService.inlineAlertMessage( m );
        } )
        .catch( () => {
        } );
    } );
    this.refreshData();
  }

  ngOnDestroy(): void {
    if (this._deleteSubscription) {
      this._deleteSubscription.unsubscribe();
    }
  }
}
