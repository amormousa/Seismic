export const DEV_API_URL = 'http://localhost:5024';
export const PUBLISHED_API_URL = 'https://correct-wolverine-majoramari-6049fd71.koyeb.app';

export function getDefaultApiUrl(isDevelopment = false): string {
  return isDevelopment ? DEV_API_URL : PUBLISHED_API_URL;
}
