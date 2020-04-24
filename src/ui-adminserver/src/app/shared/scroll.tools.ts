import { DOCUMENT } from '@angular/common';
import { Inject, Injectable } from '@angular/core';

const availablePrefixes = ['moz', 'ms', 'webkit'];
const reqAnimFrame = getRequestAnimationFrame();

type EasyingFn = (t: number, b: number, c: number, d: number) => number;

function easeInOutCubic(t: number, b: number, c: number, d: number): number {
  const cc = c - b;
  let tt = t / (d / 2);
  if (tt < 1) {
    return (cc / 2) * tt * tt * tt + b;
  } else {
    return (cc / 2) * ((tt -= 2) * tt * tt + 2) + b;
  }
}

@Injectable()
export class ScrollTools {
  private doc: Document;

  /* tslint:disable-next-line:no-any */
  constructor(@Inject(DOCUMENT) doc: any) {
    this.doc = doc;
  }

  /** Set the position of the scroll bar of `el`. */
  setScrollTop(el: Element | Window, topValue: number = 0): void {
    if (el === window) {
      this.doc.body.scrollTop = topValue;
      this.doc.documentElement!.scrollTop = topValue;
    } else {
      (el as Element).scrollTop = topValue;
    }
  }

  /** Get position of `el` against window. */
  getOffset(el: Element): { top: number; left: number } {
    const ret = {
      top: 0,
      left: 0
    };
    if (!el || !el.getClientRects().length) {
      return ret;
    }

    const rect = el.getBoundingClientRect();
    if (rect.width || rect.height) {
      const doc = el.ownerDocument!.documentElement;
      ret.top = rect.top - doc!.clientTop;
      ret.left = rect.left - doc!.clientLeft;
    } else {
      ret.top = rect.top;
      ret.left = rect.left;
    }

    return ret;
  }

  /** Get the position of the scoll bar of `el`. */
  getScroll(el?: Element | Window, top: boolean = true): number {
    const target = el ? el : window;
    const prop = top ? 'pageYOffset' : 'pageXOffset';
    const method = top ? 'scrollTop' : 'scrollLeft';
    const isWindow = target === window;
    // @ts-ignore
    let ret = isWindow ? target[prop] : target[method];
    if (isWindow && typeof ret !== 'number') {
      ret = this.doc.documentElement![method];
    }
    return ret;
  }

  /**
   * Scroll `el` to some position with animation.
   *
   * @param containerEl container, `window` by default
   * @param targetTopValue Scroll to `top`, 0 by default
   * @param easing Transition curve, `easeInOutCubic` by default
   * @param callback callback invoked when transition is done
   */
  scrollTo(containerEl: Element | Window, targetTopValue: number = 0, easing?: EasyingFn, callback?: () => void): void {
    const target = containerEl ? containerEl : window;
    const scrollTop = this.getScroll(target);
    const startTime = Date.now();
    const frameFunc = () => {
      const timestamp = Date.now();
      const time = timestamp - startTime;
      this.setScrollTop(target, (easing || easeInOutCubic)(time, scrollTop, targetTopValue, 450));
      if (time < 450) {
        reqAnimFrame(frameFunc);
      } else {
        if (callback) {
          callback();
        }
      }
    };
    reqAnimFrame(frameFunc);
  }
}

/*
 * request-animation
 */

function requestAnimationFramePolyfill(): typeof requestAnimationFrame {
  let lastTime = 0;
  return function(callback: FrameRequestCallback): any {
    const currTime = new Date().getTime();
    const timeToCall = Math.max(0, 16 - (currTime - lastTime));
    const id = setTimeout(() => {
      callback(currTime + timeToCall);
    }, timeToCall);
    lastTime = currTime + timeToCall;
    return id;
  };
}

function getRequestAnimationFrame(): typeof requestAnimationFrame {
  if (typeof window === 'undefined') {
    return () => 0;
  }

  // 如果当前页面支持 requestAnimationFrame 则绑定属性
  if (window.requestAnimationFrame) {
    return window.requestAnimationFrame.bind(window);
  }

  // 在当前页面支持的属性中获取一个可用的参数前缀
  const prefix = availablePrefixes.filter(key => `${key}RequestAnimationFrame` in window)[0];

  // 如果有可用的 requestAnimationFrame 则使用，否则使用一个通用的动画时间间隔器
  return prefix ? (window as any)[`${prefix}RequestAnimationFrame`] : requestAnimationFramePolyfill();
}

export function cancelRequestAnimationFrame(id: number): any {
  if (typeof window === 'undefined') {
    return null;
  }
  if (window.cancelAnimationFrame) {
    return window.cancelAnimationFrame(id);
  }
  const prefix = availablePrefixes.filter(
    key => `${key}CancelAnimationFrame` in window || `${key}CancelRequestAnimationFrame` in window
  )[0];

  return prefix
    ? ((window as any)[`${prefix}CancelAnimationFrame`] || (window as any)[`${prefix}CancelRequestAnimationFrame`])
        // @ts-ignore
        .call(this, id)
    : clearTimeout(id);
}
