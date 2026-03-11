"use client";

/**
 * Journal page — Daily meal timeline.
 * Displays today's meals and allows navigation between dates.
 * Includes an inline form to add a new meal entry.
 */

import { useState, useEffect, useCallback, useMemo, type FormEvent } from "react";
import { Plus, ChevronLeft, ChevronRight, Trash2 } from "lucide-react";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { useAuthContext } from "@/components/providers/AuthProvider";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { Input } from "@/components/ui/Input";
import { useCustomOptions } from "@/hooks/useCustomOptions";
import { apiFetch, type ApiError } from "@/lib/api";

/** Meal shape returned by the API. */
interface Meal {
  id: string;
  category: string;
  description: string;
  meal_time: string;
  created_at: string;
}

/** Format a Date to YYYY-MM-DD (local timezone). */
function toDateString(date: Date): string {
  const y = date.getFullYear();
  const m = String(date.getMonth() + 1).padStart(2, "0");
  const d = String(date.getDate()).padStart(2, "0");
  return `${y}-${m}-${d}`;
}

/** Format a date for display (e.g. "Lun. 10 mars 2026"). */
function formatDisplayDate(date: Date): string {
  return date.toLocaleDateString("fr-FR", {
    weekday: "short",
    day: "numeric",
    month: "long",
  });
}

/** Extract HH:MM from an ISO datetime string. */
function extractTime(isoString: string): string {
  const date = new Date(isoString);
  return date.toLocaleTimeString("fr-FR", { hour: "2-digit", minute: "2-digit" });
}

/** Format current time as HH:MM. */
function nowTimeString(): string {
  return new Date().toTimeString().slice(0, 5);
}

/** Get local timezone offset as +HH:MM or -HH:MM string. */
function localTimezoneOffset(): string {
  const offset = new Date().getTimezoneOffset();
  const sign = offset <= 0 ? "+" : "-";
  const absOffset = Math.abs(offset);
  const hours = String(Math.floor(absOffset / 60)).padStart(2, "0");
  const minutes = String(absOffset % 60).padStart(2, "0");
  return `${sign}${hours}:${minutes}`;
}

export default function JournalPage() {
  const searchParams = useSearchParams();
  const { user } = useAuthContext();
  const { options: mealCategories } = useCustomOptions("meal_category");
  const [selectedDate, setSelectedDate] = useState(() => new Date());
  const [meals, setMeals] = useState<Meal[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  // Add meal form state
  const [showForm, setShowForm] = useState(false);
  const [category, setCategory] = useState("");
  const [description, setDescription] = useState("");
  const [time, setTime] = useState(nowTimeString);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState("");

  // Open form if redirected from /ajouter with ?add=meal
  useEffect(() => {
    if (searchParams.get("add") === "meal") {
      setShowForm(true);
    }
  }, [searchParams]);

  const selectedCategory = category || (mealCategories.length > 0 ? mealCategories[0].value : "");

  /** Build label/emoji lookup maps from user's custom options. */
  const categoryLabels = useMemo(() => {
    const map: Record<string, string> = {};
    for (const cat of mealCategories) map[cat.value] = cat.label;
    return map;
  }, [mealCategories]);

  const categoryEmojis = useMemo(() => {
    const map: Record<string, string> = {};
    for (const cat of mealCategories) {
      if (cat.emoji) map[cat.value] = cat.emoji;
    }
    return map;
  }, [mealCategories]);

  const fetchMeals = useCallback(async (date: Date) => {
    setIsLoading(true);
    try {
      const data = await apiFetch<Meal[]>(
        `/v1/meals?date=${toDateString(date)}`,
      );
      setMeals(data || []);
    } catch {
      setMeals([]);
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchMeals(selectedDate);
  }, [selectedDate, fetchMeals, user?.id]);

  function goToPreviousDay() {
    setSelectedDate((prev) => {
      const d = new Date(prev);
      d.setDate(d.getDate() - 1);
      return d;
    });
  }

  function goToNextDay() {
    setSelectedDate((prev) => {
      const d = new Date(prev);
      d.setDate(d.getDate() + 1);
      return d;
    });
  }

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError("");
    setIsSubmitting(true);

    try {
      const mealTime = `${toDateString(selectedDate)}T${time}:00${localTimezoneOffset()}`;

      await apiFetch("/v1/meals", {
        method: "POST",
        body: JSON.stringify({
          category: selectedCategory,
          description,
          meal_time: mealTime,
        }),
      });
      setShowForm(false);
      setCategory("");
      setDescription("");
      setTime(nowTimeString());
      await fetchMeals(selectedDate);
    } catch (err) {
      const apiErr = err as ApiError;
      setError(apiErr.message || "Une erreur est survenue");
    } finally {
      setIsSubmitting(false);
    }
  }

  async function handleDelete(id: string) {
    if (!confirm("Supprimer ce repas ?")) return;
    try {
      await apiFetch(`/v1/meals/${id}`, { method: "DELETE" });
      await fetchMeals(selectedDate);
    } catch {
      setError("Erreur lors de la suppression");
    }
  }

  const isToday = toDateString(selectedDate) === toDateString(new Date());

  return (
    <div>
      {/* Header + Date selector */}
      <div className="mb-6 md:flex md:items-center md:justify-between">
        <div className="mb-4 md:mb-0">
          <h1 className="font-heading text-xl font-semibold">
            Bonjour{user?.first_name ? ` ${user.first_name}` : ""}
          </h1>
          <p className="text-sm text-foreground-secondary">
            Votre journal alimentaire
          </p>
        </div>

        <div className="flex items-center gap-3">
          {/* Date selector */}
          <div className="mx-auto flex w-fit items-center gap-2 rounded-[--radius-md] border border-border bg-surface px-2 py-1 md:mx-0 md:px-3">
            <button
              onClick={goToPreviousDay}
              className="rounded-[--radius-sm] p-1.5 text-foreground-secondary hover:text-foreground"
            >
              <ChevronLeft size={18} />
            </button>
            <button
              onClick={() => setSelectedDate(new Date())}
              className={`min-w-[140px] text-center font-heading text-sm font-medium ${isToday ? "text-primary" : "text-foreground"}`}
            >
              {isToday ? "Aujourd'hui" : formatDisplayDate(selectedDate)}
            </button>
            <button
              onClick={goToNextDay}
              className="rounded-[--radius-sm] p-1.5 text-foreground-secondary hover:text-foreground"
            >
              <ChevronRight size={18} />
            </button>
          </div>

          {/* Add button */}
          <button
            onClick={() => setShowForm(!showForm)}
            className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-primary text-white"
          >
            <Plus size={20} />
          </button>
        </div>
      </div>

      {error ? (
        <div className="mb-4 rounded-[--radius-md] border border-danger/30 bg-danger/10 px-4 py-3 text-sm text-danger">
          {error}
        </div>
      ) : null}

      {/* Inline add meal form */}
      {showForm ? (
        <Card padding="lg" className="mb-4">
          <form onSubmit={handleSubmit} className="flex flex-col gap-4">
            {/* Category selector */}
            <div>
              <p className="mb-2 font-heading text-sm font-medium text-foreground">
                Contexte du repas
              </p>
              <div className="grid grid-cols-3 gap-2">
                {mealCategories.map((cat) => (
                  <button
                    key={cat.value}
                    type="button"
                    onClick={() => setCategory(cat.value)}
                    className={`
                      flex flex-col items-center gap-1 rounded-[--radius-md] border px-2 py-2.5
                      font-heading text-sm transition-all
                      ${
                        selectedCategory === cat.value
                          ? "border-primary bg-primary-light text-foreground"
                          : "border-border bg-surface text-foreground-secondary hover:border-foreground-secondary"
                      }
                    `}
                  >
                    {cat.emoji ? <span className="text-xl">{cat.emoji}</span> : null}
                    <span>{cat.label}</span>
                  </button>
                ))}
              </div>
            </div>

            {/* Description */}
            <div className="flex flex-col gap-1.5">
              <label
                htmlFor="description"
                className="font-heading text-sm font-medium text-foreground"
              >
                Description
              </label>
              <textarea
                id="description"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                placeholder="Decrivez votre repas..."
                rows={3}
                required
                maxLength={1000}
                className="
                  w-full resize-none rounded-[--radius-md] border border-border
                  bg-surface px-4 py-3
                  font-body text-sm text-foreground
                  placeholder:text-foreground-secondary
                  outline-none transition-all duration-200
                  focus:border-primary focus:ring-2 focus:ring-primary-light
                "
              />
            </div>

            {/* Time */}
            <Input
              id="meal-time"
              label="Heure"
              type="time"
              value={time}
              onChange={(e) => setTime(e.target.value)}
              required
            />

            <div className="flex gap-2">
              <Button
                variant="outline"
                className="flex-1"
                type="button"
                onClick={() => setShowForm(false)}
              >
                Annuler
              </Button>
              <Button
                className="flex-1"
                type="submit"
                isLoading={isSubmitting}
              >
                Enregistrer
              </Button>
            </div>
          </form>
        </Card>
      ) : null}

      {/* Meals list */}
      {isLoading ? (
        <div className="flex justify-center py-12">
          <span className="h-6 w-6 animate-spin rounded-full border-2 border-primary border-t-transparent" />
        </div>
      ) : meals.length === 0 ? (
        <div className="flex flex-col items-center gap-3 py-12 text-center">
          <p className="text-foreground-secondary">
            Aucun repas enregistre pour cette journee
          </p>
          <button
            onClick={() => setShowForm(true)}
            className="inline-flex items-center gap-1.5 font-heading text-sm font-medium text-primary hover:underline"
          >
            <Plus size={16} />
            Ajouter un repas
          </button>
        </div>
      ) : (
        <div className="flex flex-col gap-3 md:grid md:grid-cols-2 lg:grid-cols-3">
          {meals.map((meal) => (
            <Card key={meal.id} className="flex items-center gap-3 transition-colors hover:border-primary/30">
              <Link href={`/journal/${meal.id}`} className="flex flex-1 items-center gap-3">
                <span className="text-xl">
                  {categoryEmojis[meal.category] || "🍴"}
                </span>
                <div className="flex-1">
                  <p className="line-clamp-1 font-heading text-sm font-medium">
                    {meal.description}
                  </p>
                  <p className="text-xs text-foreground-secondary">
                    {categoryLabels[meal.category] || meal.category}
                  </p>
                </div>
                <span className="text-xs text-foreground-secondary">
                  {extractTime(meal.meal_time)}
                </span>
              </Link>
              <button
                onClick={() => handleDelete(meal.id)}
                className="p-1.5 text-foreground-secondary hover:text-danger"
              >
                <Trash2 size={16} />
              </button>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
