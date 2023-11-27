import {
  ApplicationRef,
  ComponentRef,
  EnvironmentInjector,
  Injectable,
  TemplateRef,
  Type,
  ViewContainerRef,
  createComponent,
} from '@angular/core';
import { ModalComponent } from 'src/app/components/modal/modal.component';

@Injectable({
  providedIn: 'root',
})
export class ModalService {
  newModalComponent!: ComponentRef<ModalComponent>;
  title!: string;
  constructor(
    private appRef: ApplicationRef,
    private injector: EnvironmentInjector
  ) {}

  open<C>(vcrOrComponent: Type<C>, title: string, inputs: Partial<C>) {
    this.openWithComponent(vcrOrComponent, inputs);
    this.title = title;
  }

  private openWithComponent(component: Type<unknown>, inputs: Partial<any>) {
    const newComponent = createComponent(component, {
      environmentInjector: this.injector,
    });
    for (const key of Object.keys(inputs) as Array<keyof string>) {
      newComponent.setInput(key as string, inputs[key as string] as unknown);
    }

    this.newModalComponent = createComponent(ModalComponent, {
      environmentInjector: this.injector,
      projectableNodes: [[newComponent.location.nativeElement]],
    });
    document.body.appendChild(this.newModalComponent.location.nativeElement);

    this.appRef.attachView(newComponent.hostView);
    this.appRef.attachView(this.newModalComponent.hostView);
  }

  close() {
    this.newModalComponent.instance.close();
  }
}
