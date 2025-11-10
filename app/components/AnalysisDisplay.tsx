'use client';

import { useState, useRef, useEffect } from 'react';
import { PitchPoint } from '@/lib/types';
import { PitchDetector } from '@/lib/audio/pitch-detector';
import { Music, Activity, Play, Square } from 'lucide-react';
import * as Tone from 'tone';

interface AnalysisDisplayProps {
  pitches: PitchPoint[];
  bpm?: number | null;
  musicalKey?: string | null;
}

export default function AnalysisDisplay({ pitches, bpm, musicalKey }: AnalysisDisplayProps) {
  const [isPlaying, setIsPlaying] = useState(false);
  const [currentNoteIndex, setCurrentNoteIndex] = useState<number | null>(null);
  const synthRef = useRef<Tone.Synth | null>(null);
  const sequenceRef = useRef<Tone.Sequence | null>(null);

  useEffect(() => {
    // Initialize synth
    synthRef.current = new Tone.Synth({
      oscillator: { type: 'sine' },
      envelope: { attack: 0.01, decay: 0.1, sustain: 0.5, release: 0.3 }
    }).toDestination();

    return () => {
      synthRef.current?.dispose();
      sequenceRef.current?.dispose();
    };
  }, []);

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

  const playNotes = async () => {
    if (!synthRef.current || simplified.length === 0) return;

    await Tone.start();
    setIsPlaying(true);
    setCurrentNoteIndex(0);

    // Convert ALL MIDI notes to note names (not just first 20)
    const allNotes = simplified.map(midi => detector.midiToNoteName(midi));
    
    // Calculate note durations based on pitch points timing
    const noteDurations: number[] = [];
    for (let i = 0; i < simplified.length; i++) {
      const currentMidi = simplified[i];
      let duration = 0.5; // Default duration
      
      if (i < simplified.length - 1) {
        // Find the time span for this note
        const nextMidi = simplified[i + 1];
        // Find first occurrence of each note in pitches array
        const currentPitchIndex = pitches.findIndex(p => Math.round(p.midi) === currentMidi);
        const nextPitchIndex = pitches.findIndex(p => Math.round(p.midi) === nextMidi);
        
        if (currentPitchIndex >= 0 && nextPitchIndex >= 0) {
          const timeDiff = pitches[nextPitchIndex].time - pitches[currentPitchIndex].time;
          duration = Math.max(0.2, Math.min(1.0, timeDiff));
        }
      } else {
        // Last note - use default duration
        duration = 0.5;
      }
      noteDurations.push(duration);
    }

    // Create sequence to play notes
    let noteIndex = 0;
    const playNextNote = () => {
      if (noteIndex >= allNotes.length) {
        setIsPlaying(false);
        setCurrentNoteIndex(null);
        return;
      }

      setCurrentNoteIndex(noteIndex);
      const note = allNotes[noteIndex];
      const duration = noteDurations[noteIndex];
      
      if (synthRef.current) {
        synthRef.current.triggerAttackRelease(note, duration);
      }

      noteIndex++;
      setTimeout(playNextNote, duration * 1000);
    };

    playNextNote();
  };

  const stopPlayback = () => {
    if (synthRef.current) {
      synthRef.current.triggerRelease();
    }
    if (sequenceRef.current) {
      sequenceRef.current.stop();
    }
    setIsPlaying(false);
    setCurrentNoteIndex(null);
  };

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
        <div className="flex items-center justify-between mb-3">
          <h3 className="font-medium">Detected Melody (first 20 notes)</h3>
          <div className="flex gap-2">
            {!isPlaying ? (
              <button
                onClick={playNotes}
                className="flex items-center gap-2 px-4 py-2 bg-primary-500 hover:bg-primary-600 rounded-lg font-medium transition-colors text-sm"
              >
                <Play size={16} />
                Play Notes
              </button>
            ) : (
              <button
                onClick={stopPlayback}
                className="flex items-center gap-2 px-4 py-2 bg-red-500 hover:bg-red-600 rounded-lg font-medium transition-colors text-sm"
              >
                <Square size={16} />
                Stop
              </button>
            )}
          </div>
        </div>
        <div className="flex flex-wrap gap-2">
          {noteNames.map((note, i) => (
            <span
              key={i}
              className={`px-3 py-1 border rounded-md text-sm font-mono transition-colors ${
                currentNoteIndex === i && isPlaying
                  ? 'bg-primary-500 border-primary-400 text-white scale-110'
                  : 'bg-primary-500/20 border-primary-500/50 text-gray-300'
              }`}
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
        {isPlaying && (
          <p className="text-sm text-gray-400 mt-2">
            Playing detected notes... ({currentNoteIndex !== null ? currentNoteIndex + 1 : 0} / {simplified.length})
            {currentNoteIndex !== null && currentNoteIndex < noteNames.length && (
              <span className="ml-2 text-primary-400">
                Now: {noteNames[currentNoteIndex]}
              </span>
            )}
          </p>
        )}
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

