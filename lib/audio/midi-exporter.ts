import { Midi } from '@tonejs/midi';
import { PitchPoint } from '../types';

export class MidiExporter {
  export(pitches: PitchPoint[], bpm: number, key?: string): Blob {
    const midi = new Midi();
    
    // Set tempo
    midi.header.setTempo(bpm);
    
    // Add track for melody
    const track = midi.addTrack();
    track.name = 'VoxForge Melody';
    
    // Convert pitch points to MIDI notes
    pitches.forEach((pitch) => {
      const noteDuration = 0.25; // Quarter note default
      
      track.addNote({
        midi: Math.round(pitch.midi),
        time: pitch.time,
        duration: noteDuration,
        velocity: 0.8
      });
    });
    
    // Convert to binary
    const midiArray = midi.toArray();
    return new Blob([midiArray], { type: 'audio/midi' });
  }
}

