import {AfterViewInit, OnInit, Component} from '@angular/core';
import {DatePipe} from "@angular/common";

@Component({
    selector: 'dashboard',
    templateUrl: 'dashboard.component.html',
    styleUrls: ['dashboard.component.css']
})
export class DashboardComponent implements OnInit, AfterViewInit {
    serviceOptions = {};
    serviceOptionsOther ={};
    curData: any = [];
    curDataOther: any = [];
    baseDate: Date = new Date();
    oneStepTime: number = 10 * 1000;

    ngOnInit() {
        for (let i = 0; i < 11; i++) {
            let arrBuf = [this.getDate(), this.getRandomData()];
            this.curData.push(arrBuf);
            let arrOther = Array.from(arrBuf);
            arrOther[1] = this.getRandomDataOther();
            this.curDataOther.push(arrOther)
        }
    }

    ngAfterViewInit() {
        this.serviceOptionsOther = {
            tooltip: {
                trigger: 'axis',
                formatter: (params) => {
                    let xDate: Date = new Date(params[0].value[0]);
                    let pDate: DatePipe = new DatePipe("lt");
                    console.log(params);
                    return pDate.transform(xDate, "yyyy/MM/dd HH:mm:ss") +
                        `<div style='display: flex;flex-direction: column'>
                        <div style="display: flex;align-items: center">
                            <div style='width: 16px;height: 16px; background-color: red;border-radius: 50%'></div>
                            <div>pods:${params[0].value[1]}</div>
                        </div>
                        <div style="display: flex;align-items: center">
                            <div style='width: 16px;height: 16px; background-color: blue;border-radius: 50%'></div>
                            <div>containers:${params[1].value[1]}</div>
                        </div>
                    </div>`;
                },
                axisPointer: {
                    animation: false
                }
            },
            xAxis: {
                type: 'time',
                splitNumber: 10,
                splitLine: {
                    show: false
                }
            },
            yAxis: {
                type: 'value',
                show: true,
                splitLine: {
                    show: true
                }
            },
            series: [
                {
                    // id: '模拟数据',
                    type: 'line',
                    showSymbol: true,
                    smooth: true,
                    symbolSize: 10,
                    hoverAnimation: false,
                    data: this.curData
                },
                {
                    // id: '模拟数据',
                    type: 'line',
                    showSymbol: true,
                    smooth: true,
                    symbolSize: 10,
                    hoverAnimation: false,
                    data: this.curDataOther
                }
            ]
        };
        this.serviceOptions = {
            tooltip: {
                trigger: 'axis',
                formatter: (params) => {
                    let xDate: Date = new Date(params[0].value[0]);
                    let pDate: DatePipe = new DatePipe("lt");
                    console.log(params);
                    return pDate.transform(xDate, "yyyy/MM/dd HH:mm:ss") +
                        `<div style='display: flex;flex-direction: column'>
                        <div style="display: flex;align-items: center">
                            <div style='width: 16px;height: 16px; background-color: red;border-radius: 50%'></div>
                            <div>pods:${params[0].value[1]}</div>
                        </div>
                        <div style="display: flex;align-items: center">
                            <div style='width: 16px;height: 16px; background-color: blue;border-radius: 50%'></div>
                            <div>containers:${params[1].value[1]}</div>
                        </div>
                    </div>`;
                },
                axisPointer: {
                    animation: false
                }
            },
            xAxis: {
                type: 'time',
                splitNumber: 10,
                splitLine: {
                    show: false
                }
            },
            yAxis: {
                type: 'value',
                show: true,
                splitLine: {
                    show: true
                }
            },
            series: [
                {
                    // id: '模拟数据',
                    type: 'line',
                    showSymbol: true,
                    smooth: true,
                    symbolSize: 10,
                    hoverAnimation: false,
                    data: this.curData
                },
                {
                    // id: '模拟数据',
                    type: 'line',
                    showSymbol: true,
                    smooth: true,
                    symbolSize: 10,
                    hoverAnimation: false,
                    data: this.curDataOther
                }
            ]
        };
        setInterval(() => {
            this.curData.shift();
            let arrBuf = [this.getDate(), this.getRandomData()];
            this.curData.push(arrBuf);
            this.curDataOther.shift();
            let arrOther = Array.from(arrBuf);
            arrOther[1] = this.getRandomDataOther();
            this.curDataOther.push(arrOther);

            this.serviceOptions = {
                series: [
                    {
                        // id: '模拟数据',
                        type: 'line',
                        showSymbol: true,
                        smooth: true,
                        symbolSize: 10,
                        hoverAnimation: false,
                        data: this.curData
                    },
                    {
                        // id: '模拟数据',
                        type: 'line',
                        showSymbol: true,
                        smooth: true,
                        symbolSize: 10,
                        hoverAnimation: false,
                        data: this.curDataOther
                    }
                ]
            };
        }, 2000);
    }

    getRandomData(): number {//pod
        return 100 + Math.round(Math.random() * 50);
    }

    getRandomDataOther(): number {
        return 300 + Math.round(Math.random() * 100);
    }

    getDate(): Date {
        let bDate = this.curData.length > 0 ?
            this.curData[this.curData.length - 1][0] :
            this.baseDate;
        return new Date(bDate.getTime() + this.oneStepTime);
    }


    get serviceIcon(): string {
        return '../../images/service_icon.png';
    }

    get nodeIcon(): string {
        return '../../images/node_icon.png';
    }

    get storageIcon(): string {
        return '../../images/storage_icon.png'
    }


}