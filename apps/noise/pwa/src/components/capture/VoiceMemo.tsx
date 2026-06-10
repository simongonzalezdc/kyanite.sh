"use client";

import { useState, useRef, useCallback } from "react";
import { Button } from "@/components/ui/Button";
import { getSyncClient } from "@/lib/api/client";
import { addPendingIdea, addPendingMedia } from "@/lib/db";
import { MEDIA_RECORDER_TIMESLICE_MS, RECORDING_TIMER_INTERVAL_MS } from "@/lib/constants";

interface VoiceMemoProps {
  onCaptured?: () => void;
  disabled?: boolean;
}

type RecordingState = "idle" | "recording" | "processing" | "done";

export function VoiceMemo({ onCaptured, disabled }: VoiceMemoProps) {
  const [state, setState] = useState<RecordingState>("idle");
  const [duration, setDuration] = useState(0);
  const [error, setError] = useState("");
  const [audioUrl, setAudioUrl] = useState<string | null>(null);
  
  const mediaRecorderRef = useRef<MediaRecorder | null>(null);
  const chunksRef = useRef<Blob[]>([]);
  const timerRef = useRef<ReturnType<typeof setInterval> | null>(null);
  const streamRef = useRef<MediaStream | null>(null);

  const startRecording = useCallback(async () => {
    try {
      setError("");
      
      const stream = await navigator.mediaDevices.getUserMedia({
        audio: {
          echoCancellation: true,
          noiseSuppression: true,
          autoGainControl: true,
        },
      });
      
      streamRef.current = stream;
      chunksRef.current = [];

      const mediaRecorder = new MediaRecorder(stream, {
        mimeType: MediaRecorder.isTypeSupported("audio/webm;codecs=opus")
          ? "audio/webm;codecs=opus"
          : "audio/webm",
      });

      mediaRecorder.ondataavailable = (e) => {
        if (e.data.size > 0) {
          chunksRef.current.push(e.data);
        }
      };

      mediaRecorder.onstop = () => {
        const blob = new Blob(chunksRef.current, { type: "audio/webm" });
        const url = URL.createObjectURL(blob);
        setAudioUrl(url);
        setState("done");
      };

      mediaRecorderRef.current = mediaRecorder;
      mediaRecorder.start(MEDIA_RECORDER_TIMESLICE_MS);
      setState("recording");
      setDuration(0);

      // Start duration timer
      timerRef.current = setInterval(() => {
        setDuration((d) => d + 1);
      }, RECORDING_TIMER_INTERVAL_MS);
    } catch (err) {
      console.error("Failed to start recording:", err);
      setError("Could not access microphone. Please grant permission.");
    }
  }, []);

  const stopRecording = useCallback(() => {
    if (timerRef.current) {
      clearInterval(timerRef.current);
      timerRef.current = null;
    }

    if (mediaRecorderRef.current && mediaRecorderRef.current.state !== "inactive") {
      mediaRecorderRef.current.stop();
    }

    if (streamRef.current) {
      streamRef.current.getTracks().forEach((track) => track.stop());
      streamRef.current = null;
    }

    setState("processing");
  }, []);

  const cancelRecording = useCallback(() => {
    if (timerRef.current) {
      clearInterval(timerRef.current);
      timerRef.current = null;
    }

    if (streamRef.current) {
      streamRef.current.getTracks().forEach((track) => track.stop());
      streamRef.current = null;
    }

    if (audioUrl) {
      URL.revokeObjectURL(audioUrl);
    }

    mediaRecorderRef.current = null;
    chunksRef.current = [];
    setAudioUrl(null);
    setDuration(0);
    setState("idle");
    setError("");
  }, [audioUrl]);

  const saveRecording = useCallback(async () => {
    if (!audioUrl || chunksRef.current.length === 0) return;

    setState("processing");
    const localId = `voice-${Date.now()}-${Math.random().toString(36).slice(2)}`;
    const blob = new Blob(chunksRef.current, { type: "audio/webm" });
    const filename = `memo-${Date.now()}.webm`;

    try {
      const client = getSyncClient();

      if (client) {
        // Upload to server
        const { path } = await client.uploadMedia(blob, filename, "voice");
        
        // Submit idea with media path
        await client.submitIdea({
          type: "voice_memo",
          content: `Voice memo (${formatDuration(duration)})`,
          media_path: path,
        });
      } else {
        // Queue for later sync
        await addPendingIdea({
          localId,
          type: "voice_memo",
          content: `Voice memo (${formatDuration(duration)})`,
          createdAt: new Date().toISOString(),
        });

        await addPendingMedia({
          localIdeaId: localId,
          blob,
          filename,
          type: "voice",
        });
      }

      cancelRecording();
      onCaptured?.();
    } catch (err) {
      console.error("Failed to save voice memo:", err);
      
      // Try to queue locally
      try {
        await addPendingIdea({
          localId,
          type: "voice_memo",
          content: `Voice memo (${formatDuration(duration)})`,
          createdAt: new Date().toISOString(),
        });

        await addPendingMedia({
          localIdeaId: localId,
          blob,
          filename,
          type: "voice",
        });

        cancelRecording();
        onCaptured?.();
      } catch {
        setError("Failed to save voice memo");
        setState("done");
      }
    }
  }, [audioUrl, duration, cancelRecording, onCaptured]);

  return (
    <div className="flex flex-col items-center justify-center h-full p-6">
      {state === "idle" && (
        <>
          <div className="text-6xl mb-6">🎤</div>
          <p className="text-[var(--color-text-muted)] mb-8 text-center">
            Tap to start recording a voice memo
          </p>
          <Button
            onClick={startRecording}
            disabled={disabled}
            size="lg"
            className="w-48 h-48 rounded-full text-2xl"
          >
            Record
          </Button>
          {error && (
            <p className="mt-4 text-[var(--color-error)] text-sm">{error}</p>
          )}
        </>
      )}

      {state === "recording" && (
        <>
          <div className="text-6xl mb-4 animate-pulse text-[var(--color-error)]">
            ●
          </div>
          <div className="text-3xl font-mono mb-6">
            {formatDuration(duration)}
          </div>
          <p className="text-[var(--color-text-muted)] mb-8">
            Recording...
          </p>
          <div className="flex gap-4">
            <Button
              onClick={cancelRecording}
              variant="secondary"
              size="lg"
            >
              Cancel
            </Button>
            <Button
              onClick={stopRecording}
              variant="danger"
              size="lg"
            >
              Stop
            </Button>
          </div>
        </>
      )}

      {state === "processing" && (
        <>
          <div className="text-6xl mb-6 animate-spin">◐</div>
          <p className="text-[var(--color-text-muted)]">Processing...</p>
        </>
      )}

      {state === "done" && audioUrl && (
        <>
          <div className="text-6xl mb-4">✓</div>
          <div className="text-xl font-mono mb-4">
            {formatDuration(duration)}
          </div>
          
          <audio src={audioUrl} controls className="mb-6 max-w-full" />
          
          {error && (
            <p className="mb-4 text-[var(--color-error)] text-sm">{error}</p>
          )}
          
          <div className="flex gap-4">
            <Button
              onClick={cancelRecording}
              variant="secondary"
            >
              Discard
            </Button>
            <Button onClick={saveRecording}>
              Save
            </Button>
          </div>
        </>
      )}
    </div>
  );
}

function formatDuration(seconds: number): string {
  const mins = Math.floor(seconds / 60);
  const secs = seconds % 60;
  return `${mins}:${secs.toString().padStart(2, "0")}`;
}
