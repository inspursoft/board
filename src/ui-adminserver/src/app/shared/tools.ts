export function initTipsArray(num: number, flag: boolean): Array<boolean> {
    const temp: boolean[] = [];
    for (let i = 0; i < num; i++) {
        temp.push(flag);
    }
    return temp;
}

export function json2String(jsonObject: object): string {
    let resultString = '';
    for (let key in jsonObject) {
        if (jsonObject[key] instanceof Object) {
            resultString += json2String(jsonObject[key]);
        } else {
            resultString += key + ' = ' + jsonObject[key] + '\n';
        }
    }
    return resultString;
}
