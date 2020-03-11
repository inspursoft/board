import { FormControlBase } from '../form-control-base';

export class Dropdown extends FormControlBase<string> {

  options: {key: string, value: string}[] = [];

  constructor(options: {} = {}){
    super(options);
    this.options = options['options'] || [];
  }
}