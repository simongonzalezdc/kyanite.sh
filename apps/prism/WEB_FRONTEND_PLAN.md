# Prism.sh Web Frontend - Implementation Plan

## Executive Summary

This document outlines a comprehensive plan to create a browser-based version of Prism.sh, a terminal color palette design tool. The plan evaluates two primary approaches: WebAssembly (WASM) compilation of existing Go code vs. a complete JavaScript/TypeScript rewrite.

## Current State Analysis

### Existing Codebase
- **Language**: Go 1.24.7
- **UI Framework**: Bubble Tea (TUI framework)
- **Dependencies**:
  - Lipgloss (terminal styling)
  - Bubbles (text input components)
  - Minimal external dependencies for core logic
- **Architecture**: Clean separation of concerns
  - `/internal/color` - Color conversion and operations
  - `/internal/palette` - Palette generation algorithms
  - `/internal/wcag` - Accessibility calculations
  - `/internal/export` - Export formatters
  - `/internal/ui` - Terminal UI components
  - `/internal/storage` - File I/O operations

### Code Analysis
- **Reusable Core Logic** (~70% of codebase):
  - Color theory algorithms
  - WCAG calculations
  - Palette generation rules
  - Export formatters (JSON, CSS, TOML)
  - Color database (147 named colors)

- **Platform-Specific** (~30% of codebase):
  - Terminal UI rendering
  - File system operations
  - Keyboard input handling

## Approach Evaluation

### Option A: WebAssembly (WASM)

#### Advantages
- Reuse existing Go codebase (70% code reuse)
- Proven algorithms already tested
- Single codebase maintenance
- Native performance for calculations
- Type safety from Go

#### Disadvantages
- WASM bundle size (~2-5MB compressed)
- DOM manipulation awkward from Go
- Limited browser API access
- Requires JavaScript bridge layer
- Debugging complexity
- Slower initial load time

#### Technical Feasibility
```go
// Example WASM export structure
//go:build wasm

package main

import (
    "syscall/js"
    "github.com/kyanite/prism/internal/color"
    "github.com/kyanite/prism/internal/palette"
)

func generatePalette(this js.Value, args []js.Value) interface{} {
    hexColor := args[0].String()
    rule := args[1].String()

    base, _ := color.ParseHex(hexColor)
    pal := palette.Generate(base, rule)

    return js.ValueOf(exportToJS(pal))
}

func main() {
    js.Global().Set("generatePalette", js.FuncOf(generatePalette))
    <-make(chan bool)
}
```

### Option B: JavaScript/TypeScript Rewrite

#### Advantages
- Native browser integration
- Smaller bundle size with tree-shaking
- Rich ecosystem (React/Vue/Svelte)
- Easier DOM manipulation
- Better developer tools
- Progressive Web App capabilities
- Faster initial load

#### Disadvantages
- Complete rewrite required
- Port all algorithms
- Maintain two codebases
- Re-test all calculations
- Loss of Go type safety (mitigated by TypeScript)

#### Recommended Stack
- **Language**: TypeScript 5.x
- **Framework**: React 18+ or Svelte 4+
- **Build Tool**: Vite
- **Styling**: Tailwind CSS + CSS-in-JS for dynamic colors
- **State Management**: Zustand or React Context
- **Testing**: Vitest + Testing Library
- **PWA**: Workbox

## Recommendation: Hybrid Approach

**Best Strategy**: JavaScript/TypeScript rewrite with Go reference implementation

### Rationale
1. **User Experience**: Native web performance, smaller bundle, PWA capabilities
2. **Maintainability**: Web-native codebase easier to maintain and extend
3. **Algorithm Integrity**: Keep Go version as reference implementation and testing source
4. **Team Skills**: Larger pool of web developers vs. Go+WASM specialists

### Migration Path
1. **Phase 1**: Port core algorithms to TypeScript, validate against Go tests
2. **Phase 2**: Build web UI components
3. **Phase 3**: Add web-specific features
4. **Phase 4**: Deploy and iterate

## Implementation Phases

### Phase 1: Core Library (4-6 weeks)

#### 1.1 Color Module
Port `/internal/color` to TypeScript:
```typescript
// color.ts
export interface RGB {
  r: number; // 0-255
  g: number;
  b: number;
}

export interface HSL {
  h: number; // 0-360
  s: number; // 0-100
  l: number; // 0-100
}

export class Color {
  constructor(
    public rgb: RGB,
    public hsl: HSL,
    public hex: string
  ) {}

  static fromHex(hex: string): Color {
    // Port from color.ParseHex
  }

  static fromRGB(r: number, g: number, b: number): Color {
    // Port conversion logic
  }

  lighten(amount: number): Color {
    // Port from color operations
  }

  darken(amount: number): Color {}
  saturate(amount: number): Color {}
  // ... other operations
}
```

#### 1.2 Palette Generator
Port `/internal/palette`:
```typescript
// palette.ts
export type HarmonyRule =
  | 'monochromatic'
  | 'complementary'
  | 'analogous'
  | 'triadic'
  | 'tetradic'
  | 'split-complementary';

export interface Palette {
  id: string;
  name: string;
  colors: Color[];
  rule: HarmonyRule;
  createdAt: Date;
}

export class PaletteGenerator {
  generate(baseColor: Color, rule: HarmonyRule, count?: number): Palette {
    switch (rule) {
      case 'monochromatic':
        return this.generateMonochromatic(baseColor, count);
      case 'complementary':
        return this.generateComplementary(baseColor);
      // ... implement all rules
    }
  }

  private generateMonochromatic(base: Color, count: number): Palette {
    // Port from palette.Generate
  }
}
```

#### 1.3 WCAG Module
Port `/internal/wcag`:
```typescript
// wcag.ts
export interface ContrastResult {
  ratio: number;
  level: 'AAA' | 'AA' | 'FAIL';
  passedAA: boolean;
  passedAAA: boolean;
}

export class WCAGValidator {
  validate(foreground: Color, background: Color): ContrastResult {
    const fgLuminance = this.relativeLuminance(foreground.rgb);
    const bgLuminance = this.relativeLuminance(background.rgb);

    const ratio = this.contrastRatio(fgLuminance, bgLuminance);

    return {
      ratio,
      passedAA: ratio >= 4.5,
      passedAAA: ratio >= 7.0,
      level: ratio >= 7.0 ? 'AAA' : ratio >= 4.5 ? 'AA' : 'FAIL'
    };
  }

  private relativeLuminance(rgb: RGB): number {
    // Port gamma correction and luminance calculation
  }

  private contrastRatio(l1: number, l2: number): number {
    const lighter = Math.max(l1, l2);
    const darker = Math.min(l1, l2);
    return (lighter + 0.05) / (darker + 0.05);
  }
}
```

#### 1.4 Testing Strategy
```typescript
// color.test.ts
import { describe, it, expect } from 'vitest';
import { Color } from './color';

describe('Color', () => {
  it('should parse hex colors correctly', () => {
    const color = Color.fromHex('#FF5733');
    expect(color.rgb).toEqual({ r: 255, g: 87, b: 51 });
  });

  it('should convert RGB to HSL accurately', () => {
    const color = Color.fromRGB(255, 0, 0);
    expect(color.hsl.h).toBeCloseTo(0);
    expect(color.hsl.s).toBeCloseTo(100);
    expect(color.hsl.l).toBeCloseTo(50);
  });

  // Port all tests from Go test suite
});
```

### Phase 2: Web UI (6-8 weeks)

#### 2.1 Technology Stack
```json
{
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "zustand": "^4.5.0",
    "tailwindcss": "^3.4.0",
    "framer-motion": "^11.0.0",
    "lucide-react": "^0.300.0"
  },
  "devDependencies": {
    "typescript": "^5.3.0",
    "vite": "^5.0.0",
    "@vitejs/plugin-react": "^4.2.0",
    "vitest": "^1.0.0",
    "@testing-library/react": "^14.1.0"
  }
}
```

#### 2.2 Component Architecture
```
src/
├── lib/                    # Core library (ported from Go)
│   ├── color.ts
│   ├── palette.ts
│   ├── wcag.ts
│   ├── export.ts
│   └── named-colors.ts
├── components/
│   ├── ui/                 # Reusable UI components
│   │   ├── Button.tsx
│   │   ├── Input.tsx
│   │   ├── ColorSwatch.tsx
│   │   └── Modal.tsx
│   ├── ColorWheel/         # Main screens
│   │   ├── ColorWheel.tsx
│   │   ├── HueRing.tsx
│   │   └── SaturationLightness.tsx
│   ├── PaletteGenerator/
│   │   ├── PaletteGenerator.tsx
│   │   ├── RuleSelector.tsx
│   │   └── PaletteDisplay.tsx
│   ├── WCAGChecker/
│   │   ├── WCAGChecker.tsx
│   │   └── ContrastDisplay.tsx
│   ├── ColorSearch/
│   │   ├── ColorSearch.tsx
│   │   └── SearchResults.tsx
│   └── PaletteManager/
│       ├── PaletteManager.tsx
│       └── PaletteList.tsx
├── store/
│   └── useStore.ts         # Global state management
├── hooks/
│   ├── useLocalStorage.ts
│   ├── useClipboard.ts
│   └── useKeyboardShortcuts.ts
└── App.tsx
```

#### 2.3 State Management
```typescript
// store/useStore.ts
import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface AppState {
  // Current state
  currentColor: Color | null;
  currentPalette: Palette | null;
  savedPalettes: Palette[];
  theme: 'light' | 'dark' | 'amber-night';

  // Actions
  setCurrentColor: (color: Color) => void;
  generatePalette: (rule: HarmonyRule) => void;
  savePalette: (palette: Palette) => void;
  deletePalette: (id: string) => void;
  exportPalette: (format: ExportFormat) => void;
}

export const useStore = create<AppState>()(
  persist(
    (set, get) => ({
      currentColor: null,
      currentPalette: null,
      savedPalettes: [],
      theme: 'amber-night',

      setCurrentColor: (color) => set({ currentColor: color }),

      generatePalette: (rule) => {
        const { currentColor } = get();
        if (!currentColor) return;

        const generator = new PaletteGenerator();
        const palette = generator.generate(currentColor, rule);
        set({ currentPalette: palette });
      },

      savePalette: (palette) => {
        set((state) => ({
          savedPalettes: [...state.savedPalettes, palette]
        }));
      },

      // ... other actions
    }),
    {
      name: 'prism-storage',
      partialize: (state) => ({
        savedPalettes: state.savedPalettes,
        theme: state.theme
      })
    }
  )
);
```

#### 2.4 Example Component
```typescript
// components/ColorWheel/ColorWheel.tsx
import { useState } from 'react';
import { useStore } from '@/store/useStore';
import { Color } from '@/lib/color';
import { HueRing } from './HueRing';
import { SaturationLightness } from './SaturationLightness';

export function ColorWheel() {
  const [hue, setHue] = useState(0);
  const [saturation, setSaturation] = useState(100);
  const [lightness, setLightness] = useState(50);
  const { setCurrentColor } = useStore();

  const color = Color.fromHSL(hue, saturation, lightness);

  const handleColorChange = () => {
    setCurrentColor(color);
  };

  return (
    <div className="flex flex-col items-center gap-8 p-8">
      <h1 className="text-3xl font-bold">Color Wheel</h1>

      <div className="relative w-80 h-80">
        <HueRing
          value={hue}
          onChange={setHue}
        />
        <SaturationLightness
          hue={hue}
          saturation={saturation}
          lightness={lightness}
          onSaturationChange={setSaturation}
          onLightnessChange={setLightness}
        />
      </div>

      <div className="flex items-center gap-4">
        <div
          className="w-24 h-24 rounded-lg shadow-lg"
          style={{ backgroundColor: color.hex }}
        />
        <div className="text-lg font-mono">
          <div>HEX: {color.hex}</div>
          <div>RGB: {color.rgb.r}, {color.rgb.g}, {color.rgb.b}</div>
          <div>HSL: {hue}°, {saturation}%, {lightness}%</div>
        </div>
      </div>

      <button
        onClick={handleColorChange}
        className="px-6 py-3 bg-blue-600 text-white rounded-lg"
      >
        Use This Color
      </button>
    </div>
  );
}
```

### Phase 3: Web-Specific Features (4-6 weeks)

#### 3.1 Progressive Web App
```typescript
// vite.config.ts
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import { VitePWA } from 'vite-plugin-pwa';

export default defineConfig({
  plugins: [
    react(),
    VitePWA({
      registerType: 'autoUpdate',
      manifest: {
        name: 'Prism.sh - Color Palette Designer',
        short_name: 'Prism',
        description: 'Professional color palette design tool',
        theme_color: '#FF6B35',
        icons: [
          {
            src: '/icon-192.png',
            sizes: '192x192',
            type: 'image/png'
          },
          {
            src: '/icon-512.png',
            sizes: '512x512',
            type: 'image/png'
          }
        ]
      },
      workbox: {
        runtimeCaching: [
          {
            urlPattern: /^https:\/\/fonts\.googleapis\.com\/.*/i,
            handler: 'CacheFirst',
            options: {
              cacheName: 'google-fonts-cache',
              expiration: {
                maxEntries: 10,
                maxAgeSeconds: 60 * 60 * 24 * 365 // 1 year
              }
            }
          }
        ]
      }
    })
  ]
});
```

#### 3.2 Enhanced Features
- **Drag & Drop**: Upload images to extract color palettes
- **Real-time Collaboration**: Share palettes via URL
- **Color Picker**: Browser native color picker integration
- **Export to Figma/Sketch**: Plugin integration
- **Accessibility Testing**: Live preview with different vision types
- **Gradient Generator**: Smooth gradients between colors
- **Pattern Preview**: See colors in UI mockups

#### 3.3 Mobile Responsiveness
```typescript
// hooks/useResponsive.ts
import { useMediaQuery } from '@/hooks/useMediaQuery';

export function useResponsive() {
  const isMobile = useMediaQuery('(max-width: 768px)');
  const isTablet = useMediaQuery('(min-width: 769px) and (max-width: 1024px)');
  const isDesktop = useMediaQuery('(min-width: 1025px)');

  return { isMobile, isTablet, isDesktop };
}
```

### Phase 4: Deployment & Infrastructure (2-3 weeks)

#### 4.1 Hosting Options

**Option A: Vercel (Recommended)**
- Zero-config deployment
- Automatic HTTPS
- Edge network CDN
- Free tier sufficient for MVP
- GitHub integration

**Option B: Cloudflare Pages**
- Free tier with unlimited bandwidth
- Global CDN
- Built-in analytics
- DDoS protection

**Option C: Self-hosted**
- Docker container
- Nginx reverse proxy
- GitHub Actions CI/CD

#### 4.2 CI/CD Pipeline
```yaml
# .github/workflows/deploy.yml
name: Deploy

on:
  push:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: '20'
      - run: npm ci
      - run: npm test
      - run: npm run build

  deploy:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: amondnet/vercel-action@v25
        with:
          vercel-token: ${{ secrets.VERCEL_TOKEN }}
          vercel-org-id: ${{ secrets.ORG_ID }}
          vercel-project-id: ${{ secrets.PROJECT_ID }}
```

#### 4.3 Analytics & Monitoring
- **Vercel Analytics**: Page views, user engagement
- **Sentry**: Error tracking and performance monitoring
- **Google Analytics**: User behavior insights
- **Lighthouse CI**: Performance regression testing

## Feature Parity Matrix

| Feature | Terminal | Web | Priority |
|---------|----------|-----|----------|
| Color Wheel | ✅ | Phase 2 | High |
| Palette Generator | ✅ | Phase 2 | High |
| WCAG Checker | ✅ | Phase 2 | High |
| Color Search | ✅ | Phase 2 | Medium |
| Save Palettes | ✅ | Phase 2 | High |
| Export (JSON/CSS/TOML) | ✅ | Phase 2 | High |
| Theme Switching | ✅ | Phase 2 | Medium |
| Keyboard Shortcuts | ✅ | Phase 3 | Medium |
| Clipboard Integration | ✅ | Phase 2 | High |
| **Web-Only Features** | | | |
| Image Color Extraction | ❌ | Phase 3 | Medium |
| Gradient Generator | ❌ | Phase 3 | Medium |
| Share via URL | ❌ | Phase 3 | High |
| PWA/Offline Mode | ❌ | Phase 3 | Medium |
| Color Blindness Preview | ❌ | Phase 4 | Low |
| Figma/Sketch Export | ❌ | Phase 4 | Low |

## Migration Timeline

```
Month 1-2: Core Library Development
├─ Week 1-2: Color module + tests
├─ Week 3-4: Palette generator + tests
├─ Week 5-6: WCAG module + tests
└─ Week 7-8: Export formatters + integration tests

Month 3-4: UI Development
├─ Week 9-10: Component library + design system
├─ Week 11-12: Color Wheel screen
├─ Week 13-14: Palette Generator screen
└─ Week 15-16: WCAG Checker + Color Search

Month 5: Polish & Web Features
├─ Week 17-18: PWA setup + offline support
├─ Week 19: Mobile responsiveness
└─ Week 20: Performance optimization

Month 6: Launch
├─ Week 21-22: Beta testing + bug fixes
├─ Week 23: Documentation + marketing site
└─ Week 24: Public launch
```

## Risk Assessment

### Technical Risks
| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| Algorithm porting errors | Medium | High | Comprehensive test suite, reference Go impl |
| Performance issues | Low | Medium | Profiling, code splitting, lazy loading |
| Browser compatibility | Low | Low | Target modern browsers, polyfills if needed |
| Bundle size bloat | Medium | Medium | Tree-shaking, code splitting, lazy components |

### Business Risks
| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| Maintenance overhead (2 codebases) | High | Medium | Share test cases, consider eventual deprecation of CLI |
| User migration | Low | Low | Keep both versions, cross-promotion |
| Feature divergence | Medium | Medium | Shared roadmap, regular sync |

## Success Metrics

### Phase 1 (Core Library)
- ✅ 100% test coverage on core algorithms
- ✅ <1% deviation from Go implementation results
- ✅ All Go tests ported and passing

### Phase 2 (UI Launch)
- 🎯 <3s initial page load (3G network)
- 🎯 >90 Lighthouse score
- 🎯 100 active users within first month
- 🎯 <5% bounce rate

### Phase 3 (Growth)
- 🎯 1000+ monthly active users
- 🎯 >50% mobile usage
- 🎯 <0.1% error rate
- 🎯 PWA installation rate >10%

## Budget Estimate

### Development Costs
- Core Library Development: 200-300 hours
- UI Development: 300-400 hours
- Testing & QA: 100-150 hours
- Documentation: 50-75 hours
- **Total**: 650-925 hours

### Infrastructure Costs (Annual)
- Hosting (Vercel Pro): $240/year
- Domain: $15/year
- Monitoring (Sentry): $26/month = $312/year
- **Total**: ~$570/year (can start with free tiers)

## Open Questions

1. **Branding**: Same name (Prism.sh) or separate brand?
2. **Monetization**: Free forever, freemium, or paid?
3. **API**: Should we expose a public API?
4. **Backend**: Do we need a backend for sync/collaboration features?
5. **Mobile App**: Native mobile apps (React Native) in future?

## Conclusion

**Recommended Path**: TypeScript/React rewrite with phased rollout

The web version of Prism.sh presents an opportunity to reach a much larger audience while maintaining the quality and accuracy of the terminal version. By porting the core algorithms to TypeScript and building a modern web UI, we can deliver a superior user experience while keeping the CLI version as a reference implementation.

The estimated timeline of 6 months for a full-featured launch is realistic with dedicated development effort. Starting with a free Vercel deployment keeps initial costs minimal while we validate market fit.

Next steps:
1. Set up TypeScript project structure
2. Port color conversion algorithms with comprehensive tests
3. Build basic color wheel component as proof of concept
4. Iterate based on user feedback
