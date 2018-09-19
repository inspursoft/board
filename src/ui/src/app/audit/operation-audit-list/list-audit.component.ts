import { Component, OnInit } from '@angular/core';
import { MessageService } from "../../shared/message-service/message.service";
import { Audit, AuditQueryData } from "../audit";
import { ClrDatagridSortOrder, ClrDatagridStateInterface } from "@clr/angular";
import { OperationAuditService } from "../audit-service";
import { HttpErrorResponse } from "@angular/common/http";
import { User } from "../../user-center/user";

@Component({
  selector: 'list-audit',
  templateUrl: './list-audit.component.html',
  styleUrls: ['./list-audit.component.css']
})
export class ListAuditComponent implements OnInit {
  endDate: Date;
  beginDate: Date;
  userNames: Array<User>;
  isInLoading: boolean = false;
  auditsListData: Array<Audit>;
  auditQueryData: AuditQueryData = new AuditQueryData();
  totalRecordCount: number = 0;
  descSort = ClrDatagridSortOrder.DESC;
  oldStateInfo: ClrDatagridStateInterface;
  objectQueryMap: Array<{ key: string, title: string, isSpecial?: boolean }>;
  actionQueryMap: Array<{ key: string, title: string, isSpecial?: boolean }>;
  actionStatusMap: Array<{ key: string, title: string, isSpecial?: boolean }>;


  constructor(private messageService: MessageService,
              private auditService: OperationAuditService) {
    let now: Date = new Date();
    this.beginDate = new Date(now.getFullYear(), now.getMonth(), now.getDate(), 0, 0, 0, 0);
    this.endDate = new Date(now.getFullYear(), now.getMonth(), now.getDate(), 23, 59, 59, 0);
    this.auditsListData = Array<Audit>();
    this.userNames = Array<User>();
    this.auditQueryData = new AuditQueryData();
    this.objectQueryMap = Array<{ key: string, title: string }>();
    this.actionQueryMap = Array<{ key: string, title: string }>();
    this.actionStatusMap = Array<{ key: string, title: string }>();
  }

  ngOnInit() {
    this.initObjectQueryMap();
    this.initActionQueryMap();
    this.initStatusQueryMap();
    this.getUserList();
  }

  initObjectQueryMap() {
    this.objectQueryMap.push({key: "ALL", title: "AUDIT.ALL", isSpecial: true});
    this.objectQueryMap.push({key: "user", title: "AUDIT.AUDIT_USER"});
    this.objectQueryMap.push({key: "sign-in", title: "AUDIT.AUDIT_SIGN_IN"});
    this.objectQueryMap.push({key: "log-out", title: "AUDIT.AUDIT_LOG_OUT"});
    this.objectQueryMap.push({key: "dashboard", title: "AUDIT.AUDIT_DASHBOARD"});
    this.objectQueryMap.push({key: "nodegroup", title: "AUDIT.AUDIT_NODEGROUP"});
    this.objectQueryMap.push({key: "node", title: "AUDIT.AUDIT_NODE"});
    this.objectQueryMap.push({key: "projects", title: "AUDIT.AUDIT_PROJECTS"});
    this.objectQueryMap.push({key: "services", title: "AUDIT.AUDIT_SERVICES"});
    this.objectQueryMap.push({key: "images", title: "AUDIT.AUDIT_IMAGES"});
    this.objectQueryMap.push({key: "file", title: "AUDIT.AUDIT_FILE"});
    this.objectQueryMap.push({key: "system", title: "AUDIT.AUDIT_SYSTEM"});
    this.objectQueryMap.push({key: "reset-password", title: "AUDIT.AUDIT_RESET_PASSWORD"});
  }

  initActionQueryMap() {
    this.actionQueryMap.push({key: "ALL", title: "AUDIT.ALL", isSpecial: true});
    this.actionQueryMap.push({key: "get", title: "AUDIT.AUDIT_GET"});
    this.actionQueryMap.push({key: "create", title: "AUDIT.AUDIT_CREATE"});
    this.actionQueryMap.push({key: "delete", title: "AUDIT.AUDIT_DELETE"});
    this.actionQueryMap.push({key: "update", title: "AUDIT.AUDIT_UPDATE"});
  }

  initStatusQueryMap() {
    this.actionStatusMap.push({key: "ALL", title: "AUDIT.ALL", isSpecial: true});
    this.actionStatusMap.push({key: "Unknown", title: "AUDIT.AUDIT_UNKNOWN"});
    this.actionStatusMap.push({key: "Success", title: "AUDIT.AUDIT_SUCCESS"});
    this.actionStatusMap.push({key: "Failed", title: "AUDIT.AUDIT_FAILED"});
    this.actionStatusMap.push({key: "Error", title: "AUDIT.AUDIT_ERROR"});
  }

  getUserList(): void {
    this.isInLoading = true;
    this.auditService.getUserList().subscribe((res: Array<User>) => {
      let user = new User();
      user.user_name = "AUDIT.ALL";
      user["isSpecial"] = true;
      this.userNames.push(user);
      this.userNames = this.userNames.concat(res);
      this.isInLoading = false;
    }, (err: HttpErrorResponse) => this.isInLoading = false)
  }

  changeObjectQuery(event: { key: string, title: string }) {
    this.auditQueryData.object_name = event.key == "ALL" ? "" : event.key;
  }

  changeActionQuery(event: { key: string, title: string }) {
    this.auditQueryData.action = event.key == "ALL" ? "" : event.key;
  }

  changeStatusQuery(event: { key: string, title: string }) {
    this.auditQueryData.status = event.key == "ALL" ? "" : event.key;
  }

  changeUsernameQuery(user: User) {
    this.auditQueryData.user_name = user.user_name == "AUDIT.ALL" ? "" : user.user_name;
  }

  changeEndData(event: Date) {
    this.endDate = new Date(event.getFullYear(), event.getMonth(), event.getDate(), 23, 59, 59, 0);
  }

  retrieve(state: ClrDatagridStateInterface): void {
    this.oldStateInfo = state;
    this.auditQueryData.sortBy = state.sort.by as string;
    this.auditQueryData.isReverse = state.sort.reverse;
    this.queryListData();
  }

  queryListData() {
    if (this.beginDate.getTime() < this.endDate.getTime()) {
      setTimeout(() => {
        this.isInLoading = true;
        this.auditQueryData.beginTimestamp = this.beginDate ? this.beginDate.getTime() : 0;
        this.auditQueryData.endTimestamp = this.endDate ? this.endDate.getTime() : 0;
        this.auditService.getAuditList(this.auditQueryData).subscribe(paginatedProjects => {
          this.totalRecordCount = paginatedProjects.pagination.total_count;
          this.auditsListData = paginatedProjects['operation_list'];
          this.isInLoading = false;
        }, () => this.isInLoading = false)
      });
    } else {
      this.messageService.showAlert('AUDIT.AUDIT_QUERY_DATE_ERROR', {alertType: 'alert-warning'})
    }
  }

  getObjectTitle(key: string): string {
    let query = this.objectQueryMap.find(value => value.key.toUpperCase() == key.toUpperCase());
    return query ? query.title : key;
  }

  getActionTitle(key: string): string {
    let query = this.actionQueryMap.find(value => value.key.toUpperCase() == key.toUpperCase());
    return query ? query.title : key;
  }

  getStatusTitle(key: string): string {
    let query = this.actionStatusMap.find(value => value.key.toUpperCase() == key.toUpperCase());
    return query ? query.title : key;
  }
}
