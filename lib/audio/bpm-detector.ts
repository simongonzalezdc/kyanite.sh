import { createRealTimeBpmProcessor } from 'realtime-bpm-analyzer';
import { BPMAnalysis } from '../types';

export class BPMDetector {
  async analyze(audioBuffer: AudioBuffer): Promise<BPMAnalysis> {
    const audioContext = new AudioContext({ sampleRate: audioBuffer.sampleRate });
    const processor = await createRealTimeBpmProcessor(audioContext, {
      continuousAnalysis: false,
      stabilizationTime: Math.min(audioBuffer.duration * 1000, 10000) // Max 10 seconds
    });
    
    return new Promise((resolve) => {
      let detectedBPM: number | null = null;
      let confidence = 0.5;
      
      processor.port.onmessage = (event) => {
        if (event.data.message === 'BPM') {
          detectedBPM = Math.round(event.data.data.bpm);
          confidence = event.data.data.confidence || 0.7;
        }
        
        if (event.data.message === 'BPM_STABLE') {
          detectedBPM = Math.round(event.data.data.bpm);
          confidence = 1.0;
        }
      };

      const source = audioContext.createBufferSource();
      source.buffer = audioBuffer;
      source.connect(processor);
      source.connect(audioContext.destination);
      
      source.onended = () => {
        // Wait a bit for final BPM calculation
        setTimeout(() => {
          audioContext.close();
          resolve({
            bpm: detectedBPM || 120, // Default to 120 if not detected
            confidence: confidence,
            stable: detectedBPM !== null
          });
        }, 500);
      };
      
      source.start();
    });
  }
}

