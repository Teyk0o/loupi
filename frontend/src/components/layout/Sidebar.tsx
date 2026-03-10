"use client";

/**
 * Sidebar navigation for tablet and desktop viewports.
 * Hidden on mobile, displayed as a fixed left sidebar on md+ screens.
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
import { Logo } from "@/components/ui/Logo";

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

export function Sidebar() {
  const pathname = usePathname();

  return (
    <aside className="fixed left-0 top-0 z-50 hidden h-screen w-60 flex-col border-r border-border bg-surface md:flex">
      {/* Logo */}
      <div className="flex h-16 items-center px-6">
        <Logo width={100} height={33} />
      </div>

      {/* Navigation */}
      <nav className="flex flex-1 flex-col gap-1 px-3 pt-2">
        {navItems.map((item) => {
          const isActive =
            pathname === item.href || pathname.startsWith(item.href + "/");
          const Icon = item.icon;

          return (
            <Link
              key={item.href}
              href={item.href}
              className={`
                flex items-center gap-3 rounded-[--radius-md] px-3 py-2.5
                font-heading text-sm transition-colors duration-200
                ${
                  isActive
                    ? "bg-primary-light font-medium text-primary"
                    : "text-foreground-secondary hover:bg-background hover:text-foreground"
                }
              `}
            >
              <Icon size={20} strokeWidth={isActive ? 2.2 : 1.8} />
              <span>{item.label}</span>
            </Link>
          );
        })}
      </nav>

      {/* Footer */}
      <div className="border-t border-border px-6 py-4">
        <p className="text-[10px] text-foreground-secondary">
          Loupi — Open Source
        </p>
      </div>
    </aside>
  );
}
