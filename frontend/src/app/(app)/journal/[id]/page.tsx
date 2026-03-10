"use client";

/**
 * Meal detail page.
 * Shows a single meal with its description and associated symptom check-ins.
 * Allows adding new check-ins and navigating back to the journal.
 */

import { useState, useEffect, useCallback } from "react";
import { useParams, useRouter } from "next/navigation";
import {
  ArrowLeft,
  Clock,
  Plus,
  AlertTriangle,
  Trash2,
} from "lucide-react";
import Link from "next/link";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { apiFetch, type ApiError } from "@/lib/api";

interface Meal {
  id: string;
  category: string;
  description: string;
  meal_time: string;
  created_at: string;
}

interface SymptomDetail {
  type: string;
  severity: number;
}

interface SymptomCheckin {
  id: string;
  meal_id: string;
  delay_hours: number;
  symptoms: SymptomDetail[];
  notes?: string;
  created_at: string;
}

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

const delayOptions = [
  { value: 6, label: "Après 6h" },
  { value: 8, label: "Après 8h" },
  { value: 12, label: "Après 12h" },
];

/** Symptom types matching the API oneof validation. */
const symptomLabels: Record<string, string> = {
  diarrhea: "Diarrhée",
  stomach_ache: "Maux de ventre",
  nausea: "Nausée",
  bloating: "Ballonnements",
  heartburn: "Brûlures d'estomac",
  cramps: "Crampes",
  constipation: "Constipation",
  gas: "Gaz",
  reflux: "Reflux",
  fatigue: "Fatigue",
};

/** Severity display with colored dots. */
function SeverityDots({ severity }: { severity: number }) {
  return (
    <div className="flex gap-0.5">
      {[1, 2, 3, 4, 5].map((n) => (
        <span
          key={n}
          className={`h-2 w-2 rounded-full ${
            n <= severity ? "bg-danger" : "bg-border"
          }`}
        />
      ))}
    </div>
  );
}

export default function MealDetailPage() {
  const params = useParams();
  const router = useRouter();
  const mealId = params.id as string;

  const [meal, setMeal] = useState<Meal | null>(null);
  const [checkins, setCheckins] = useState<SymptomCheckin[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isDeleting, setIsDeleting] = useState(false);
  const [error, setError] = useState("");

  // Check-in form state
  const [showCheckinForm, setShowCheckinForm] = useState(false);
  const [checkinDelay, setCheckinDelay] = useState(6);
  const [checkinSymptoms, setCheckinSymptoms] = useState<SymptomDetail[]>([]);
  const [checkinNotes, setCheckinNotes] = useState("");
  const [isSubmittingCheckin, setIsSubmittingCheckin] = useState(false);

  const fetchData = useCallback(async () => {
    setIsLoading(true);
    try {
      const [mealData, checkinData] = await Promise.all([
        apiFetch<Meal>(`/v1/meals/${mealId}`),
        apiFetch<SymptomCheckin[]>(`/v1/meals/${mealId}/check-ins`),
      ]);
      setMeal(mealData);
      setCheckins(checkinData || []);
    } catch {
      setError("Impossible de charger le repas");
    } finally {
      setIsLoading(false);
    }
  }, [mealId]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  async function handleDeleteMeal() {
    if (!confirm("Supprimer ce repas ?")) return;
    setIsDeleting(true);
    try {
      await apiFetch(`/v1/meals/${mealId}`, { method: "DELETE" });
      router.push("/journal");
    } catch (err) {
      const apiErr = err as ApiError;
      setError(apiErr.message || "Erreur lors de la suppression");
      setIsDeleting(false);
    }
  }

  function toggleSymptom(type: string) {
    setCheckinSymptoms((prev) => {
      const existing = prev.find((s) => s.type === type);
      if (existing) return prev.filter((s) => s.type !== type);
      return [...prev, { type, severity: 3 }];
    });
  }

  function updateSymptomSeverity(type: string, severity: number) {
    setCheckinSymptoms((prev) =>
      prev.map((s) => (s.type === type ? { ...s, severity } : s)),
    );
  }

  async function handleSubmitCheckin() {
    if (checkinSymptoms.length === 0) return;
    setIsSubmittingCheckin(true);
    setError("");

    try {
      await apiFetch(`/v1/meals/${mealId}/check-ins`, {
        method: "POST",
        body: JSON.stringify({
          delay_hours: checkinDelay,
          symptoms: checkinSymptoms,
          notes: checkinNotes || undefined,
        }),
      });
      setShowCheckinForm(false);
      setCheckinSymptoms([]);
      setCheckinNotes("");
      await fetchData();
    } catch (err) {
      const apiErr = err as ApiError;
      setError(apiErr.message || "Erreur lors de l'enregistrement");
    } finally {
      setIsSubmittingCheckin(false);
    }
  }

  if (isLoading) {
    return (
      <div className="flex justify-center py-12">
        <span className="h-6 w-6 animate-spin rounded-full border-2 border-primary border-t-transparent" />
      </div>
    );
  }

  if (!meal) {
    return (
      <div className="py-12 text-center">
        <p className="text-foreground-secondary">Repas introuvable</p>
        <Link href="/journal" className="mt-2 text-sm text-primary hover:underline">
          Retour au journal
        </Link>
      </div>
    );
  }

  const mealDate = new Date(meal.meal_time);

  return (
    <div>
      {/* Header */}
      <div className="mb-6 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Link
            href="/journal"
            className="rounded-[--radius-sm] p-1.5 text-foreground-secondary hover:text-foreground"
          >
            <ArrowLeft size={20} />
          </Link>
          <h1 className="font-heading text-xl font-semibold">
            {categoryLabels[meal.category] || meal.category}
          </h1>
        </div>
        <Button
          variant="ghost"
          className="text-danger"
          onClick={handleDeleteMeal}
          isLoading={isDeleting}
        >
          <Trash2 size={16} />
        </Button>
      </div>

      {error ? (
        <div className="mb-4 rounded-[--radius-md] border border-danger/30 bg-danger/10 px-4 py-3 text-sm text-danger">
          {error}
        </div>
      ) : null}

      {/* Meal info */}
      <Card padding="lg" className="mb-4">
        <p className="mb-2 text-sm text-foreground">{meal.description}</p>
        <div className="flex items-center gap-1.5 text-xs text-foreground-secondary">
          <Clock size={12} />
          <span>
            {mealDate.toLocaleDateString("fr-FR")} à{" "}
            {mealDate.toLocaleTimeString("fr-FR", { hour: "2-digit", minute: "2-digit" })}
          </span>
        </div>
      </Card>

      {/* Symptom check-ins */}
      <div className="mb-4 flex items-center justify-between">
        <h2 className="font-heading text-sm font-semibold">
          Symptômes signalés
        </h2>
        <button
          onClick={() => setShowCheckinForm(!showCheckinForm)}
          className="inline-flex items-center gap-1 text-sm font-medium text-primary hover:underline"
        >
          <Plus size={14} />
          Ajouter
        </button>
      </div>

      {/* Check-in form */}
      {showCheckinForm ? (
        <Card padding="lg" className="mb-4">
          {/* Delay selection */}
          <p className="mb-2 font-heading text-xs font-medium text-foreground-secondary">
            Délai après le repas
          </p>
          <div className="mb-4 flex gap-2">
            {delayOptions.map((d) => (
              <button
                key={d.value}
                type="button"
                onClick={() => setCheckinDelay(d.value)}
                className={`
                  rounded-[--radius-sm] border px-3 py-1.5 text-xs font-medium transition-all
                  ${
                    checkinDelay === d.value
                      ? "border-primary bg-primary-light text-foreground"
                      : "border-border text-foreground-secondary"
                  }
                `}
              >
                {d.label}
              </button>
            ))}
          </div>

          {/* Symptom type toggles */}
          <p className="mb-2 font-heading text-xs font-medium text-foreground-secondary">
            Symptômes
          </p>
          <div className="mb-4 flex flex-wrap gap-1.5">
            {Object.entries(symptomLabels).map(([type, label]) => {
              const isSelected = checkinSymptoms.some((s) => s.type === type);
              return (
                <button
                  key={type}
                  type="button"
                  onClick={() => toggleSymptom(type)}
                  className={`
                    rounded-[--radius-full] border px-3 py-1 text-xs transition-all
                    ${
                      isSelected
                        ? "border-primary bg-primary-light text-foreground"
                        : "border-border text-foreground-secondary"
                    }
                  `}
                >
                  {label}
                </button>
              );
            })}
          </div>

          {/* Severity sliders for selected symptoms */}
          {checkinSymptoms.length > 0 ? (
            <div className="mb-4 flex flex-col gap-3">
              {checkinSymptoms.map((symptom) => (
                <div key={symptom.type} className="flex items-center justify-between">
                  <span className="text-xs text-foreground">
                    {symptomLabels[symptom.type] || symptom.type}
                  </span>
                  <div className="flex gap-1">
                    {[1, 2, 3, 4, 5].map((n) => (
                      <button
                        key={n}
                        type="button"
                        onClick={() => updateSymptomSeverity(symptom.type, n)}
                        className={`
                          h-7 w-7 rounded-[--radius-sm] text-xs font-medium transition-all
                          ${
                            symptom.severity === n
                              ? "bg-danger text-white"
                              : "border border-border text-foreground-secondary"
                          }
                        `}
                      >
                        {n}
                      </button>
                    ))}
                  </div>
                </div>
              ))}
            </div>
          ) : null}

          {/* Notes */}
          <textarea
            value={checkinNotes}
            onChange={(e) => setCheckinNotes(e.target.value)}
            placeholder="Notes (optionnel)..."
            rows={2}
            className="
              mb-3 w-full resize-none rounded-[--radius-md] border border-border
              bg-surface px-3 py-2 text-xs text-foreground outline-none
              transition-all focus:border-primary focus:ring-2 focus:ring-primary-light
            "
          />

          <div className="flex gap-2">
            <Button
              variant="outline"
              className="flex-1"
              onClick={() => setShowCheckinForm(false)}
            >
              Annuler
            </Button>
            <Button
              className="flex-1"
              onClick={handleSubmitCheckin}
              isLoading={isSubmittingCheckin}
              disabled={checkinSymptoms.length === 0}
            >
              Enregistrer
            </Button>
          </div>
        </Card>
      ) : null}

      {/* Existing check-ins */}
      {checkins.length === 0 && !showCheckinForm ? (
        <Card padding="md" className="flex items-center gap-3 text-foreground-secondary">
          <AlertTriangle size={16} />
          <p className="text-xs">
            Aucun symptôme signalé pour ce repas
          </p>
        </Card>
      ) : (
        <div className="flex flex-col gap-2">
          {checkins.map((checkin) => (
            <Card key={checkin.id} padding="md">
              <div className="mb-2 flex items-center justify-between">
                <span className="rounded-[--radius-full] bg-primary-light px-2 py-0.5 text-xs font-medium text-foreground">
                  Après {checkin.delay_hours}h
                </span>
              </div>
              <div className="flex flex-col gap-1.5">
                {(Array.isArray(checkin.symptoms) ? checkin.symptoms : []).map((s, i) => (
                  <div key={i} className="flex items-center justify-between">
                    <span className="text-xs text-foreground">
                      {symptomLabels[s.type] || s.type}
                    </span>
                    <SeverityDots severity={s.severity} />
                  </div>
                ))}
              </div>
              {checkin.notes ? (
                <p className="mt-2 text-xs text-foreground-secondary">
                  {checkin.notes}
                </p>
              ) : null}
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
