import { ComponentFactoryResolver, ComponentRef, Injectable, Type, ViewContainerRef } from '@angular/core';
import { HttpErrorResponse } from '@angular/common/http';
import { Observable, TimeoutError } from 'rxjs';
import { tap } from 'rxjs/operators';
import { AlertMessage, AlertType, BUTTON_STYLE, GlobalAlertMessage, GlobalAlertType, Message } from '../shared/shared.types';
import { CsAlertComponent } from './cs-alert/cs-alert.component';
import { CsGlobalAlertComponent } from './cs-global-alert/cs-global-alert.component';
import { CsDialogComponent } from './cs-dialog/cs-dialog.component';

@Injectable()
export class MessageService {
  private dialogView: ViewContainerRef;
  private dialogResolver: ComponentFactoryResolver;

  private createComponent<T>(component: Type<T>, view: ViewContainerRef): ComponentRef<T> {
    view.clear();
    const factory = this.dialogResolver.resolveComponentFactory(component);
    return view.createComponent<T>(factory);
  }

  public registerDialogHandle(view: ViewContainerRef, resolver: ComponentFactoryResolver) {
    this.dialogView = view;
    this.dialogResolver = resolver;
  }

  public cleanNotification() {
    this.dialogView.clear();
  }

  public showAlert(msg: string, optional?: {alertType?: AlertType, view?: ViewContainerRef}): void {
    this.dialogView.clear();
    const alertView: ViewContainerRef = optional ? optional.view || this.dialogView : this.dialogView;
    const message: AlertMessage = new AlertMessage();
    message.message = msg;
    message.alertType = optional ? optional.alertType || 'success' : 'success';
    const componentRef = this.createComponent(CsAlertComponent, alertView);
    componentRef.instance.openAlert(message).subscribe(() => alertView.remove(alertView.indexOf(componentRef.hostView)));
  }

  public showGlobalMessage(msg: string,
                           optional?: {
                             alertType?: AlertType,
                             globalAlertType?: GlobalAlertType,
                             errorObject?: HttpErrorResponse | Type<Error> | TimeoutError,
                             view?: ViewContainerRef,
                             endMessage?: string
                           }): void {
    const globalView: ViewContainerRef = optional ? optional.view || this.dialogView : this.dialogView;
    const message: GlobalAlertMessage = new GlobalAlertMessage();
    message.message = msg;
    message.alertType = optional ? optional.alertType || 'danger' : 'danger';
    message.type = optional ? optional.globalAlertType || GlobalAlertType.gatNormal : GlobalAlertType.gatNormal;
    message.errorObject = optional ? optional.errorObject : null;
    message.endMessage = optional && optional.endMessage ? `:${optional.endMessage}` : '';
    const componentRef = this.createComponent(CsGlobalAlertComponent, globalView);
    componentRef.instance.openAlert(message).subscribe(() => globalView.remove(globalView.indexOf(componentRef.hostView)));
  }

  public showDialog(msg: string,
                    optional?: {
                      title?: string,
                      buttonStyle?: BUTTON_STYLE,
                      data?: any,
                      view?: ViewContainerRef
                    }): Observable<Message> {
    const dialogView = optional ? optional.view || this.dialogView : this.dialogView;
    const message: Message = new Message();
    message.message = msg;
    message.title = optional ? optional.title || 'GLOBAL_ALERT.TITLE' : 'GLOBAL_ALERT.TITLE';
    message.buttonStyle = optional ? optional.buttonStyle || BUTTON_STYLE.ONLY_CONFIRM : BUTTON_STYLE.ONLY_CONFIRM;
    message.data = optional ? optional.data : null;
    const componentRef = this.createComponent(CsDialogComponent, dialogView);
    return componentRef.instance.openDialog(message)
      .pipe(tap(() => dialogView.remove(dialogView.indexOf(componentRef.hostView))));
  }

  public showOnlyOkDialog(msg: string, title?: string): void {
    this.showDialog(msg, {title: title || 'GLOBAL_ALERT.HINT', buttonStyle: BUTTON_STYLE.ONLY_CONFIRM}).subscribe();
  }

  public showOnlyOkDialogObservable(msg: string, title?: string): Observable<Message> {
    return this.showDialog(msg, {title: title || 'GLOBAL_ALERT.HINT', buttonStyle: BUTTON_STYLE.ONLY_CONFIRM});
  }

  public showYesNoDialog(msg: string, title?: string): Observable<Message> {
    return this.showDialog(msg, {title: title || 'GLOBAL_ALERT.ASK', buttonStyle: BUTTON_STYLE.YES_NO});
  }

  public showConfirmationDialog(msg: string, title?: string): Observable<Message> {
    return this.showDialog(msg, {title: title || 'GLOBAL_ALERT.ASK', buttonStyle: BUTTON_STYLE.CONFIRMATION});
  }

  public showDeleteDialog(msg: string, title?: string): Observable<Message> {
    return this.showDialog(msg, {title: title || 'GLOBAL_ALERT.DELETE', buttonStyle: BUTTON_STYLE.DELETION});
  }
}
