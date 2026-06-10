"use client";

import { useState, useRef, useCallback } from "react";
import { Button } from "@/components/ui/Button";
import { getSyncClient } from "@/lib/api/client";
import { addPendingIdea, addPendingMedia } from "@/lib/db";
import { CAMERA_WIDTH, CAMERA_HEIGHT, JPEG_QUALITY } from "@/lib/constants";

interface PhotoCaptureProps {
  onCaptured?: () => void;
  disabled?: boolean;
}

type CaptureState = "idle" | "camera" | "preview" | "processing";

export function PhotoCapture({ onCaptured, disabled }: PhotoCaptureProps) {
  const [state, setState] = useState<CaptureState>("idle");
  const [error, setError] = useState("");
  const [imageUrl, setImageUrl] = useState<string | null>(null);
  const [imageBlob, setImageBlob] = useState<Blob | null>(null);
  
  const videoRef = useRef<HTMLVideoElement>(null);
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const streamRef = useRef<MediaStream | null>(null);

  const startCamera = useCallback(async () => {
    try {
      setError("");
      
      const stream = await navigator.mediaDevices.getUserMedia({
        video: {
          facingMode: "environment", // Prefer back camera
          width: { ideal: CAMERA_WIDTH },
          height: { ideal: CAMERA_HEIGHT },
        },
      });
      
      streamRef.current = stream;
      
      if (videoRef.current) {
        videoRef.current.srcObject = stream;
        await videoRef.current.play();
      }
      
      setState("camera");
    } catch (err) {
      console.error("Failed to start camera:", err);
      setError("Could not access camera. Please grant permission.");
    }
  }, []);

  const stopCamera = useCallback(() => {
    if (streamRef.current) {
      streamRef.current.getTracks().forEach((track) => track.stop());
      streamRef.current = null;
    }
    
    if (videoRef.current) {
      videoRef.current.srcObject = null;
    }
  }, []);

  const capturePhoto = useCallback(() => {
    if (!videoRef.current || !canvasRef.current) return;

    const video = videoRef.current;
    const canvas = canvasRef.current;
    
    canvas.width = video.videoWidth;
    canvas.height = video.videoHeight;
    
    const ctx = canvas.getContext("2d");
    if (!ctx) return;
    
    ctx.drawImage(video, 0, 0);
    
    canvas.toBlob(
      (blob) => {
        if (blob) {
          const url = URL.createObjectURL(blob);
          setImageUrl(url);
          setImageBlob(blob);
          setState("preview");
          stopCamera();
        }
      },
      "image/jpeg",
      JPEG_QUALITY
    );
  }, [stopCamera]);

  const cancelCapture = useCallback(() => {
    stopCamera();
    
    if (imageUrl) {
      URL.revokeObjectURL(imageUrl);
    }
    
    setImageUrl(null);
    setImageBlob(null);
    setState("idle");
    setError("");
  }, [imageUrl, stopCamera]);

  const savePhoto = useCallback(async () => {
    if (!imageBlob) return;

    setState("processing");
    const localId = `photo-${Date.now()}-${Math.random().toString(36).slice(2)}`;
    const filename = `photo-${Date.now()}.jpg`;

    try {
      const client = getSyncClient();

      if (client) {
        // Upload to server
        const { path } = await client.uploadMedia(imageBlob, filename, "photo");
        
        // Submit idea with media path
        await client.submitIdea({
          type: "photo",
          content: "Photo capture",
          media_path: path,
        });
      } else {
        // Queue for later sync
        await addPendingIdea({
          localId,
          type: "photo",
          content: "Photo capture",
          createdAt: new Date().toISOString(),
        });

        await addPendingMedia({
          localIdeaId: localId,
          blob: imageBlob,
          filename,
          type: "photo",
        });
      }

      cancelCapture();
      onCaptured?.();
    } catch (err) {
      console.error("Failed to save photo:", err);
      
      // Try to queue locally
      try {
        await addPendingIdea({
          localId,
          type: "photo",
          content: "Photo capture",
          createdAt: new Date().toISOString(),
        });

        await addPendingMedia({
          localIdeaId: localId,
          blob: imageBlob,
          filename,
          type: "photo",
        });

        cancelCapture();
        onCaptured?.();
      } catch {
        setError("Failed to save photo");
        setState("preview");
      }
    }
  }, [imageBlob, cancelCapture, onCaptured]);

  return (
    <div className="flex flex-col items-center justify-center h-full p-6">
      {/* Hidden canvas for capturing */}
      <canvas ref={canvasRef} className="hidden" />
      
      {state === "idle" && (
        <>
          <div className="text-6xl mb-6">📷</div>
          <p className="text-[var(--color-text-muted)] mb-8 text-center">
            Capture a photo of lyrics, setlists, or inspiration
          </p>
          <Button
            onClick={startCamera}
            disabled={disabled}
            size="lg"
          >
            Open Camera
          </Button>
          {error && (
            <p className="mt-4 text-[var(--color-error)] text-sm">{error}</p>
          )}
        </>
      )}

      {state === "camera" && (
        <div className="w-full max-w-md">
          <div className="relative aspect-[4/3] bg-black rounded-lg overflow-hidden mb-4">
            <video
              ref={videoRef}
              autoPlay
              playsInline
              muted
              className="w-full h-full object-cover"
            />
          </div>
          
          <div className="flex gap-4 justify-center">
            <Button
              onClick={cancelCapture}
              variant="secondary"
            >
              Cancel
            </Button>
            <Button
              onClick={capturePhoto}
              size="lg"
              className="w-20 h-20 rounded-full"
            >
              📷
            </Button>
          </div>
        </div>
      )}

      {state === "preview" && imageUrl && (
        <div className="w-full max-w-md">
          <div className="relative aspect-[4/3] bg-black rounded-lg overflow-hidden mb-4">
            {/* eslint-disable-next-line @next/next/no-img-element */}
            <img
              src={imageUrl}
              alt="Captured photo"
              className="w-full h-full object-cover"
            />
          </div>
          
          {error && (
            <p className="mb-4 text-center text-[var(--color-error)] text-sm">{error}</p>
          )}
          
          <div className="flex gap-4 justify-center">
            <Button
              onClick={cancelCapture}
              variant="secondary"
            >
              Retake
            </Button>
            <Button onClick={savePhoto}>
              Save
            </Button>
          </div>
        </div>
      )}

      {state === "processing" && (
        <>
          <div className="text-6xl mb-6 animate-spin">◐</div>
          <p className="text-[var(--color-text-muted)]">Saving...</p>
        </>
      )}
    </div>
  );
}
