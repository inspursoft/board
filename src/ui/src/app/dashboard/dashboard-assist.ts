
import { DatePipe } from "@angular/common";

export class Assist {
	static datePipe: DatePipe = new DatePipe("lt");
	static getHoverValue(first: string, second: string, firstHint: string, secondHint: string): string {
		return `<div style="display: flex;flex-direction: column">
    <div style="display: flex;align-items: center">
        <div style='width: 16px;height: 16px; background-color: #c23531;border-radius: 50%'></div>
        <div>${firstHint}:${first}</div>
    </div>
    <div style="display: flex;align-items: center">
        <div style='width: 16px;height: 16px; background-color: #2f4554;border-radius: 50%'></div>
        <div>${secondHint}:${second}</div>
    </div>
</div>`
	}

	static getTooltip(hint1: string, hint2: string): object {
		return {
			trigger: "axis",
			formatter: (params) => {
				let xDate: Date = new Date(params[0].value[0]);
				return this.datePipe.transform(xDate, "yyyy/MM/dd HH:mm:ss") +
					this.getHoverValue(params[0].value[1], params[1].value[1], hint1, hint2);
			},
			axisPointer: { animation: false }
		}
	}

	static getBaseOptions(): object {
		return {
			xAxis: {
				type: "time",
				splitNumber: 10,
				splitLine: { show: false }
			},
			yAxis: {
				type: "value",
				splitLine: { show: true }
			}
		};
	}

	static getBaseSeries(): object {
		return {
			type: 'line',
			showSymbol: true,
			smooth: true,
			symbolSize: 10,
			hoverAnimation: false
		}
	}
}