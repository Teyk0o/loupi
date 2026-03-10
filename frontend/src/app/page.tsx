/**
 * Root page — redirects to the journal.
 */

import { redirect } from "next/navigation";

export default function Home() {
  redirect("/journal");
}
