import { Component, inject, signal, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { ApiService } from '../../../core/api/api.service';

interface VerifyResponse {
  newUser: boolean;
  accessToken?: string;
  signupToken?: string;
  email?: string;
}

interface SignupResponse {
  accessToken: string;
  user: { id: string; username: string; email: string };
}

@Component({
  selector: 'app-verify',
  standalone: true,
  imports: [FormsModule],
  templateUrl: './verify.html',
})
export class Verify implements OnInit {
  private api = inject(ApiService);
  private route = inject(ActivatedRoute);
  private router = inject(Router);

  loading = signal(true);
  error = signal('');

  // Set to true only if the backend says this is a new user
  isNewUser = signal(false);
  signupToken = '';
  email = signal('');

  username = signal('');
  displayName = signal('');
  signupLoading = signal(false);
  signupError = signal('');

  ngOnInit() {
    const token = this.route.snapshot.queryParamMap.get('token');
    if (!token) {
      this.error.set('Missing token');
      this.loading.set(false);
      return;
    }

    this.api.get<VerifyResponse>('/api/auth/verify', { token }).subscribe({
      next: (data) => {
        this.loading.set(false);

        if (data.newUser && data.signupToken) {
          this.isNewUser.set(true);
          this.signupToken = data.signupToken;
          this.email.set(data.email ?? '');
        } else if (data.accessToken) {
          localStorage.setItem('accessToken', data.accessToken);
          this.router.navigate(['/dashboard']);
        }
      },
      error: () => {
        this.error.set('Invalid or expired link');
        this.loading.set(false);
      },
    });
  }

  completeSignup() {
    this.signupLoading.set(true);
    this.signupError.set('');

    this.api
      .post<SignupResponse>('/api/auth/complete-signup', {
        signupToken: this.signupToken,
        username: this.username(),
        displayName: this.displayName(),
      })
      .subscribe({
        next: (data) => {
          localStorage.setItem('accessToken', data.accessToken);
          this.router.navigate(['/dashboard']);
        },
        error: (err) => {
          this.signupLoading.set(false);
          this.signupError.set(err.error?.message || 'Something went wrong');
        },
      });
  }
}
