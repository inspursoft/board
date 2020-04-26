/**
 * @fileName: sign-in.component.spec
 * @author: liyanqing
 * @dateTime: 25/01/2018 13:16
 * @desc:
 */
import { By } from '@angular/platform-browser';
import { DebugElement, Injectable } from '@angular/core';
import { ActivatedRoute, NavigationExtras, Router } from '@angular/router';
import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { HttpErrorResponse } from '@angular/common/http';
import { defer } from 'rxjs';
import { SignInComponent } from './sign-in.component';
import { AccountService } from '../account.service';
import { AppModule } from '../../app.module';
import { AppComponent } from '../../app.component';
import { AppInitService } from '../../shared.service/app-init.service';
import { HeaderComponent } from '../../shared/header/header.component';
import { AccountModule } from '../account.module';
import Spy = jasmine.Spy;

export function newEvent(eventName: string, bubbles = false, cancelable = false) {
  const evt = document.createEvent('CustomEvent');  // MUST be 'CustomEvent'
  evt.initCustomEvent(eventName, bubbles, cancelable, null);
  return evt;
}

export function asyncData<T>(data: T) {
  return defer(() => Promise.resolve(data));
}

@Injectable()
export class RouterStub {
  navigate(commands: any[], extras?: NavigationExtras): Promise<boolean> {
    return Promise.resolve(true);
  }
}

describe('SignInComponent', () => {
  const fakeActivatedRoute = {
    snapshot: {
      data: {
        systeminfo: {
          auth_mode: 'ldap_auth',
          board_host: '10.0.0.0',
          board_version: 'dev',
          init_project_repo: 'created',
          redirection_url: 'http://redirection.mydomain.com',
          set_auth_password: 'updated',
          sync_k8s: 'created',
          dns_suffix: '.cluster.local',
          kubernetes_version: 'v1.18.0',
          processor_type: 'x86_64'
        }
      }
    }
  };
  let fixture: ComponentFixture<SignInComponent>;
  let component: SignInComponent;
  let accountService: AccountService;
  let signInSpy: Spy;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [AppModule, AccountModule],
      providers: [
        {provide: Router, useClass: RouterStub},
        {provide: ActivatedRoute, useValue: fakeActivatedRoute}]
    }).overrideTemplate(HeaderComponent, `<div></div>`)
      .compileComponents()
      .then(() => {
      /*The AppComponent must be first create,
       because set currentLang on it's constructor function*/
      const appFixture = TestBed.createComponent(AppComponent);
      appFixture.detectChanges();
      appFixture.componentInstance.ngAfterViewInit();
      fixture = TestBed.createComponent(SignInComponent);
      fixture.detectChanges();
      component = fixture.componentInstance;
      accountService = fixture.debugElement.injector.get(AccountService);
      signInSpy = spyOn(accountService, 'signIn');
    });
  }));

  // Todo:Test HttpClient Service:https://angular.cn/guide/testing#component-with-async-service
  function testErrorHandle(status: number): void {
    const btnSubmitDebug: DebugElement = fixture.debugElement.query(By.css('.btn-primary-outline'));
    const btnSubmit: HTMLButtonElement = btnSubmitDebug.nativeElement as HTMLButtonElement;
    const appInitService = fixture.debugElement.injector.get(AppInitService);
    component.signInUser.username = 'hello';
    component.signInUser.password = 'world';
    const errorResponse = new HttpErrorResponse({
      error: 'test error',
      status
    });
    signInSpy.and.returnValue(errorResponse);
    btnSubmit.click();
    fixture.whenStable().then(() => {
      fixture.detectChanges();
      setTimeout(() => {
        expect(signInSpy.calls.any()).toBe(true, '/v1/api/sign-in restApi called');
        expect(appInitService.token).toBeUndefined();
      }, 500);
    });
  }

  it('should create the SignInComponent', () => {
    expect(SignInComponent).toBeTruthy();
  });

  it('should submit button is disabled', () => {
    const btnSubmitDebug: DebugElement = fixture.debugElement.query(By.css('.btn-primary-outline'));
    const btnSubmit: HTMLButtonElement = btnSubmitDebug.nativeElement as HTMLButtonElement;
    fixture.detectChanges();
    expect<boolean>(btnSubmit.disabled).toBeTruthy('submit button is disabled');
  });

  it('should submit button is enable', () => {
    const btnSubmitDebug: DebugElement = fixture.debugElement.query(By.css('.btn-primary-outline'));
    const btnSubmit: HTMLButtonElement = btnSubmitDebug.nativeElement as HTMLButtonElement;
    component.signInUser.username = 'hello';
    component.signInUser.password = 'world';
    fixture.detectChanges();
    expect<boolean>(btnSubmit.disabled).toBeFalsy('submit button is enabled');
  });

  // it("test sign-in click and catch 400 error code", async(() => {
  //   testErrorHandle(400);
  // }));
  //
  // it("test sign-in click and catch 409 error code", async(() => {
  //   testErrorHandle(409);
  // }));

  // it("should sign-in with {username:'admin',password:'xxxxxx' success}", async(() => {
  //   let btnSubmitDebug: DebugElement = fixture.debugElement.query(By.css("[type='submit']"));
  //   let btnSubmit: HTMLInputElement = btnSubmitDebug.nativeElement as HTMLInputElement;
  //   let appInitService = fixture.debugElement.injector.get(AppInitService);
  //   let router = fixture.debugElement.injector.get(Router);
  //   // let navSpy = spyOn(router, 'navigate');
  //   component.signInUser.username = "admin";
  //   component.signInUser.password = "xxxxxx";
  //   // navSpy.and.returnValue(Promise.resolve(true));
  //   const response = new HttpResponse({
  //     body:{token:'tokenString'},
  //     status: 200,
  //   });
  //   signInSpy.and.returnValue(response);
  //   btnSubmit.click();
  //   fixture.whenStable().then(() => {
  //     fixture.detectChanges();
  //     // const navArgs = navSpy.calls.first().args[0];
  //     expect(signInSpy.calls.any()).toBe(true, '/v1/api/sign-in restApi called');
  //     // expect(navSpy.calls.any()).toBe(true, 'navigate called');
  //     // expect(navArgs[0]).toContain('dashboard', 'nav to dashboard detail URL');
  //     expect(appInitService.token).toBe("tokenString", 'appInitService.token is undefined')
  //   });
  // }))
});
