"use client";

import * as React from "react";
import { Loader2 } from "lucide-react";

export function LoadingPage({ ...props }) {
  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-background/80"
      {...props}
    >
      <Loader2 className="h-8 w-8 animate-spin" />
    </div>
  );
}
