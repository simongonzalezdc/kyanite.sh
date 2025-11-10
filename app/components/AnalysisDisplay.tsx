'use client';

import { PitchPoint } from '@/lib/types';
import { PitchDetector } from '@/lib/audio/pitch-detector';
import { Music, Activity } from 'lucide-react';

interface AnalysisDisplayProps {
  pitches: PitchPoint[];
  bpm?: number | null;
  musicalKey?: string | null;
}

export default function AnalysisDisplay({ pitches, bpm, musicalKey }: AnalysisDisplayProps) {
  if (pitches.length === 0) {
    return (
      <div className="bg-gray-900 rounded-xl p-8 border border-gray-800">
        <p className="text-gray-400 text-center">No analysis data yet</p>
      </div>
    );
  }

  const detector = new PitchDetector();
  const stats = detector.getStats(pitches);
  const simplified = detector.getSimplifiedMelody(pitches);
  const noteNames = simplified.slice(0, 20).map(midi => detector.midiToNoteName(midi));

  return (
    <div className="bg-gray-900 rounded-xl p-8 border border-gray-800 space-y-6">
      <h2 className="text-xl font-semibold flex items-center gap-2">
        <Activity size={24} className="text-primary-500" />
        Audio Analysis
      </h2>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        {/* Pitch Info */}
        <div className="bg-black/30 rounded-lg p-4 border border-gray-800">
          <div className="flex items-center gap-2 mb-2">
            <Music size={20} className="text-secondary-500" />
            <h3 className="font-medium">Pitch Range</h3>
          </div>
          <p className="text-2xl font-bold text-secondary-500 capitalize">{stats.range}</p>
          <p className="text-sm text-gray-400 mt-1">
            {stats.averageFrequency.toFixed(1)} Hz average
          </p>
        </div>

        {/* BPM */}
        <div className="bg-black/30 rounded-lg p-4 border border-gray-800">
          <h3 className="font-medium mb-2">Tempo</h3>
          <p className="text-2xl font-bold text-primary-500">
            {bpm ? `${bpm} BPM` : 'Detecting...'}
          </p>
          <p className="text-sm text-gray-400 mt-1">
            {bpm ? 'Auto-detected' : 'Analysis pending'}
          </p>
        </div>

        {/* Key */}
        <div className="bg-black/30 rounded-lg p-4 border border-gray-800">
          <h3 className="font-medium mb-2">Musical Key</h3>
          <p className="text-2xl font-bold text-primary-500">
            {musicalKey || 'Detecting...'}
          </p>
          <p className="text-sm text-gray-400 mt-1">
            {musicalKey ? 'Auto-detected' : 'Analysis pending'}
          </p>
        </div>
      </div>

      {/* Detected Notes */}
      <div>
        <h3 className="font-medium mb-3">Detected Melody (first 20 notes)</h3>
        <div className="flex flex-wrap gap-2">
          {noteNames.map((note, i) => (
            <span
              key={i}
              className="px-3 py-1 bg-primary-500/20 border border-primary-500/50 rounded-md text-sm font-mono"
            >
              {note}
            </span>
          ))}
          {simplified.length > 20 && (
            <span className="px-3 py-1 text-gray-400 text-sm">
              +{simplified.length - 20} more
            </span>
          )}
        </div>
      </div>

      {/* Stats */}
      <div className="text-sm text-gray-400 space-y-1">
        <p>Total pitch points detected: {pitches.length}</p>
        <p>Unique notes: {simplified.length}</p>
        <p>Frequency range: {stats.minFrequency.toFixed(1)} - {stats.maxFrequency.toFixed(1)} Hz</p>
      </div>
    </div>
  );
}

