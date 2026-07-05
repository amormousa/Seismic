import * as fs from 'node:fs';
import * as os from 'node:os';
import * as path from 'node:path';
import { PUBLISHED_API_URL } from './apiUrls';

/**
 * Reads and writes a shared config file at ~/.seismic.cfg
 * This lets the Editors all use the same
 * API key without configuring each editor separately
 */

const CONFIG_PATH = path.join(os.homedir(), '.seismic.cfg');

interface GlobalConfig {
  apiKey: string;
  apiUrl: string;
}

export function readGlobalConfig(): GlobalConfig | null {
  if (!fs.existsSync(CONFIG_PATH)) return null;

  const content = fs.readFileSync(CONFIG_PATH, 'utf-8');
  const apiKey = content.match(/api_key\s*=\s*(.+)/)?.[1]?.trim() ?? '';
  const apiUrl = content.match(/api_url\s*=\s*(.+)/)?.[1]?.trim() ?? '';

  if (!apiKey) return null;

  return { apiKey, apiUrl: apiUrl || PUBLISHED_API_URL };
}

export function writeGlobalConfig(config: GlobalConfig): void {
  const content = `[settings]\napi_key = ${config.apiKey}\napi_url = ${config.apiUrl}\n`;
  fs.writeFileSync(CONFIG_PATH, content, 'utf-8');
}
