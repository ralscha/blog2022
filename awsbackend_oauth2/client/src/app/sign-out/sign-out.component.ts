import {Component} from '@angular/core';
import {RouterLink} from '@angular/router';
import {IonContent, IonHeader, IonRouterLinkWithHref, IonTitle, IonToolbar} from "@ionic/angular/standalone";

@Component({
  selector: 'app-sign-out',
  templateUrl: './sign-out.component.html',
  styleUrls: ['./sign-out.component.scss'],
  imports: [RouterLink, IonRouterLinkWithHref, IonHeader, IonToolbar, IonTitle, IonContent]
})
export class SignOutComponent {
}
