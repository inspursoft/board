import { AfterViewChecked, AfterViewInit, Component, ElementRef, OnInit, ViewChild } from '@angular/core';

const CharsResource = [
  '1', '2', '3', '4', '5', '6', '7', '8', '9', '0',
  'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J',
  'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T',
  'U', 'V', 'W', 'X', 'Y', 'Z', 'a', 'b', 'c', 'd',
  'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n',
  'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x',
  'y', 'z'
];

@Component({
  selector: 'app-verify-code',
  templateUrl: './verify-code.component.html',
  styleUrls: ['./verify-code.component.css']
})
export class VerifyCodeComponent implements OnInit, AfterViewInit {
  @ViewChild('verifyCanvas') verifyCanvas: ElementRef;
  ctx: CanvasRenderingContext2D;

  constructor() {
  }

  ngOnInit() {
  }

  ngAfterViewInit(): void {
    this.ctx = this.canvasElement.getContext('2d');
    this.drawCode();
  }

  get canvasElement(): HTMLCanvasElement {
    return this.verifyCanvas.nativeElement;
  }

  drawCode() {
    this.ctx.fillStyle = 'cornflowerblue';
    this.ctx.fillRect(0, 0, this.canvasElement.width, this.canvasElement.height);
    const gradient = this.ctx.createLinearGradient(0, 0, this.canvasElement.width, 0);
    gradient.addColorStop(0, 'magenta');
    gradient.addColorStop(0.5, 'blue');
    gradient.addColorStop(1.0, 'red');

    this.ctx.fillStyle = gradient;
    this.ctx.font = '25px Arial';
    const rand = Array();
    const x = Array();
    const y = Array();
    for (let i = 0; i < 4; i++) {
      rand[i] = CharsResource[Math.floor(Math.random() * CharsResource.length)];
      x[i] = i * 16 + 10;
      y[i] = Math.random() * 20 + 20;
      this.ctx.fillText(rand[i], x[i], y[i]);
    }

    for (let i = 0; i < 3; i++) {
      this.drawLine();
    }

    for (let i = 0; i < 30; i++) {
      this.drawDot();
    }

    this.convertCanvasToImage();

    // $('#submit').click((e) => {
    //   var newRand = rand.join('').toUpperCase();
    //   console.log(newRand);
    //
    //   //下面就是判断是否== 的代码，无需解释
    //   var oValue = $('#verify').val().toUpperCase();
    //   console.log(oValue);
    //   if (oValue == 0) {
    //   } else if (oValue != newRand) {
    //     oValue = ' ';
    //   } else {
    //     window.open('http://www.baidu.com', '_self');
    //   }
    //
    // });
  }


  drawLine() {
    this.ctx.moveTo(Math.floor(Math.random() * this.canvasElement.width), Math.floor(Math.random() * this.canvasElement.height));
    this.ctx.lineTo(Math.floor(Math.random() * this.canvasElement.width), Math.floor(Math.random() * this.canvasElement.height));
    this.ctx.lineWidth = 0.5;
    this.ctx.strokeStyle = 'rgba(50,50,50,0.3)';
    this.ctx.stroke();
  }

  drawDot() {
    const px = Math.floor(Math.random() * this.canvasElement.width);
    const py = Math.floor(Math.random() * this.canvasElement.height);
    this.ctx.moveTo(px, py);
    this.ctx.lineTo(px + 1, py + 1);
    this.ctx.lineWidth = 0.2;
    this.ctx.stroke();
  }

  convertCanvasToImage() {
    document.getElementById('verifyCanvas').style.display = 'none';
    const image = document.getElementById('code_img') as HTMLImageElement;
    image.src = this.canvasElement.toDataURL('image/png');
    return image;
  }

}
