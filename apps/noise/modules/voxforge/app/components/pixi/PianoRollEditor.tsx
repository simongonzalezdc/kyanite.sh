'use client';

import { useEffect, useRef, useState } from 'react';
import { InteractivePianoRoll } from '@/lib/pixi/piano-roll';
import { PitchPoint } from '@/lib/types';
import { Download, Trash2, Grid3x3 } from 'lucide-react';

interface PianoRollEditorProps {
  initialPitches?: PitchPoint[];
  bpm: number;
  onExport?: (pitches: PitchPoint[]) => void;
}

export default function PianoRollEditor({ 
  initialPitches, 
  bpm,
  onExport 
}: PianoRollEditorProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const pianoRollRef = useRef<InteractivePianoRoll | null>(null);
  const [snapToGrid, setSnapToGrid] = useState(true);

  useEffect(() => {
    if (!containerRef.current) return;

    // Create piano roll
    const pianoRoll = new InteractivePianoRoll({
      width: containerRef.current.clientWidth,
      height: 600,
      snapToGrid
    });

    // Wait for initialization, then add canvas to DOM
    pianoRoll.waitForInit().then(() => {
      const canvas = pianoRoll.getCanvas();
      canvas.style.display = 'block';
      canvas.style.width = '100%';
      canvas.style.height = '100%';
      containerRef.current?.appendChild(canvas);
      pianoRollRef.current = pianoRoll;

      // Load initial notes
      if (initialPitches) {
        pianoRoll.loadNotes(initialPitches);
      }
    }).catch(console.error);

    // Handle resize
    const handleResize = () => {
      if (containerRef.current) {
        pianoRoll.resize(containerRef.current.clientWidth, 600);
      }
    };

    window.addEventListener('resize', handleResize);

    return () => {
      window.removeEventListener('resize', handleResize);
      pianoRoll.destroy();
    };
  }, [snapToGrid, initialPitches]);

  const handleExport = () => {
    if (!pianoRollRef.current) return;
    const pitches = pianoRollRef.current.exportToPitches();
    onExport?.(pitches);
  };

  const handleExportMIDI = () => {
    if (!pianoRollRef.current) return;
    const midi = pianoRollRef.current.exportToMIDI(bpm);
    const midiArray = midi.toArray();
    const blob = new Blob([midiArray.buffer as ArrayBuffer], { type: 'audio/midi' });
    const url = URL.createObjectURL(blob);
    
    const a = document.createElement('a');
    a.href = url;
    a.download = `voxforge-edited-${Date.now()}.mid`;
    a.click();
    URL.revokeObjectURL(url);
  };

  const handleClear = () => {
    if (!pianoRollRef.current) return;
    if (confirm('Clear all notes?')) {
      pianoRollRef.current.clear();
    }
  };

  const toggleSnapToGrid = () => {
    setSnapToGrid(!snapToGrid);
  };

  return (
    <div className="space-y-4">
      {/* Toolbar */}
      <div className="flex items-center justify-between bg-gray-900 p-4 rounded-lg border border-gray-800">
        <div className="flex items-center gap-4">
          <h3 className="text-lg font-semibold">Piano Roll Editor</h3>
          <button
            onClick={toggleSnapToGrid}
            className={`flex items-center gap-2 px-3 py-2 rounded-lg transition-colors ${
              snapToGrid 
                ? 'bg-primary-500 text-white' 
                : 'bg-gray-800 text-gray-400'
            }`}
          >
            <Grid3x3 size={16} />
            Snap to Grid
          </button>
        </div>

        <div className="flex gap-2">
          <button
            onClick={handleClear}
            className="flex items-center gap-2 px-4 py-2 bg-red-500/20 hover:bg-red-500/30 text-red-400 rounded-lg transition-colors"
          >
            <Trash2 size={16} />
            Clear
          </button>
          
          <button
            onClick={handleExport}
            className="flex items-center gap-2 px-4 py-2 bg-primary-500 hover:bg-primary-600 rounded-lg transition-colors"
          >
            <Download size={16} />
            Apply Changes
          </button>
          
          <button
            onClick={handleExportMIDI}
            className="flex items-center gap-2 px-4 py-2 bg-secondary-500 hover:bg-secondary-600 rounded-lg transition-colors"
          >
            <Download size={16} />
            Export MIDI
          </button>
        </div>
      </div>

      {/* Instructions */}
      <div className="bg-gray-900/50 p-4 rounded-lg border border-gray-800">
        <p className="text-sm text-gray-400">
          <strong>Controls:</strong> Click to add notes • Drag to move • Right-click to delete • Scroll to zoom
        </p>
      </div>

      {/* Piano Roll Canvas */}
      <div 
        ref={containerRef}
        className="rounded-lg overflow-hidden border border-gray-800 relative"
        style={{ minHeight: '600px', maxHeight: '600px' }}
      />
    </div>
  );
}

