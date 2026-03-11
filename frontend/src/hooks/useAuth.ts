"use client";

/**
 * Authentication hook for managing user state across the application.
 * Uses cookie-based auth — no tokens stored in localStorage.
 * Checks session validity on mount via /v1/auth/me.
 */

import { useState, useEffect, useCallback } from "react";
import {
  apiFetch,
  apiPublicFetch,
  type CookieAuthResponse,
  type UserResponse,
} from "@/lib/api";

interface UseAuthReturn {
  user: UserResponse | null;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (email: string, password: string, firstName?: string) => Promise<void>;
  logout: () => Promise<void>;
  deleteAccount: () => Promise<void>;
}

export function useAuth(): UseAuthReturn {
  const [user, setUser] = useState<UserResponse | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  // Check session validity on mount
  useEffect(() => {
    let cancelled = false;
    async function checkSession() {
      try {
        const userData = await apiFetch<UserResponse>("/v1/auth/me");
        if (!cancelled) setUser(userData);
      } catch {
        if (!cancelled) setUser(null);
      } finally {
        if (!cancelled) setIsLoading(false);
      }
    }
    checkSession();
    return () => { cancelled = true; };
  }, []);

  const login = useCallback(async (email: string, password: string) => {
    const response = await apiPublicFetch<CookieAuthResponse>("/v1/auth/login", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    });
    setUser(response.user);
    // Full page reload so the browser sends the newly-set cookie to Next.js middleware
    window.location.href = "/journal";
  }, []);

  const register = useCallback(async (email: string, password: string, firstName?: string) => {
    const response = await apiPublicFetch<CookieAuthResponse>("/v1/auth/register", {
      method: "POST",
      body: JSON.stringify({ email, password, first_name: firstName || undefined }),
    });
    setUser(response.user);
    window.location.href = "/journal";
  }, []);

  const logout = useCallback(async () => {
    try {
      await apiFetch("/v1/auth/logout", { method: "POST" });
    } catch {
      // Best effort — cookies will be cleared server-side
    }
    setUser(null);
    window.location.href = "/login";
  }, []);

  const deleteAccount = useCallback(async () => {
    await apiFetch("/v1/auth/account", { method: "DELETE" });
    setUser(null);
    window.location.href = "/login";
  }, []);

  return { user, isLoading, login, register, logout, deleteAccount };
}
