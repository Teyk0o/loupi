/**
 * API client for communicating with the Loupi backend.
 * Handles authentication tokens, automatic refresh, and error formatting.
 */

const API_BASE = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

/** Standard error shape returned by the API. */
export interface ApiError {
  error: string;
  message: string;
}

/** Token pair returned by auth endpoints. */
export interface TokenResponse {
  access_token: string;
  refresh_token: string;
  expires_in: number;
  user: UserResponse;
}

/** Public user data. */
export interface UserResponse {
  id: string;
  email: string;
  first_name?: string;
  email_verified: boolean;
  created_at: string;
}

/** Stored auth state. */
interface AuthState {
  accessToken: string;
  refreshToken: string;
  user: UserResponse;
}

const AUTH_KEY = "loupi_auth";

/** Retrieve stored auth state from localStorage. */
export function getAuth(): AuthState | null {
  if (typeof window === "undefined") return null;
  const raw = localStorage.getItem(AUTH_KEY);
  if (!raw) return null;
  try {
    return JSON.parse(raw);
  } catch {
    return null;
  }
}

/** Persist auth state to localStorage. */
export function setAuth(tokens: TokenResponse): void {
  const state: AuthState = {
    accessToken: tokens.access_token,
    refreshToken: tokens.refresh_token,
    user: tokens.user,
  };
  localStorage.setItem(AUTH_KEY, JSON.stringify(state));
}

/** Clear stored auth state. */
export function clearAuth(): void {
  localStorage.removeItem(AUTH_KEY);
}

/**
 * Attempt to refresh the access token using the stored refresh token.
 * Returns true if successful, false otherwise.
 */
async function refreshAccessToken(): Promise<boolean> {
  const auth = getAuth();
  if (!auth) return false;

  try {
    const res = await fetch(`${API_BASE}/v1/auth/refresh`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ refresh_token: auth.refreshToken }),
    });

    if (!res.ok) {
      clearAuth();
      return false;
    }

    const tokens: TokenResponse = await res.json();
    setAuth(tokens);
    return true;
  } catch {
    clearAuth();
    return false;
  }
}

/**
 * Make an authenticated request to the API.
 * Automatically attaches the access token and retries once on 401 (token refresh).
 */
export async function apiFetch<T>(
  path: string,
  options: RequestInit = {},
): Promise<T> {
  const auth = getAuth();

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(options.headers as Record<string, string>),
  };

  if (auth) {
    headers["Authorization"] = `Bearer ${auth.accessToken}`;
  }

  let res = await fetch(`${API_BASE}${path}`, { ...options, headers });

  // Retry once with refreshed token on 401
  if (res.status === 401 && auth) {
    const refreshed = await refreshAccessToken();
    if (refreshed) {
      const newAuth = getAuth();
      if (newAuth) {
        headers["Authorization"] = `Bearer ${newAuth.accessToken}`;
        res = await fetch(`${API_BASE}${path}`, { ...options, headers });
      }
    }
  }

  if (!res.ok) {
    const err: ApiError = await res.json().catch(() => ({
      error: "unknown",
      message: "An unexpected error occurred",
    }));
    throw err;
  }

  return res.json();
}

/**
 * Make an unauthenticated request to the API (for login/register).
 */
export async function apiPublicFetch<T>(
  path: string,
  options: RequestInit = {},
): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...(options.headers as Record<string, string>),
    },
  });

  if (!res.ok) {
    const err: ApiError = await res.json().catch(() => ({
      error: "unknown",
      message: "An unexpected error occurred",
    }));
    throw err;
  }

  return res.json();
}
