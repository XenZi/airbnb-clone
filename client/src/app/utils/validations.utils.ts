import { AbstractControl, ValidationErrors, ValidatorFn } from '@angular/forms';

export const customPasswordStrengthValidator = (): ValidatorFn => {
  return (control: AbstractControl): ValidationErrors | null => {
    const value = control.value;
    if (!value) {
      return null;
    }
    const hasUpperCase = /[A-Z]+/.test(value);
    const hasNumeric = /[0-9]+/.test(value);
    const hasSpecialChars = /[!@#$%^&*()_+\?><>:';\]\[']/.test(value);
    const minLength = (value as string).length >= 8;
    const passwordValid =
      hasUpperCase && hasNumeric && hasSpecialChars && minLength;
    return !passwordValid ? { passwordStrength: true } : null;
  };
};
