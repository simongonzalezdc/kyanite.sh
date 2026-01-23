"use client";

import { useState } from "react";
import { Button } from "@/components/ui/Button";

// =============================================================================
// Constants
// =============================================================================

const ONBOARDING_STORAGE_KEY = "noise-onboarding-complete";

interface OnboardingStep {
  title: string;
  description: string;
  icon: string;
  tip?: string;
}

const ONBOARDING_STEPS: OnboardingStep[] = [
  {
    title: "Welcome to noise.sh",
    description:
      "Capture song ideas on the go. Voice memos, photos, tap tempo, and quick notes - all synced to your desktop.",
    icon: "♪",
    tip: "Your ideas will sync automatically when connected to your computer.",
  },
  {
    title: "Voice Memos",
    description:
      "Record melody ideas, lyric snippets, or production notes. Tap the microphone to start recording.",
    icon: "🎤",
    tip: "Try humming your melody idea - it's often easier than describing it!",
  },
  {
    title: "Photo Capture",
    description:
      "Snap photos of handwritten lyrics, chord charts, studio setups, or inspiration for album art.",
    icon: "📷",
    tip: "Photos work great for capturing chord progressions you see at a gig.",
  },
  {
    title: "Tap Tempo",
    description:
      "Find your groove by tapping the beat. We'll calculate the BPM and save it with your idea.",
    icon: "♪",
    tip: "Tap at least 4 times for an accurate BPM reading.",
  },
  {
    title: "Quick Notes",
    description:
      "Type lyrics, song structure ideas, or production notes. Everything syncs to your desktop project.",
    icon: "📝",
    tip: "Use Cmd/Ctrl+Enter to quickly save your note.",
  },
  {
    title: "Swipe to Navigate",
    description:
      "Swipe left or right on the capture area to quickly switch between capture tools.",
    icon: "👆",
    tip: "Works best on touchscreen devices.",
  },
];

// =============================================================================
// Component
// =============================================================================

interface OnboardingWizardProps {
  onComplete: () => void;
  /** Force show even if previously completed */
  forceShow?: boolean;
}

export function OnboardingWizard({ onComplete, forceShow = false }: OnboardingWizardProps) {
  const [currentStep, setCurrentStep] = useState(0);

  // Initialize visibility based on forceShow or localStorage
  const [isVisible, setIsVisible] = useState(() => {
    if (forceShow) return true;
    if (typeof window === "undefined") return false;
    return localStorage.getItem(ONBOARDING_STORAGE_KEY) !== "true";
  });

  const handleNext = () => {
    if (currentStep < ONBOARDING_STEPS.length - 1) {
      setCurrentStep(currentStep + 1);
    } else {
      handleComplete();
    }
  };

  const handlePrevious = () => {
    if (currentStep > 0) {
      setCurrentStep(currentStep - 1);
    }
  };

  const handleComplete = () => {
    localStorage.setItem(ONBOARDING_STORAGE_KEY, "true");
    setIsVisible(false);
    onComplete();
  };

  const handleSkip = () => {
    handleComplete();
  };

  if (!isVisible) {
    return null;
  }

  const step = ONBOARDING_STEPS[currentStep];
  const isLastStep = currentStep === ONBOARDING_STEPS.length - 1;
  const isFirstStep = currentStep === 0;

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-[var(--color-background)]/95 backdrop-blur-sm p-4"
      role="dialog"
      aria-modal="true"
      aria-labelledby="onboarding-title"
    >
      <div className="w-full max-w-md bg-[var(--color-surface)] rounded-xl p-6 shadow-2xl">
        {/* Progress dots */}
        <div className="flex justify-center gap-2 mb-6">
          {ONBOARDING_STEPS.map((_, index) => (
            <button
              key={index}
              onClick={() => setCurrentStep(index)}
              className={`w-2 h-2 rounded-full transition-all ${
                index === currentStep
                  ? "w-6 bg-[var(--color-primary)]"
                  : index < currentStep
                  ? "bg-[var(--color-success)]"
                  : "bg-[var(--color-text-muted)]"
              }`}
              aria-label={`Go to step ${index + 1}`}
            />
          ))}
        </div>

        {/* Icon */}
        <div className="text-6xl text-center mb-4" aria-hidden="true">
          {step.icon}
        </div>

        {/* Title */}
        <h2
          id="onboarding-title"
          className="text-xl font-bold text-[var(--color-text)] text-center mb-3"
        >
          {step.title}
        </h2>

        {/* Description */}
        <p className="text-[var(--color-text-muted)] text-center mb-4">
          {step.description}
        </p>

        {/* Tip */}
        {step.tip && (
          <div className="bg-[var(--color-primary)]/10 border border-[var(--color-primary)]/30 rounded-lg p-3 mb-6">
            <p className="text-sm text-[var(--color-text)]">
              <span className="font-medium text-[var(--color-primary)]">💡 Tip: </span>
              {step.tip}
            </p>
          </div>
        )}

        {/* Navigation buttons */}
        <div className="flex gap-3">
          {!isFirstStep && (
            <Button variant="ghost" onClick={handlePrevious} className="flex-1">
              Back
            </Button>
          )}

          {isFirstStep && (
            <Button variant="ghost" onClick={handleSkip} className="flex-1">
              Skip
            </Button>
          )}

          <Button onClick={handleNext} className="flex-1">
            {isLastStep ? "Get Started" : "Next"}
          </Button>
        </div>

        {/* Step counter */}
        <p className="text-xs text-[var(--color-text-muted)] text-center mt-4">
          Step {currentStep + 1} of {ONBOARDING_STEPS.length}
        </p>
      </div>
    </div>
  );
}

// =============================================================================
// Utility Functions
// =============================================================================

/**
 * Check if user has completed onboarding
 */
export function hasCompletedOnboarding(): boolean {
  if (typeof window === "undefined") return true;
  return localStorage.getItem(ONBOARDING_STORAGE_KEY) === "true";
}

/**
 * Reset onboarding (for testing or re-showing)
 */
export function resetOnboarding(): void {
  if (typeof window !== "undefined") {
    localStorage.removeItem(ONBOARDING_STORAGE_KEY);
  }
}
