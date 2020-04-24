import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  Inject,
  Input,
  OnDestroy,
  OnInit,
  TemplateRef,
  ViewEncapsulation,
} from '@angular/core';
import { fadeMotion } from '../animations';
import { Subscription, fromEvent } from 'rxjs';
import { DOCUMENT } from '@angular/common';
import { throttleTime, distinctUntilChanged } from 'rxjs/operators';
import { ScrollTools } from '../scroll.tools';

@Component({
  selector: 'app-back-top',
  animations: [fadeMotion],
  templateUrl: './back-top.component.html',
  styleUrls: ['./back-top.component.css'],
  // 当组件实例化之后，Angular 就会创建一个变更检测器，它负责传播组件各个绑定值的变化。
  // 该策略是下列值之一：
  // ChangeDetectionStrategy#OnPush(0) 把策略设置为 CheckOnce（按需）
  // ChangeDetectionStrategy#Default(1) 把策略设置为 CheckAlways
  changeDetection: ChangeDetectionStrategy.OnPush,
  encapsulation: ViewEncapsulation.None,
  preserveWhitespaces: false
})
export class BackTopComponent implements OnInit, OnDestroy {
  private scroll$: Subscription | null = null;
  private target: HTMLElement | null = null;
  private myVisibilityHeight = 400;

  visible = false;

  // 如果自定义了样式，则使用自定义的样式
  @Input() backTopTemplate: TemplateRef<void>;

  // 如果传入了高度，则使用传入的高度值
  @Input()
  set visibilityHeight(height: string | number) {
    if (typeof height === 'number') {
      this.myVisibilityHeight = height;
    } else {
      this.myVisibilityHeight = parseInt(height, 10) >= 0 ? parseInt(height, 10) : 400;
    }
  }

  // 如果传入了目标元素，则为目标元素注册订阅信息
  @Input()
  set topTarget(el: string | HTMLElement) {
    this.target = typeof el === 'string' ? this.doc.querySelector(el) : el;
    this.registerScrollEvent();
  }

  constructor(
    private scrollTools: ScrollTools,
    // tslint:disable-next-line:no-any
    @Inject(DOCUMENT) private doc: any,
    // 与changeDetection一起使用，手动调用声明视图刷新的时间
    private cd: ChangeDetectorRef
  ) { }

  ngOnInit(): void {
    if (!this.scroll$) {
      this.registerScrollEvent();
    }
  }

  // 点击返回顶部
  clickBackTop(): void {
    this.scrollTools.scrollTo(this.getTarget(), 0);
  }

  private getTarget(): HTMLElement | Window {
    return this.target || window;
  }

  // 当滑动的距离超过 myVisibilityHeight 设定的高度时才显示该组件
  private handleScroll(): void {
    if (this.visible === this.scrollTools.getScroll(this.getTarget()) > this.myVisibilityHeight) {
      return;
    }
    this.visible = !this.visible;
    this.cd.markForCheck();
  }

  // （如果之前订阅过的话）清除之前的订阅缓存
  private removeListen(): void {
    if (this.scroll$) {
      this.scroll$.unsubscribe();
    }
  }

  private registerScrollEvent(): void {
    // 先清除缓存
    this.removeListen();
    // 处理滚动操作
    this.handleScroll();
    // 监听组件的滚动操作
    this.scroll$ = fromEvent(this.getTarget(), 'scroll')
      .pipe(
        throttleTime(50),
        distinctUntilChanged(),
      )
      .subscribe(() => this.handleScroll());
  }

  // 销毁页面时清除缓存
  ngOnDestroy(): void {
    this.removeListen();
  }
}
