"use client";

/**
 * Skip to main content link for keyboard accessibility
 */
export function SkipLink() {
  return (
    <a
      href="#main-content"
      className="sr-only focus:not-sr-only focus:fixed focus:top-2 focus:left-2 focus:z-50 focus:px-4 focus:py-2 focus:bg-[var(--color-primary)] focus:text-[var(--color-background)] focus:rounded"
    >
      Skip to main content
    </a>
  );
}
