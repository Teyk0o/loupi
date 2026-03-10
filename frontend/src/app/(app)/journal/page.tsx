"use client";

/**
 * Journal page — Daily meal timeline and quick wellness overview.
 * Displays today's meals and allows navigation between dates.
 */

import { useState, useEffect, useCallback, useMemo } from "react";
import { Plus, ChevronLeft, ChevronRight } from "lucide-react";
import Link from "next/link";
import { useAuthContext } from "@/components/providers/AuthProvider";
import { Card } from "@/components/ui/Card";
import { useCustomOptions } from "@/hooks/useCustomOptions";
import { apiFetch } from "@/lib/api";

/** Meal shape returned by the API. */
interface Meal {
  id: string;
  category: string;
  description: string;
  meal_time: string;
  created_at: string;
}

/** Format a Date to YYYY-MM-DD. */
function toDateString(date: Date): string {
  return date.toISOString().split("T")[0];
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

export default function JournalPage() {
  const { user } = useAuthContext();
  const { options: mealCategories } = useCustomOptions("meal_category");
  const [selectedDate, setSelectedDate] = useState(() => new Date());
  const [meals, setMeals] = useState<Meal[]>([]);
  const [isLoading, setIsLoading] = useState(true);

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

  const isToday = toDateString(selectedDate) === toDateString(new Date());

  return (
    <div>
      {/* Header + Date selector */}
      <div className="mb-6 md:flex md:items-center md:justify-between">
        <div className="mb-4 md:mb-0">
          <h1 className="font-heading text-xl font-semibold">
            Bonjour{user?.first_name ? ` ${user.first_name}` : ""} 👋
          </h1>
          <p className="text-sm text-foreground-secondary">
            Votre journal alimentaire
          </p>
        </div>

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
      </div>

      {/* Meals list */}
      {isLoading ? (
        <div className="flex justify-center py-12">
          <span className="h-6 w-6 animate-spin rounded-full border-2 border-primary border-t-transparent" />
        </div>
      ) : meals.length === 0 ? (
        <div className="flex flex-col items-center gap-3 py-12 text-center">
          <p className="text-foreground-secondary">
            Aucun repas enregistré pour cette journée
          </p>
          <Link
            href="/ajouter"
            className="inline-flex items-center gap-1.5 font-heading text-sm font-medium text-primary hover:underline"
          >
            <Plus size={16} />
            Ajouter un repas
          </Link>
        </div>
      ) : (
        <div className="flex flex-col gap-3 md:grid md:grid-cols-2 lg:grid-cols-3">
          {meals.map((meal) => (
            <Link key={meal.id} href={`/journal/${meal.id}`}>
              <Card className="flex items-center gap-3 transition-colors hover:border-primary/30">
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
              </Card>
            </Link>
          ))}
        </div>
      )}

    </div>
  );
}
