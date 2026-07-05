import { Injectable, inject, signal } from '@angular/core';
import { ApiService } from '../api/api.service';

interface CurrentUser {
  username: string;
  country: string | null;
}

interface StatsSummary {
  currentStreak: number;
}

@Injectable({ providedIn: 'root' })
export class UserService {
  private api = inject(ApiService);

  user = signal<CurrentUser | null>(null);
  streak = signal<number>(0);

  load() {
    this.api.get<CurrentUser>('/api/auth/me').subscribe({
      next: (data) => this.user.set(data),
      error: () => {},
    });

    this.api.get<StatsSummary>('/api/stats/summary', { range: 'today' }).subscribe({
      next: (data) => this.streak.set(data.currentStreak),
      error: () => {},
    });
  }
}
