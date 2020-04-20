import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';
import { FormGroup, FormControl, FormBuilder, ValidatorFn, Validators, AbstractControl, ValidationErrors } from '@angular/forms';
import { Subscription } from 'rxjs';

@Component({
  selector: 'app-my-input-template',
  templateUrl: './my-input-template.component.html',
  styleUrls: ['./my-input-template.component.css']
})
export class MyInputTemplateComponent implements OnInit {
  @Input() label = '';
  @Input() type: 'text' | 'number' | 'boolean' | 'enum' | 'password' = 'text';
  @Input() required = false;
  @Input() pattern: RegExp;
  @Input() min: number;
  @Input() max: number;
  @Input() minLength: number;
  @Input() maxLength: number;
  @Input() value: string | number | boolean;
  @Input() enumType: 'text' | 'number' = 'text';
  @Input() enumItem: Array<any>;
  @Input() enumInline = false;
  @Input() enableCustomItem = false;
  @Input() customItemPlaceholder = 'Custom Input';
  @Input() customItemWidth = '';
  customItem: any;
  @Input() confirmLabel = 'confirm Password';
  confirmPasswordValue = '';
  @Input() placeholder = '';
  @Input() helper = '';
  @Input() showHelper: 'onfocus' | 'always' | 'never' = 'onfocus';
  @Input() requiredMsg = 'Required!';
  @Input() patternMsg = 'Error value!';
  @Input() minMsg = 'Too small!';
  @Input() maxMsg = 'Too big!';
  @Input() minLengthMsg = 'Too short!';
  @Input() maxLengthMsg = 'Too long!';
  @Input() confirmMsg = 'Inconsistent password!';
  @Input() inputUpdateOn: 'change' | 'blur' | 'submit' = 'blur';
  @Input() firstItem = false;
  @Output() editEvent: EventEmitter<any>;
  @Output() commitEvent: EventEmitter<any>;
  @Output() valueChanges: EventEmitter<any>;

  enableHelper = false;
  showError = false;
  showPassword = false;
  showPasswordconfirm = false;

  inputFormGroup: FormGroup;
  inputControl: FormControl;
  confirmControl: FormControl;

  inputValidatorFns: Array<ValidatorFn>;
  valueSubscription: Subscription;
  statusSubscription: Subscription;

  constructor(private fb: FormBuilder) {
    this.editEvent = new EventEmitter();
    this.commitEvent = new EventEmitter();
    this.valueChanges = new EventEmitter();
    this.inputControl = this.fb.control({ value: '', disabled: false });
    this.confirmControl = this.fb.control({ value: '', disabled: true });
    this.inputValidatorFns = new Array<ValidatorFn>();
  }

  static passwordValidator(source: string, target: string): ValidatorFn {
    return (self: AbstractControl): { [key: string]: any } => {    // 这里严格按照ValidatorFn的声明来
      const _form = self.parent;
      if (_form) {
        const sourceControl: AbstractControl = _form.controls[source];
        const targetControl: AbstractControl = _form.controls[target];
        if (targetControl.value && sourceControl.value && targetControl.value !== sourceControl.value) {   // 如果两个值不一致
          console.log('###' + targetControl.value + '###' + sourceControl.value)
          return { verifyPassword: 'verify-password' };
        }
      }
    };
  }

  ngOnInit() {
    this.enableHelper = this.showHelper === 'always' ? true : false;
    // if (this.type === 'boolean') {
    //   this.inputControl.setValue(false);
    // }
    this.inputFormGroup = this.fb.group({
      inputControl: this.inputControl,
      confirmControl: this.confirmControl
    }, {
      updateOn: this.inputUpdateOn
    });
    this.installValidators();
  }

  @Input()
  set disabled(value: boolean) {
    if (value) {
      this.inputControl.disable();
    } else {
      this.inputControl.enable();
    }
  }

  get disabled(): boolean {
    return this.inputControl.disabled;
  }

  @Input()
  set confirmPassword(value: boolean) {
    if (value) {
      this.confirmControl.enable();
    }
  }

  get confirmPassword(): boolean {
    return this.confirmControl.enabled;
  }

  @Input()
  set defaultValue(value: string | number | boolean) {
    console.log(value)
    this.inputControl.setValue(value);
  }

  installValidators() {
    this.valueSubscription = this.inputControl.valueChanges.subscribe((value: any) => {
      this.valueChanges.next(value);
      console.log('change' + value);
    });
    this.statusSubscription = this.inputControl.statusChanges.subscribe(() => {
      if (this.inputControl.valid) {
        this.commitEvent.emit(this.inputControl.value);
      }
    });
    if (this.required) {
      this.inputValidatorFns.push(Validators.required);
    }
    if (this.type === 'number') {
      if (this.max) {
        this.inputValidatorFns.push(Validators.max(this.max));
      }
      if (this.min) {
        this.inputValidatorFns.push(Validators.min(this.min));
      }
    } else if (this.type === 'text' || this.type === 'password') {
      if (this.maxLength > 0) {
        this.inputValidatorFns.push(Validators.maxLength(this.maxLength));
      }
      if (this.minLength > 0) {
        this.inputValidatorFns.push(Validators.minLength(this.minLength));
      }
      if (this.pattern) {
        this.inputValidatorFns.push(Validators.pattern(this.pattern));
      }
      if (this.type === 'password' && this.confirmPassword) {
        console.log("confirmPassword")
        this.inputValidatorFns.push(MyInputTemplateComponent.passwordValidator('inputControl', 'confirmControl'));
        this.confirmControl.setValidators(this.inputValidatorFns);
      }
    }
    this.inputControl.setValidators(this.inputValidatorFns);
  }

  unInstallValidators() {
    if (this.valueSubscription) {
      this.valueSubscription.unsubscribe();
    }
    if (this.statusSubscription) {
      this.statusSubscription.unsubscribe();
    }
    this.inputControl.clearValidators();
    this.inputControl.clearAsyncValidators();
  }

  getValidatorMessage(errors: ValidationErrors): string {
    let result = '';
    if (errors) {
      if (Reflect.has(errors, 'required')) {
        result = this.requiredMsg;
      } else if (Reflect.has(errors, 'pattern')) {
        result = this.patternMsg;
      } else if (Reflect.has(errors, 'maxlength')) {
        result = this.maxLengthMsg;
      } else if (Reflect.has(errors, 'minlength')) {
        result = this.minLengthMsg;
      } else if (Reflect.has(errors, 'max')) {
        result = this.maxMsg;
      } else if (Reflect.has(errors, 'min')) {
        result = this.minMsg;
      } else if (Reflect.has(errors, 'verifyPassword')) {
        console.log("confirmPassword Error")
        result = this.confirmMsg;
      } else if (Object.keys(errors).length > 0) {
        result = errors[Object.keys(errors)[0]];
      }
    }
    return result;
  }

  onInputFocus() {
    this.showError = false;
    if (this.showHelper === 'onfocus' && this.inputControl.enabled) {
      this.enableHelper = true;
    }
  }

  onInputBlur() {
    this.inputControl.updateValueAndValidity();
    if (this.confirmPassword) {
      this.confirmControl.updateValueAndValidity();
    }
    if (this.showHelper === 'onfocus' && (!this.inputControl.value || this.inputControl.valid)) {
      if (this.confirmPassword && this.confirmControl.invalid) {
        this.enableHelper = true;
      } else {
        this.enableHelper = false;
      }
    }
  }

  onCustomFocus(value: string | number) {
    console.log(value);
    if (this.showHelper === 'onfocus' && this.inputControl.enabled) {
      this.enableHelper = true;
    }
    // this.inputControl.setValue(value);
    // this.checkSelf();
  }

  onCustomBlur(value: string | number) {
    if (typeof (value) === 'string' && this.enumType === 'text') {
      console.log('string');
      value = value.trim();
      this.customItem = value;
    }
    this.inputControl.setValue(value);
    this.checkSelf();
  }

  togglePassword() {
    console.log('toggle')
    this.showPassword = !this.showPassword;
  }

  togglePasswordconfirm() {
    console.log('toggle C')
    this.showPasswordconfirm = !this.showPasswordconfirm;
  }

  public checkSelf() {
    if (this.inputControl.enabled) {
      this.inputControl.markAsTouched({ onlySelf: true });
      this.inputControl.updateValueAndValidity();
      this.showError = this.inputControl.errors ? true : false;
      if (this.showHelper === 'onfocus') {
        this.enableHelper = this.inputControl.errors ? true : false;
      }
    }
    if (this.confirmControl.enabled) {
      this.confirmControl.markAsTouched({ onlySelf: true });
      this.confirmControl.updateValueAndValidity();
    }
  }

  public get isValid(): boolean {
    return this.inputControl.valid;
  }
}
