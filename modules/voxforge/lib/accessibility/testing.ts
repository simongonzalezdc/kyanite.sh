// Note: axe-core is available through eslint-plugin-jsx-a11y
// For production use, install: npm install --save-dev axe-core @axe-core/react
declare global {
  interface Window {
    axe: any;
  }
}

// Accessibility testing configuration
export const accessibilityConfig = {
  rules: {
    // WCAG 2.1 AA compliance rules
    'color-contrast': { enabled: true },
    'keyboard-navigation': { enabled: true },
    'focus-order-semantics': { enabled: true },
    'aria-labels': { enabled: true },
    'heading-order': { enabled: true },
    'landmark-roles': { enabled: true },
    'link-name': { enabled: true },
    'button-name': { enabled: true },
    'image-alt': { enabled: true },
    'form-field-multiple-labels': { enabled: true },
    'input-button-name': { enabled: true },
    'label-title-only': { enabled: true },
    'object-alt': { enabled: true },
    'select-name': { enabled: true },
    'textarea-name': { enabled: true },
    'video-caption': { enabled: true },
    'audio-caption': { enabled: true },
    'duplicate-id': { enabled: true },
    'tabindex': { enabled: true },
    'skip-link': { enabled: true },
    'frame-title': { enabled: true },
    'html-has-lang': { enabled: true },
    'page-has-heading-one': { enabled: true },
    'region': { enabled: true }
  }
};

// Initialize axe-core for accessibility testing
export function initializeAccessibilityTesting() {
  if (typeof window !== 'undefined' && window.axe) {
    // axe-core is available
    if (typeof window.axe.configure === 'function') {
      window.axe.configure(accessibilityConfig);
    }
  }
}

// Run accessibility audit on a specific element
export async function runAccessibilityAudit(
  context?: string | Element,
  options?: any
): Promise<{
  violations: any[];
  passes: any[];
  incomplete: any[];
  url: string;
  timestamp: string;
}> {
  if (typeof window === 'undefined') {
    return {
      violations: [],
      passes: [],
      incomplete: [],
      url: '',
      timestamp: new Date().toISOString()
    };
  }

  try {
    if (typeof window !== 'undefined' && window.axe && typeof window.axe.run === 'function') {
      const results = await window.axe.run(context || document, options);
      
      return {
        violations: results.violations || [],
        passes: results.passes || [],
        incomplete: results.incomplete || [],
        url: window.location.href,
        timestamp: new Date().toISOString()
      };
    } else {
      // axe-core not available, return empty results
      console.warn('axe-core not available for accessibility testing');
      return {
        violations: [],
        passes: [],
        incomplete: [],
        url: window.location.href,
        timestamp: new Date().toISOString()
      };
    }
  } catch (error) {
    console.error('Accessibility audit failed:', error);
    return {
      violations: [],
      passes: [],
      incomplete: [],
      url: window.location.href,
      timestamp: new Date().toISOString()
    };
  }
}

// Generate accessibility report
export function generateAccessibilityReport(auditResults: any) {
  const { violations, passes, incomplete, url, timestamp } = auditResults;
  
  const report = {
    summary: {
      total: violations.length + passes.length + incomplete.length,
      violations: violations.length,
      passes: passes.length,
      incomplete: incomplete.length,
      wcagCompliance: violations.length === 0 ? 'AA' : 'Non-compliant'
    },
    url,
    timestamp,
    violations: violations.map((violation: any) => ({
      id: violation.id,
      impact: violation.impact,
      description: violation.description,
      help: violation.help,
      helpUrl: violation.helpUrl,
      nodes: violation.nodes.map((node: any) => ({
        target: node.target,
        html: node.html,
        failureSummary: node.failureSummary
      }))
    })),
    passes: passes.map((pass: any) => ({
      id: pass.id,
      description: pass.description,
      nodes: pass.nodes.map((node: any) => ({
        target: node.target,
        html: node.html
      }))
    })),
    incomplete: incomplete.map((item: any) => ({
      id: item.id,
      description: item.description,
      nodes: item.nodes.map((node: any) => ({
        target: node.target,
        html: node.html
      }))
    }))
  };

  return report;
}

// Manual testing checklist
export const manualTestingChecklist = {
  keyboardNavigation: [
    'Can all interactive elements be reached with Tab key?',
    'Is the focus order logical and intuitive?',
    'Can all actions be performed with keyboard only?',
    'Is focus clearly visible on all elements?',
    'Can users skip to main content?',
    'Are there keyboard traps that prevent navigation?',
    'Do modals trap focus appropriately?',
    'Can users navigate within forms using keyboard?'
  ],
  screenReader: [
    'Are all images have meaningful alt text?',
    'Are form fields properly labeled?',
    'Are headings used hierarchically?',
    'Are links descriptive when read out of context?',
    'Are dynamic content changes announced?',
    'Are error messages accessible?',
    'Are tables properly marked up?',
    'Are ARIA landmarks used appropriately?',
    'Is page language specified?',
    'Are custom controls properly announced?'
  ],
  visualAccessibility: [
    'Is text contrast at least 4.5:1 for normal text?',
    'Is large text contrast at least 3:1?',
    'Can text be resized to 200% without breaking layout?',
    'Is content usable in high contrast mode?',
    'Is content usable with reduced motion?',
    'Are color blind friendly palettes used?',
    'Is focus indication visible?',
    'Are status indicators not color-only?'
  ],
  audioAccessibility: [
    'Are visual indicators provided for audio events?',
    'Are captions available for audio content?',
    'Are transcripts provided for audio content?',
    'Are volume controls accessible?',
    'Are audio descriptions available for visual content?',
    'Is vibration feedback available for touch devices?',
    'Are alternative input methods available?'
  ],
  cognitiveAccessibility: [
    'Is content written in clear, simple language?',
    'Are instructions clear and easy to follow?',
    'Is consistent navigation provided?',
    'Are error messages helpful and specific?',
    'Is sufficient time provided to read and use content?',
    'Are complex processes broken into steps?',
    'Is help available when needed?',
    'Are predictable page layouts used?'
  ]
};

// Screen reader testing setup
export function setupScreenReaderTesting() {
  const testInstructions = {
    nvda: {
      name: 'NVDA (NonVisual Desktop Access)',
      platform: 'Windows',
      setup: [
        'Download and install NVDA from nvaccess.org',
        'Restart computer after installation',
        'Launch NVDA (Ctrl+Alt+N)',
        'Test with web browser (recommended: Firefox)'
      ],
      shortcuts: [
        'H: Move to next heading',
        '1-6: Move to heading level 1-6',
        'Tab: Move to next interactive element',
        'Enter: Activate button or link',
        'Space: Toggle checkbox or activate button',
        'Arrow keys: Navigate within content',
        'Ctrl+Home: Move to top of page',
        'Ctrl+End: Move to bottom of page',
        'B: Move to next button',
        'L: Move to next list',
        'T: Move to next table',
        'G: Move to next graphic',
        'Insert+F7: List all elements',
        'Insert+Q: Quit NVDA'
      ]
    },
    jaws: {
      name: 'JAWS (Job Access With Speech)',
      platform: 'Windows',
      setup: [
        'Download and install JAWS from freedomscientific.com',
        'Restart computer after installation',
        'Launch JAWS',
        'Test with web browser (recommended: Internet Explorer or Firefox)'
      ],
      shortcuts: [
        'H: Move to next heading',
        '1-6: Move to heading level 1-6',
        'Tab: Move to next interactive element',
        'Enter: Activate button or link',
        'Space: Toggle checkbox or activate button',
        'Arrow keys: Navigate within content',
        'Ctrl+Home: Move to top of page',
        'Ctrl+End: Move to bottom of page',
        'B: Move to next button',
        'Insert+F7: List all elements',
        'Insert+F5: Refresh screen',
        'Insert+Q: Quit JAWS'
      ]
    },
    voiceOver: {
      name: 'VoiceOver',
      platform: 'macOS/iOS',
      setup: [
        'On Mac: System Preferences > Accessibility > VoiceOver > Enable VoiceOver',
        'On iOS: Settings > Accessibility > VoiceOver > Enable VoiceOver',
        'Use Cmd+F5 to toggle on Mac',
        'Use triple-click Home button to toggle on iOS'
      ],
      shortcuts: [
        'VO+Right Arrow: Move to next element',
        'VO+Left Arrow: Move to previous element',
        'VO+Down Arrow: Enter element (interact)',
        'VO+Up Arrow: Exit element (stop interacting)',
        'VO+Space: Activate selected element',
        'VO+Shift+Down Arrow: Select text',
        'VO+Command+Space: Speak selected text',
        'VO+U: Rotor (quick navigation)',
        'VO+I: Item chooser',
        'VO+Command+F: Find',
        'VO+F8: VoiceOver utility'
      ]
    }
  };

  return testInstructions;
}

// Keyboard-only navigation testing
export function testKeyboardNavigation() {
  const testSteps = [
    {
      name: 'Tab Navigation',
      description: 'Test if all interactive elements can be reached with Tab',
      steps: [
        'Press Tab repeatedly to navigate through the page',
        'Verify all buttons, links, form fields are reachable',
        'Check that focus moves in logical order',
        'Ensure no elements are skipped'
      ]
    },
    {
      name: 'Shift+Tab Navigation',
      description: 'Test backward navigation',
      steps: [
        'Press Shift+Tab to navigate backwards',
        'Verify focus moves to previous element',
        'Check that navigation is consistent'
      ]
    },
    {
      name: 'Enter/Space Activation',
      description: 'Test element activation',
      steps: [
        'Navigate to buttons and press Enter',
        'Navigate to checkboxes and press Space',
        'Navigate to links and press Enter',
        'Verify all elements activate correctly'
      ]
    },
    {
      name: 'Arrow Key Navigation',
      description: 'Test arrow key navigation where applicable',
      steps: [
        'Test arrow keys in menus',
        'Test arrow keys in radio button groups',
        'Test arrow keys in sliders',
        'Test arrow keys in custom components'
      ]
    },
    {
      name: 'Escape Key',
      description: 'Test Escape key functionality',
      steps: [
        'Open modals and press Escape to close',
        'Test Escape in dropdown menus',
        'Test Escape in custom components',
        'Verify Escape cancels operations'
      ]
    },
    {
      name: 'Focus Management',
      description: 'Test focus visibility and management',
      steps: [
        'Verify focus is clearly visible on all elements',
        'Test focus in modals (should be trapped)',
        'Test focus returns after modal closes',
        'Check focus moves to new content after navigation'
      ]
    }
  ];

  return testSteps;
}

// Color contrast validation
export function validateColorContrast() {
  const contrastRequirements = {
    WCAG_AA: {
      normalText: 4.5,
      largeText: 3.0,
      graphicalObjects: 3.0
    },
    WCAG_AAA: {
      normalText: 7.0,
      largeText: 4.5,
      graphicalObjects: 4.5
    }
  };

  const calculateContrast = (foreground: string, background: string): number => {
    const getLuminance = (color: string) => {
      const rgb = parseInt(color.slice(1), 16);
      const r = (rgb >> 16) & 0xff;
      const g = (rgb >> 8) & 0xff;
      const b = rgb & 0xff;
      
      const [rs, gs, bs] = [r, g, b].map(c => {
        c = c / 255;
        return c <= 0.03928 ? c / 12.92 : Math.pow((c + 0.055) / 1.055, 2.4);
      });
      
      return 0.2126 * rs + 0.7152 * gs + 0.0722 * bs;
    };

    const l1 = getLuminance(foreground);
    const l2 = getLuminance(background);
    const contrast = (Math.max(l1, l2) + 0.05) / (Math.min(l1, l2) + 0.05);
    
    return Math.round(contrast * 100) / 100;
  };

  const validateContrast = (foreground: string, background: string, isLargeText = false) => {
    const ratio = calculateContrast(foreground, background);
    
    return {
      ratio,
      wcagAA: ratio >= (isLargeText ? contrastRequirements.WCAG_AA.largeText : contrastRequirements.WCAG_AA.normalText),
      wcagAAA: ratio >= (isLargeText ? contrastRequirements.WCAG_AAA.largeText : contrastRequirements.WCAG_AAA.normalText)
    };
  };

  return {
    contrastRequirements,
    calculateContrast,
    validateContrast
  };
}

// Automated accessibility test runner
export class AccessibilityTestRunner {
  private results: any[] = [];
  private isRunning = false;

  async runFullAudit() {
    if (this.isRunning) return;
    
    this.isRunning = true;
    
    try {
      // Run axe-core audit
      const axeResults = await runAccessibilityAudit();
      
      // Run custom tests
      const customResults = await this.runCustomTests();
      
      // Combine results
      const fullResults = {
        ...axeResults,
        customTests: customResults,
        timestamp: new Date().toISOString()
      };
      
      this.results.push(fullResults);
      
      return fullResults;
    } finally {
      this.isRunning = false;
    }
  }

  private async runCustomTests() {
    const customTests = [
      this.testSkipLinks(),
      this.testFocusManagement(),
      this.testAriaLabels(),
      this.testHeadingStructure(),
      this.testFormAccessibility()
    ];

    const results = await Promise.all(customTests);
    return results;
  }

  private async testSkipLinks() {
    const skipLinks = document.querySelectorAll('a[href^="#"]');
    const results = [];

    for (const link of skipLinks) {
      const href = link.getAttribute('href');
      if (href && href !== '#') {
        const target = document.querySelector(href);
        results.push({
          element: link.outerHTML,
          hasTarget: !!target,
          targetAccessible: target && this.isElementAccessible(target)
        });
      }
    }

    return { test: 'skip-links', results };
  }

  private async testFocusManagement() {
    const focusableElements = document.querySelectorAll(
      'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
    );

    const results: any[] = [];
    focusableElements.forEach(element => {
      results.push({
        element: element.outerHTML,
        hasTabIndex: element.hasAttribute('tabindex'),
        tabIndex: element.getAttribute('tabindex'),
        isFocusable: this.isElementFocusable(element)
      });
    });

    return { test: 'focus-management', results };
  }

  private async testAriaLabels() {
    const interactiveElements = document.querySelectorAll('button, input, select, textarea, a');
    const results: any[] = [];

    interactiveElements.forEach(element => {
      const hasLabel = element.hasAttribute('aria-label') ||
                     element.hasAttribute('aria-labelledby') ||
                     element.getAttribute('title') ||
                     (element as HTMLInputElement).labels?.length;

      results.push({
        element: element.outerHTML,
        hasLabel,
        labelType: this.getLabelType(element)
      });
    });

    return { test: 'aria-labels', results };
  }

  private async testHeadingStructure() {
    const headings = document.querySelectorAll('h1, h2, h3, h4, h5, h6');
    const results: any[] = [];

    let previousLevel = 0;
    headings.forEach(heading => {
      const level = parseInt(heading.tagName.substring(1));
      const isProperNesting = level <= previousLevel + 1;
      
      results.push({
        element: heading.outerHTML,
        level,
        isProperNesting,
        text: heading.textContent
      });
      
      previousLevel = level;
    });

    return { test: 'heading-structure', results };
  }

  private async testFormAccessibility() {
    const forms = document.querySelectorAll('form');
    const results: any[] = [];

    forms.forEach(form => {
      const inputs = form.querySelectorAll('input, select, textarea');
      const labels = form.querySelectorAll('label');
      const hasSubmitButton = form.querySelector('button[type="submit"], input[type="submit"]');
      
      results.push({
        form: form.outerHTML,
        inputCount: inputs.length,
        labelCount: labels.length,
        hasSubmitButton,
        inputsHaveLabels: Array.from(inputs).every(input => 
          this.hasAssociatedLabel(input as HTMLElement)
        )
      });
    });

    return { test: 'form-accessibility', results };
  }

  private isElementAccessible(element: Element): boolean {
    const style = window.getComputedStyle(element);
    return style.display !== 'none' && 
           style.visibility !== 'hidden' && 
           !element.hasAttribute('aria-hidden');
  }

  private isElementFocusable(element: Element): boolean {
    const style = window.getComputedStyle(element);
    return style.display !== 'none' && 
           style.visibility !== 'hidden' && 
           !element.hasAttribute('disabled') &&
           !element.hasAttribute('aria-disabled');
  }

  private getLabelType(element: Element): string {
    if (element.hasAttribute('aria-label')) return 'aria-label';
    if (element.hasAttribute('aria-labelledby')) return 'aria-labelledby';
    if (element.getAttribute('title')) return 'title';
    if ((element as HTMLInputElement).labels?.length) return 'label';
    return 'none';
  }

  private hasAssociatedLabel(input: HTMLElement): boolean {
    const id = input.getAttribute('id');
    if (id) {
      const label = document.querySelector(`label[for="${id}"]`);
      if (label) return true;
    }
    
    return !!input.hasAttribute('aria-label') ||
           !!input.hasAttribute('aria-labelledby') ||
           !!input.getAttribute('title') ||
           !!((input as HTMLInputElement).labels?.length);
  }

  getResults() {
    return this.results;
  }

  clearResults() {
    this.results = [];
  }
}

// Export test runner instance
export const accessibilityTestRunner = new AccessibilityTestRunner();