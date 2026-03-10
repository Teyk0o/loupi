"use client";

/**
 * Registration page.
 * Allows new users to create an account with email and password.
 */

import { useState, type FormEvent } from "react";
import Link from "next/link";
import { useAuthContext } from "@/components/providers/AuthProvider";
import { GuestGuard } from "@/components/guards/GuestGuard";
import { Button } from "@/components/ui/Button";
import { Input } from "@/components/ui/Input";
import { Logo } from "@/components/ui/Logo";
import { mapApiError, type ApiError } from "@/lib/api";

const EMAIL_REGEX = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

function validatePassword(pw: string): string | null {
  if (pw.length < 8) return "Le mot de passe doit contenir au moins 8 caractères.";
  if (!/[A-Z]/.test(pw)) return "Le mot de passe doit contenir au moins une majuscule.";
  if (!/[a-z]/.test(pw)) return "Le mot de passe doit contenir au moins une minuscule.";
  if (!/[0-9]/.test(pw)) return "Le mot de passe doit contenir au moins un chiffre.";
  if (!/[!@#$%^&*()_+\-=[\]{};':"\\|,.<>/?`~]/.test(pw))
    return "Le mot de passe doit contenir au moins un caractère spécial.";
  return null;
}

export default function RegisterPage() {
  const { register } = useAuthContext();
  const [firstName, setFirstName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [error, setError] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError("");

    if (!EMAIL_REGEX.test(email)) {
      setError("Veuillez entrer un email valide.");
      return;
    }

    if (password !== confirmPassword) {
      setError("Les mots de passe ne correspondent pas.");
      return;
    }

    const pwError = validatePassword(password);
    if (pwError) {
      setError(pwError);
      return;
    }

    setIsSubmitting(true);

    try {
      await register(email, password, firstName || undefined);
    } catch (err) {
      const apiErr = err as ApiError;
      setError(mapApiError(apiErr));
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
            Créer un compte
          </h1>
          <p className="mb-8 text-center text-sm text-foreground-secondary">
            Commencez à suivre votre alimentation et bien-être
          </p>

          {error ? (
            <div className="mb-4 rounded-[--radius-md] border border-danger/30 bg-danger/10 px-4 py-3 text-sm text-danger">
              {error}
            </div>
          ) : null}

          <form onSubmit={handleSubmit} className="flex flex-col gap-4">
            <Input
              id="first-name"
              label="Prénom"
              type="text"
              placeholder="Votre prénom"
              value={firstName}
              onChange={(e) => setFirstName(e.target.value)}
              autoComplete="given-name"
              maxLength={100}
            />
            <Input
              id="email"
              label="Email"
              type="email"
              placeholder="vous@exemple.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              autoComplete="email"
              maxLength={254}
              required
            />
            <Input
              id="password"
              label="Mot de passe"
              type="password"
              placeholder="Min. 8 car., majuscule, chiffre, spécial"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              autoComplete="new-password"
              maxLength={128}
              required
            />
            <Input
              id="confirm-password"
              label="Confirmer le mot de passe"
              type="password"
              placeholder="Retapez votre mot de passe"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              autoComplete="new-password"
              maxLength={128}
              required
            />

            <Button
              type="submit"
              isLoading={isSubmitting}
              className="mt-2 w-full"
            >
              Créer mon compte
            </Button>
          </form>

          <p className="mt-6 text-center text-sm text-foreground-secondary">
            Déjà un compte ?{" "}
            <Link
              href="/login"
              className="font-medium text-primary hover:underline"
            >
              Se connecter
            </Link>
          </p>
        </div>
      </div>
    </GuestGuard>
  );
}
