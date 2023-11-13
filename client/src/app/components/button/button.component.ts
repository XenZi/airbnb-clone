import { Component, EventEmitter, Input, Output } from '@angular/core';

@Component({
  selector: 'app-button',
  templateUrl: './button.component.html',
  styleUrls: ['./button.component.scss'],
})
export class ButtonComponent {
  @Input() type: string | undefined;
  @Input() size!: string;
  @Input() color!: string;
  @Input() class: string | undefined;
  @Output() onClick: EventEmitter<void> = new EventEmitter<void>();
}
