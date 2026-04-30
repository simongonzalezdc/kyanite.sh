'use client';

import { useEffect, useRef, useState, useCallback } from 'react';
import { PitchPoint } from '@/lib/types';

interface VisualizerCanvasProps {
  pitches?: PitchPoint[];
  isRecording?: boolean;
  onVisualizerReady?: (visualizer: any) => void;
}

export default function VisualizerCanvas({ 
  pitches, 
  isRecording,
  onVisualizerReady 
}: VisualizerCanvasProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const canvasRef = useRef<HTMLCanvasElement | null>(null);
  const [isInitialized, setIsInitialized] = useState(false);
  const [initError, setInitError] = useState<string | null>(null);
  const [showVisualizer, setShowVisualizer] = useState(true);

  // Simple canvas-based visualizer that doesn't depend on PixiJS
  const createSimpleVisualizer = useCallback(() => {
    if (!containerRef.current) return null;

    try {
      // Create a simple canvas visualizer
      const canvas = document.createElement('canvas');
      const ctx = canvas.getContext('2d');
      
      if (!ctx) {
        throw new Error('Could not get 2D context');
      }

      // Set canvas properties
      canvas.width = containerRef.current.clientWidth;
      canvas.height = 400;
      canvas.style.width = '100%';
      canvas.style.height = '100%';
      canvas.style.display = 'block';

      // Clear existing content
      while (containerRef.current.firstChild) {
        containerRef.current.removeChild(containerRef.current.firstChild);
      }
      
      // Add canvas to container
      containerRef.current.appendChild(canvas);
      canvasRef.current = canvas;

      // Return simple API
      return {
        loadNotes: (notes: PitchPoint[]) => {
          if (!ctx || !canvasRef.current) return;
          
          // Clear canvas
          ctx.clearRect(0, 0, canvasRef.current.width, canvasRef.current.height);
          
          // Draw grid
          ctx.strokeStyle = '#333333';
          ctx.lineWidth = 1;
          
          // Horizontal lines (octaves)
          for (let i = 0; i <= 4; i++) {
            const y = (i * canvasRef.current.height) / 4;
            ctx.beginPath();
            ctx.moveTo(0, y);
            ctx.lineTo(canvasRef.current.width, y);
            ctx.stroke();
          }
          
          // Vertical lines (beats)
          for (let i = 0; i <= 8; i++) {
            const x = (i * canvasRef.current.width) / 8;
            ctx.beginPath();
            ctx.moveTo(x, 0);
            ctx.lineTo(x, canvasRef.current.height);
            ctx.stroke();
          }
          
          // Draw notes
          notes.forEach((note, index) => {
            const x = (index * canvasRef.current!.width) / Math.max(notes.length, 1);
            const y = ((84 - note.midi) / 48) * canvasRef.current!.height;
            
            ctx.fillStyle = '#3B82F6';
            ctx.fillRect(x - 2, y - 2, 4, 4);
          });
        },
        
        scroll: (deltaTime: number) => {
          // Simple scroll implementation
          if (canvasRef.current) {
            // Redraw if needed during recording
            if (isRecording && pitches) {
              // You can implement scrolling logic here
            }
          }
        },
        
        clear: () => {
          if (ctx && canvasRef.current) {
            ctx.clearRect(0, 0, canvasRef.current.width, canvasRef.current.height);
          }
        },
        
        resize: (width: number, height: number) => {
          if (canvasRef.current) {
            canvasRef.current.width = width;
            canvasRef.current.height = height;
          }
        },
        
        destroy: () => {
          if (canvasRef.current && containerRef.current) {
            containerRef.current.removeChild(canvasRef.current);
            canvasRef.current = null;
          }
        }
      };
    } catch (error) {
      console.error('Failed to create simple visualizer:', error);
      return null;
    }
  }, [isRecording, pitches]);

  useEffect(() => {
    if (!containerRef.current || !showVisualizer) return;

    console.log('Initializing simple visualizer...');
    
    try {
      const visualizer = createSimpleVisualizer();
      
      if (visualizer) {
        setIsInitialized(true);
        setInitError(null);
        onVisualizerReady?.(visualizer);
        console.log('Simple visualizer initialized successfully');
      } else {
        setInitError('Failed to create visualizer');
      }
    } catch (error) {
      console.error('Visualizer initialization failed:', error);
      setInitError(error instanceof Error ? error.message : 'Unknown error');
    }
  }, [createSimpleVisualizer, onVisualizerReady, showVisualizer]);

  // Handle resize
  useEffect(() => {
    const handleResize = () => {
      if (canvasRef.current && containerRef.current) {
        canvasRef.current.width = containerRef.current.clientWidth;
        // Redraw if needed
        if (pitches && !isRecording) {
          // Trigger redraw
        }
      }
    };

    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, [pitches, isRecording]);

  // Load pitches when they change
  useEffect(() => {
    if (isInitialized && pitches && !isRecording && canvasRef.current) {
      // This would be handled by the visualizer.loadNotes method
    }
  }, [pitches, isRecording, isInitialized]);

  // Auto-scroll during recording
  useEffect(() => {
    if (!isRecording || !isInitialized) return;

    let animationFrame: number;
    let lastTime = performance.now();

    const scroll = (currentTime: number) => {
      if (!isRecording || !isInitialized) return;
      
      const deltaTime = (currentTime - lastTime) / 1000;
      lastTime = currentTime;

      // Auto-scroll logic would go here
      animationFrame = requestAnimationFrame(scroll);
    };

    animationFrame = requestAnimationFrame(scroll);
    return () => cancelAnimationFrame(animationFrame);
  }, [isRecording, isInitialized]);

  // Disable visualizer to fix issues
  useEffect(() => {
    const timer = setTimeout(() => {
      console.log('Disabling visualizer to fix audio issues');
      setShowVisualizer(false);
    }, 1000);

    return () => clearTimeout(timer);
  }, []);

  // Show error state
  if (initError) {
    return (
      <div 
        ref={containerRef} 
        className="w-full rounded-lg overflow-hidden border border-gray-800 relative flex items-center justify-center bg-red-900/20"
        style={{ minHeight: '400px', maxHeight: '400px' }}
      >
        <div className="text-center p-4">
          <p className="text-red-400 mb-2">Visualizer Disabled</p>
          <p className="text-gray-400 text-sm">Focusing on audio functionality</p>
        </div>
      </div>
    );
  }

  // Show disabled state
  if (!showVisualizer) {
    return (
      <div 
        ref={containerRef} 
        className="w-full rounded-lg overflow-hidden border border-gray-800 relative flex items-center justify-center bg-gray-900/50"
        style={{ minHeight: '400px', maxHeight: '400px' }}
      >
        <div className="text-center p-4">
          <p className="text-gray-400 mb-2">Audio Recording Active</p>
          <p className="text-gray-500 text-sm">Visualization disabled for better audio performance</p>
        </div>
      </div>
    );
  }

  // Show loading state
  if (!isInitialized) {
    return (
      <div 
        ref={containerRef} 
        className="w-full rounded-lg overflow-hidden border border-gray-800 relative flex items-center justify-center bg-gray-900/50"
        style={{ minHeight: '400px', maxHeight: '400px' }}
      >
        <div className="text-center p-4">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-500 mx-auto mb-2"></div>
          <p className="text-gray-400">Initializing audio system...</p>
        </div>
      </div>
    );
  }

  // Main visualizer (minimal implementation)
  return (
    <div 
      ref={containerRef} 
      className="w-full rounded-lg overflow-hidden border border-gray-800 relative"
      style={{ minHeight: '400px', maxHeight: '400px' }}
    >
      <canvas
        ref={canvasRef}
        className="w-full h-full block"
        style={{ display: 'block' }}
      />
    </div>
  );
}
