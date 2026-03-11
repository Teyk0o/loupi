"use client";

/**
 * Bottom navigation bar for the main application (mobile only).
 * Hidden on md+ screens where the Sidebar is used instead.
 */

import { usePathname } from "next/navigation";
import Link from "next/link";
import {
  BookOpen,
  PlusCircle,
  Heart,
  Settings,
  Stethoscope,
  type LucideIcon,
} from "lucide-react";

interface NavItem {
  href: string;
  label: string;
  icon: LucideIcon;
}

const navItems: NavItem[] = [
  { href: "/journal", label: "Journal", icon: BookOpen },
  { href: "/symptomes", label: "Symptômes", icon: Stethoscope },
  { href: "/ajouter", label: "Ajouter", icon: PlusCircle },
  { href: "/bien-etre", label: "Bien-être", icon: Heart },
  { href: "/parametres", label: "Paramètres", icon: Settings },
];

export function BottomNav() {
  const pathname = usePathname();

  return (
    <nav className="fixed bottom-0 left-0 right-0 z-50 border-t border-border bg-surface/95 pb-[env(safe-area-inset-bottom)] backdrop-blur-sm md:hidden">
      <div className="mx-auto flex items-center justify-around px-2 pb-6 pt-3">
        {navItems.map((item) => {
          const isActive = pathname.startsWith(item.href);
          const Icon = item.icon;

          return (
            <Link
              key={item.href}
              href={item.href}
              className={`
                flex flex-col items-center gap-0.5 px-2 py-1.5
                transition-colors duration-200
                ${isActive ? "text-primary" : "text-foreground-secondary hover:text-foreground"}
              `}
            >
              <Icon size={20} strokeWidth={isActive ? 2.2 : 1.8} />
              <span className="font-heading text-[10px] font-medium leading-tight">
                {item.label}
              </span>
            </Link>
          );
        })}
      </div>
    </nav>
  );
}
