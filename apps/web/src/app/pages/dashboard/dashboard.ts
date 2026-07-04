import { Component, inject, OnInit, signal } from '@angular/core';
import { ApiService } from '../../core/api/api.service';
import { AuthService } from '../../core/auth/auth.service';
import { Router } from '@angular/router';

interface StatsSummary {
  totalSeconds: number;
  topLanguage: string | null;
  topProject: string | null;
  dailyAverage: number;
  currentStreak: number;
}

@Component({
  selector: 'app-dashboard',
  standalone: true,
  templateUrl: './dashboard.html',
})
export class Dashboard implements OnInit {
  private api = inject(ApiService);
  private auth = inject(AuthService);
  private router = inject(Router);

  stats = signal<StatsSummary | null>(null);
  loading = signal(true);

  ngOnInit() {
    this.api.get<StatsSummary>('/api/stats/summary', { range: 'today' }).subscribe({
      next: (data) => {
        this.stats.set(data);
        this.loading.set(false);
      },
      error: () => this.loading.set(false),
    });
  }

  formatSeconds(seconds: number): string {
    if (seconds < 60) return '< 1m';
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    return hours > 0 ? `${hours}h ${minutes}m` : `${minutes}m`;
  }

  logout() {
    this.auth.logout();
    void this.router.navigate(['/login']);
  }
}
