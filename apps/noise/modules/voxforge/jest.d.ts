/// <reference types="jest" />
/// <reference types="@testing-library/jest-dom" />

import '@testing-library/jest-dom';
import 'jest-axe';

declare global {
  namespace jest {
    interface Expect {
      extend(matchers: Record<string, any>): void;
      objectContaining(sample: Record<string, unknown>): any;
    }
    interface Matchers<R = void> {
      toHaveNoViolations(): R;
      toBeAudioBuffer(): R;
      toBePitchPoint(): R;
      toBeBPMAnalysis(): R;
    }
  }
}