"use client";

/**
 * Authentication hook for managing user state across the application.
 * Provides login, register, logout, and current user access.
 */

import { useState, useEffect, useCallback } from "react";
import {
  getAuth,
  setAuth,
  clearAuth,
  apiPublicFetch,
  apiFetch,
  type TokenResponse,
  type UserResponse,
} from "@/lib/api";

interface UseAuthReturn {
  user: UserResponse | null;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (email: string, password: string, firstName?: string) => Promise<void>;
  logout: () => void;
  deleteAccount: () => Promise<void>;
}

export function useAuth(): UseAuthReturn {
  const [user, setUser] = useState<UserResponse | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  // Restore auth state from localStorage on mount
  useEffect(() => {
    const auth = getAuth();
    if (auth) {
      setUser(auth.user);
    }
    setIsLoading(false);
  }, []);

  const login = useCallback(async (email: string, password: string) => {
    const tokens = await apiPublicFetch<TokenResponse>("/v1/auth/login", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    });
    setAuth(tokens);
    setUser(tokens.user);
  }, []);

  const register = useCallback(async (email: string, password: string, firstName?: string) => {
    const tokens = await apiPublicFetch<TokenResponse>("/v1/auth/register", {
      method: "POST",
      body: JSON.stringify({ email, password, first_name: firstName || undefined }),
    });
    setAuth(tokens);
    setUser(tokens.user);
  }, []);

  const logout = useCallback(() => {
    clearAuth();
    setUser(null);
    window.location.href = "/login";
  }, []);

  const deleteAccount = useCallback(async () => {
    await apiFetch("/v1/auth/account", { method: "DELETE" });
    clearAuth();
    setUser(null);
  }, []);

  return { user, isLoading, login, register, logout, deleteAccount };
}
