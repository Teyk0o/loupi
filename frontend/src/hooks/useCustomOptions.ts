/**
 * Hook for fetching and managing user-configurable options.
 * Provides CRUD operations for symptom types, meal categories, and sport types.
 */

import { useState, useEffect, useCallback } from "react";
import { apiFetch, isValidUUID, type ApiError } from "@/lib/api";

export type OptionCategory = "symptom_type" | "meal_category" | "sport_type";

const VALID_CATEGORIES: OptionCategory[] = ["symptom_type", "meal_category", "sport_type"];

export interface CustomOption {
  id: string;
  category: OptionCategory;
  value: string;
  label: string;
  emoji?: string;
  sort_order: number;
  created_at: string;
}

interface UseCustomOptionsReturn {
  options: CustomOption[];
  isLoading: boolean;
  error: string;
  refresh: () => Promise<void>;
  addOption: (value: string, label: string, emoji?: string) => Promise<void>;
  updateOption: (id: string, label: string, emoji?: string) => Promise<void>;
  deleteOption: (id: string) => Promise<void>;
}

export function useCustomOptions(category: OptionCategory): UseCustomOptionsReturn {
  const [options, setOptions] = useState<CustomOption[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState("");

  const refresh = useCallback(async () => {
    if (!VALID_CATEGORIES.includes(category)) {
      setError("Catégorie invalide");
      setIsLoading(false);
      return;
    }
    setIsLoading(true);
    setError("");
    try {
      const data = await apiFetch<CustomOption[]>(`/v1/options/${category}`);
      setOptions(data || []);
    } catch {
      setError("Impossible de charger les options");
    } finally {
      setIsLoading(false);
    }
  }, [category]);

  useEffect(() => {
    refresh();
  }, [refresh]);

  const addOption = useCallback(
    async (value: string, label: string, emoji?: string) => {
      setError("");
      try {
        await apiFetch(`/v1/options/${category}`, {
          method: "POST",
          body: JSON.stringify({ value, label, emoji: emoji || undefined }),
        });
        await refresh();
      } catch (err) {
        const apiErr = err as ApiError;
        throw new Error(apiErr.message || "Erreur lors de l'ajout");
      }
    },
    [category, refresh],
  );

  const updateOption = useCallback(
    async (id: string, label: string, emoji?: string) => {
      if (!isValidUUID(id)) throw new Error("ID invalide");
      setError("");
      try {
        await apiFetch(`/v1/options/item/${id}`, {
          method: "PUT",
          body: JSON.stringify({ label, emoji: emoji || undefined }),
        });
        await refresh();
      } catch (err) {
        const apiErr = err as ApiError;
        throw new Error(apiErr.message || "Erreur lors de la mise à jour");
      }
    },
    [refresh],
  );

  const deleteOption = useCallback(
    async (id: string) => {
      if (!isValidUUID(id)) throw new Error("ID invalide");
      setError("");
      try {
        await apiFetch(`/v1/options/item/${id}`, { method: "DELETE" });
        await refresh();
      } catch (err) {
        const apiErr = err as ApiError;
        throw new Error(apiErr.message || "Erreur lors de la suppression");
      }
    },
    [refresh],
  );

  return { options, isLoading, error, refresh, addOption, updateOption, deleteOption };
}
