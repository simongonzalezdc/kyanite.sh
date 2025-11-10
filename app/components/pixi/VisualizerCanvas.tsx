'use client';

import { useEffect, useRef } from 'react';
import { MusicVisualizer } from '@/lib/pixi/visualizer';
import { PitchPoint } from '@/lib/types';

interface VisualizerCanvasProps {
  pitches?: PitchPoint[];
  isRecording?: boolean;
  onVisualizerReady?: (visualizer: MusicVisualizer) => void;
}

export default function VisualizerCanvas({ 
  pitches, 
  isRecording,
  onVisualizerReady 
}: VisualizerCanvasProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const visualizerRef = useRef<MusicVisualizer | null>(null);

  useEffect(() => {
    if (!containerRef.current) return;

    // Create visualizer
    const visualizer = new MusicVisualizer({
      width: containerRef.current.clientWidth,
      height: 400
    });

    // Wait for initialization, then add canvas to DOM
    visualizer.waitForInit().then(() => {
      const canvas = visualizer.getCanvas();
      canvas.style.display = 'block';
      canvas.style.width = '100%';
      canvas.style.height = '100%';
      containerRef.current?.appendChild(canvas);
      visualizerRef.current = visualizer;

      // Notify parent
      onVisualizerReady?.(visualizer);
    }).catch(console.error);

    // Handle resize
    const handleResize = () => {
      if (containerRef.current) {
        visualizer.resize(
          containerRef.current.clientWidth,
          400
        );
      }
    };

    window.addEventListener('resize', handleResize);

    return () => {
      window.removeEventListener('resize', handleResize);
      visualizer.destroy();
    };
  }, []);

  // Load pitches when they change
  useEffect(() => {
    if (visualizerRef.current && pitches && !isRecording) {
      visualizerRef.current.loadNotes(pitches);
    }
  }, [pitches, isRecording]);

  // Auto-scroll during recording
  useEffect(() => {
    if (!visualizerRef.current || !isRecording) return;

    let animationFrame: number;
    let lastTime = performance.now();

    const scroll = (currentTime: number) => {
      const deltaTime = (currentTime - lastTime) / 1000; // Convert to seconds
      lastTime = currentTime;

      visualizerRef.current?.scroll(deltaTime);
      animationFrame = requestAnimationFrame(scroll);
    };

    animationFrame = requestAnimationFrame(scroll);

    return () => {
      cancelAnimationFrame(animationFrame);
    };
  }, [isRecording]);

  return (
    <div 
      ref={containerRef} 
      className="w-full rounded-lg overflow-hidden border border-gray-800 relative"
      style={{ minHeight: '400px', maxHeight: '400px' }}
    />
  );
}

