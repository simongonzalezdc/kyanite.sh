"use client";

import { useState, useRef, useEffect } from "react";
import { useTheme } from "@/lib/theme/provider";

interface ThemeSelectorProps {
  compact?: boolean;
}

export function ThemeSelector({ compact = false }: ThemeSelectorProps) {
  const { theme, themeId, setTheme, availableThemes } = useTheme();
  const [isOpen, setIsOpen] = useState(false);
  const containerRef = useRef<HTMLDivElement>(null);

  // Close on click outside
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (containerRef.current && !containerRef.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    }

    if (isOpen) {
      document.addEventListener("mousedown", handleClickOutside);
      return () => document.removeEventListener("mousedown", handleClickOutside);
    }
  }, [isOpen]);

  // Close on escape
  useEffect(() => {
    function handleEscape(event: KeyboardEvent) {
      if (event.key === "Escape") {
        setIsOpen(false);
      }
    }

    if (isOpen) {
      document.addEventListener("keydown", handleEscape);
      return () => document.removeEventListener("keydown", handleEscape);
    }
  }, [isOpen]);

  return (
    <div ref={containerRef} className="relative">
      {/* Trigger button */}
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="flex items-center gap-2 min-h-[44px] px-3 rounded-lg hover:bg-[var(--color-surface)] transition-colors focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]"
        aria-label="Select theme"
        aria-expanded={isOpen}
        aria-haspopup="listbox"
      >
        {/* Theme color preview */}
        <div
          className="w-5 h-5 rounded-full border-2 border-[var(--color-text-muted)]"
          style={{ backgroundColor: theme.primary }}
        />
        {!compact && (
          <span className="text-sm text-[var(--color-text)]">{theme.name}</span>
        )}
        <svg
          className={`w-4 h-4 text-[var(--color-text-muted)] transition-transform ${isOpen ? "rotate-180" : ""}`}
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
        </svg>
      </button>

      {/* Dropdown */}
      {isOpen && (
        <div
          role="listbox"
          aria-label="Available themes"
          className="absolute right-0 top-full mt-2 w-64 py-2 rounded-lg border border-[var(--color-surface)] shadow-xl z-50"
          style={{ backgroundColor: "var(--color-background)" }}
        >
          <div className="px-3 py-2 text-xs font-medium text-[var(--color-text-muted)] uppercase tracking-wider">
            Choose Theme
          </div>
          <div className="max-h-[300px] overflow-y-auto">
            {availableThemes.map((t) => (
              <button
                key={t.id}
                role="option"
                aria-selected={t.id === themeId}
                onClick={() => {
                  setTheme(t.id);
                  setIsOpen(false);
                }}
                className={`w-full flex items-center gap-3 px-3 py-3 min-h-[44px] text-left transition-colors hover:bg-[var(--color-surface)] focus:outline-none focus:bg-[var(--color-surface)] ${
                  t.id === themeId ? "bg-[var(--color-surface)]" : ""
                }`}
              >
                {/* Theme color swatches */}
                <div className="flex gap-1">
                  <div
                    className="w-4 h-4 rounded-full"
                    style={{ backgroundColor: t.primary }}
                    title="Primary"
                  />
                  <div
                    className="w-4 h-4 rounded-full"
                    style={{ backgroundColor: t.secondary }}
                    title="Secondary"
                  />
                  <div
                    className="w-4 h-4 rounded-full"
                    style={{ backgroundColor: t.accent }}
                    title="Accent"
                  />
                </div>

                {/* Theme name */}
                <span className="flex-1 text-sm text-[var(--color-text)]">{t.name}</span>

                {/* Selected indicator */}
                {t.id === themeId && (
                  <svg
                    className="w-5 h-5 text-[var(--color-success)]"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                  </svg>
                )}
              </button>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
