"use client";

/**
 * Add meal page.
 * Form to create a new meal entry with category, description, and meal time.
 */

import { useState, type FormEvent } from "react";
import { useRouter } from "next/navigation";
import { ArrowLeft } from "lucide-react";
import Link from "next/link";
import { Button } from "@/components/ui/Button";
import { Input } from "@/components/ui/Input";
import { apiFetch, type ApiError } from "@/lib/api";

type MealCategory =
  | "homemade"
  | "restaurant"
  | "takeout"
  | "snack"
  | "fast_food"
  | "cafeteria"
  | "family"
  | "friends"
  | "other";

interface CategoryOption {
  value: MealCategory;
  label: string;
  emoji: string;
}

const categories: CategoryOption[] = [
  { value: "homemade", label: "Fait maison", emoji: "🏠" },
  { value: "restaurant", label: "Restaurant", emoji: "🍽️" },
  { value: "takeout", label: "À emporter", emoji: "🥡" },
  { value: "snack", label: "Collation", emoji: "🍪" },
  { value: "fast_food", label: "Fast-food", emoji: "🍔" },
  { value: "cafeteria", label: "Cantine", emoji: "🏫" },
  { value: "family", label: "En famille", emoji: "👨‍👩‍👧" },
  { value: "friends", label: "Entre amis", emoji: "👫" },
  { value: "other", label: "Autre", emoji: "🍴" },
];

/** Format current date as YYYY-MM-DD. */
function todayString(): string {
  return new Date().toISOString().split("T")[0];
}

/** Format current time as HH:MM. */
function nowTimeString(): string {
  return new Date().toTimeString().slice(0, 5);
}

export default function AddMealPage() {
  const router = useRouter();
  const [category, setCategory] = useState<MealCategory>("homemade");
  const [description, setDescription] = useState("");
  const [date, setDate] = useState(todayString);
  const [time, setTime] = useState(nowTimeString);
  const [error, setError] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError("");
    setIsSubmitting(true);

    try {
      // Combine date + time into ISO datetime string for the API
      const mealTime = `${date}T${time}:00Z`;

      await apiFetch("/v1/meals", {
        method: "POST",
        body: JSON.stringify({
          category,
          description,
          meal_time: mealTime,
        }),
      });
      router.push("/journal");
    } catch (err) {
      const apiErr = err as ApiError;
      setError(apiErr.message || "Une erreur est survenue");
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <div>
      {/* Header */}
      <div className="mb-6 flex items-center gap-3">
        <Link
          href="/journal"
          className="rounded-[--radius-sm] p-1.5 text-foreground-secondary hover:text-foreground"
        >
          <ArrowLeft size={20} />
        </Link>
        <h1 className="font-heading text-xl font-semibold">Ajouter un repas</h1>
      </div>

      {error ? (
        <div className="mb-4 rounded-[--radius-md] border border-danger/30 bg-danger/10 px-4 py-3 text-sm text-danger">
          {error}
        </div>
      ) : null}

      <form onSubmit={handleSubmit} className="flex flex-col gap-5">
        {/* Category selector */}
        <div>
          <p className="mb-2 font-heading text-sm font-medium text-foreground">
            Contexte du repas
          </p>
          <div className="grid grid-cols-3 gap-2">
            {categories.map((cat) => (
              <button
                key={cat.value}
                type="button"
                onClick={() => setCategory(cat.value)}
                className={`
                  flex flex-col items-center gap-1 rounded-[--radius-md] border px-2 py-2.5
                  font-heading text-xs transition-all
                  ${
                    category === cat.value
                      ? "border-primary bg-primary-light text-foreground"
                      : "border-border bg-surface text-foreground-secondary hover:border-foreground-secondary"
                  }
                `}
              >
                <span className="text-lg">{cat.emoji}</span>
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
            placeholder="Décrivez votre repas..."
            rows={3}
            required
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

        {/* Date & Time */}
        <div className="grid grid-cols-2 gap-3">
          <Input
            id="date"
            label="Date"
            type="date"
            value={date}
            onChange={(e) => setDate(e.target.value)}
            required
          />
          <Input
            id="time"
            label="Heure"
            type="time"
            value={time}
            onChange={(e) => setTime(e.target.value)}
            required
          />
        </div>

        <Button type="submit" isLoading={isSubmitting} className="mt-2 w-full">
          Enregistrer
        </Button>
      </form>
    </div>
  );
}
