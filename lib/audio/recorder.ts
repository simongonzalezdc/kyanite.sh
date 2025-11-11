export class AudioRecorder {
  private mediaRecorder: MediaRecorder | null = null;
  private audioChunks: Blob[] = [];
  private audioContext: AudioContext;
  private analyser: AnalyserNode;
  private stream: MediaStream | null = null;

  constructor() {
    this.audioContext = new AudioContext({ sampleRate: 44100 });
    this.analyser = this.audioContext.createAnalyser();
    this.analyser.fftSize = 2048;
  }

  async requestPermission(): Promise<boolean> {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({
        audio: {
          sampleRate: 44100,
          channelCount: 1,
          echoCancellation: false,
          noiseSuppression: false,
          autoGainControl: false
        }
      });
      
      // Test and close immediately
      stream.getTracks().forEach(track => track.stop());
      return true;
    } catch (error) {
      console.error('Microphone permission denied:', error);
      return false;
    }
  }

  async startRecording(): Promise<void> {
    // Ensure AudioContext is resumed (required in browsers)
    if (this.audioContext.state === 'suspended') {
      await this.audioContext.resume();
    }

    this.stream = await navigator.mediaDevices.getUserMedia({
      audio: {
        sampleRate: 44100,
        channelCount: 1,
        echoCancellation: false,
        noiseSuppression: false,
        autoGainControl: false
      }
    });

    this.mediaRecorder = new MediaRecorder(this.stream, {
      mimeType: 'audio/webm'
    });
    
    this.audioChunks = [];

    this.mediaRecorder.ondataavailable = (event) => {
      if (event.data.size > 0) {
        this.audioChunks.push(event.data);
      }
    };

    // Connect to analyser for real-time visualization
    const source = this.audioContext.createMediaStreamSource(this.stream);
    source.connect(this.analyser);

    this.mediaRecorder.start(100); // Collect data every 100ms
  }

  stopRecording(): Promise<AudioBuffer> {
    return new Promise((resolve, reject) => {
      if (!this.mediaRecorder) {
        reject(new Error('No recording in progress'));
        return;
      }

      this.mediaRecorder.onstop = async () => {
        try {
          // Ensure AudioContext is resumed for decoding
          if (this.audioContext.state === 'suspended') {
            await this.audioContext.resume();
          }

          // Create blob from chunks
          const audioBlob = new Blob(this.audioChunks, { type: 'audio/webm' });
          
          // Convert to array buffer
          const arrayBuffer = await audioBlob.arrayBuffer();
          
          // Decode to audio buffer
          const audioBuffer = await this.audioContext.decodeAudioData(arrayBuffer);
          
          // Cleanup
          this.cleanup();
          
          resolve(audioBuffer);
        } catch (error) {
          reject(error);
        }
      };

      this.mediaRecorder.stop();
    });
  }

  getWaveformData(): Uint8Array {
    const dataArray = new Uint8Array(this.analyser.frequencyBinCount);
    this.analyser.getByteTimeDomainData(dataArray);
    return dataArray;
  }

  getFrequencyData(): Uint8Array {
    const dataArray = new Uint8Array(this.analyser.frequencyBinCount);
    this.analyser.getByteFrequencyData(dataArray);
    return dataArray;
  }

  private cleanup(): void {
    if (this.stream) {
      this.stream.getTracks().forEach(track => track.stop());
      this.stream = null;
    }
  }

  dispose(): void {
    this.cleanup();
    this.audioContext.close();
  }
}
