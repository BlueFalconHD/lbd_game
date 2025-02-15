"use client";

import Image from "next/image";
import WithAuth from "@/components/withAuth";
import { getCurrentPhrase } from "@/lib/api/phrase";
import { useState } from "react";
import { AuthContext, AuthProvider } from "@/contexts/AuthContext";
import { PhraseDisplay } from "@/components/phrase-display";
import PickUser from "@/components/verify-usage/pick-user";
import { Verifications } from "@/components/verifications";

const Home = () => {
  const [error, setError] = useState<string | null>(null);

  // font-[family-name:var(--font-geist-sans)]"

  return (
    <div className="flex flex-col items-center justify-center min-h-screen max-w-screen gap-8">
      <PhraseDisplay />
      <Verifications />
    </div>
  );
};

export default WithAuth(Home);
