"use client";

import { ReactNode } from "react";
import { ToastProvider } from "@/components/ui/Toast";
import { ThemeProvider } from "@/lib/theme/provider";
import { OfflineBanner } from "@/components/ui/OfflineBanner";

interface ProvidersProps {
  children: ReactNode;
}

/**
 * Client-side providers wrapper for the application.
 * Combines all context providers that need to wrap the app.
 */
export function Providers({ children }: ProvidersProps) {
  return (
    <ThemeProvider>
      <ToastProvider>
        <OfflineBanner />
        {children}
      </ToastProvider>
    </ThemeProvider>
  );
}
