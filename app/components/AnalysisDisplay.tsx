'use client';

import { useState, useRef, useEffect } from 'react';
import { PitchPoint, TimeSignature } from '@/lib/types';
import { PitchDetector } from '@/lib/audio/pitch-detector';
import { TimeSignatureDetector } from '@/lib/audio/time-signature-detector';
import { Music, Activity, Play, Square, Edit2, Check, X } from 'lucide-react';
import * as Tone from 'tone';

interface AnalysisDisplayProps {
  pitches: PitchPoint[];
  bpm?: number | null;
  musicalKey?: string | null;
  timeSignature?: TimeSignature | null;
  onBPMChange?: (bpm: number) => void;
  onKeyChange?: (key: string) => void;
  onTimeSignatureChange?: (ts: TimeSignature) => void;
}

export default function AnalysisDisplay({ 
  pitches, 
  bpm, 
  musicalKey, 
  timeSignature,
  onBPMChange,
  onKeyChange,
  onTimeSignatureChange
}: AnalysisDisplayProps) {
  const [isPlaying, setIsPlaying] = useState(false);
  const [currentNoteIndex, setCurrentNoteIndex] = useState<number | null>(null);
  const [editingBPM, setEditingBPM] = useState(false);
  const [editingKey, setEditingKey] = useState(false);
  const [editingTimeSig, setEditingTimeSig] = useState(false);
  const [tempBPM, setTempBPM] = useState(bpm?.toString() || '120');
  const [tempKey, setTempKey] = useState(musicalKey || 'C Major');
  const [tempTimeSig, setTempTimeSig] = useState(timeSignature?.display || '4/4');
  const synthRef = useRef<Tone.Synth | null>(null);
  const sequenceRef = useRef<Tone.Sequence | null>(null);
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);
  const isPlayingRef = useRef<boolean>(false);

  useEffect(() => {
    // Initialize synth
    synthRef.current = new Tone.Synth({
      oscillator: { type: 'sine' },
      envelope: { attack: 0.01, decay: 0.1, sustain: 0.5, release: 0.3 }
    }).toDestination();

    return () => {
      // Clean up on unmount
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
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
    isPlayingRef.current = true;
    setCurrentNoteIndex(0);

    // Convert ALL MIDI notes to note names (not just first 20)
    const allNotes = simplified.map(midi => detector.midiToNoteName(midi));
    
    // Calculate note durations based on pitch points timing
    // Track sequential pitch point indices to handle repeated MIDI values correctly
    const noteDurations: number[] = [];
    let currentPitchIndex = 0;
    
    for (let i = 0; i < simplified.length; i++) {
      const currentMidi = simplified[i];
      let duration = 0.5; // Default duration
      
      if (i < simplified.length - 1) {
        // Find the next occurrence of the next MIDI value starting from current position
        const nextMidi = simplified[i + 1];
        let nextPitchIndex = -1;
        
        // Search forward from current position to find next note transition
        for (let j = currentPitchIndex; j < pitches.length; j++) {
          if (Math.round(pitches[j].midi) === nextMidi) {
            nextPitchIndex = j;
            break;
          }
        }
        
        // If next note found, calculate time difference
        if (nextPitchIndex > currentPitchIndex) {
          const timeDiff = pitches[nextPitchIndex].time - pitches[currentPitchIndex].time;
          duration = Math.max(0.2, Math.min(1.0, timeDiff));
          currentPitchIndex = nextPitchIndex; // Update position for next iteration
        } else {
          // If next note not found, advance to next pitch point
          currentPitchIndex = Math.min(currentPitchIndex + 1, pitches.length - 1);
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
      // Check if playback was stopped (use ref to avoid stale closure)
      if (!isPlayingRef.current || noteIndex >= allNotes.length) {
        setIsPlaying(false);
        isPlayingRef.current = false;
        setCurrentNoteIndex(null);
        timeoutRef.current = null;
        return;
      }

      setCurrentNoteIndex(noteIndex);
      const note = allNotes[noteIndex];
      const duration = noteDurations[noteIndex];
      
      if (synthRef.current) {
        synthRef.current.triggerAttackRelease(note, duration);
      }

      noteIndex++;
      // Store timeout ID so it can be cleared
      timeoutRef.current = setTimeout(playNextNote, duration * 1000);
    };

    playNextNote();
  };

  const stopPlayback = () => {
    // Set flag to stop playback (checked in playNextNote closure)
    isPlayingRef.current = false;
    
    // Clear any pending timeout to prevent stale callbacks
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
      timeoutRef.current = null;
    }
    
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

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
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
          <div className="flex items-center justify-between mb-2">
            <h3 className="font-medium">Tempo</h3>
            {!editingBPM ? (
              <button
                onClick={() => {
                  setTempBPM(bpm?.toString() || '120');
                  setEditingBPM(true);
                }}
                className="text-gray-400 hover:text-primary-500 transition-colors"
              >
                <Edit2 size={16} />
              </button>
            ) : (
              <div className="flex gap-1">
                <button
                  onClick={() => {
                    const newBPM = parseInt(tempBPM) || 120;
                    onBPMChange?.(newBPM);
                    setEditingBPM(false);
                  }}
                  className="text-green-400 hover:text-green-300"
                >
                  <Check size={16} />
                </button>
                <button
                  onClick={() => {
                    setTempBPM(bpm?.toString() || '120');
                    setEditingBPM(false);
                  }}
                  className="text-red-400 hover:text-red-300"
                >
                  <X size={16} />
                </button>
              </div>
            )}
          </div>
          {editingBPM ? (
            <input
              type="number"
              value={tempBPM}
              onChange={(e) => setTempBPM(e.target.value)}
              className="w-full bg-gray-800 border border-gray-700 rounded px-2 py-1 text-2xl font-bold text-primary-500"
              min="60"
              max="200"
              autoFocus
            />
          ) : (
            <p className="text-2xl font-bold text-primary-500">
              {bpm ? `${bpm} BPM` : 'Detecting...'}
            </p>
          )}
          <p className="text-sm text-gray-400 mt-1">
            {bpm ? 'Auto-detected' : 'Analysis pending'}
          </p>
        </div>

        {/* Key */}
        <div className="bg-black/30 rounded-lg p-4 border border-gray-800">
          <div className="flex items-center justify-between mb-2">
            <h3 className="font-medium">Musical Key</h3>
            {!editingKey ? (
              <button
                onClick={() => {
                  setTempKey(musicalKey || 'C Major');
                  setEditingKey(true);
                }}
                className="text-gray-400 hover:text-primary-500 transition-colors"
              >
                <Edit2 size={16} />
              </button>
            ) : (
              <div className="flex gap-1">
                <button
                  onClick={() => {
                    onKeyChange?.(tempKey);
                    setEditingKey(false);
                  }}
                  className="text-green-400 hover:text-green-300"
                >
                  <Check size={16} />
                </button>
                <button
                  onClick={() => {
                    setTempKey(musicalKey || 'C Major');
                    setEditingKey(false);
                  }}
                  className="text-red-400 hover:text-red-300"
                >
                  <X size={16} />
                </button>
              </div>
            )}
          </div>
          {editingKey ? (
            <select
              value={tempKey}
              onChange={(e) => setTempKey(e.target.value)}
              className="w-full bg-gray-800 border border-gray-700 rounded px-2 py-1 text-2xl font-bold text-primary-500"
              autoFocus
            >
              {['C', 'C#', 'D', 'D#', 'E', 'F', 'F#', 'G', 'G#', 'A', 'A#', 'B'].map(note => (
                <option key={`${note} Major`} value={`${note} Major`}>{note} Major</option>
              ))}
              {['C', 'C#', 'D', 'D#', 'E', 'F', 'F#', 'G', 'G#', 'A', 'A#', 'B'].map(note => (
                <option key={`${note} Minor`} value={`${note} Minor`}>{note} Minor</option>
              ))}
            </select>
          ) : (
            <p className="text-2xl font-bold text-primary-500">
              {musicalKey || 'Detecting...'}
            </p>
          )}
          <p className="text-sm text-gray-400 mt-1">
            {musicalKey ? 'Auto-detected' : 'Analysis pending'}
          </p>
        </div>

        {/* Time Signature */}
        <div className="bg-black/30 rounded-lg p-4 border border-gray-800">
          <div className="flex items-center justify-between mb-2">
            <h3 className="font-medium">Time Signature</h3>
            {!editingTimeSig ? (
              <button
                onClick={() => {
                  setTempTimeSig(timeSignature?.display || '4/4');
                  setEditingTimeSig(true);
                }}
                className="text-gray-400 hover:text-primary-500 transition-colors"
              >
                <Edit2 size={16} />
              </button>
            ) : (
              <div className="flex gap-1">
                <button
                  onClick={() => {
                    const [num, den] = tempTimeSig.split('/').map(Number);
                    onTimeSignatureChange?.({ numerator: num, denominator: den, display: tempTimeSig });
                    setEditingTimeSig(false);
                  }}
                  className="text-green-400 hover:text-green-300"
                >
                  <Check size={16} />
                </button>
                <button
                  onClick={() => {
                    setTempTimeSig(timeSignature?.display || '4/4');
                    setEditingTimeSig(false);
                  }}
                  className="text-red-400 hover:text-red-300"
                >
                  <X size={16} />
                </button>
              </div>
            )}
          </div>
          {editingTimeSig ? (
            <select
              value={tempTimeSig}
              onChange={(e) => setTempTimeSig(e.target.value)}
              className="w-full bg-gray-800 border border-gray-700 rounded px-2 py-1 text-2xl font-bold text-primary-500"
              autoFocus
            >
              <option value="4/4">4/4</option>
              <option value="3/4">3/4</option>
              <option value="2/4">2/4</option>
              <option value="6/8">6/8</option>
              <option value="12/8">12/8</option>
              <option value="5/4">5/4</option>
              <option value="7/8">7/8</option>
            </select>
          ) : (
            <p className="text-2xl font-bold text-primary-500">
              {timeSignature?.display || '4/4'}
            </p>
          )}
          <p className="text-sm text-gray-400 mt-1">
            {timeSignature ? 'Auto-detected' : 'Default'}
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
            {currentNoteIndex !== null && (
              <span className="ml-2 text-primary-400 font-mono">
                Now: {detector.midiToNoteName(simplified[currentNoteIndex])}
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

