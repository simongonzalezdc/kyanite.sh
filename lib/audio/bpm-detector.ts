import { BPMAnalysis } from '../types';

export class BPMDetector {
  async analyze(audioBuffer: AudioBuffer): Promise<BPMAnalysis> {
    try {
      console.log('Starting BPM analysis...');
      console.log('Audio duration:', audioBuffer.duration, 'seconds');
      console.log('Sample rate:', audioBuffer.sampleRate);
      
      // Use onset-based detection as primary method
      const onsetBPM = this.detectBPMFromOnsets(audioBuffer);
      
      if (onsetBPM) {
        console.log('BPM detected from onsets:', onsetBPM);
        return {
          bpm: onsetBPM,
          confidence: 0.8,
          stable: true
        };
      }
      
      // Fallback to energy-based analysis
      const energyBPM = this.detectBPMFromEnergy(audioBuffer);
      if (energyBPM) {
        console.log('BPM detected from energy:', energyBPM);
        return {
          bpm: energyBPM,
          confidence: 0.6,
          stable: false
        };
      }
      
      // Default fallback
      console.log('Using default BPM: 120');
      return {
        bpm: 120,
        confidence: 0.1,
        stable: false
      };
      
    } catch (error) {
      console.error('BPM detection error:', error);
      return {
        bpm: 120,
        confidence: 0.1,
        stable: false
      };
    }
  }
  
  private detectBPMFromOnsets(audioBuffer: AudioBuffer): number | null {
    const channelData = audioBuffer.getChannelData(0);
    const sampleRate = audioBuffer.sampleRate;
    
    // Find note onsets using peak detection
    const onsets: number[] = [];
    const threshold = 0.05; // Lower threshold for better detection
    const minInterval = Math.floor(sampleRate * 0.2); // Minimum 0.2s between onsets
    
    let lastOnset = -minInterval;
    
    // First pass: find potential onsets
    for (let i = 1; i < channelData.length - 1; i++) {
      const current = Math.abs(channelData[i]);
      const prev = Math.abs(channelData[i - 1]);
      const next = Math.abs(channelData[i + 1]);
      
      // Peak detection
      if (current > threshold && current > prev && current > next) {
        if (i - lastOnset >= minInterval) {
          onsets.push(i);
          lastOnset = i;
        }
      }
    }
    
    console.log('Found onsets:', onsets.length);
    
    if (onsets.length < 4) return null;
    
    // Calculate intervals between onsets
    const intervals: number[] = [];
    for (let i = 1; i < onsets.length; i++) {
      const interval = (onsets[i] - onsets[i - 1]) / sampleRate;
      intervals.push(interval);
    }
    
    // Find most common interval using histogram
    const intervalHist = new Map<number, number>();
    const tolerance = 0.1; // 100ms tolerance
    
    for (const interval of intervals) {
      let found = false;
      for (const [key, count] of intervalHist.entries()) {
        if (Math.abs(key - interval) < tolerance) {
          intervalHist.set(key, count + 1);
          found = true;
          break;
        }
      }
      if (!found) {
        intervalHist.set(interval, 1);
      }
    }
    
    // Find interval with most votes
    let maxVotes = 0;
    let bestInterval = 0;
    intervalHist.forEach((votes, interval) => {
      if (votes > maxVotes && interval > 0.3 && interval < 2.0) {
        maxVotes = votes;
        bestInterval = interval;
      }
    });
    
    if (bestInterval > 0) {
      const bpm = Math.round(60 / bestInterval);
      if (bpm >= 60 && bpm <= 200) {
        return bpm;
      }
    }
    
    return null;
  }
  
  private detectBPMFromEnergy(audioBuffer: AudioBuffer): number | null {
    const channelData = audioBuffer.getChannelData(0);
    const sampleRate = audioBuffer.sampleRate;
    
    // Calculate energy in small windows
    const windowSize = Math.floor(sampleRate * 0.05); // 50ms windows
    const hopSize = Math.floor(windowSize / 2); // 50% overlap
    const energies: number[] = [];
    
    for (let i = 0; i < channelData.length - windowSize; i += hopSize) {
      let energy = 0;
      for (let j = 0; j < windowSize; j++) {
        energy += channelData[i + j] * channelData[i + j];
      }
      energies.push(energy / windowSize);
    }
    
    // Find peaks in energy
    const peaks: number[] = [];
    const threshold = 0.01;
    const minPeakInterval = Math.floor(sampleRate * 0.3 / hopSize); // 300ms minimum
    
    let lastPeak = -minPeakInterval;
    for (let i = 1; i < energies.length - 1; i++) {
      if (energies[i] > threshold && 
          energies[i] > energies[i - 1] && 
          energies[i] > energies[i + 1]) {
        if (i - lastPeak >= minPeakInterval) {
          peaks.push(i);
          lastPeak = i;
        }
      }
    }
    
    if (peaks.length < 4) return null;
    
    // Calculate intervals between peaks
    const intervals: number[] = [];
    for (let i = 1; i < peaks.length; i++) {
      const interval = (peaks[i] - peaks[i - 1]) * hopSize / sampleRate;
      intervals.push(interval);
    }
    
    // Find most common interval
    const intervalCounts = new Map<number, number>();
    intervals.forEach(interval => {
      const rounded = Math.round(interval * 10) / 10; // Round to 0.1s
      intervalCounts.set(rounded, (intervalCounts.get(rounded) || 0) + 1);
    });
    
    let maxCount = 0;
    let bestInterval = 0;
    intervalCounts.forEach((count, interval) => {
      if (count > maxCount && interval > 0.3 && interval < 2.0) {
        maxCount = count;
        bestInterval = interval;
      }
    });
    
    if (bestInterval > 0) {
      const bpm = Math.round(60 / bestInterval);
      if (bpm >= 60 && bpm <= 200) {
        return bpm;
      }
    }
    
    return null;
  }
}
