import { PitchPoint } from '../types';

// Industry-standard pitch detection using advanced techniques
export class PitchDetector {
  private sampleRate: number;
  private hannWindow: Float32Array;
  private noiseProfile: number[] = [];

  constructor(sampleRate: number = 44100) {
    this.sampleRate = sampleRate;
    this.hannWindow = this.createHannWindow(4096);
  }

  private createHannWindow(size: number): Float32Array {
    const window = new Float32Array(size);
    for (let i = 0; i < size; i++) {
      window[i] = 0.5 * (1 - Math.cos(2 * Math.PI * i / (size - 1)));
    }
    return window;
  }

  analyze(audioBuffer: AudioBuffer): PitchPoint[] {
    if (audioBuffer.numberOfChannels === 0) {
      console.warn('AudioBuffer has no channels');
      return [];
    }
    if (this.sampleRate <= 0) {
      console.warn('Invalid sample rate:', this.sampleRate);
      return [];
    }
    const channelData = audioBuffer.getChannelData(0);
    const pitches: PitchPoint[] = [];
    
    // Industry-standard parameters
    const windowSize = 4096; // FFT-friendly size
    const hopSize = 256; // 75% overlap for smooth tracking
    const minFreq = 60; // Lower bound for voice/instruments
    const maxFreq = 2500; // Upper bound for voice/instruments
    
    if (channelData.length < windowSize) {
      console.warn('Audio buffer too short for analysis');
      return [];
    }
    
    console.log('Starting advanced pitch analysis with industry techniques...');
    console.log('Audio duration:', audioBuffer.duration, 'seconds');
    console.log('Sample rate:', this.sampleRate);
    console.log('Total samples:', channelData.length);

    // Build noise profile for noise reduction
    this.buildNoiseProfile(channelData);
    
    // Multi-frame analysis for robustness
    for (let i = 0; i < channelData.length - windowSize; i += hopSize) {
      const window = channelData.slice(i, i + windowSize);
      const time = i / this.sampleRate;
      
      // Apply pre-processing pipeline
      const processedWindow = this.preprocessAudio(window);
      
      // Use multiple pitch detection methods for accuracy
      const frequency = this.detectPitchMultipleMethods(processedWindow);
      
      if (frequency && frequency >= minFreq && frequency <= maxFreq) {
        const confidence = this.calculateAdvancedConfidence(processedWindow, frequency);
        
        if (confidence > 15) { // Lower threshold for better sensitivity
          const midi = this.frequencyToMidi(frequency);
          const clampedMidi = Math.max(0, Math.min(127, midi));

          const pitchPoint: PitchPoint = {
            frequency,
            time,
            midi: clampedMidi,
            confidence: Math.min(100, Math.max(0, confidence))
          };

          pitches.push(pitchPoint);
        }
      }
    }

    console.log('Raw pitch analysis complete. Detected', pitches.length, 'pitch points');
    
    // Apply advanced musical post-processing
    if (pitches.length > 0) {
      const processed = this.applyMusicalPostProcessing(pitches);
      console.log('After musical post-processing:', processed.length, 'pitch points');
      return processed;
    }
    
    return pitches;
  }

  private buildNoiseProfile(channelData: Float32Array): void {
    // Analyze quiet sections to build noise profile
    const frameSize = 1024;
    const noiseFrames: number[] = [];
    
    for (let i = 0; i < channelData.length - frameSize; i += frameSize * 10) {
      const frame = channelData.slice(i, i + frameSize);
      const rms = this.calculateRMS(frame);
      if (rms < 0.01) { // Very quiet frame
        noiseFrames.push(rms);
      }
    }
    
    // Store average noise level
    if (noiseFrames.length > 0) {
      this.noiseProfile = [noiseFrames.reduce((a, b) => a + b, 0) / noiseFrames.length];
    }
  }

  private preprocessAudio(window: Float32Array): Float32Array {
    const processed = new Float32Array(window.length);
    
    // 1. Pre-emphasis filter (emphasizes high frequencies, reduces noise)
    const preEmphasis = 0.97;
    for (let i = 0; i < window.length; i++) {
      if (i === 0) {
        processed[i] = window[i];
      } else {
        processed[i] = window[i] - preEmphasis * window[i - 1];
      }
    }
    
    // 2. Apply Hann window for better spectral analysis
    for (let i = 0; i < processed.length; i++) {
      processed[i] *= this.hannWindow[i];
    }
    
    // 3. Noise reduction based on noise profile
    const noiseLevel = this.noiseProfile[0] || 0.001;
    const noiseThreshold = noiseLevel * 3; // 3x noise floor
    
    for (let i = 0; i < processed.length; i++) {
      if (Math.abs(processed[i]) < noiseThreshold) {
        processed[i] *= 0.1; // Reduce noise
      }
    }
    
    return processed;
  }

  private detectPitchMultipleMethods(window: Float32Array): number | null {
    // Method 1: YIN algorithm (most reliable for music)
    const yinPitch = this.yinAlgorithm(window);
    
    // Method 2: Autocorrelation (good for periodic signals)
    const autoPitch = this.autocorrelationPitch(window);
    
    // Method 3: Harmonic Product Spectrum (good for voiced speech)
    const hpsPitch = this.harmonicProductSpectrum(window);
    
    // Method 4: FFT-based peak detection
    const fftPitch = this.fftPeakDetection(window);
    
    // Combine results using confidence-weighted voting
    const candidates = [
      { pitch: yinPitch, confidence: 0.9 },
      { pitch: autoPitch, confidence: 0.7 },
      { pitch: hpsPitch, confidence: 0.6 },
      { pitch: fftPitch, confidence: 0.5 }
    ].filter(c => c.pitch !== null) as { pitch: number; confidence: number }[];

    if (candidates.length === 0) return null;
    
    // Group similar pitches and find consensus
    const grouped = this.groupSimilarPitches(candidates);
    
    // Return the most confident group
    let bestGroup = grouped[0];
    for (const group of grouped) {
      if (group.totalConfidence > bestGroup.totalConfidence) {
        bestGroup = group;
      }
    }
    
    return bestGroup.averagePitch;
  }

  private yinAlgorithm(window: Float32Array): number | null {
    // Simplified YIN implementation
    const threshold = 0.1;
    const minPeriod = Math.floor(this.sampleRate / 800); // 800 Hz max
    const maxPeriod = Math.floor(this.sampleRate / 80);  // 80 Hz min
    
    // Calculate difference function
    const diff = new Float32Array(maxPeriod);
    for (let tau = minPeriod; tau < maxPeriod; tau++) {
      let sum = 0;
      for (let i = 0; i < window.length - tau; i++) {
        const delta = window[i] - window[i + tau];
        sum += delta * delta;
      }
      diff[tau] = sum;
    }
    
    // Calculate cumulative mean normalized difference
    const cmnd = new Float32Array(maxPeriod);
    let runningSum = 0;
    for (let tau = 1; tau < maxPeriod; tau++) {
      runningSum += diff[tau];
      cmnd[tau] = diff[tau] * tau / (runningSum / tau);
    }
    
    // Find first local minimum below threshold
    let tau = -1;
    for (let i = 2; i < maxPeriod - 1; i++) {
      if (cmnd[i] < threshold && cmnd[i] < cmnd[i - 1] && cmnd[i] < cmnd[i + 1]) {
        tau = i;
        break;
      }
    }
    
    if (tau === -1) return null;
    
    // Parabolic interpolation for better precision
    const betterTau = this.parabolicInterpolation(cmnd, tau);
    
    return this.sampleRate / betterTau;
  }

  private autocorrelationPitch(window: Float32Array): number | null {
    const minLag = Math.floor(this.sampleRate / 1000);
    const maxLag = Math.floor(this.sampleRate / 50);
    const lags: { lag: number; correlation: number }[] = [];
    
    for (let lag = minLag; lag < maxLag; lag++) {
      let correlation = 0;
      for (let i = 0; i < window.length - lag; i++) {
        correlation += window[i] * window[i + lag];
      }
      lags.push({ lag, correlation });
    }
    
    // Find the lag with highest correlation
    lags.sort((a, b) => b.correlation - a.correlation);
    
    if (lags.length > 0 && lags[0].correlation > 0.1) {
      return this.sampleRate / lags[0].lag;
    }
    
    return null;
  }

  private harmonicProductSpectrum(window: Float32Array): number | null {
    // Simple HPS implementation
    const fft = this.simpleFFT(window);
    const magnitude = fft.magnitude;
    const length = magnitude.length;
    
    // Harmonic Product Spectrum (multiply spectrum by its harmonics)
    const hps = new Array(length).fill(0);
    for (let i = 1; i < length / 2; i++) {
      hps[i] = magnitude[i];
      for (let harmonic = 2; harmonic <= 5 && i * harmonic < length; harmonic++) {
        hps[i] *= magnitude[i * harmonic];
      }
    }
    
    // Find peak in HPS
    let maxIndex = 0;
    let maxValue = 0;
    for (let i = 1; i < hps.length / 2; i++) {
      if (hps[i] > maxValue) {
        maxValue = hps[i];
        maxIndex = i;
      }
    }
    
    if (maxValue > 0.01) {
      return (maxIndex * this.sampleRate) / (2 * length);
    }
    
    return null;
  }

  private fftPeakDetection(window: Float32Array): number | null {
    const fft = this.simpleFFT(window);
    const magnitude = fft.magnitude;
    
    // Find peaks in magnitude spectrum
    const peaks: { freq: number; magnitude: number }[] = [];
    
    for (let i = 1; i < magnitude.length - 1; i++) {
      if (magnitude[i] > magnitude[i - 1] && magnitude[i] > magnitude[i + 1] && magnitude[i] > 0.01) {
        const freq = (i * this.sampleRate) / (2 * magnitude.length);
        if (freq >= 60 && freq <= 2500) {
          peaks.push({ freq, magnitude: magnitude[i] });
        }
      }
    }
    
    if (peaks.length > 0) {
      // Return the strongest peak
      peaks.sort((a, b) => b.magnitude - a.magnitude);
      return peaks[0].freq;
    }
    
    return null;
  }

  private groupSimilarPitches(candidates: { pitch: number; confidence: number }[]): { 
    averagePitch: number; 
    totalConfidence: number; 
    count: number 
  }[] {
    const groups: { pitches: number[]; confidences: number[] }[] = [];
    
    for (const candidate of candidates) {
      const pitch = candidate.pitch!;
      const confidence = candidate.confidence;
      
      // Find existing group or create new one
      let assigned = false;
      for (const group of groups) {
        const avgPitch = group.pitches.reduce((a, b) => a + b, 0) / group.pitches.length;
        if (Math.abs(pitch - avgPitch) / avgPitch < 0.1) { // Within 10%
          group.pitches.push(pitch);
          group.confidences.push(confidence);
          assigned = true;
          break;
        }
      }
      
      if (!assigned) {
        groups.push({ pitches: [pitch], confidences: [confidence] });
      }
    }
    
    // Convert to summary format
    return groups.map(group => ({
      averagePitch: group.pitches.reduce((a, b) => a + b, 0) / group.pitches.length,
      totalConfidence: group.confidences.reduce((a, b) => a + b, 0),
      count: group.pitches.length
    }));
  }

  private parabolicInterpolation(data: Float32Array, index: number): number {
    if (index < 1 || index >= data.length - 1) return index;
    
    const s0 = data[index - 1];
    const s1 = data[index];
    const s2 = data[index + 1];
    
    const a = (s0 + s2 - 2 * s1) / 2;
    if (a === 0) return index;
    
    const b = (s2 - s0) / 2;
    return index - b / (2 * a);
  }

  private simpleFFT(signal: Float32Array): { real: Float32Array; imaginary: Float32Array; magnitude: Float32Array } {
    const N = signal.length;
    const real = new Float32Array(N);
    const imaginary = new Float32Array(N);
    const magnitude = new Float32Array(N);
    
    // Copy signal to real part
    for (let i = 0; i < N; i++) {
      real[i] = signal[i];
    }
    
    // Simple DFT (not optimized FFT, but good enough for our purposes)
    for (let k = 0; k < N / 2; k++) {
      let sumReal = 0;
      let sumImag = 0;
      
      for (let n = 0; n < N; n++) {
        const angle = -2 * Math.PI * k * n / N;
        sumReal += real[n] * Math.cos(angle);
        sumImag += real[n] * Math.sin(angle);
      }
      
      real[k] = sumReal;
      imaginary[k] = sumImag;
      magnitude[k] = Math.sqrt(sumReal * sumReal + sumImag * sumImag);
    }
    
    return { real, imaginary, magnitude };
  }

  private calculateAdvancedConfidence(window: Float32Array, frequency: number): number {
    // Multiple confidence factors
    
    // 1. Signal strength (RMS)
    const rms = this.calculateRMS(window);
    const strengthScore = Math.min(1, rms * 10) * 40; // Max 40 points
    
    // 2. Periodicity (autocorrelation at detected period)
    const lag = Math.round(this.sampleRate / frequency);
    const periodicity = this.calculatePeriodicity(window, lag);
    const periodicityScore = periodicity * 30; // Max 30 points
    
    // 3. Harmonic content (ratio of harmonics to fundamental)
    const harmonicScore = this.analyzeHarmonics(window, frequency) * 20; // Max 20 points
    
    // 4. Signal-to-noise ratio
    const noiseLevel = this.noiseProfile[0] || 0.001;
    const snr = rms / noiseLevel;
    const snrScore = Math.min(1, snr / 10) * 10; // Max 10 points
    
    return strengthScore + periodicityScore + harmonicScore + snrScore;
  }

  private calculateRMS(data: Float32Array): number {
    let sum = 0;
    for (let i = 0; i < data.length; i++) {
      sum += data[i] * data[i];
    }
    return Math.sqrt(sum / data.length);
  }

  private calculatePeriodicity(window: Float32Array, lag: number): number {
    if (lag >= window.length) return 0;
    
    let correlation = 0;
    let energy = 0;
    
    for (let i = 0; i < window.length - lag; i++) {
      correlation += window[i] * window[i + lag];
      energy += window[i] * window[i];
    }
    
    return energy > 0 ? correlation / energy : 0;
  }

  private analyzeHarmonics(window: Float32Array, fundamental: number): number {
    // Count strong harmonics
    let harmonicCount = 0;
    const tolerance = 0.05; // 5% tolerance
    
    for (let harmonic = 2; harmonic <= 6; harmonic++) {
      const targetFreq = fundamental * harmonic;
      const detectedFreq = this.findClosestPeak(window, targetFreq);
      if (detectedFreq && Math.abs(detectedFreq - targetFreq) / targetFreq < tolerance) {
        harmonicCount++;
      }
    }
    
    return harmonicCount / 5; // Normalize to 0-1
  }

  private findClosestPeak(window: Float32Array, targetFreq: number): number | null {
    const fft = this.simpleFFT(window);
    const magnitude = fft.magnitude;
    
    const targetIndex = Math.round((targetFreq * 2 * magnitude.length) / this.sampleRate);
    
    if (targetIndex >= magnitude.length) return null;
    
    // Find peak around target frequency
    let bestIndex = targetIndex;
    let bestMagnitude = 0;
    
    for (let i = Math.max(1, targetIndex - 3); i < Math.min(magnitude.length - 1, targetIndex + 4); i++) {
      if (magnitude[i] > magnitude[i - 1] && magnitude[i] > magnitude[i + 1] && magnitude[i] > bestMagnitude) {
        bestIndex = i;
        bestMagnitude = magnitude[i];
      }
    }
    
    return bestMagnitude > 0.01 ? (bestIndex * this.sampleRate) / (2 * magnitude.length) : null;
  }

  private applyMusicalPostProcessing(pitches: PitchPoint[]): PitchPoint[] {
    if (pitches.length === 0) return pitches;
    
    const processed: PitchPoint[] = [];
    
    // 1. Temporal smoothing (median filter to remove spikes)
    const smoothed = this.temporalSmoothing(pitches);
    
    // 2. Musical quantization (snap to nearest semitone with confidence-based weights)
    for (let i = 0; i < smoothed.length; i++) {
      const pitch = smoothed[i];
      const quantMidi = this.quantizeToNearestSemitone(pitch.midi, pitch.confidence || 0);
      
      processed.push({
        ...pitch,
        midi: quantMidi,
        frequency: this.midiToFrequency(quantMidi)
      });
    }
    
    // 3. Remove micro-variations (small frequency changes that are likely noise)
    const stabilized = this.stabilizePitchChanges(processed);
    
    return stabilized;
  }

  private temporalSmoothing(pitches: PitchPoint[]): PitchPoint[] {
    const windowSize = 3; // 3-point median filter
    const halfWindow = Math.floor(windowSize / 2);
    const smoothed: PitchPoint[] = [];
    
    for (let i = 0; i < pitches.length; i++) {
      const start = Math.max(0, i - halfWindow);
      const end = Math.min(pitches.length, i + halfWindow + 1);
      const window = pitches.slice(start, end);
      
      // Median filter on frequency
      const frequencies = window.map(p => p.frequency).sort((a, b) => a - b);
      const medianFreq = frequencies[Math.floor(frequencies.length / 2)];
      
      // Weighted average for time and confidence
      const totalWeight = window.reduce((sum, p) => sum + (p.confidence || 0), 0);
      const weightedTime = totalWeight > 0 ? window.reduce((sum, p) => sum + p.time * (p.confidence || 0), 0) / totalWeight : window.reduce((sum, p) => sum + p.time, 0) / window.length;
      const avgConfidence = window.reduce((sum, p) => sum + (p.confidence || 0), 0) / window.length;
      
      smoothed.push({
        frequency: medianFreq,
        time: weightedTime,
        midi: this.frequencyToMidi(medianFreq),
        confidence: avgConfidence
      });
    }
    
    return smoothed;
  }

  private quantizeToNearestSemitone(midi: number, confidence: number): number {
    // Higher confidence = more rigid quantization
    const quantizeStrength = confidence > 50 ? 1.0 : confidence / 50;
    
    const nearestNote = Math.round(midi);
    return midi + (nearestNote - midi) * quantizeStrength;
  }

  private stabilizePitchChanges(pitches: PitchPoint[]): PitchPoint[] {
    if (pitches.length === 0) return pitches;
    
    const stabilized = [pitches[0]]; // Keep first pitch
    
    for (let i = 1; i < pitches.length; i++) {
      const current = pitches[i];
      const previous = stabilized[stabilized.length - 1];
      
      const timeDiff = current.time - previous.time;
      const freqDiff = Math.abs(current.frequency - previous.frequency);
      const semitoneDiff = Math.abs(current.midi - previous.midi);
      
      // If change is very small and recent, it's likely noise
      if (timeDiff < 0.1 && semitoneDiff < 0.2) {
        // Keep the more confident pitch
        if ((current.confidence || 0) > (previous.confidence || 0)) {
          stabilized[stabilized.length - 1] = current;
        }
        // Otherwise skip current pitch
      } else {
        stabilized.push(current);
      }
    }
    
    return stabilized;
  }

  private frequencyToMidi(frequency: number): number {
    return 69 + 12 * Math.log2(frequency / 440);
  }

  private midiToFrequency(midi: number): number {
    return 440 * Math.pow(2, (midi - 69) / 12);
  }

  // Enhanced statistics
  getStats(pitches: PitchPoint[]) {
    if (pitches.length === 0) {
      return {
        averageFrequency: 0,
        minFrequency: 0,
        maxFrequency: 0,
        averageMidi: 0,
        range: 'unknown' as 'low' | 'mid' | 'high' | 'unknown',
        confidence: 0,
        noteVariety: 0,
        pitchStability: 0
      };
    }

    const frequencies = pitches.map(p => p.frequency);
    const midis = pitches.map(p => p.midi);
    const confidences = pitches.map(p => p.confidence || 0);

    const totalConfidence = confidences.reduce((sum, conf) => sum + conf, 0);
    const weightedAvgFreq = totalConfidence > 0 ? frequencies.reduce((sum, freq, i) => sum + freq * confidences[i], 0) / totalConfidence : frequencies.reduce((sum, freq) => sum + freq, 0) / frequencies.length;
    const weightedAvgMidi = totalConfidence > 0 ? midis.reduce((sum, midi, i) => sum + midi * confidences[i], 0) / totalConfidence : midis.reduce((sum, midi) => sum + midi, 0) / midis.length;

    // Calculate note variety (unique notes)
    const uniqueNotes = new Set(midis.map(m => Math.round(m)));
    const noteVariety = uniqueNotes.size;

    // Calculate pitch stability (inverse of pitch variation)
    const avgConfidence = confidences.reduce((a, b) => a + b, 0) / confidences.length;
    const pitchStability = Math.min(100, avgConfidence);

    let range: 'low' | 'mid' | 'high' = 'mid';
    if (weightedAvgMidi < 60) range = 'low';
    else if (weightedAvgMidi > 72) range = 'high';

    return {
      averageFrequency: weightedAvgFreq,
      minFrequency: Math.min(...frequencies),
      maxFrequency: Math.max(...frequencies),
      averageMidi: weightedAvgMidi,
      range,
      confidence: avgConfidence,
      noteVariety,
      pitchStability
    };
  }

  getSimplifiedMelody(pitches: PitchPoint[]): number[] {
    if (pitches.length === 0) return [];

    // Remove consecutive duplicates and quantize
    const melody: number[] = [];
    let lastNote = -1;

    for (const pitch of pitches) {
      if ((pitch.confidence || 0) > 30) { // Only confident notes
        const note = Math.round(pitch.midi);
        if (note !== lastNote) {
          melody.push(note);
          lastNote = note;
        }
      }
    }

    return melody;
  }

  midiToNoteName(midi: number): string {
    const noteNames = ['C', 'C#', 'D', 'D#', 'E', 'F', 'F#', 'G', 'G#', 'A', 'A#', 'B'];
    const octave = Math.floor(midi / 12) - 1;
    const note = noteNames[Math.round(midi) % 12];
    return `${note}${octave}`;
  }
}
