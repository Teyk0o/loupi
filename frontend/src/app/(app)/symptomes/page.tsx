"use client";

/**
 * Standalone symptoms page.
 * Allows users to log symptoms that are not linked to a specific meal.
 */

import { useState, useEffect, useCallback, useMemo, type FormEvent } from "react";
import { Plus, Trash2 } from "lucide-react";
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

function todayString(): string {
  return new Date().toISOString().split("T")[0];
}

function nowTimeString(): string {
  return new Date().toTimeString().slice(0, 5);
}

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

export default function SymptomsPage() {
  const { options: symptomTypes } = useCustomOptions("symptom_type");
  const [entries, setEntries] = useState<SymptomEntry[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [selectedDate, setSelectedDate] = useState(todayString);

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

  const fetchEntries = useCallback(async (date: string) => {
    setIsLoading(true);
    try {
      const data = await apiFetch<SymptomEntry[]>(
        `/v1/symptoms?date=${date}`,
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
      const entryTime = `${selectedDate}T${time}:00Z`;
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
      // Silent fail
    }
  }

  return (
    <div>
      <div className="mb-6 flex items-center justify-between">
        <div>
          <h1 className="font-heading text-xl font-semibold">Symptômes</h1>
          <p className="text-sm text-foreground-secondary">
            Signalez des symptômes indépendants d&apos;un repas
          </p>
        </div>
        <button
          onClick={() => setShowForm(!showForm)}
          className="flex h-10 w-10 items-center justify-center rounded-full bg-primary text-white"
        >
          <Plus size={20} />
        </button>
      </div>

      {/* Date filter */}
      <div className="mb-4">
        <Input
          id="date-filter"
          type="date"
          value={selectedDate}
          onChange={(e) => setSelectedDate(e.target.value)}
        />
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
              <p className="mb-2 font-heading text-xs font-medium text-foreground-secondary">
                Symptômes ressentis
              </p>
              <div className="flex flex-wrap gap-1.5">
                {symptomTypes.map((st) => {
                  const isSelected = symptoms.some((s) => s.type === st.value);
                  return (
                    <button
                      key={st.value}
                      type="button"
                      onClick={() => toggleSymptom(st.value)}
                      className={`
                        rounded-[--radius-full] border px-3 py-1 text-xs transition-all
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
                    <span className="text-xs text-foreground">
                      {symptomLabels[symptom.type]}
                    </span>
                    <div className="flex gap-1">
                      {[1, 2, 3, 4, 5].map((n) => (
                        <button
                          key={n}
                          type="button"
                          onClick={() => updateSeverity(symptom.type, n)}
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
        <p className="py-8 text-center text-sm text-foreground-secondary">
          Aucun symptôme enregistré pour cette date
        </p>
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
                  className="text-foreground-secondary hover:text-danger"
                >
                  <Trash2 size={14} />
                </button>
              </div>
              <p className="mb-2 text-xs text-foreground-secondary">
                {new Date(entry.entry_time).toLocaleTimeString("fr-FR", {
                  hour: "2-digit",
                  minute: "2-digit",
                })}
              </p>
              <div className="flex flex-col gap-1.5">
                {(Array.isArray(entry.symptoms) ? entry.symptoms : []).map(
                  (s, i) => (
                    <div key={i} className="flex items-center justify-between">
                      <span className="text-xs text-foreground">
                        {symptomLabels[s.type] || s.type}
                      </span>
                      <SeverityDots severity={s.severity} />
                    </div>
                  ),
                )}
              </div>
              {entry.notes ? (
                <p className="mt-2 text-xs text-foreground-secondary">
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
