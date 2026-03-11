"use client";

/**
 * Standalone symptoms page.
 * Allows users to log symptoms that are not linked to a specific meal.
 */

import { useState, useEffect, useCallback, useMemo, type FormEvent } from "react";
import { Plus, Trash2, ChevronLeft, ChevronRight } from "lucide-react";
import { useSearchParams } from "next/navigation";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { Input } from "@/components/ui/Input";
import { useCustomOptions } from "@/hooks/useCustomOptions";
import { apiFetch, type ApiError } from "@/lib/api";

interface SymptomDetail {
  type: string;
  severity: number;
}

interface SymptomEntry {
  id: string;
  symptoms: SymptomDetail[];
  notes?: string;
  entry_time: string;
  created_at: string;
}

/** Format a Date to YYYY-MM-DD. */
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

/** Severity display with colored dots. */
function SeverityDots({ severity }: { severity: number }) {
  return (
    <div className="flex gap-1">
      {[1, 2, 3, 4, 5].map((n) => (
        <span
          key={n}
          className={`h-2.5 w-2.5 rounded-full ${
            n <= severity ? "bg-danger" : "bg-border"
          }`}
        />
      ))}
    </div>
  );
}

export default function SymptomsPage() {
  const searchParams = useSearchParams();
  const { options: symptomTypes } = useCustomOptions("symptom_type");
  const [entries, setEntries] = useState<SymptomEntry[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [selectedDate, setSelectedDate] = useState(() => new Date());

  /** Build label lookup from user's custom options. */
  const symptomLabels = useMemo(() => {
    const map: Record<string, string> = {};
    for (const s of symptomTypes) map[s.value] = s.label;
    return map;
  }, [symptomTypes]);

  // Form state
  const [showForm, setShowForm] = useState(false);
  const [symptoms, setSymptoms] = useState<SymptomDetail[]>([]);
  const [notes, setNotes] = useState("");
  const [time, setTime] = useState(nowTimeString);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState("");

  // Open form if redirected from /ajouter with ?add=true
  useEffect(() => {
    if (searchParams.get("add") === "true") {
      setShowForm(true);
    }
  }, [searchParams]);

  const fetchEntries = useCallback(async (date: Date) => {
    setIsLoading(true);
    try {
      const data = await apiFetch<SymptomEntry[]>(
        `/v1/symptoms?date=${toDateString(date)}`,
      );
      setEntries(data || []);
    } catch {
      setEntries([]);
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchEntries(selectedDate);
  }, [selectedDate, fetchEntries]);

  function toggleSymptom(type: string) {
    setSymptoms((prev) => {
      const existing = prev.find((s) => s.type === type);
      if (existing) return prev.filter((s) => s.type !== type);
      return [...prev, { type, severity: 3 }];
    });
  }

  function updateSeverity(type: string, severity: number) {
    setSymptoms((prev) =>
      prev.map((s) => (s.type === type ? { ...s, severity } : s)),
    );
  }

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    if (symptoms.length === 0) return;

    setIsSubmitting(true);
    setError("");

    try {
      const entryTime = `${toDateString(selectedDate)}T${time}:00${localTimezoneOffset()}`;
      await apiFetch("/v1/symptoms", {
        method: "POST",
        body: JSON.stringify({
          symptoms,
          notes: notes || undefined,
          entry_time: entryTime,
        }),
      });
      setShowForm(false);
      setSymptoms([]);
      setNotes("");
      setTime(nowTimeString());
      await fetchEntries(selectedDate);
    } catch (err) {
      const apiErr = err as ApiError;
      setError(apiErr.message || "Une erreur est survenue");
    } finally {
      setIsSubmitting(false);
    }
  }

  async function handleDelete(id: string) {
    if (!confirm("Supprimer cette entrée ?")) return;
    try {
      await apiFetch(`/v1/symptoms/${id}`, { method: "DELETE" });
      await fetchEntries(selectedDate);
    } catch {
      setError("Erreur lors de la suppression");
    }
  }

  return (
    <div>
      {/* Header + Date selector (same layout as journal) */}
      <div className="mb-6 md:flex md:items-center md:justify-between">
        <div className="mb-4 md:mb-0">
          <h1 className="font-heading text-xl font-semibold">Symptômes</h1>
          <p className="text-sm text-foreground-secondary">
            Signalez des symptômes indépendants d&apos;un repas
          </p>
        </div>

        <div className="flex items-center gap-3">
          {/* Date selector */}
          <div className="mx-auto flex w-fit items-center gap-2 rounded-[--radius-md] border border-border bg-surface px-2 py-1 md:mx-0 md:px-3">
            <button
              onClick={() =>
                setSelectedDate((prev) => {
                  const d = new Date(prev);
                  d.setDate(d.getDate() - 1);
                  return d;
                })
              }
              className="rounded-[--radius-sm] p-1.5 text-foreground-secondary hover:text-foreground"
            >
              <ChevronLeft size={18} />
            </button>
            <button
              onClick={() => setSelectedDate(new Date())}
              className={`min-w-[140px] text-center font-heading text-sm font-medium ${
                toDateString(selectedDate) === toDateString(new Date())
                  ? "text-primary"
                  : "text-foreground"
              }`}
            >
              {toDateString(selectedDate) === toDateString(new Date())
                ? "Aujourd'hui"
                : formatDisplayDate(selectedDate)}
            </button>
            <button
              onClick={() =>
                setSelectedDate((prev) => {
                  const d = new Date(prev);
                  d.setDate(d.getDate() + 1);
                  return d;
                })
              }
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

      {/* Add form */}
      {showForm ? (
        <Card padding="lg" className="mb-4">
          <form onSubmit={handleSubmit} className="flex flex-col gap-4">
            {/* Symptom type toggles */}
            <div>
              <p className="mb-2 font-heading text-sm font-medium text-foreground-secondary">
                Symptômes ressentis
              </p>
              <div className="flex flex-wrap gap-2">
                {symptomTypes.map((st) => {
                  const isSelected = symptoms.some((s) => s.type === st.value);
                  return (
                    <button
                      key={st.value}
                      type="button"
                      onClick={() => toggleSymptom(st.value)}
                      className={`
                        rounded-[--radius-full] border px-3.5 py-1.5 text-sm transition-all
                        ${
                          isSelected
                            ? "border-primary bg-primary-light text-foreground"
                            : "border-border text-foreground-secondary"
                        }
                      `}
                    >
                      {st.label}
                    </button>
                  );
                })}
              </div>
            </div>

            {/* Severity per symptom */}
            {symptoms.length > 0 ? (
              <div className="flex flex-col gap-3">
                {symptoms.map((symptom) => (
                  <div key={symptom.type} className="flex items-center justify-between">
                    <span className="text-sm text-foreground">
                      {symptomLabels[symptom.type]}
                    </span>
                    <div className="flex gap-1.5">
                      {[1, 2, 3, 4, 5].map((n) => (
                        <button
                          key={n}
                          type="button"
                          onClick={() => updateSeverity(symptom.type, n)}
                          className={`
                            h-9 w-9 rounded-[--radius-sm] text-sm font-medium transition-all
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

            {/* Time */}
            <Input
              id="symptom-time"
              label="Heure"
              type="time"
              value={time}
              onChange={(e) => setTime(e.target.value)}
              required
            />

            {/* Notes */}
            <textarea
              value={notes}
              onChange={(e) => setNotes(e.target.value)}
              placeholder="Notes (optionnel)..."
              rows={2}
              maxLength={1000}
              className="
                w-full resize-none rounded-[--radius-md] border border-border
                bg-surface px-3 py-2 text-sm text-foreground outline-none
                transition-all focus:border-primary focus:ring-2 focus:ring-primary-light
              "
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
                disabled={symptoms.length === 0}
              >
                Enregistrer
              </Button>
            </div>
          </form>
        </Card>
      ) : null}

      {/* Existing entries */}
      {isLoading ? (
        <div className="flex justify-center py-12">
          <span className="h-6 w-6 animate-spin rounded-full border-2 border-primary border-t-transparent" />
        </div>
      ) : entries.length === 0 ? (
        <div className="flex flex-col items-center gap-3 py-12 text-center">
          <p className="text-foreground-secondary">
            Aucun symptôme enregistré pour cette date
          </p>
          <button
            onClick={() => setShowForm(true)}
            className="inline-flex items-center gap-1.5 font-heading text-sm font-medium text-primary hover:underline"
          >
            <Plus size={16} />
            Ajouter un symptôme
          </button>
        </div>
      ) : (
        <div className="flex flex-col gap-3 md:grid md:grid-cols-2 lg:grid-cols-3">
          {entries.map((entry) => (
            <Card key={entry.id} padding="md">
              <div className="mb-2 flex items-center justify-between">
                <span className="font-heading text-sm font-medium">
                  {(Array.isArray(entry.symptoms) ? entry.symptoms : [])
                    .map((s) => symptomLabels[s.type] || s.type)
                    .join(", ")}
                </span>
                <button
                  onClick={() => handleDelete(entry.id)}
                  className="p-1.5 text-foreground-secondary hover:text-danger"
                >
                  <Trash2 size={16} />
                </button>
              </div>
              <p className="mb-2 text-sm text-foreground-secondary">
                {new Date(entry.entry_time).toLocaleTimeString("fr-FR", {
                  hour: "2-digit",
                  minute: "2-digit",
                })}
              </p>
              <div className="flex flex-col gap-2">
                {(Array.isArray(entry.symptoms) ? entry.symptoms : []).map(
                  (s, i) => (
                    <div key={i} className="flex items-center justify-between">
                      <span className="text-sm text-foreground">
                        {symptomLabels[s.type] || s.type}
                      </span>
                      <SeverityDots severity={s.severity} />
                    </div>
                  ),
                )}
              </div>
              {entry.notes ? (
                <p className="mt-2 text-sm text-foreground-secondary">
                  {entry.notes}
                </p>
              ) : null}
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
