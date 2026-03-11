"use client";

/**
 * Layout for authenticated application routes.
 * Mobile: logo header + bottom navigation + centered content.
 * Desktop: sidebar navigation + wider content area.
 */

import { AuthGuard } from "@/components/guards/AuthGuard";
import { BottomNav } from "@/components/layout/BottomNav";
import { FloatingAddButton } from "@/components/layout/FloatingAddButton";
import { Sidebar } from "@/components/layout/Sidebar";
import { Logo } from "@/components/ui/Logo";

export default function AppLayout({ children }: { children: React.ReactNode }) {
  return (
    <AuthGuard>
      <Sidebar />
      <div className="md:ml-60">
        {/* Mobile header with logo */}
        <header className="flex h-14 items-center border-b border-border px-4 md:hidden">
          <Logo width={90} height={30} />
        </header>
        <main className="mx-auto max-w-lg px-4 pb-24 pt-4 md:max-w-none md:px-8 md:pb-6 md:pt-6 lg:px-12">
          {children}
        </main>
      </div>
      <FloatingAddButton />
      <BottomNav />
    </AuthGuard>
  );
}
