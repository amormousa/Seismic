import { Routes } from '@angular/router';
import { Login } from './pages/auth/login/login';
import { Verify } from './pages/auth/verify/verify';
import { Dashboard } from './pages/dashboard/dashboard';

export const routes: Routes = [
  { path: 'login', component: Login },
  { path: 'verify', component: Verify },
  { path: 'dashboard', component: Dashboard },
];
