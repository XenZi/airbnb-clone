import {
  Component,
  ElementRef,
  HostListener,
  Input,
  ViewChild,
} from '@angular/core';
import { ModalService } from 'src/app/services/modal/modal.service';
@Component({
  selector: 'app-modal',
  templateUrl: './modal.component.html',
  styleUrls: ['./modal.component.scss'],
})
export class ModalComponent {
  @ViewChild('modal') modal!: ElementRef<HTMLDivElement>;
  title!: string;
  constructor(
    private modalService: ModalService,
    private element: ElementRef
  ) {}

  ngOnInit() {
    this.title = this.modalService.title;
  }

  removeElement(element: HTMLDivElement) {
    element.remove();
  }

  @HostListener('document:keydown.escape')
  onEscape() {
    this.modalService.close();
  }

  @HostListener('document:click', ['$event'])
  onClickOutside(event: any) {
    if (event.target.classList.contains('modal')) {
      this.modalService.close();
    }
  }

  onClose() {
    this.modalService.close();
  }

  close() {
    this.removeElement;
    this.element.nativeElement.remove();
  }
}
