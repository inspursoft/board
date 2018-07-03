import {Component, OnInit} from '@angular/core';
import {MessageService} from "../../shared/message-service/message.service";
import {Message} from "../../shared/message-service/message";
import {BUTTON_STYLE, MESSAGE_TARGET} from "../../shared/shared.const";
import {AuditService} from "../audit-service/audit-service";
import {Audit, Query} from "../audit";
import {ClrDatagridSortOrder, ClrDatagridStateInterface} from "@clr/angular";

@Component({
  selector: 'app-list-audit',
  templateUrl: './list-audit.component.html',
  styleUrls: ['./list-audit.component.css']
})
export class ListAuditComponent implements OnInit {
  endDate: Date = new Date();
  beginDate: Date = new Date(this.endDate.getTime() - 86400000);
  usernames: string[] = [];
  isInLoading: boolean = false;
  auditsListData: Array<Audit> = Array<Audit>();
  querydata: Query = new Query();
  totalRecordCount: number;
  descSort = ClrDatagridSortOrder.DESC;
  oldStateInfo: ClrDatagridStateInterface;

  constructor(
    private messageService: MessageService,
    private auditService: AuditService,
  ) {
  }

  ngOnInit() {
    this.refreshDataUser();
  }

  refreshDataUser(): void {
    this.isInLoading = true;
    this.auditService.getUserList()
      .then(res => {
        for (var index in res) {
          this.usernames[index] = res[index]["user_name"];
        }
        if (this.usernames.length == 0) {
          this.usernames[0] = "nobody";
        }
        this.isInLoading = false;
      })
      .catch(err => {
        this.messageService.dispatchError(err, '');
        this.isInLoading = false;
      });
  }

  dateTest() {
    if (!this.beginDate || !this.endDate) {
      return
    } else if (this.beginDate > this.endDate) {
      let msg: Message = new Message();
      msg.title = "AUDIT.ILLEGAL_DATE_TITLE";
      msg.message = "AUDIT.ILLEGAL_DATE_MSG";
      msg.buttons = BUTTON_STYLE.ONLY_CONFIRM;
      this.messageService.announceMessage(msg);
      this.endDate = new Date();
      this.beginDate = new Date(this.endDate.getTime() - 86400000);
      this.querydata.endDate = new Date().getTime().toString();
      this.querydata.beginDate = new Date(this.endDate.getTime() - 86400000).getTime().toString();
      return false;
    } else {
      this.querydata.beginDate = this.beginDate.getTime().toString();
      this.querydata.endDate = this.endDate.getTime().toString();
      return
    }
  }

  retrieve(state: ClrDatagridStateInterface): void {
    if (state) {
      this.isInLoading = true;
      this.oldStateInfo = state;
      this.querydata.sortBy = state.sort.by as string;
      this.querydata.isReverse = state.sort.reverse;
      this.auditService
        .getAuditList(this.querydata)
        .then(paginatedProjects => {
          this.totalRecordCount = paginatedProjects.pagination.total_count;
          this.auditsListData = paginatedProjects['operation_list'];
          this.isInLoading = false;
        })
        .catch(err => {
          this.messageService.dispatchError(err, 'PROJECT.FAILED_TO_RETRIEVE_PROJECTS');
          this.isInLoading = false;
        });
    }
  }

  query() {
      this.isInLoading = true;
      this.auditService
        .getAuditList(this.querydata)
        .then(paginatedProjects => {
          this.totalRecordCount = paginatedProjects.pagination.total_count;
          this.auditsListData = paginatedProjects['operation_list'];
          this.isInLoading = false;
        })
        .catch(err => {
          this.messageService.dispatchError(err, 'PROJECT.FAILED_TO_RETRIEVE_PROJECTS');
          this.isInLoading = false;
        });
  }

  ChangeObject(e) {
    switch (e.toString()) {
      case "登录":
        this.querydata.object_name = "sign";
        break;
      case "用户":
        this.querydata.object_name = "user";
        break;
      case "监控":
        this.querydata.object_name = "dashboard";
        break;
      case "组标签":
        this.querydata.object_name = "nodegroup";
        break;
      case "节点":
        this.querydata.object_name = "node";
        break;
      case "项目":
        this.querydata.object_name = "projects";
        break;
      case "服务":
        this.querydata.object_name = "services";
        break;
      case "镜像":
        this.querydata.object_name = "images";
        break;
      case "文件":
        this.querydata.object_name = "file";
        break;
      case "系统":
        this.querydata.object_name = "system";
        break;
      default:
        this.querydata.object_name = e.toString();
    }
  }

  ChangeAction(e) {
    switch (e.toString()) {
      case "创建":
        this.querydata.action = "Create";
        break;
      case "删除":
        this.querydata.action = "Delete";
        break;
      case "修改":
        this.querydata.action = "Update";
        break;
      default:
        this.querydata.action = e.toString();
    }
  }

  ChangeStatus(e) {
    switch (e.toString()) {
      case "未知":
        this.querydata.status = "Unknown";
        break;
      case "成功":
        this.querydata.status = "Success";
        break;
      case "失败":
        this.querydata.status = "Fail";
        break;
      case "错误":
        this.querydata.status = "Error";
        break;
      default:
        this.querydata.status = e.toString();
    }
  }

  ChangeUsername(e) {
    this.querydata.user_name = e.toString();
  }

}
