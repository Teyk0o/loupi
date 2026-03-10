"use client";

/**
 * Layout for authenticated application routes.
 * Wraps children with AuthGuard and includes bottom navigation.
 */

import { AuthGuard } from "@/components/guards/AuthGuard";
import { BottomNav } from "@/components/layout/BottomNav";

export default function AppLayout({ children }: { children: React.ReactNode }) {
  return (
    <AuthGuard>
      <main className="mx-auto min-h-screen max-w-lg px-4 pb-20 pt-6">
        {children}
      </main>
      <BottomNav />
    </AuthGuard>
  );
}
