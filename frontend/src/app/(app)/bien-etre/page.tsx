"use client";

/**
 * Wellness tracking page.
 * Allows users to log daily wellness metrics: stress, mood, energy, sleep, hydration, sport.
 */

import { useState, useEffect, type FormEvent } from "react";
import { useAuthContext } from "@/components/providers/AuthProvider";
import {
  Brain,
  Smile,
  Zap,
  Moon,
  Droplets,
  Dumbbell,
  Check,
  Plus,
  X,
} from "lucide-react";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { Input } from "@/components/ui/Input";
import { useCustomOptions } from "@/hooks/useCustomOptions";
import { apiFetch, type ApiError } from "@/lib/api";

interface SportSession {
  type: string;
  duration_minutes: number;
  intensity: number;
}

interface WellnessData {
  stress?: number;
  mood?: number;
  energy?: number;
  sleep_hours?: number;
  sleep_quality?: number;
  hydration?: number;
  sport?: SportSession[];
  notes?: string;
}


/** Scale labels for each metric (0 = none/best, 5 = worst/highest). */
const metricScaleLabels: Record<string, string[]> = {
  stress: ["Aucun", "Léger", "Modéré", "Élevé", "Intense", "Extrême"],
  mood: ["Très mauvaise", "Mauvaise", "Neutre", "Bonne", "Très bonne", "Excellente"],
  energy: ["Aucune", "Très faible", "Faible", "Correcte", "Bonne", "Débordante"],
  sleep_quality: ["Aucune", "Très mauvaise", "Mauvaise", "Correcte", "Bonne", "Excellente"],
};

/** Metric selector component (0-5 scale with descriptive labels). */
function MetricSelector({
  label,
  icon: Icon,
  value,
  metricKey,
  onChange,
}: {
  label: string;
  icon: React.ComponentType<{ size: number }>;
  value: number | undefined;
  metricKey: string;
  onChange: (v: number) => void;
}) {
  const labels = metricScaleLabels[metricKey] || [];
  return (
    <div>
      <div className="mb-2 flex items-center gap-2">
        <Icon size={16} />
        <span className="font-heading text-sm font-medium">{label}</span>
        {value !== undefined ? (
          <span className="text-xs text-foreground-secondary">
            — {labels[value] || value}
          </span>
        ) : null}
      </div>
      <div className="flex gap-1.5">
        {[0, 1, 2, 3, 4, 5].map((n) => (
          <button
            key={n}
            type="button"
            onClick={() => onChange(n)}
            className={`
              flex h-10 w-10 items-center justify-center rounded-[--radius-sm]
              text-sm font-medium transition-all
              ${
                value === n
                  ? "bg-primary text-white"
                  : "border border-border bg-surface text-foreground-secondary hover:border-primary"
              }
            `}
          >
            {n}
          </button>
        ))}
      </div>
    </div>
  );
}

function todayString(): string {
  return new Date().toISOString().split("T")[0];
}

export default function WellnessPage() {
  const { options: sportTypes } = useCustomOptions("sport_type");
  const [data, setData] = useState<WellnessData>({});
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [saved, setSaved] = useState(false);
  const [error, setError] = useState("");

  // Sport form state
  const [showSportForm, setShowSportForm] = useState(false);
  const [sportType, setSportType] = useState("");
  const [sportDuration, setSportDuration] = useState(30);
  const [sportIntensity, setSportIntensity] = useState(3);

  const { user } = useAuthContext();

  useEffect(() => {
    setData({});
    setIsLoading(true);
    async function load() {
      try {
        const existing = await apiFetch<WellnessData>(
          `/v1/wellness?date=${todayString()}`,
        );
        if (existing) setData(existing);
      } catch {
        // No entry yet
      } finally {
        setIsLoading(false);
      }
    }
    load();
  }, [user?.id]);

  function addSportSession() {
    const selectedType = sportType || (sportTypes.length > 0 ? sportTypes[0].value : "other");
    const session: SportSession = {
      type: selectedType,
      duration_minutes: sportDuration,
      intensity: sportIntensity,
    };
    setData((prev) => ({
      ...prev,
      sport: [...(prev.sport || []), session],
    }));
    setShowSportForm(false);
    setSportType("");
    setSportDuration(30);
    setSportIntensity(3);
  }

  function removeSportSession(index: number) {
    setData((prev) => ({
      ...prev,
      sport: (prev.sport || []).filter((_, i) => i !== index),
    }));
  }

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError("");
    setIsSaving(true);
    setSaved(false);

    try {
      await apiFetch("/v1/wellness", {
        method: "POST",
        body: JSON.stringify({
          date: todayString(),
          ...data,
          sport: data.sport && data.sport.length > 0 ? data.sport : undefined,
        }),
      });
      setSaved(true);
      setTimeout(() => setSaved(false), 2000);
    } catch (err) {
      const apiErr = err as ApiError;
      setError(apiErr.message || "Une erreur est survenue");
    } finally {
      setIsSaving(false);
    }
  }

  if (isLoading) {
    return (
      <div className="flex justify-center py-12">
        <span className="h-6 w-6 animate-spin rounded-full border-2 border-primary border-t-transparent" />
      </div>
    );
  }

  return (
    <div>
      <div className="mb-6">
        <h1 className="font-heading text-xl font-semibold">Bien-être</h1>
        <p className="text-sm text-foreground-secondary">
          Comment vous sentez-vous aujourd&apos;hui ?
        </p>
      </div>

      {error ? (
        <div className="mb-4 rounded-[--radius-md] border border-danger/30 bg-danger/10 px-4 py-3 text-sm text-danger">
          {error}
        </div>
      ) : null}

      <form onSubmit={handleSubmit} className="flex flex-col gap-5 md:grid md:grid-cols-2 md:items-start">
        {/* Mood & Energy metrics */}
        <Card padding="lg" className="flex flex-col gap-5">
          <MetricSelector
            label="Stress"
            icon={Brain}
            value={data.stress}
            metricKey="stress"
            onChange={(v) => setData((prev) => ({ ...prev, stress: v }))}
          />
          <MetricSelector
            label="Humeur"
            icon={Smile}
            value={data.mood}
            metricKey="mood"
            onChange={(v) => setData((prev) => ({ ...prev, mood: v }))}
          />
          <MetricSelector
            label="Énergie"
            icon={Zap}
            value={data.energy}
            metricKey="energy"
            onChange={(v) => setData((prev) => ({ ...prev, energy: v }))}
          />
          <MetricSelector
            label="Qualité du sommeil"
            icon={Moon}
            value={data.sleep_quality}
            metricKey="sleep_quality"
            onChange={(v) => setData((prev) => ({ ...prev, sleep_quality: v }))}
          />
        </Card>

        {/* Sleep, Hydration & Notes */}
        <div className="flex flex-col gap-5">
          <Card padding="lg" className="flex flex-col gap-4">
            <div>
              <div className="mb-2 flex items-center gap-2">
                <Moon size={16} />
                <span className="font-heading text-sm font-medium">
                  Heures de sommeil
                </span>
              </div>
              <input
                type="number"
                min={0}
                max={24}
                step={0.5}
                value={data.sleep_hours ?? ""}
                onChange={(e) =>
                  setData((prev) => ({
                    ...prev,
                    sleep_hours: e.target.value ? parseFloat(e.target.value) : undefined,
                  }))
                }
                placeholder="7.5"
                className="
                  w-full rounded-[--radius-md] border border-border bg-surface
                  px-4 py-3 text-sm text-foreground outline-none
                  transition-all focus:border-primary focus:ring-2 focus:ring-primary-light
                "
              />
            </div>

            <div>
              <div className="mb-2 flex items-center gap-2">
                <Droplets size={16} />
                <span className="font-heading text-sm font-medium">
                  Verres d&apos;eau
                </span>
              </div>
              <input
                type="number"
                min={0}
                max={50}
                value={data.hydration ?? ""}
                onChange={(e) =>
                  setData((prev) => ({
                    ...prev,
                    hydration: e.target.value ? parseInt(e.target.value) : undefined,
                  }))
                }
                placeholder="8"
                className="
                  w-full rounded-[--radius-md] border border-border bg-surface
                  px-4 py-3 text-sm text-foreground outline-none
                  transition-all focus:border-primary focus:ring-2 focus:ring-primary-light
                "
              />
            </div>
          </Card>

          <Card padding="lg">
            <label className="mb-2 block font-heading text-sm font-medium">
              Notes
            </label>
            <textarea
              value={data.notes ?? ""}
              onChange={(e) =>
                setData((prev) => ({ ...prev, notes: e.target.value || undefined }))
              }
              placeholder="Quelque chose à noter sur votre journée..."
              rows={3}
              maxLength={1000}
              className="
                w-full resize-none rounded-[--radius-md] border border-border
                bg-surface px-4 py-3 text-sm text-foreground
                placeholder:text-foreground-secondary outline-none
                transition-all focus:border-primary focus:ring-2 focus:ring-primary-light
              "
            />
          </Card>
        </div>

        {/* Sport sessions */}
        <Card padding="lg">
          <div className="mb-3 flex items-center justify-between">
            <div className="flex items-center gap-2">
              <Dumbbell size={16} />
              <span className="font-heading text-sm font-medium">Sport</span>
            </div>
            <button
              type="button"
              onClick={() => setShowSportForm(true)}
              className="inline-flex items-center gap-1 text-sm font-medium text-primary"
            >
              <Plus size={14} />
              Ajouter
            </button>
          </div>

          {/* Existing sport sessions */}
          {data.sport && data.sport.length > 0 ? (
            <div className="mb-3 flex flex-col gap-2">
              {data.sport.map((session, i) => {
                const sportInfo = sportTypes.find((s) => s.value === session.type);
                return (
                  <div
                    key={i}
                    className="flex items-center justify-between rounded-[--radius-sm] border border-border px-3 py-2"
                  >
                    <div className="flex items-center gap-2">
                      <span>{sportInfo?.emoji || "🏅"}</span>
                      <span className="text-xs font-medium">
                        {sportInfo?.label || session.type}
                      </span>
                      <span className="text-xs text-foreground-secondary">
                        {session.duration_minutes} min
                      </span>
                      <span className="text-xs text-foreground-secondary">
                        — Intensité {session.intensity}/5
                      </span>
                    </div>
                    <button
                      type="button"
                      onClick={() => removeSportSession(i)}
                      className="text-foreground-secondary hover:text-danger"
                    >
                      <X size={14} />
                    </button>
                  </div>
                );
              })}
            </div>
          ) : !showSportForm ? (
            <p className="mb-3 text-xs text-foreground-secondary">
              Aucune session enregistrée
            </p>
          ) : null}

          {/* Sport add form */}
          {showSportForm ? (
            <div className="flex flex-col gap-3 rounded-[--radius-md] border border-border p-3">
              {/* Sport type */}
              <div className="flex flex-wrap gap-1.5">
                {sportTypes.map((st) => (
                  <button
                    key={st.value}
                    type="button"
                    onClick={() => setSportType(st.value)}
                    className={`
                      rounded-[--radius-full] border px-2.5 py-1 text-xs transition-all
                      ${
                        sportType === st.value
                          ? "border-primary bg-primary-light text-foreground"
                          : "border-border text-foreground-secondary"
                      }
                    `}
                  >
                    {st.emoji} {st.label}
                  </button>
                ))}
              </div>

              {/* Duration & Intensity */}
              <div className="grid grid-cols-2 gap-3">
                <Input
                  id="sport-duration"
                  label="Durée (min)"
                  type="number"
                  min={1}
                  max={600}
                  value={sportDuration}
                  onChange={(e) => setSportDuration(parseInt(e.target.value) || 0)}
                />
                <div>
                  <p className="mb-1.5 font-heading text-sm font-medium text-foreground">
                    Intensité
                  </p>
                  <div className="flex gap-1">
                    {[1, 2, 3, 4, 5].map((n) => (
                      <button
                        key={n}
                        type="button"
                        onClick={() => setSportIntensity(n)}
                        className={`
                          flex h-10 flex-1 items-center justify-center rounded-[--radius-sm]
                          text-xs font-medium transition-all
                          ${
                            sportIntensity === n
                              ? "bg-primary text-white"
                              : "border border-border text-foreground-secondary"
                          }
                        `}
                      >
                        {n}
                      </button>
                    ))}
                  </div>
                </div>
              </div>

              <div className="flex gap-2">
                <Button
                  variant="outline"
                  className="flex-1"
                  type="button"
                  onClick={() => setShowSportForm(false)}
                >
                  Annuler
                </Button>
                <Button
                  className="flex-1"
                  type="button"
                  onClick={addSportSession}
                >
                  Ajouter
                </Button>
              </div>
            </div>
          ) : null}
        </Card>

        <Button type="submit" isLoading={isSaving} className="w-full md:col-span-2">
          {saved ? (
            <span className="inline-flex items-center gap-1.5">
              <Check size={16} /> Enregistré
            </span>
          ) : (
            "Enregistrer"
          )}
        </Button>
      </form>
    </div>
  );
}
