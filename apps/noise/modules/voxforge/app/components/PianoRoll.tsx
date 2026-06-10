'use client';

import { useState, useRef, useEffect } from 'react';
import { PitchPoint } from '@/lib/types';
import { PitchDetector } from '@/lib/audio/pitch-detector';
import * as Tone from 'tone';
import { Play, Square } from 'lucide-react';

interface PianoRollProps {
  pitches: PitchPoint[];
  duration: number;
  bpm?: number | null;
}

export default function PianoRoll({ pitches, duration, bpm }: PianoRollProps) {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const [isPlaying, setIsPlaying] = useState(false);
  const synthRef = useRef<Tone.Synth | null>(null);
  const sequenceRef = useRef<Tone.Sequence | null>(null);
  const isPlayingRef = useRef<boolean>(false);
  
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

  useEffect(() => {
    if (!canvasRef.current || pitches.length === 0) return;

    const canvas = canvasRef.current;
    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    // Set canvas size
    const width = Math.max(1, canvas.offsetWidth);
    const height = Math.max(1, canvas.offsetHeight);
    canvas.width = width * 2; // Retina
    canvas.height = height * 2;
    ctx.scale(2, 2);

    // Piano roll settings
    const detector = new PitchDetector();
    const simplified = detector.getSimplifiedMelody(pitches);
    
    if (simplified.length === 0) return;
    
    // Get note range
    const midiNotes = simplified.map(midi => Math.round(midi));
    const minMidi = Math.max(0, Math.min(...midiNotes));
    const maxMidi = Math.min(127, Math.max(...midiNotes));
    
    // Add padding to range
    const paddedMin = Math.max(0, minMidi - 2);
    const paddedMax = Math.min(127, maxMidi + 2);
    const paddedRange = Math.max(1, paddedMax - paddedMin + 1);

    // Calculate dimensions - ensure positive values
    const noteHeight = Math.max(1, height / paddedRange);
    const timeScale = duration > 0 ? Math.max(0.001, width / duration) : 1;

    // Clear canvas
    ctx.fillStyle = '#0a0a0a';
    ctx.fillRect(0, 0, width, height);

    // Draw grid lines (octaves)
    ctx.strokeStyle = '#1a1a1a';
    ctx.lineWidth = 1;
    
    for (let midi = paddedMin; midi <= paddedMax; midi++) {
      const y = height - ((midi - paddedMin) * noteHeight);
      
      // Highlight C notes (every 12 semitones)
      if (midi % 12 === 0) {
        ctx.fillStyle = '#1a1a1a';
        ctx.fillRect(0, y - noteHeight / 2, width, noteHeight);
      }
      
      ctx.beginPath();
      ctx.moveTo(0, y);
      ctx.lineTo(width, y);
      ctx.stroke();
    }

    // Draw note labels (C notes only)
    ctx.fillStyle = '#666';
    ctx.font = '10px monospace';
    ctx.textAlign = 'right';
    ctx.textBaseline = 'middle';
    
    for (let midi = paddedMin; midi <= paddedMax; midi++) {
      if (midi % 12 === 0) {
        const y = height - ((midi - paddedMin) * noteHeight);
        const noteName = detector.midiToNoteName(midi);
        ctx.fillText(noteName, width - 5, y);
      }
    }

    // Draw notes - group by actual note events
    const noteMap = new Map<number, { start: number; end: number }[]>();
    
    // Group consecutive same notes into note events
    let currentMidi = -1;
    let currentStart = 0;
    let lastTime = 0;
    
    pitches.forEach((pitch, i) => {
      const midi = Math.round(pitch.midi);
      const time = pitch.time;
      
      if (midi !== currentMidi || (time - lastTime) > 0.2) {
        // Save previous note if it existed
        if (currentMidi >= 0 && currentMidi >= 0 && currentMidi <= 127) {
          const noteEnd = lastTime || time;
          if (!noteMap.has(currentMidi)) {
            noteMap.set(currentMidi, []);
          }
          noteMap.get(currentMidi)!.push({ 
            start: currentStart, 
            end: noteEnd 
          });
        }
        
        // Start new note
        currentMidi = midi;
        currentStart = time;
      }
      
      lastTime = time;
    });
    
    // Save last note
    if (currentMidi >= 0 && currentMidi <= 127) {
      if (!noteMap.has(currentMidi)) {
        noteMap.set(currentMidi, []);
      }
      noteMap.get(currentMidi)!.push({ 
        start: currentStart, 
        end: duration 
      });
    }

    // Draw note rectangles - handle multiple note events per MIDI note
    noteMap.forEach((noteEvents, midi) => {
      if (midi < paddedMin || midi > paddedMax) return;
      
      noteEvents.forEach((timeRange) => {
        const y = height - ((midi - paddedMin) * noteHeight);
        const x = Math.max(0, timeRange.start * timeScale);
        const noteWidth = Math.max(1, (timeRange.end - timeRange.start) * timeScale);
        const noteY = Math.max(0, y - noteHeight * 0.8);
        const noteHeightDraw = Math.max(1, noteHeight * 0.6);
        
        // Ensure values are within canvas bounds
        if (x < 0 || noteY < 0 || x + noteWidth > width || noteY + noteHeightDraw > height) {
          return; // Skip notes outside canvas
        }
        
        // Draw note
        ctx.fillStyle = '#3B82F6';
        ctx.fillRect(x, noteY, noteWidth, noteHeightDraw);
        
        // Draw border
        ctx.strokeStyle = '#2563EB';
        ctx.lineWidth = 1;
        ctx.strokeRect(x, noteY, noteWidth, noteHeightDraw);
      });
    });

    // Draw time grid (beats if BPM available)
    if (bpm && bpm > 0) {
      const beatInterval = Math.max(0.01, 60 / bpm);
      ctx.strokeStyle = '#333';
      ctx.lineWidth = 1;
      
      for (let time = 0; time <= duration; time += beatInterval) {
        const x = Math.max(0, Math.min(width, time * timeScale));
        ctx.beginPath();
        ctx.moveTo(x, 0);
        ctx.lineTo(x, height);
        ctx.stroke();
      }
    }

    // Draw time axis
    ctx.strokeStyle = '#666';
    ctx.lineWidth = 1;
    ctx.beginPath();
    ctx.moveTo(0, height);
    ctx.lineTo(width, height);
    ctx.stroke();

    // Draw time labels
    ctx.fillStyle = '#999';
    ctx.font = '9px monospace';
    ctx.textAlign = 'center';
    ctx.textBaseline = 'top';
    
    const timeLabels = Math.min(10, Math.max(1, Math.ceil(duration)));
    for (let i = 0; i <= timeLabels; i++) {
      const time = (i / timeLabels) * duration;
      const x = Math.max(0, Math.min(width, time * timeScale));
      ctx.fillText(time.toFixed(1) + 's', x, height + 3);
    }
  }, [pitches, duration, bpm]);

  if (pitches.length === 0) {
    return (
      <div className="bg-gray-900 rounded-xl p-8 border border-gray-800">
        <p className="text-gray-400 text-center">No notes to display</p>
      </div>
    );
  }

  const playPianoRoll = async () => {
    if (!synthRef.current || pitches.length === 0) return;
    
    await Tone.start();
    setIsPlaying(true);
    isPlayingRef.current = true;
    
    const detector = new PitchDetector();
    
    // Create note events with proper timing from pitches
    const noteEvents: Array<{ time: number; note: string; duration: number }> = [];
    let currentMidi = -1;
    let noteStart = 0;
    
    pitches.forEach((pitch, i) => {
      const midi = Math.round(pitch.midi);
      
      if (midi !== currentMidi) {
        // Save previous note if it existed
        if (currentMidi >= 0 && i > 0) {
          const noteEnd = pitch.time;
          const noteDuration = Math.max(0.1, Math.min(2.0, noteEnd - noteStart));
          const noteName = detector.midiToNoteName(currentMidi);
          noteEvents.push({
            time: noteStart,
            note: noteName,
            duration: noteDuration
          });
        }
        
        // Start new note
        currentMidi = midi;
        noteStart = pitch.time;
      }
    });
    
    // Add last note
    if (currentMidi >= 0) {
      const lastNoteDuration = Math.max(0.1, Math.min(2.0, duration - noteStart));
      const noteName = detector.midiToNoteName(currentMidi);
      noteEvents.push({
        time: noteStart,
        note: noteName,
        duration: lastNoteDuration
      });
    }
    
    if (noteEvents.length === 0) {
      setIsPlaying(false);
      isPlayingRef.current = false;
      return;
    }
    
    // Play notes in sequence using Tone.js scheduling
    const now = Tone.now();
    noteEvents.forEach((event) => {
      if (synthRef.current && isPlayingRef.current) {
        synthRef.current.triggerAttackRelease(
          event.note,
          event.duration,
          now + event.time
        );
      }
    });
    
    // Stop after all notes
    const lastEvent = noteEvents[noteEvents.length - 1];
    const totalDuration = lastEvent.time + lastEvent.duration;
    
    setTimeout(() => {
      if (isPlayingRef.current) {
        setIsPlaying(false);
        isPlayingRef.current = false;
      }
    }, totalDuration * 1000 + 100);
  };
  
  const stopPlayback = () => {
    isPlayingRef.current = false;
    if (synthRef.current) {
      synthRef.current.triggerRelease();
    }
    if (sequenceRef.current) {
      sequenceRef.current.stop();
    }
    setIsPlaying(false);
  };

  return (
    <div className="bg-gray-900 rounded-xl p-8 border border-gray-800 space-y-4">
      <div className="flex items-center justify-between">
        <h3 className="font-medium">MIDI Piano Roll</h3>
        <div className="flex gap-2">
          {!isPlaying ? (
            <button
              onClick={playPianoRoll}
              className="flex items-center gap-2 px-4 py-2 bg-primary-500 hover:bg-primary-600 rounded-lg font-medium transition-colors text-sm"
            >
              <Play size={16} />
              Play Piano Roll
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
      <div className="relative">
        <canvas
          ref={canvasRef}
          className="w-full h-96 rounded-lg bg-black border border-gray-800"
        />
      </div>
      <p className="text-xs text-gray-400">
        Notes are shown as blue rectangles. C notes are highlighted in gray. Click "Play Piano Roll" to hear the detected melody.
      </p>
    </div>
  );
}

