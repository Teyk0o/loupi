"use client";

/**
 * Floating action button with expandable menu.
 * Allows quick access to add a meal, symptom, or wellness entry from any page.
 * Positioned bottom-right on all viewports.
 */

import { useState, useRef, useEffect } from "react";
import { useRouter } from "next/navigation";
import {
  Plus,
  X,
  UtensilsCrossed,
  Stethoscope,
  Heart,
} from "lucide-react";

interface QuickAction {
  label: string;
  href: string;
  icon: React.ComponentType<{ size: number }>;
  color: string;
}

const actions: QuickAction[] = [
  { label: "Repas", href: "/journal?add=meal", icon: UtensilsCrossed, color: "bg-primary" },
  { label: "Symptômes", href: "/symptomes?add=true", icon: Stethoscope, color: "bg-amber-500" },
  { label: "Bien-être", href: "/bien-etre", icon: Heart, color: "bg-rose-400" },
];

export function FloatingAddButton() {
  const [isOpen, setIsOpen] = useState(false);
  const menuRef = useRef<HTMLDivElement>(null);
  const router = useRouter();

  /** Close menu when clicking outside. */
  useEffect(() => {
    if (!isOpen) return;
    function handleClickOutside(e: MouseEvent) {
      if (menuRef.current && !menuRef.current.contains(e.target as Node)) {
        setIsOpen(false);
      }
    }
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, [isOpen]);

  function handleAction(href: string) {
    setIsOpen(false);
    router.push(href);
  }

  return (
    <div ref={menuRef} className="fixed bottom-24 right-4 z-40 h-14 w-14 md:bottom-6">
      {/* Expandable menu — positioned absolutely so it doesn't shift the FAB */}
      {isOpen ? (
        <div className="absolute bottom-16 right-0 flex flex-col gap-2 pb-2">
          {actions.map((action) => {
            const Icon = action.icon;
            return (
              <button
                key={action.href}
                onClick={() => handleAction(action.href)}
                className="flex items-center gap-3 self-end whitespace-nowrap"
              >
                <span className="rounded-[--radius-md] bg-surface px-3 py-1.5 text-sm font-medium text-foreground shadow-md">
                  {action.label}
                </span>
                <span
                  className={`flex h-10 w-10 shrink-0 items-center justify-center rounded-full text-white shadow-md ${action.color}`}
                >
                  <Icon size={18} />
                </span>
              </button>
            );
          })}
        </div>
      ) : null}

      {/* Main FAB */}
      <button
        onClick={() => setIsOpen((prev) => !prev)}
        className={`
          flex h-14 w-14 items-center justify-center rounded-full
          bg-primary text-white shadow-lg
          transition-transform duration-200
          hover:scale-105 active:scale-95
          ${isOpen ? "rotate-45" : ""}
        `}
        aria-label={isOpen ? "Fermer le menu" : "Ajouter"}
      >
        {isOpen ? <X size={28} /> : <Plus size={28} />}
      </button>
    </div>
  );
}
