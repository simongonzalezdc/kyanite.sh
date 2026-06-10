"use client";

import { useState, useRef, useEffect } from "react";
import { Button } from "@/components/ui/Button";
import { getSyncClient } from "@/lib/api/client";
import { addPendingIdea } from "@/lib/db";
import { useToastActions } from "@/components/ui/Toast";

interface TextCaptureProps {
  onCaptured?: () => void;
  disabled?: boolean;
}

export function TextCapture({ onCaptured, disabled }: TextCaptureProps) {
  const [text, setText] = useState("");
  const [sending, setSending] = useState(false);
  const textareaRef = useRef<HTMLTextAreaElement>(null);
  const toast = useToastActions();

  // Auto-resize textarea
  useEffect(() => {
    if (textareaRef.current) {
      textareaRef.current.style.height = "auto";
      textareaRef.current.style.height = `${textareaRef.current.scrollHeight}px`;
    }
  }, [text]);

  const handleSubmit = async () => {
    const content = text.trim();
    if (!content) return;

    setSending(true);

    const localId = `text-${Date.now()}-${Math.random().toString(36).slice(2)}`;

    try {
      const client = getSyncClient();

      if (client) {
        // Try to send directly to server
        await client.submitIdea({
          type: "text",
          content,
        });
        toast.success("Idea captured!");
      } else {
        // Queue for later sync
        await addPendingIdea({
          localId,
          type: "text",
          content,
          createdAt: new Date().toISOString(),
        });
        toast.info("Saved offline - will sync when connected");
      }

      setText("");
      onCaptured?.();
    } catch (err) {
      console.error("Failed to submit text idea:", err);

      // Save to offline queue on failure
      try {
        await addPendingIdea({
          localId,
          type: "text",
          content,
          createdAt: new Date().toISOString(),
        });
        setText("");
        toast.info("Saved offline - will sync when connected");
        onCaptured?.();
      } catch {
        toast.error("Failed to save idea. Please try again.");
      }
    } finally {
      setSending(false);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    // Submit on Cmd/Ctrl+Enter
    if ((e.metaKey || e.ctrlKey) && e.key === "Enter") {
      e.preventDefault();
      handleSubmit();
    }
  };

  return (
    <div className="flex flex-col h-full">
      <div className="flex-1 p-4">
        <textarea
          ref={textareaRef}
          value={text}
          onChange={(e) => setText(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder="Type your lyric idea, melody description, or quick note..."
          className="w-full min-h-[120px] max-h-[60vh] resize-none bg-transparent border-none outline-none text-lg placeholder:text-[var(--color-text-muted)]"
          disabled={disabled || sending}
          autoFocus
        />
      </div>

      <div className="p-4 border-t border-[var(--color-surface)] flex items-center justify-between">
        <span className="text-sm text-[var(--color-text-muted)]">
          {text.length > 0 && `${text.length} characters`}
        </span>
        <Button
          onClick={handleSubmit}
          disabled={!text.trim() || disabled || sending}
          loading={sending}
        >
          Capture
        </Button>
      </div>
    </div>
  );
}
