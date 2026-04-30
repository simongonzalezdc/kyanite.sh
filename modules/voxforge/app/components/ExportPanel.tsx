'use client';

import { useState } from 'react';
import { Download } from 'lucide-react';
import { StemExporter } from '@/lib/audio/stem-exporter';
import { MidiExporter } from '@/lib/audio/midi-exporter';
import { MusicGenerator } from '@/lib/audio/music-generator';
import { PitchPoint, BPMAnalysis } from '@/lib/types';
import Tooltip from './Tooltip';
import { isMobile, triggerHaptic } from '@/lib/utils';

interface ExportPanelProps {
  audioBuffer: AudioBuffer | null;
  pitches: PitchPoint[];
  bpm: BPMAnalysis | null;
  musicalKey: string | null;
  generator: MusicGenerator | null;
}

export default function ExportPanel({ audioBuffer, pitches, bpm, musicalKey, generator }: ExportPanelProps) {
  const [exporting, setExporting] = useState<string | null>(null);

  const downloadBlob = (blob: Blob, filename: string) => {
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = filename;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  };

  const handleExportStem = async (type: 'vocal' | 'drums' | 'bass' | 'chords' | 'mix') => {
    if (!generator) {
      alert('Please generate music first');
      return;
    }

    setExporting(type);
    try {
      const exporter = new StemExporter(generator);
      const duration = audioBuffer?.duration || 30;
      const blob = await exporter.exportStem(type, audioBuffer || undefined, duration);
      const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
      downloadBlob(blob, `voxforge-${type}-${timestamp}.wav`);
    } catch (error) {
      console.error('Export failed:', error);
      alert('Failed to export stem. Please try again.');
    } finally {
      setExporting(null);
    }
  };

  const handleExportMIDI = () => {
    if (pitches.length === 0) {
      alert('No melody to export');
      return;
    }

    try {
      const exporter = new MidiExporter();
      const bpmValue = bpm?.bpm || 120;
      const blob = exporter.export(pitches, bpmValue, musicalKey || undefined);
      const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
      downloadBlob(blob, `voxforge-melody-${timestamp}.mid`);
    } catch (error) {
      console.error('MIDI export failed:', error);
      alert('Failed to export MIDI. Please try again.');
    }
  };

  const exportButtons = [
    { type: 'vocal' as const, label: 'Vocal' },
    { type: 'drums' as const, label: 'Drums' },
    { type: 'bass' as const, label: 'Bass' },
    { type: 'chords' as const, label: 'Chords' },
    { type: 'mix' as const, label: 'Full Mix' },
  ];

  return (
    <div className="bg-gray-900 rounded-xl p-6 sm:p-8 border border-gray-800 space-y-6" id="export">
      <h2 className="fluid-xl font-semibold">Export</h2>
      
      <div className="space-y-4">
        <h3 className="font-medium flex items-center gap-2">
          Stems
          <Tooltip content="Individual audio tracks for each instrument. Perfect for mixing in your favorite DAW or audio editor." />
        </h3>
        <div className="grid grid-cols-2 sm:grid-cols-3 gap-3">
          {exportButtons.map(({ type, label }) => (
            <button
              key={type}
              onClick={() => {
                handleExportStem(type);
                triggerHaptic('light');
              }}
              disabled={exporting === type || !generator}
              className={`
                flex items-center justify-center gap-2 px-4 py-3 sm:py-2
                bg-primary-500/20 border border-primary-500/50 rounded-lg
                hover:bg-primary-500/30 transition-all duration-200
                disabled:opacity-50 disabled:cursor-not-allowed
                min-h-[44px] min-w-[44px] sm:min-w-0
                ${isMobile() ? 'text-sm' : ''}
                active:scale-95 touch-manipulation
              `}
            >
              <Download size={isMobile() ? 18 : 16} />
              <span className="truncate">
                {exporting === type ? 'Exporting...' : label}
              </span>
            </button>
          ))}
        </div>
      </div>

      <div className="space-y-4">
        <h3 className="font-medium flex items-center gap-2">
          MIDI
          <Tooltip content="Standard MIDI file containing your melody and chord progression. Can be imported into any music software for further editing." />
        </h3>
        <button
          onClick={() => {
            handleExportMIDI();
            triggerHaptic('light');
          }}
          disabled={pitches.length === 0}
          className={`
            flex items-center justify-center gap-2 px-6 py-3 sm:py-2
            bg-secondary-500/20 border border-secondary-500/50 rounded-lg
            hover:bg-secondary-500/30 transition-all duration-200
            disabled:opacity-50 disabled:cursor-not-allowed
            min-h-[44px] min-w-[44px] sm:min-w-0
            ${isMobile() ? 'text-base' : ''}
            active:scale-95 touch-manipulation
          `}
        >
          <Download size={isMobile() ? 20 : 16} />
          <span>Export MIDI</span>
        </button>
      </div>

      {/* Mobile-specific export info */}
      {isMobile() && (
        <div className="bg-gray-800/50 rounded-lg p-4 border border-gray-700">
          <h4 className="font-medium text-sm mb-2">Export Tips</h4>
          <ul className="text-xs text-gray-400 space-y-1">
            <li>• Exported files will be saved to your device</li>
            <li>• For best results, export all stems individually</li>
            <li>• MIDI files are smaller and easier to share</li>
            <li>• Check your Downloads folder after export</li>
          </ul>
        </div>
      )}
    </div>
  );
}

