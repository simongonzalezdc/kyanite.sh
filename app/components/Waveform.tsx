'use client';

import { useEffect, useRef } from 'react';

interface WaveformProps {
  isRecording: boolean;
  getWaveformData?: () => Uint8Array;
  audioBuffer?: AudioBuffer | null;
}

export default function Waveform({ isRecording, getWaveformData, audioBuffer }: WaveformProps) {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const animationRef = useRef<number>();

  useEffect(() => {
    if (!canvasRef.current) return;

    const canvas = canvasRef.current;
    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    // Set canvas size
    canvas.width = canvas.offsetWidth * 2; // Retina
    canvas.height = canvas.offsetHeight * 2;
    ctx.scale(2, 2);

    if (isRecording && getWaveformData) {
      // Real-time visualization
      const draw = () => {
        const data = getWaveformData();
        const width = canvas.offsetWidth;
        const height = canvas.offsetHeight;

        // Clear canvas
        ctx.fillStyle = '#0a0a0a';
        ctx.fillRect(0, 0, width, height);

        // Draw waveform
        ctx.lineWidth = 2;
        ctx.strokeStyle = '#3B82F6';
        ctx.beginPath();

        const sliceWidth = width / data.length;
        let x = 0;

        for (let i = 0; i < data.length; i++) {
          const v = data[i] / 128.0; // Normalize to 0-2
          const y = (v * height) / 2;

          if (i === 0) {
            ctx.moveTo(x, y);
          } else {
            ctx.lineTo(x, y);
          }

          x += sliceWidth;
        }

        ctx.stroke();

        animationRef.current = requestAnimationFrame(draw);
      };

      draw();

      return () => {
        if (animationRef.current) {
          cancelAnimationFrame(animationRef.current);
        }
      };
    } else if (audioBuffer) {
      // Static visualization of recorded audio
      const data = audioBuffer.getChannelData(0);
      const width = canvas.offsetWidth;
      const height = canvas.offsetHeight;

      ctx.fillStyle = '#0a0a0a';
      ctx.fillRect(0, 0, width, height);

      ctx.lineWidth = 1;
      ctx.strokeStyle = '#8B5CF6';
      ctx.beginPath();

      const step = Math.ceil(data.length / width);
      const amp = height / 2;

      for (let i = 0; i < width; i++) {
        let min = 1.0;
        let max = -1.0;

        for (let j = 0; j < step; j++) {
          const datum = data[(i * step) + j];
          if (datum < min) min = datum;
          if (datum > max) max = datum;
        }

        ctx.moveTo(i, (1 + min) * amp);
        ctx.lineTo(i, (1 + max) * amp);
      }

      ctx.stroke();
    }
  }, [isRecording, getWaveformData, audioBuffer]);

  return (
    <canvas
      ref={canvasRef}
      className="w-full h-32 rounded-lg bg-black/50 border border-gray-800"
    />
  );
}

