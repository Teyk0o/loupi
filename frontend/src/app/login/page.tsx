"use client";

/**
 * Login page.
 * Allows existing users to authenticate with email and password.
 */

import { useState, type FormEvent } from "react";
import Link from "next/link";
import { useAuthContext } from "@/components/providers/AuthProvider";
import { GuestGuard } from "@/components/guards/GuestGuard";
import { Button } from "@/components/ui/Button";
import { Input } from "@/components/ui/Input";
import { Logo } from "@/components/ui/Logo";
import type { ApiError } from "@/lib/api";

export default function LoginPage() {
  const { login } = useAuthContext();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError("");
    setIsSubmitting(true);

    try {
      await login(email, password);
    } catch (err) {
      const apiErr = err as ApiError;
      setError(apiErr.message || "Une erreur est survenue");
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <GuestGuard>
      <div className="flex min-h-screen flex-col items-center justify-center px-6">
        <div className="w-full max-w-sm">
          <div className="mb-10 flex justify-center">
            <Logo width={140} height={46} />
          </div>

          <h1 className="mb-1 text-center font-heading text-2xl font-semibold text-foreground">
            Content de vous revoir
          </h1>
          <p className="mb-8 text-center text-sm text-foreground-secondary">
            Connectez-vous pour accéder à votre journal
          </p>

          {error ? (
            <div className="mb-4 rounded-[--radius-md] border border-danger/30 bg-danger/10 px-4 py-3 text-sm text-danger">
              {error}
            </div>
          ) : null}

          <form onSubmit={handleSubmit} className="flex flex-col gap-4">
            <Input
              id="email"
              label="Email"
              type="email"
              placeholder="vous@exemple.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              autoComplete="email"
              required
            />
            <Input
              id="password"
              label="Mot de passe"
              type="password"
              placeholder="••••••••"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              autoComplete="current-password"
              required
            />

            <Button
              type="submit"
              isLoading={isSubmitting}
              className="mt-2 w-full"
            >
              Se connecter
            </Button>
          </form>

          <p className="mt-6 text-center text-sm text-foreground-secondary">
            Pas encore de compte ?{" "}
            <Link
              href="/register"
              className="font-medium text-primary hover:underline"
            >
              Créer un compte
            </Link>
          </p>
        </div>
      </div>
    </GuestGuard>
  );
}
