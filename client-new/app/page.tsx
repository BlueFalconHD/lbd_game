"use client";

import WithAuth from "@/components/withAuth";
import { useContext, useState } from "react";
import { AuthContext } from "@/contexts/AuthContext";
import { PhraseDisplay } from "@/components/phrase-display";
import { Verifications } from "@/components/verifications";
import { Button } from "@/components/ui/button";

const Home = () => {
  const [error, setError] = useState<string | null>(null);
  const { logout } = useContext(AuthContext);

  // logout button at top right
  return (
    <div className="relative min-h-screen max-w-screen">
      <div className="absolute top-0 right-0 p-4">
        <Button onClick={logout} variant="link">
          Log Out
        </Button>
      </div>
      <div className="flex flex-col items-center justify-center min-h-screen gap-8">
        <PhraseDisplay />
        <Verifications />
      </div>
    </div>
  );
};

export default WithAuth(Home);
