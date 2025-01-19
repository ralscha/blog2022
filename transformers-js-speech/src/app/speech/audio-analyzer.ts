interface AudioAnalyzerOptions {
  energyThreshold?: number;
  frameSize?: number;
  minDuration?: number;
  sampleRate?: number;
}

interface AnalysisDetails {
  totalFrames: number;
  consecutiveSpeechSamples: number;
  threshold: number;
  minDuration: number;
  averageEnergy: number;
}

interface AnalysisResult {
  hasSpeech: boolean;
  details: AnalysisDetails;
}

export class AudioAnalyzer {
  private readonly energyThreshold: number;
  private readonly frameSize: number;
  private readonly minDuration: number;
  private readonly sampleRate: number;

  constructor(options: AudioAnalyzerOptions = {}) {
    this.energyThreshold = options.energyThreshold ?? 0.01;
    this.frameSize = options.frameSize ?? 2048;
    this.minDuration = options.minDuration ?? 0.1;
    this.sampleRate = options.sampleRate ?? 44100;
  }

  public analyzeAudioData(channelData: Float32Array): AnalysisResult {
    const frames: number = Math.floor(channelData.length / this.frameSize);
    const samplesNeededForMinDuration: number = Math.floor(this.sampleRate * this.minDuration);
    let consecutiveSpeechSamples: number = 0;
    let hasSpeech: boolean = false;
    let totalEnergy: number = 0;

    for (let i = 0; i < frames && !hasSpeech; i++) {
      const startIdx: number = i * this.frameSize;
      const frame: Float32Array = channelData.subarray(startIdx, startIdx + this.frameSize);

      const frameEnergy: number = this.calculateFrameEnergy(frame);
      totalEnergy += frameEnergy;

      if (frameEnergy > this.energyThreshold) {
        consecutiveSpeechSamples += this.frameSize;
        if (consecutiveSpeechSamples >= samplesNeededForMinDuration) {
          hasSpeech = true;
        }
      } else {
        consecutiveSpeechSamples = 0;
      }
    }

    return {
      hasSpeech,
      details: {
        totalFrames: frames,
        consecutiveSpeechSamples,
        threshold: this.energyThreshold,
        minDuration: this.minDuration,
        averageEnergy: totalEnergy / frames
      }
    };
  }

  private calculateFrameEnergy(frame: Float32Array): number {
    let sum: number = 0;
    for (let i = 0; i < frame.length; i++) {
      sum += frame[i] * frame[i];
    }
    return sum / frame.length;
  }
}
