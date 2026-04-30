import * as Tone from 'tone';
import { MusicGenerator } from './music-generator';
import { InstrumentType } from '../types';

export class StemExporter {
  private generator: MusicGenerator;

  constructor(generator: MusicGenerator) {
    this.generator = generator;
  }

  async exportStem(
    type: 'vocal' | 'drums' | 'bass' | 'chords' | 'mix',
    audioBuffer?: AudioBuffer,
    duration: number = 30
  ): Promise<Blob> {
    const recorder = new Tone.Recorder();
    
    // Connect appropriate sources to recorder
    if (type === 'vocal' && audioBuffer) {
      const player = new Tone.Player(audioBuffer);
      player.connect(recorder);
      recorder.start();
      player.start();
      
      await new Promise(resolve => setTimeout(resolve, duration * 1000));
      
      player.stop();
      const recording = await recorder.stop();
      return recording;
    } else if (type === 'drums' || type === 'bass' || type === 'chords' || type === 'mix') {
      // For generated instruments, we need to record from the generator
      // This is a simplified version - in production, you'd route specific instruments
      recorder.start();
      
      // Start the generator
      this.generator.start();
      
      await new Promise(resolve => setTimeout(resolve, duration * 1000));
      
      this.generator.stop();
      const recording = await recorder.stop();
      return recording;
    }

    throw new Error(`Unknown stem type: ${type}`);
  }
}

