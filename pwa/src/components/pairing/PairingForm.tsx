"use client";

import { useState } from "react";
import { Button } from "@/components/ui/Button";
import { initSyncClient } from "@/lib/api/client";
import { setServerConfig } from "@/lib/db";

interface PairingFormProps {
  onPaired: (serverUrl: string, deviceId: string) => void;
}

export function PairingForm({ onPaired }: PairingFormProps) {
  const [serverUrl, setServerUrl] = useState("");
  const [pairingCode, setPairingCode] = useState("");
  const [deviceName, setDeviceName] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const [step, setStep] = useState<"url" | "code">("url");

  const handleUrlSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      // Normalize URL
      let url = serverUrl.trim();
      if (!url.startsWith("http://") && !url.startsWith("https://")) {
        url = `http://${url}`;
      }

      // Test connection
      const client = initSyncClient(url);
      const reachable = await client.ping();

      if (!reachable) {
        throw new Error("Could not reach server. Check the URL and make sure noise.sh is running with sync enabled.");
      }

      setServerUrl(url);
      setStep("code");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Connection failed");
    } finally {
      setLoading(false);
    }
  };

  const handleCodeSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      const client = initSyncClient(serverUrl);
      const name = deviceName.trim() || getDefaultDeviceName();
      const result = await client.pair(pairingCode, name);

      if (!result.success) {
        throw new Error(result.error || "Pairing failed");
      }

      // Save config to IndexedDB
      await setServerConfig({
        url: serverUrl,
        deviceId: result.device_id!,
        deviceName: name,
        pairedAt: new Date().toISOString(),
      });

      onPaired(serverUrl, result.device_id!);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Pairing failed");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="w-full max-w-md mx-auto p-6">
      <div className="text-center mb-8">
        <div className="text-5xl mb-4">♪</div>
        <h1 className="text-2xl font-bold text-[var(--color-primary)] mb-2">
          noise.sh
        </h1>
        <p className="text-[var(--color-text-muted)]">
          Connect to your desktop to start capturing ideas
        </p>
      </div>

      {step === "url" ? (
        <form onSubmit={handleUrlSubmit} className="space-y-4">
          <div>
            <label htmlFor="serverUrl" className="block text-sm mb-2 text-[var(--color-text-muted)]">
              Server URL
            </label>
            <input
              id="serverUrl"
              type="text"
              value={serverUrl}
              onChange={(e) => setServerUrl(e.target.value)}
              placeholder="192.168.1.100:8765"
              className="w-full"
              required
              autoFocus
              autoComplete="off"
            />
            <p className="text-xs text-[var(--color-text-muted)] mt-1">
              Find this in noise.sh → Settings → Sync
            </p>
          </div>

          {error && (
            <div className="p-3 bg-[var(--color-error)]/10 border border-[var(--color-error)]/30 rounded text-[var(--color-error)] text-sm">
              {error}
            </div>
          )}

          <Button type="submit" className="w-full" loading={loading}>
            Connect
          </Button>
        </form>
      ) : (
        <form onSubmit={handleCodeSubmit} className="space-y-4">
          <div>
            <label htmlFor="pairingCode" className="block text-sm mb-2 text-[var(--color-text-muted)]">
              Pairing Code
            </label>
            <input
              id="pairingCode"
              type="text"
              value={pairingCode}
              onChange={(e) => setPairingCode(e.target.value.replace(/\D/g, "").slice(0, 6))}
              placeholder="000000"
              className="w-full text-center text-2xl tracking-[0.5em] font-mono"
              maxLength={6}
              pattern="[0-9]{6}"
              required
              autoFocus
              autoComplete="off"
              inputMode="numeric"
            />
            <p className="text-xs text-[var(--color-text-muted)] mt-1">
              Generate a code in noise.sh → Settings → Sync → Generate Pairing Code
            </p>
          </div>

          <div>
            <label htmlFor="deviceName" className="block text-sm mb-2 text-[var(--color-text-muted)]">
              Device Name (optional)
            </label>
            <input
              id="deviceName"
              type="text"
              value={deviceName}
              onChange={(e) => setDeviceName(e.target.value)}
              placeholder={getDefaultDeviceName()}
              className="w-full"
              autoComplete="off"
            />
          </div>

          {error && (
            <div className="p-3 bg-[var(--color-error)]/10 border border-[var(--color-error)]/30 rounded text-[var(--color-error)] text-sm">
              {error}
            </div>
          )}

          <div className="flex gap-3">
            <Button
              type="button"
              variant="secondary"
              onClick={() => {
                setStep("url");
                setError("");
              }}
            >
              Back
            </Button>
            <Button type="submit" className="flex-1" loading={loading}>
              Pair Device
            </Button>
          </div>
        </form>
      )}

      <div className="mt-8 pt-6 border-t border-[var(--color-surface)] text-center">
        <p className="text-xs text-[var(--color-text-muted)]">
          Syncs over your local network only.
          <br />
          No data is sent to the cloud.
        </p>
      </div>
    </div>
  );
}

function getDefaultDeviceName(): string {
  if (typeof navigator === "undefined") return "Mobile Device";
  
  const ua = navigator.userAgent;
  if (/iPhone/i.test(ua)) return "iPhone";
  if (/iPad/i.test(ua)) return "iPad";
  if (/Android/i.test(ua)) return "Android";
  return "Mobile Device";
}
