import { Component, ViewChild, OnInit } from '@angular/core';

@Component({
  selector: 'custom-chart',
  templateUrl: './chart.component.html'
})
export class ChartComponent implements OnInit {
  
  options: {[key: string]: any};

  ngOnInit() {
    this.options = {
        title: {
            text: 'ECharts 入门示例'
        },
        tooltip: {},
        legend: {
            data:['销量']
        },
        xAxis: {
            data: ["衬衫","羊毛衫","雪纺衫","裤子","高跟鞋","袜子"]
        },
        yAxis: {},
        series: [{
            name: '销量',
            type: 'bar',
            data: [5, 20, 36, 10, 10, 20]
        }]
    };
  }
}