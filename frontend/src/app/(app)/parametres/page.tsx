"use client";

/**
 * Settings page.
 * Provides account info, logout, account deletion, and custom option management.
 * Desktop: two-column layout. Mobile: single column.
 */

import { useState } from "react";
import { LogOut, Trash2, User } from "lucide-react";
import { useAuthContext } from "@/components/providers/AuthProvider";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { Logo } from "@/components/ui/Logo";
import { OptionListEditor } from "@/components/settings/OptionListEditor";
import type { ApiError } from "@/lib/api";

export default function SettingsPage() {
  const { user, logout, deleteAccount } = useAuthContext();
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);
  const [error, setError] = useState("");

  async function handleDelete() {
    setIsDeleting(true);
    setError("");

    try {
      await deleteAccount();
    } catch (err) {
      const apiErr = err as ApiError;
      setError(apiErr.message || "Une erreur est survenue");
      setIsDeleting(false);
    }
  }

  return (
    <div>
      <div className="mb-6">
        <h1 className="font-heading text-xl font-semibold">Paramètres</h1>
      </div>

      {/* Account & Actions */}
      <div className="mb-8 flex flex-col gap-6 md:grid md:grid-cols-2 md:gap-8">
        {/* Left: Account info + About */}
        <div className="flex flex-col gap-4">
          <Card padding="lg" className="flex items-center gap-4">
            <div className="flex h-14 w-14 items-center justify-center rounded-full bg-primary-light">
              <User size={24} className="text-primary" />
            </div>
            <div>
              <p className="font-heading text-sm font-semibold">
                {user?.first_name || "Utilisateur"}
              </p>
              <p className="text-sm text-foreground-secondary">{user?.email}</p>
              <p className="text-xs text-foreground-secondary">
                Membre depuis{" "}
                {user?.created_at
                  ? new Date(user.created_at).toLocaleDateString("fr-FR", {
                      month: "long",
                      year: "numeric",
                    })
                  : ""}
              </p>
            </div>
          </Card>

          <Card padding="lg" className="flex flex-col items-center gap-2">
            <Logo width={80} height={26} />
            <p className="text-xs text-foreground-secondary">
              Application gratuite et open-source
            </p>
          </Card>
        </div>

        {/* Right: Actions */}
        <div className="flex flex-col gap-3">
          <h2 className="font-heading text-sm font-semibold text-foreground-secondary">
            Actions
          </h2>

          <Button variant="outline" className="w-full" onClick={logout}>
            <LogOut size={16} />
            Se déconnecter
          </Button>

          {!showDeleteConfirm ? (
            <Button
              variant="ghost"
              className="w-full text-danger"
              onClick={() => setShowDeleteConfirm(true)}
            >
              <Trash2 size={16} />
              Supprimer mon compte
            </Button>
          ) : (
            <Card padding="md" className="border-danger/30">
              <p className="mb-3 text-sm text-foreground">
                Cette action est irréversible. Toutes vos données seront
                définitivement supprimées.
              </p>
              {error ? (
                <p className="mb-2 text-xs text-danger">{error}</p>
              ) : null}
              <div className="flex gap-2">
                <Button
                  variant="outline"
                  className="flex-1"
                  onClick={() => setShowDeleteConfirm(false)}
                >
                  Annuler
                </Button>
                <Button
                  variant="danger"
                  className="flex-1"
                  isLoading={isDeleting}
                  onClick={handleDelete}
                >
                  Confirmer
                </Button>
              </div>
            </Card>
          )}
        </div>
      </div>

      {/* Customizable options */}
      <div className="mb-6">
        <h2 className="mb-1 font-heading text-lg font-semibold">Personnalisation</h2>
        <p className="text-sm text-foreground-secondary">
          Gérez vos listes personnalisées de symptômes, catégories de repas et types de sport.
        </p>
      </div>

      <div className="flex flex-col gap-6 md:grid md:grid-cols-2 lg:grid-cols-3 md:gap-6">
        <OptionListEditor
          title="Types de symptômes"
          category="symptom_type"
        />
        <OptionListEditor
          title="Catégories de repas"
          category="meal_category"
          showEmoji
        />
        <OptionListEditor
          title="Types de sport"
          category="sport_type"
          showEmoji
        />
      </div>
    </div>
  );
}
