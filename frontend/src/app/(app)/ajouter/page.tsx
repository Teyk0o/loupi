"use client";

/**
 * Add entry hub page.
 * Presents three choices that redirect to the appropriate page
 * with the add form automatically opened.
 */

import Link from "next/link";
import {
  UtensilsCrossed,
  Stethoscope,
  Heart,
  ChevronRight,
} from "lucide-react";
import { Card } from "@/components/ui/Card";

interface QuickAction {
  href: string;
  label: string;
  description: string;
  icon: React.ComponentType<{ size: number; className?: string }>;
  iconBg: string;
}

const actions: QuickAction[] = [
  {
    href: "/journal?add=meal",
    label: "Repas",
    description: "Enregistrez ce que vous avez mange",
    icon: UtensilsCrossed,
    iconBg: "bg-primary/15 text-primary",
  },
  {
    href: "/symptomes?add=true",
    label: "Symptome",
    description: "Signalez un symptome digestif",
    icon: Stethoscope,
    iconBg: "bg-amber-500/15 text-amber-600",
  },
  {
    href: "/bien-etre",
    label: "Bien-etre",
    description: "Notez votre humeur, stress, sommeil...",
    icon: Heart,
    iconBg: "bg-rose-400/15 text-rose-500",
  },
];

export default function AddEntryPage() {
  return (
    <div>
      <div className="mb-6">
        <h1 className="font-heading text-xl font-semibold">Ajouter</h1>
        <p className="text-sm text-foreground-secondary">
          Que souhaitez-vous enregistrer ?
        </p>
      </div>

      <div className="flex flex-col gap-3">
        {actions.map((action) => {
          const Icon = action.icon;
          return (
            <Link key={action.href} href={action.href}>
              <Card className="flex items-center gap-4 transition-colors hover:border-primary/30">
                <div
                  className={`flex h-12 w-12 shrink-0 items-center justify-center rounded-[--radius-md] ${action.iconBg}`}
                >
                  <Icon size={24} />
                </div>
                <div className="flex-1">
                  <p className="font-heading text-base font-medium">
                    {action.label}
                  </p>
                  <p className="text-sm text-foreground-secondary">
                    {action.description}
                  </p>
                </div>
                <ChevronRight
                  size={20}
                  className="shrink-0 text-foreground-secondary"
                />
              </Card>
            </Link>
          );
        })}
      </div>
    </div>
  );
}
