import { Component, inject, signal, OnInit } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ApiService } from '../../core/api/api.service';
import { ToastService } from '../../core/toast/toast.service';

interface ApiKeyResponse {
  apiKey: string;
}

interface PrivacySettings {
  hideProjects: boolean;
  hideTime: boolean;
  hideLanguages: boolean;
  hideLeaderboard: boolean;
  profilePublic: boolean;
}

type SettingsSection = 'apikey' | 'privacy' | 'account';
type PrivacySubTab = 'toggles';
type AccountSubTab = 'reset' | 'delete';

@Component({
  selector: 'app-settings',
  standalone: true,
  imports: [FormsModule],
  templateUrl: './settings.html',
})
export class Settings implements OnInit {
  private api = inject(ApiService);
  private toast = inject(ToastService);

  activeSection = signal<SettingsSection>('apikey');
  accountSubTab = signal<AccountSubTab>('reset');

  apiKey = signal('');
  apiKeyVisible = signal(false);
  regenerating = signal(false);

  privacy = signal<PrivacySettings | null>(null);
  privacyDirty = signal(false);
  savingPrivacy = signal(false);

  showDeleteConfirm = signal(false);
  deleteConfirmText = signal('');

  ngOnInit() {
    this.loadApiKey();
    this.loadPrivacy();
  }

  setSection(section: SettingsSection) {
    this.activeSection.set(section);
  }

  setAccountSubTab(tab: AccountSubTab) {
    this.accountSubTab.set(tab);
  }

  private loadApiKey() {
    this.api.get<ApiKeyResponse>('/api/auth/apikey').subscribe({
      next: (data) => this.apiKey.set(data.apiKey),
      error: () => this.toast.error('Failed to load API key'),
    });
  }

  toggleApiKeyVisible() {
    this.apiKeyVisible.set(!this.apiKeyVisible());
  }

  copyApiKey() {
    navigator.clipboard.writeText(this.apiKey());
    this.toast.success('API key copied to clipboard');
  }

  regenerateApiKey() {
    this.regenerating.set(true);
    this.api.post<ApiKeyResponse>('/api/auth/apikey/regenerate', {}).subscribe({
      next: (data) => {
        this.apiKey.set(data.apiKey);
        this.regenerating.set(false);
        this.toast.success('API key regenerated. Update it in your editors.');
      },
      error: () => {
        this.regenerating.set(false);
        this.toast.error('Failed to regenerate API key');
      },
    });
  }

  private loadPrivacy() {
    this.api.get<PrivacySettings>('/api/settings/privacy').subscribe({
      next: (data) => this.privacy.set(data),
      error: () => this.toast.error('Failed to load privacy settings'),
    });
  }

  toggleLocal(key: keyof PrivacySettings, value: boolean) {
    const current = this.privacy();
    if (!current) return;
    this.privacy.set({ ...current, [key]: value });
    this.privacyDirty.set(true);
  }

  savePrivacy() {
    const current = this.privacy();
    if (!current) return;

    this.savingPrivacy.set(true);
    this.api.post('/api/settings/privacy', current).subscribe({
      next: () => {
        this.savingPrivacy.set(false);
        this.privacyDirty.set(false);
        this.toast.success('Privacy settings saved');
      },
      error: () => {
        this.savingPrivacy.set(false);
        this.toast.error('Failed to save settings');
      },
    });
  }

  resetTimers() {
    const confirmed = confirm(
      'This will permanently delete all your tracked coding time. This cannot be undone. Continue?',
    );
    if (!confirmed) return;

    this.api.post('/api/settings/reset-timers', {}).subscribe({
      next: () => this.toast.success('All coding stats have been reset'),
      error: () => this.toast.error('Failed to reset timers'),
    });
  }

  confirmDelete() {
    if (this.deleteConfirmText() !== 'delete') return;

    this.api.post('/api/settings/account', {}).subscribe({
      next: () => {
        localStorage.clear();
        window.location.href = '/login';
      },
      error: () => this.toast.error('Failed to delete account'),
    });
  }
}
