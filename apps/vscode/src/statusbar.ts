import * as vscode from 'vscode';
import * as config from './config';

/**
 * Manages the small "⏱ 3h 42m" indicator shown in the
 * bottom status bar, and keeps it updated every minute.
 */
export class StatusBarManager {
  private item: vscode.StatusBarItem;
  private timer: ReturnType<typeof setInterval> | undefined;

  constructor() {
    this.item = vscode.window.createStatusBarItem(vscode.StatusBarAlignment.Left, 10);
    this.item.show();
    this.refresh().catch(console.error);
  }

  startUpdating(): void {
    this.timer = setInterval(() => {
      this.refresh().catch(console.error);
    }, 60 * 1000);
  }

  dispose(): void {
    if (this.timer) clearInterval(this.timer);
    this.item.dispose();
  }

  async refresh(): Promise<void> {
    if (!config.isStatusBarEnabled()) {
      this.item.hide();
      return;
    }
    this.item.show();

    if (!config.hasApiKey()) {
      this.item.text = '$(clock) Seismic: Set API Key';
      this.item.tooltip = 'Click to set your Seismic API key';
      this.item.command = 'seismic.setApiKey';
      return;
    }

    if (!config.isEnabled()) {
      this.item.text = '$(clock) Seismic: Paused';
      this.item.tooltip = 'Seismic tracking is disabled';
      this.item.command = 'seismic.openDashboard';
      return;
    }

    try {
      const seconds = await this.fetchTodaySeconds();
      this.item.text = `$(clock) ${this.formatSeconds(seconds)}`;
      this.item.tooltip = "Today's coding time on Seismic\nClick to open dashboard";
      this.item.command = 'seismic.openDashboard';
    } catch {
      this.item.text = '$(warning) Seismic';
      this.item.tooltip = 'Could not connect to Seismic';
    }
  }

  private async fetchTodaySeconds(): Promise<number> {
    const res = await fetch(`${config.getApiUrl()}/api/stats/summary?range=today`, {
      headers: { Authorization: `Bearer ${config.getApiKey()}` },
    });
    if (!res.ok) throw new Error('failed to fetch stats');
    const data = (await res.json()) as { totalSeconds: number };
    return data.totalSeconds;
  }

  private formatSeconds(seconds: number): string {
    if (seconds < 60) return '< 1m';
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    return hours > 0 ? `${hours}h ${minutes}m` : `${minutes}m`;
  }
}
