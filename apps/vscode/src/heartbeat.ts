import * as vscode from 'vscode';
import * as config from './config';
import * as detector from './detector';
import { HeartbeatQueue } from './queue';

/**
 * The data sent to the Seismic API every time a heartbeat fires.
 */
export interface HeartbeatPayload {
  file: string;
  project: string;
  language: string;
  editor: 'vscode';
  branch?: string;
  os?: string;
  machine?: string;
  lines?: number;
  cursorLine?: number;
  timezone?: string;
  time: number;
}

const TWO_MINUTES = 2 * 60 * 1000;

export class HeartbeatService {
  private lastHeartbeatTime = 0;
  private lastFile = '';
  private queue = new HeartbeatQueue();
  private hasShownInvalidKeyWarning = false;

  /**
   * Called on every relevant editor event. Decides whether
   * this specific event should actually trigger a network call,
   * based on the 2 minute rule and whether the file changed.
   */
  async handleActivity(document: vscode.TextDocument, forced = false): Promise<void> {
    if (!config.isEnabled()) return;
    if (!config.hasApiKey()) return;
    if (!detector.shouldTrack(document)) return;

    const now = Date.now();
    const fileChanged = document.fileName !== this.lastFile;

    if (!forced && !fileChanged && now - this.lastHeartbeatTime < TWO_MINUTES) {
      return; // too soon, and nothing important changed
    }

    this.lastHeartbeatTime = now;
    this.lastFile = document.fileName;

    const payload = await this.buildPayload(document);
    await this.send(payload);
  }

  private async buildPayload(document: vscode.TextDocument): Promise<HeartbeatPayload> {
    const editor = vscode.window.activeTextEditor;

    return {
      file: document.fileName,
      project: detector.detectProject(document),
      language: document.languageId,
      editor: 'vscode',
      branch: await detector.detectBranch(),
      os: detector.detectOS(),
      machine: detector.detectMachine(),
      lines: document.lineCount,
      cursorLine: editor ? editor.selection.active.line + 1 : undefined,
      timezone: detector.detectTimezone(),
      time: Date.now(),
    };
  }

  private async send(payload: HeartbeatPayload): Promise<void> {
    const apiKey = config.getApiKey();
    const apiUrl = config.getApiUrl();

    try {
      const res = await fetch(`${apiUrl}/api/heartbeat`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${apiKey}`,
        },
        body: JSON.stringify(payload),
      });

      if (res.status === 401) {
        this.notifyInvalidKey();
        return;
      }

      if (res.ok) {
        // Since we're online right now, also try to clear
        // out anything stuck in the offline queue.
        await this.queue.flush(apiKey, apiUrl);
      } else {
        this.queue.enqueue(payload);
      }
    } catch {
      // Network error (offline, DNS failure, etc) — queue it
      // silently and try again later. Never bother the user.
      this.queue.enqueue(payload);
    }
  }

  private notifyInvalidKey(): void {
    if (this.hasShownInvalidKeyWarning) return;
    this.hasShownInvalidKeyWarning = true;
    vscode.window.showWarningMessage(
      'Seismic: Invalid API key. Run "Seismic: Set API Key" to update it.',
    );
  }

  /**
   * Called periodically in the background to retry any
   * heartbeats that failed to send earlier.
   */
  async flushQueue(): Promise<void> {
    await this.queue.flush(config.getApiKey(), config.getApiUrl());
  }
}
