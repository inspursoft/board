import {AfterViewInit, Component} from '@angular/core';
import {DatePipe} from "@angular/common";

@Component({
    selector: 'dashboard',
    templateUrl: 'dashboard.component.html',
    styleUrls: ['dashboard.component.css']
})
export class DashboardComponent implements AfterViewInit {
    serviceOptions={};
    ngAfterViewInit() {
        this.serviceOptions = {
            tooltip: {
                trigger: 'axis',
                formatter: (params) => {
                    let xDate: Date = new Date(params[0].value[0]);
                    let pDate: DatePipe = new DatePipe("lt");
                    return pDate.transform(xDate, "yyyy/MM/dd HH:mm:ss") +
                        `<div style='display: flex;flex-direction: column'>
                        <div style="display: flex;align-items: center">
                            <div style='width: 16px;height: 16px; background-color: red;border-radius: 50%'></div>
                            <div>pods:${params[0].value[1]}</div>
                        </div>
                        <div style="display: flex;align-items: center">
                            <div style='width: 16px;height: 16px; background-color: blue;border-radius: 50%'></div>
                            <div>containers:${params[0].value[2]}</div>
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
                    data: [
                        [new Date(2017, 6, 29, 0, 0, 0), 4, 6],
                        [new Date(2017, 6, 29, 2, 0, 0), 10, 9],
                        [new Date(2017, 6, 29, 4, 0, 0), 20, 2],
                        [new Date(2017, 6, 29, 5, 0, 0), 40, 10],
                        [new Date(2017, 6, 29, 6, 0, 0), 50, 60],
                        [new Date(2017, 6, 29, 8, 0, 0), 30, 50],
                        [new Date(2017, 6, 29, 10, 0, 0), 70, 35],
                        [new Date(2017, 6, 29, 12, 0, 0), 65, 40]]
                },
                {
                    // id: '模拟数据',
                    type: 'line',
                    showSymbol: true,
                    smooth: true,
                    symbolSize: 10,
                    hoverAnimation: false,
                    data: [
                        [new Date(2017, 6, 29, 0, 0, 0), 6],
                        [new Date(2017, 6, 29, 2, 0, 0), 9],
                        [new Date(2017, 6, 29, 4, 0, 0), 2],
                        [new Date(2017, 6, 29, 5, 0, 0), 10],
                        [new Date(2017, 6, 29, 6, 0, 0), 60],
                        [new Date(2017, 6, 29, 8, 0, 0), 50],
                        [new Date(2017, 6, 29, 10, 0, 0), 35],
                        [new Date(2017, 6, 29, 12, 0, 0), 40]]
                }
            ]
        };

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