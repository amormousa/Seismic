import { Routes } from '@angular/router';
import { Login } from './pages/auth/login/login';
import { Verify } from './pages/auth/verify/verify';
import { Dashboard } from './pages/dashboard/dashboard';
import { Leaderboard } from './pages/leaderboard/leaderboard';
import { authGuard } from './core/auth/auth.guard';
import { guestGuard } from './core/auth/guest.guard';
import { Settings } from './pages/settings/settings';

export const routes: Routes = [
  { path: 'login', component: Login, canActivate: [guestGuard] },
  { path: 'verify', component: Verify },
  { path: 'dashboard', component: Dashboard, canActivate: [authGuard] },
  { path: 'leaderboard', component: Leaderboard },
  { path: 'settings', component: Settings },
];
