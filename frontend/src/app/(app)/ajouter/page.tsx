"use client";

/**
 * Add meal page.
 * Form to create a new meal entry with category, description, and meal time.
 * Categories are loaded dynamically from user's custom options.
 * Desktop: two-column layout. Mobile: single column.
 */

import { useState, type FormEvent } from "react";
import { useRouter } from "next/navigation";
import { ArrowLeft } from "lucide-react";
import Link from "next/link";
import { Button } from "@/components/ui/Button";
import { Input } from "@/components/ui/Input";
import { useCustomOptions } from "@/hooks/useCustomOptions";
import { apiFetch, type ApiError } from "@/lib/api";

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
  const { options: categories, isLoading: categoriesLoading } = useCustomOptions("meal_category");

  const [category, setCategory] = useState("");
  const [description, setDescription] = useState("");
  const [date, setDate] = useState(todayString);
  const [time, setTime] = useState(nowTimeString);
  const [error, setError] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  // Auto-select first category once loaded.
  const selectedCategory = category || (categories.length > 0 ? categories[0].value : "");

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError("");
    setIsSubmitting(true);

    try {
      const mealTime = `${date}T${time}:00Z`;

      await apiFetch("/v1/meals", {
        method: "POST",
        body: JSON.stringify({
          category: selectedCategory,
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

      <form onSubmit={handleSubmit} className="flex flex-col gap-5 md:grid md:grid-cols-2 md:gap-6">
        {/* Left column: Category selector */}
        <div>
          <p className="mb-2 font-heading text-sm font-medium text-foreground">
            Contexte du repas
          </p>
          {categoriesLoading ? (
            <div className="flex justify-center py-8">
              <span className="h-5 w-5 animate-spin rounded-full border-2 border-primary border-t-transparent" />
            </div>
          ) : (
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
                      selectedCategory === cat.value
                        ? "border-primary bg-primary-light text-foreground"
                        : "border-border bg-surface text-foreground-secondary hover:border-foreground-secondary"
                    }
                  `}
                >
                  {cat.emoji ? <span className="text-lg">{cat.emoji}</span> : null}
                  <span>{cat.label}</span>
                </button>
              ))}
            </div>
          )}
        </div>

        {/* Right column: Description + Date/Time */}
        <div className="flex flex-col gap-5">
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
              rows={4}
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
        </div>

        {/* Submit button full width */}
        <Button type="submit" isLoading={isSubmitting} className="mt-2 w-full md:col-span-2">
          Enregistrer
        </Button>
      </form>
    </div>
  );
}
