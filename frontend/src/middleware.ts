import { NextRequest, NextResponse } from "next/server";

const PROTECTED_PATHS = ["/journal", "/ajouter", "/symptomes", "/bien-etre", "/parametres"];
const GUEST_PATHS = ["/login", "/register"];

/**
 * Next.js middleware for defense-in-depth auth routing.
 * Checks for the presence of the loupi_access cookie:
 * - If absent on protected routes → redirect to /login
 * - If present on guest routes → redirect to /journal
 *
 * Note: actual token validation happens server-side in the Go API.
 */
export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;
  const hasAccessCookie = request.cookies.has("loupi_access");

  // Protected routes: redirect to login if no cookie
  if (PROTECTED_PATHS.some((p) => pathname.startsWith(p))) {
    if (!hasAccessCookie) {
      return NextResponse.redirect(new URL("/login", request.url));
    }
  }

  // Guest routes: redirect to journal if already authenticated
  if (GUEST_PATHS.some((p) => pathname.startsWith(p))) {
    if (hasAccessCookie) {
      return NextResponse.redirect(new URL("/journal", request.url));
    }
  }

  return NextResponse.next();
}

export const config = {
  matcher: [
    "/journal/:path*",
    "/ajouter/:path*",
    "/symptomes/:path*",
    "/bien-etre/:path*",
    "/parametres/:path*",
    "/login",
    "/register",
  ],
};
