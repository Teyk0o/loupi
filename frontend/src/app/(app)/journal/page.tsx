"use client";

/**
 * Journal page — Daily meal timeline and quick wellness overview.
 * Displays today's meals and allows navigation between dates.
 */

import { useState, useEffect, useCallback } from "react";
import { Plus, ChevronLeft, ChevronRight } from "lucide-react";
import Link from "next/link";
import { useAuthContext } from "@/components/providers/AuthProvider";
import { Card } from "@/components/ui/Card";
import { apiFetch } from "@/lib/api";

/** Meal shape returned by the API. */
interface Meal {
  id: string;
  category: string;
  description: string;
  meal_time: string;
  created_at: string;
}

/** Category display config (matches API oneof values). */
const categoryLabels: Record<string, string> = {
  homemade: "Fait maison",
  restaurant: "Restaurant",
  takeout: "À emporter",
  snack: "Collation",
  fast_food: "Fast-food",
  cafeteria: "Cantine",
  family: "En famille",
  friends: "Entre amis",
  other: "Autre",
};

const categoryEmojis: Record<string, string> = {
  homemade: "🏠",
  restaurant: "🍽️",
  takeout: "🥡",
  snack: "🍪",
  fast_food: "🍔",
  cafeteria: "🏫",
  family: "👨‍👩‍👧",
  friends: "👫",
  other: "🍴",
};

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
  const [selectedDate, setSelectedDate] = useState(() => new Date());
  const [meals, setMeals] = useState<Meal[]>([]);
  const [isLoading, setIsLoading] = useState(true);

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
      {/* Header */}
      <div className="mb-6">
        <h1 className="font-heading text-xl font-semibold">
          Bonjour{user?.first_name ? ` ${user.first_name}` : ""} 👋
        </h1>
        <p className="text-sm text-foreground-secondary">
          Votre journal alimentaire
        </p>
      </div>

      {/* Date selector */}
      <div className="mb-6 flex items-center justify-between">
        <button
          onClick={goToPreviousDay}
          className="rounded-[--radius-sm] p-2 text-foreground-secondary hover:bg-surface hover:text-foreground"
        >
          <ChevronLeft size={20} />
        </button>
        <button
          onClick={() => setSelectedDate(new Date())}
          className={`font-heading text-sm font-medium ${isToday ? "text-primary" : "text-foreground"}`}
        >
          {isToday ? "Aujourd'hui" : formatDisplayDate(selectedDate)}
        </button>
        <button
          onClick={goToNextDay}
          className="rounded-[--radius-sm] p-2 text-foreground-secondary hover:bg-surface hover:text-foreground"
        >
          <ChevronRight size={20} />
        </button>
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
        <div className="flex flex-col gap-3">
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

      {/* Floating add button */}
      <Link
        href="/ajouter"
        className="fixed bottom-20 right-4 z-40 flex h-14 w-14 items-center justify-center rounded-full bg-primary text-white shadow-lg transition-transform hover:scale-105 active:scale-95"
      >
        <Plus size={28} />
      </Link>
    </div>
  );
}
