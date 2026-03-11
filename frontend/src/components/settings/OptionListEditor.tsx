"use client";

/**
 * Editable list component for user-configurable options.
 * Allows adding, editing, and deleting options within a category.
 */

import { useState } from "react";
import { Pencil, Trash2, Plus, Check, X } from "lucide-react";
import { Card } from "@/components/ui/Card";
import { Button } from "@/components/ui/Button";
import type { CustomOption, OptionCategory } from "@/hooks/useCustomOptions";
import { useCustomOptions } from "@/hooks/useCustomOptions";

interface OptionListEditorProps {
  title: string;
  category: OptionCategory;
  showEmoji?: boolean;
}

export function OptionListEditor({ title, category, showEmoji = false }: OptionListEditorProps) {
  const { options, isLoading, error, addOption, updateOption, deleteOption } =
    useCustomOptions(category);

  const [showAddForm, setShowAddForm] = useState(false);
  const [newLabel, setNewLabel] = useState("");
  const [newEmoji, setNewEmoji] = useState("");
  const [editingId, setEditingId] = useState<string | null>(null);
  const [editLabel, setEditLabel] = useState("");
  const [editEmoji, setEditEmoji] = useState("");
  const [actionError, setActionError] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  async function handleAdd() {
    if (!newLabel.trim()) return;
    setIsSubmitting(true);
    setActionError("");
    try {
      const value = newLabel
        .trim()
        .toLowerCase()
        .replace(/[^a-z0-9àâäéèêëïîôùûüÿçœæ]+/gi, "_")
        .replace(/^_|_$/g, "");
      await addOption(value, newLabel.trim(), showEmoji && newEmoji ? newEmoji : undefined);
      setNewLabel("");
      setNewEmoji("");
      setShowAddForm(false);
    } catch (err) {
      setActionError((err as Error).message);
    } finally {
      setIsSubmitting(false);
    }
  }

  function startEdit(option: CustomOption) {
    setEditingId(option.id);
    setEditLabel(option.label);
    setEditEmoji(option.emoji || "");
    setActionError("");
  }

  async function handleUpdate() {
    if (!editingId || !editLabel.trim()) return;
    setIsSubmitting(true);
    setActionError("");
    try {
      await updateOption(editingId, editLabel.trim(), showEmoji && editEmoji ? editEmoji : undefined);
      setEditingId(null);
    } catch (err) {
      setActionError((err as Error).message);
    } finally {
      setIsSubmitting(false);
    }
  }

  async function handleDelete(id: string) {
    setActionError("");
    try {
      await deleteOption(id);
    } catch (err) {
      setActionError((err as Error).message);
    }
  }

  if (isLoading) {
    return (
      <Card padding="lg">
        <h3 className="mb-4 font-heading text-base font-semibold">{title}</h3>
        <div className="flex justify-center py-4">
          <span className="h-5 w-5 animate-spin rounded-full border-2 border-primary border-t-transparent" />
        </div>
      </Card>
    );
  }

  return (
    <Card padding="lg">
      <div className="mb-4 flex items-center justify-between">
        <h3 className="font-heading text-base font-semibold">{title}</h3>
        <button
          onClick={() => { setShowAddForm(!showAddForm); setActionError(""); }}
          className="inline-flex items-center gap-1.5 text-sm font-medium text-primary hover:underline"
        >
          <Plus size={18} />
          Ajouter
        </button>
      </div>

      {error || actionError ? (
        <p className="mb-3 text-sm text-danger">{error || actionError}</p>
      ) : null}

      {/* Add form */}
      {showAddForm ? (
        <div className="mb-4 flex flex-col gap-3 rounded-[--radius-md] border border-border p-4">
          <div className="flex gap-2">
            {showEmoji ? (
              <input
                value={newEmoji}
                onChange={(e) => setNewEmoji(e.target.value)}
                placeholder="🔹"
                className="w-14 rounded-[--radius-sm] border border-border bg-surface px-2 py-2.5 text-center text-base outline-none focus:border-primary"
              />
            ) : null}
            <input
              value={newLabel}
              onChange={(e) => setNewLabel(e.target.value)}
              placeholder="Nom de l'option..."
              className="flex-1 rounded-[--radius-sm] border border-border bg-surface px-3 py-2.5 text-sm outline-none focus:border-primary"
              onKeyDown={(e) => e.key === "Enter" && handleAdd()}
              autoFocus
            />
          </div>
          <div className="flex gap-2">
            <Button
              variant="outline"
              className="flex-1"
              onClick={() => { setShowAddForm(false); setNewLabel(""); setNewEmoji(""); }}
            >
              Annuler
            </Button>
            <Button
              className="flex-1"
              onClick={handleAdd}
              isLoading={isSubmitting}
              disabled={!newLabel.trim()}
            >
              Ajouter
            </Button>
          </div>
        </div>
      ) : null}

      {/* Options list */}
      <div className="flex flex-col gap-1">
        {options.map((option) => (
          <div
            key={option.id}
            className="flex items-center gap-3 rounded-[--radius-sm] px-3 py-2.5 hover:bg-background"
          >
            {editingId === option.id ? (
              /* Edit mode */
              <>
                {showEmoji ? (
                  <input
                    value={editEmoji}
                    onChange={(e) => setEditEmoji(e.target.value)}
                    className="w-12 rounded-[--radius-sm] border border-border bg-surface px-2 py-2 text-center text-base outline-none focus:border-primary"
                  />
                ) : null}
                <input
                  value={editLabel}
                  onChange={(e) => setEditLabel(e.target.value)}
                  className="flex-1 rounded-[--radius-sm] border border-border bg-surface px-3 py-2 text-sm outline-none focus:border-primary"
                  onKeyDown={(e) => e.key === "Enter" && handleUpdate()}
                  autoFocus
                />
                <button
                  onClick={handleUpdate}
                  disabled={isSubmitting}
                  className="rounded-[--radius-sm] p-2 text-primary hover:bg-primary-light"
                >
                  <Check size={18} />
                </button>
                <button
                  onClick={() => setEditingId(null)}
                  className="rounded-[--radius-sm] p-2 text-foreground-secondary hover:text-foreground"
                >
                  <X size={18} />
                </button>
              </>
            ) : (
              /* Display mode */
              <>
                {option.emoji ? (
                  <span className="w-8 text-center text-lg">{option.emoji}</span>
                ) : null}
                <span className="flex-1 text-sm text-foreground">{option.label}</span>
                <span className="text-xs text-foreground-secondary">{option.value}</span>
                <button
                  onClick={() => startEdit(option)}
                  className="rounded-[--radius-sm] p-2 text-foreground-secondary hover:text-foreground"
                >
                  <Pencil size={16} />
                </button>
                <button
                  onClick={() => handleDelete(option.id)}
                  className="rounded-[--radius-sm] p-2 text-foreground-secondary hover:text-danger"
                >
                  <Trash2 size={16} />
                </button>
              </>
            )}
          </div>
        ))}
        {options.length === 0 ? (
          <p className="py-3 text-center text-sm text-foreground-secondary">
            Aucune option configurée
          </p>
        ) : null}
      </div>
    </Card>
  );
}
