import { createRealTimeBpmProcessor } from 'realtime-bpm-analyzer';
import { BPMAnalysis } from '../types';

export class BPMDetector {
  async analyze(audioBuffer: AudioBuffer): Promise<BPMAnalysis> {
    // First try the realtime-bpm-analyzer
    try {
      const audioContext = new AudioContext({ sampleRate: audioBuffer.sampleRate });
      const processor = await createRealTimeBpmProcessor(audioContext, {
        continuousAnalysis: false,
        stabilizationTime: Math.min(audioBuffer.duration * 1000, 10000) // Max 10 seconds
      });
      
      return new Promise((resolve) => {
        let detectedBPM: number | null = null;
        let confidence = 0.5;
        let bpmValues: number[] = [];
        
        processor.port.onmessage = (event) => {
          if (event.data.message === 'BPM') {
            const bpm = Math.round(event.data.data.bpm);
            if (bpm >= 60 && bpm <= 200) { // Valid BPM range
              bpmValues.push(bpm);
              detectedBPM = bpm;
              confidence = event.data.data.confidence || 0.7;
            }
          }
          
          if (event.data.message === 'BPM_STABLE') {
            const bpm = Math.round(event.data.data.bpm);
            if (bpm >= 60 && bpm <= 200) {
              detectedBPM = bpm;
              confidence = 1.0;
            }
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
            
            // If we got multiple BPM readings, use the median
            if (bpmValues.length > 0) {
              bpmValues.sort((a, b) => a - b);
              const medianIndex = Math.floor(bpmValues.length / 2);
              detectedBPM = bpmValues[medianIndex];
            }
            
            // Fallback to onset-based detection if realtime analyzer failed
            if (!detectedBPM || detectedBPM < 60 || detectedBPM > 200) {
              const onsetBPM = this.detectBPMFromOnsets(audioBuffer);
              if (onsetBPM) {
                detectedBPM = onsetBPM;
                confidence = 0.6;
              }
            }
            
            const finalBPM = detectedBPM || 120;
            
            // Debug logging (only in development)
            if (process.env.NODE_ENV === 'development') {
              console.log('BPM detection:', {
                detectedBPM: finalBPM,
                confidence,
                stable: detectedBPM !== null && confidence > 0.7,
                method: detectedBPM ? 'realtime-analyzer' : 'onset-fallback'
              });
            }
            
            resolve({
              bpm: finalBPM,
              confidence: confidence,
              stable: detectedBPM !== null && confidence > 0.7
            });
          }, 1000); // Increased wait time
        };
        
        source.start();
      });
    } catch (error) {
      console.error('BPM detection error:', error);
      // Fallback to onset-based detection
      const onsetBPM = this.detectBPMFromOnsets(audioBuffer);
      return {
        bpm: onsetBPM || 120,
        confidence: onsetBPM ? 0.5 : 0.1,
        stable: false
      };
    }
  }
  
  private detectBPMFromOnsets(audioBuffer: AudioBuffer): number | null {
    const channelData = audioBuffer.getChannelData(0);
    const sampleRate = audioBuffer.sampleRate;
    
    // Find note onsets (peaks in amplitude)
    const onsets: number[] = [];
    const threshold = 0.1; // Amplitude threshold
    const minInterval = sampleRate * 0.1; // Minimum 0.1s between onsets
    
    let lastOnset = -minInterval;
    for (let i = 1; i < channelData.length - 1; i++) {
      const current = Math.abs(channelData[i]);
      const prev = Math.abs(channelData[i - 1]);
      const next = Math.abs(channelData[i + 1]);
      
      // Peak detection
      if (current > threshold && current > prev && current > next) {
        const time = i / sampleRate;
        if (time - lastOnset >= 0.1) {
          onsets.push(time);
          lastOnset = time;
        }
      }
    }
    
    if (onsets.length < 4) return null;
    
    // Calculate intervals between onsets
    const intervals: number[] = [];
    for (let i = 1; i < onsets.length; i++) {
      intervals.push(onsets[i] - onsets[i - 1]);
    }
    
    // Find most common interval
    const intervalCounts = new Map<number, number>();
    intervals.forEach(interval => {
      // Round to nearest 0.1s
      const rounded = Math.round(interval * 10) / 10;
      intervalCounts.set(rounded, (intervalCounts.get(rounded) || 0) + 1);
    });
    
    let maxCount = 0;
    let mostCommonInterval = 0;
    intervalCounts.forEach((count, interval) => {
      if (count > maxCount && interval > 0.2 && interval < 2.0) {
        maxCount = count;
        mostCommonInterval = interval;
      }
    });
    
    if (mostCommonInterval > 0) {
      const bpm = Math.round(60 / mostCommonInterval);
      if (bpm >= 60 && bpm <= 200) {
        return bpm;
      }
    }
    
    return null;
  }
}

