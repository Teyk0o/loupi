/**
 * API client for communicating with the Loupi backend.
 * Uses httpOnly cookies for authentication (no localStorage).
 * Handles automatic token refresh and CSRF protection.
 */

const API_BASE = process.env.NEXT_PUBLIC_API_URL ?? "";

/** Standard error shape returned by the API. */
export interface ApiError {
  error: string;
  message: string;
}

/** Response from login/register (no tokens — they're in cookies). */
export interface CookieAuthResponse {
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

/** UUID v4 validation. */
const UUID_REGEX =
  /^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i;

export function isValidUUID(s: string): boolean {
  return UUID_REGEX.test(s);
}

/** Read a cookie value by name (client-side). */
function getCookie(name: string): string | null {
  if (typeof document === "undefined") return null;
  const match = document.cookie.match(new RegExp(`(?:^|; )${name}=([^;]*)`));
  return match ? decodeURIComponent(match[1]) : null;
}

/** Get CSRF token from cookie. */
function getCSRFToken(): string | null {
  return getCookie("loupi_csrf");
}

/** Map API errors to French user-friendly messages. */
export function mapApiError(err: ApiError): string {
  switch (err.error) {
    case "validation_error":
      return "Données invalides. Veuillez vérifier votre saisie.";
    case "invalid_credentials":
      return "Email ou mot de passe incorrect.";
    case "conflict":
      return "Un compte avec cet email existe déjà.";
    case "rate_limit_exceeded":
      return "Trop de tentatives. Veuillez réessayer plus tard.";
    case "account_locked":
      return "Compte temporairement verrouillé. Veuillez réessayer dans 15 minutes.";
    case "csrf_error":
      return "Session expirée. Veuillez recharger la page.";
    case "not_found":
      return "Ressource introuvable.";
    case "unauthorized":
      return "Vous devez vous connecter.";
    case "payload_too_large":
      return "La requête est trop volumineuse.";
    default:
      return err.message || "Une erreur est survenue.";
  }
}

/**
 * Attempt to refresh the access token using the refresh cookie.
 * Returns true if successful, false otherwise.
 */
async function refreshAccessToken(): Promise<boolean> {
  try {
    const res = await fetch(`${API_BASE}/v1/auth/refresh`, {
      method: "POST",
      credentials: "include",
      headers: { "Content-Type": "application/json" },
    });

    return res.ok;
  } catch {
    return false;
  }
}

/**
 * Build headers for a request, including CSRF token for state-changing methods.
 */
function buildHeaders(
  method: string,
  extra?: Record<string, string>,
): Record<string, string> {
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...extra,
  };

  // Attach CSRF token on state-changing requests
  if (method !== "GET" && method !== "HEAD") {
    const csrf = getCSRFToken();
    if (csrf) {
      headers["X-CSRF-Token"] = csrf;
    }
  }

  return headers;
}

/**
 * Make an authenticated request to the API.
 * Cookies are sent automatically. Retries once on 401 (token refresh).
 */
export async function apiFetch<T>(
  path: string,
  options: RequestInit = {},
): Promise<T> {
  const method = (options.method || "GET").toUpperCase();
  const headers = buildHeaders(method, options.headers as Record<string, string>);

  let res = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers,
    credentials: "include",
  });

  // Retry once with refreshed token on 401
  if (res.status === 401) {
    const refreshed = await refreshAccessToken();
    if (refreshed) {
      res = await fetch(`${API_BASE}${path}`, {
        ...options,
        headers: buildHeaders(method, options.headers as Record<string, string>),
        credentials: "include",
      });
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
  const method = (options.method || "GET").toUpperCase();
  const headers = buildHeaders(method, options.headers as Record<string, string>);

  const res = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers,
    credentials: "include",
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
