"use client";

import { useState, useCallback, useRef } from "react";
import { Button } from "@/components/ui/Button";
import { getSyncClient } from "@/lib/api/client";
import { addPendingIdea } from "@/lib/db";
import {
  TAP_TEMPO_MIN_TAPS,
  TAP_TEMPO_MAX_TAPS,
  TAP_TEMPO_TIMEOUT_MS,
  BPM_MIN,
  BPM_MAX,
  MS_PER_MINUTE,
  TAP_HAPTIC_DURATION_MS,
} from "@/lib/constants";

interface TapTempoProps {
  onCaptured?: () => void;
  disabled?: boolean;
}

export function TapTempo({ onCaptured, disabled }: TapTempoProps) {
  const [taps, setTaps] = useState<number[]>([]);
  const [bpm, setBpm] = useState<number | null>(null);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [saved, setSaved] = useState(false);
  
  const timeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const calculateBpm = useCallback((tapTimes: number[]): number | null => {
    if (tapTimes.length < TAP_TEMPO_MIN_TAPS) return null;

    const intervals: number[] = [];
    for (let i = 1; i < tapTimes.length; i++) {
      intervals.push(tapTimes[i] - tapTimes[i - 1]);
    }

    const averageInterval = intervals.length > 0 
      ? intervals.reduce((a, b) => a + b, 0) / intervals.length 
      : 0;
    
    if (averageInterval <= 0) return null;
    
    const calculatedBpm = Math.round(MS_PER_MINUTE / averageInterval);

    // Clamp to reasonable BPM range
    return Math.max(BPM_MIN, Math.min(BPM_MAX, calculatedBpm));
  }, []);

  const handleTap = useCallback(() => {
    // Haptic feedback if available
    if (navigator.vibrate) {
      navigator.vibrate(TAP_HAPTIC_DURATION_MS);
    }

    const now = Date.now();

    // Reset timeout
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }
    
    timeoutRef.current = setTimeout(() => {
      // Don't reset if we have enough taps
      if (taps.length < TAP_TEMPO_MIN_TAPS) {
        setTaps([]);
        setBpm(null);
      }
    }, TAP_TEMPO_TIMEOUT_MS);

    setTaps((prevTaps) => {
      const newTaps = [...prevTaps, now].slice(-TAP_TEMPO_MAX_TAPS);
      const newBpm = calculateBpm(newTaps);
      setBpm(newBpm);
      return newTaps;
    });
    
    setSaved(false);
    setError("");
  }, [taps.length, calculateBpm]);

  const handleReset = useCallback(() => {
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }
    setTaps([]);
    setBpm(null);
    setSaved(false);
    setError("");
  }, []);

  const handleSave = useCallback(async () => {
    if (!bpm) return;

    setSaving(true);
    setError("");
    const localId = `tempo-${Date.now()}-${Math.random().toString(36).slice(2)}`;

    try {
      const client = getSyncClient();

      if (client) {
        await client.submitIdea({
          type: "tempo",
          content: `Tap tempo: ${bpm} BPM`,
          bpm,
        });
      } else {
        await addPendingIdea({
          localId,
          type: "tempo",
          content: `Tap tempo: ${bpm} BPM`,
          bpm,
          createdAt: new Date().toISOString(),
        });
      }

      setSaved(true);
      onCaptured?.();
    } catch (err) {
      console.error("Failed to save tempo:", err);
      
      try {
        await addPendingIdea({
          localId,
          type: "tempo",
          content: `Tap tempo: ${bpm} BPM`,
          bpm,
          createdAt: new Date().toISOString(),
        });
        setSaved(true);
        onCaptured?.();
      } catch {
        setError("Failed to save tempo");
      }
    } finally {
      setSaving(false);
    }
  }, [bpm, onCaptured]);

  const confidence = Math.min(100, Math.round((taps.length / TAP_TEMPO_MIN_TAPS) * 100));
  const hasEnoughTaps = taps.length >= TAP_TEMPO_MIN_TAPS;

  return (
    <div className="flex flex-col items-center justify-center h-full p-6">
      <div className="text-center mb-8">
        <div 
          className={`text-7xl font-mono font-bold mb-2 transition-all ${
            bpm ? "text-[var(--color-primary)]" : "text-[var(--color-text-muted)]"
          }`}
        >
          {bpm || "---"}
        </div>
        <div className="text-[var(--color-text-muted)] text-lg">BPM</div>
      </div>

      {/* Confidence bar */}
      <div className="w-48 h-2 bg-[var(--color-surface)] rounded-full mb-8 overflow-hidden">
        <div
          className={`h-full transition-all ${
            hasEnoughTaps ? "bg-[var(--color-success)]" : "bg-[var(--color-primary)]"
          }`}
          style={{ width: `${confidence}%` }}
        />
      </div>

      {/* Tap area */}
      <button
        onClick={handleTap}
        disabled={disabled || saving}
        className={`w-48 h-48 rounded-full flex items-center justify-center text-2xl font-bold transition-all active:scale-95 ${
          hasEnoughTaps
            ? "bg-[var(--color-success)] text-[var(--color-background)]"
            : "bg-[var(--color-surface)] text-[var(--color-text)]"
        }`}
      >
        {hasEnoughTaps ? "♪" : `Tap (${taps.length}/${TAP_TEMPO_MIN_TAPS})`}
      </button>

      <p className="mt-4 text-[var(--color-text-muted)] text-sm text-center">
        {hasEnoughTaps
          ? "Keep tapping to refine, or save the tempo"
          : "Tap the circle to the beat of your song"}
      </p>

      {error && (
        <p className="mt-2 text-[var(--color-error)] text-sm">{error}</p>
      )}

      {saved && (
        <p className="mt-2 text-[var(--color-success)] text-sm">
          ✓ Tempo saved!
        </p>
      )}

      {/* Actions */}
      <div className="flex gap-4 mt-8">
        <Button
          onClick={handleReset}
          variant="secondary"
          disabled={taps.length === 0 || saving}
        >
          Reset
        </Button>
        <Button
          onClick={handleSave}
          disabled={!hasEnoughTaps || saving || saved}
          loading={saving}
        >
          Save {bpm} BPM
        </Button>
      </div>
    </div>
  );
}
