import { trigger, state, style, transition, animate } from '@angular/animations';

// 卡片滑动动画
export const cardSlide = trigger('cardSlide', [
    state('open', style({
        opacity: 1,
    })),
    state('closed', style({
        height: 0,
        opacity: 0,
        padding: 0
    })),
    transition('open <=> closed', [
        animate(250)
    ]),
]);

// 逆时针旋转90度
export const rotateNega90 = trigger('rotateNega90', [
    state('open', style({

    })),
    state('closed', style({
        transform: 'rotate(-90deg)'
    })),
    transition('open <=> closed', [
        animate(250)
    ])
]);

// 渐入和渐出
export const fadeMotion = trigger('fadeMotion', [
    transition(':enter', [style({ opacity: 0 }), animate(200, style({ opacity: 1 }))]),
    transition(':leave', [style({ opacity: 1 }), animate(200, style({ opacity: 0 }))])
]);
