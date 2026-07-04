import { HttpClient } from '@angular/common/http';
import { Injectable, inject, signal } from '@angular/core';
import { map, Observable, tap } from 'rxjs';
import { environment } from '../../../environments/environment';

interface RefreshResponse {
  accessToken: string;
}

interface ApiEnvelope<T> {
  success: boolean;
  message: string;
  data: T;
}

@Injectable({ providedIn: 'root' })
export class AuthService {
  private readonly TOKEN_KEY = 'accessToken';
  private http = inject(HttpClient);
  private baseUrl = environment.apiUrl;

  isLoggedIn = signal(this.hasToken());

  private hasToken(): boolean {
    return !!localStorage.getItem(this.TOKEN_KEY);
  }

  getToken(): string | null {
    return localStorage.getItem(this.TOKEN_KEY);
  }

  setToken(token: string): void {
    localStorage.setItem(this.TOKEN_KEY, token);
    this.isLoggedIn.set(true);
  }

  logout(): void {
    localStorage.removeItem(this.TOKEN_KEY);
    this.isLoggedIn.set(false);
  }

  // Calls the refresh endpoint. The refresh token itself is
  // sent automatically as an httpOnly cookie by the browser
  refreshToken(): Observable<RefreshResponse> {
    return this.http
      .post<ApiEnvelope<RefreshResponse>>(
        `${this.baseUrl}/api/auth/refresh`,
        {},
        { withCredentials: true },
      )
      .pipe(
        map((res) => res.data),
        tap((data) => this.setToken(data.accessToken)),
      );
  }
}
