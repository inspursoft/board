import { FormControlBase } from '../form-control-base';

export class TextInput extends FormControlBase<string> {
  
  constructor(options: {} = {}){
    super(options);
    this.type = 'text';
  }
}