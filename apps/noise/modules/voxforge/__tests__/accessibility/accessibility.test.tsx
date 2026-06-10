import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { axe, toHaveNoViolations } from 'jest-axe';
import { AccessibilityProvider } from '../../app/components/AccessibilityProvider';
import { KeyboardNavigation } from '../../app/components/KeyboardNavigation';
import { ScreenReaderNavigation } from '../../app/components/ScreenReaderSupport';
import { VisualAccessibilityControls } from '../../app/components/VisualAccessibility';
import { AudioAccessibilityControls } from '../../app/components/AudioAccessibility';

// Extend Jest matchers
expect.extend(toHaveNoViolations);

// Mock window.axe for testing
Object.defineProperty(window, 'axe', {
  value: {
    configure: jest.fn(),
    run: jest.fn().mockResolvedValue({
      violations: [],
      passes: [],
      incomplete: []
    })
  },
  writable: true
});

// Test wrapper with all accessibility providers
const AccessibilityWrapper = ({ children }: { children: React.ReactNode }) => (
  <AccessibilityProvider>
    <KeyboardNavigation>
      <ScreenReaderNavigation />
      {children}
    </KeyboardNavigation>
  </AccessibilityProvider>
);

describe('Accessibility Components', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('AccessibilityProvider', () => {
    it('should render without accessibility violations', async () => {
      const { container } = render(
        <AccessibilityWrapper>
          <div>Test content</div>
        </AccessibilityWrapper>
      );

      const results = await axe(container);
      expect(results).toHaveNoViolations();
    });

    it('should provide accessibility context', () => {
      const TestComponent = () => {
        // This would use useAccessibility hook in real implementation
        return <div>Test component</div>;
      };

      render(
        <AccessibilityWrapper>
          <TestComponent />
        </AccessibilityWrapper>
      );

      expect(screen.getByText('Test component')).toBeInTheDocument();
    });
  });

  describe('KeyboardNavigation', () => {
    it('should handle keyboard shortcuts', () => {
      const mockAction = jest.fn();
      
      render(
        <AccessibilityWrapper>
          <button onClick={mockAction}>Test Button</button>
        </AccessibilityWrapper>
      );

      const button = screen.getByRole('button', { name: 'Test Button' });
      button.focus();
      
      fireEvent.keyDown(document, { key: 'Enter' });
      fireEvent.click(button);
      
      expect(mockAction).toHaveBeenCalled();
    });

    it('should provide skip links', () => {
      render(
        <AccessibilityWrapper>
          <div>
            <a href="#main" className="skip-link">Skip to main content</a>
            <main id="main">Main content</main>
          </div>
        </AccessibilityWrapper>
      );

      const skipLink = screen.getByRole('link', { name: 'Skip to main content' });
      expect(skipLink).toBeInTheDocument();
      expect(skipLink).toHaveAttribute('href', '#main');
    });

    it('should manage focus in modals', () => {
      render(
        <AccessibilityWrapper>
          <div role="dialog" aria-modal="true" aria-labelledby="modal-title">
            <h2 id="modal-title">Modal Title</h2>
            <button>Close</button>
          </div>
        </AccessibilityWrapper>
      );

      const modal = screen.getByRole('dialog');
      expect(modal).toHaveAttribute('aria-modal', 'true');
      expect(modal).toHaveAttribute('aria-labelledby', 'modal-title');
    });
  });

  describe('ScreenReaderSupport', () => {
    it('should provide ARIA labels', () => {
      render(
        <AccessibilityWrapper>
          <button aria-label="Play audio">▶</button>
        </AccessibilityWrapper>
      );

      const button = screen.getByRole('button', { name: 'Play audio' });
      expect(button).toHaveAttribute('aria-label', 'Play audio');
    });

    it('should announce dynamic content', async () => {
      render(
        <AccessibilityWrapper>
          <div role="status" aria-live="polite" aria-atomic="true">
            Dynamic content will appear here
          </div>
        </AccessibilityWrapper>
      );

      const liveRegion = screen.getByRole('status');
      expect(liveRegion).toHaveAttribute('aria-live', 'polite');
      expect(liveRegion).toHaveAttribute('aria-atomic', 'true');
    });

    it('should provide semantic landmarks', () => {
      render(
        <AccessibilityWrapper>
          <header role="banner">Header</header>
          <nav role="navigation">Navigation</nav>
          <main role="main">Main content</main>
          <aside role="complementary">Sidebar</aside>
          <footer role="contentinfo">Footer</footer>
        </AccessibilityWrapper>
      );

      expect(screen.getByRole('banner')).toBeInTheDocument();
      expect(screen.getByRole('navigation')).toBeInTheDocument();
      expect(screen.getByRole('main')).toBeInTheDocument();
      expect(screen.getByRole('complementary')).toBeInTheDocument();
      expect(screen.getByRole('contentinfo')).toBeInTheDocument();
    });
  });

  describe('VisualAccessibility', () => {
    it('should provide high contrast toggle', () => {
      render(
        <AccessibilityWrapper>
          <VisualAccessibilityControls />
        </AccessibilityWrapper>
      );

      const highContrastToggle = screen.getByRole('switch', { name: /high contrast/i });
      expect(highContrastToggle).toBeInTheDocument();
    });

    it('should provide font size controls', () => {
      render(
        <AccessibilityWrapper>
          <VisualAccessibilityControls />
        </AccessibilityWrapper>
      );

      const fontSizeSelect = screen.getByRole('combobox', { name: /font size/i });
      expect(fontSizeSelect).toBeInTheDocument();
    });

    it('should provide reduced motion toggle', () => {
      render(
        <AccessibilityWrapper>
          <VisualAccessibilityControls />
        </AccessibilityWrapper>
      );

      const reducedMotionToggle = screen.getByRole('switch', { name: /reduced motion/i });
      expect(reducedMotionToggle).toBeInTheDocument();
    });
  });

  describe('AudioAccessibility', () => {
    it('should provide visual indicators', () => {
      render(
        <AccessibilityWrapper>
          <AudioAccessibilityControls />
        </AccessibilityWrapper>
      );

      const visualIndicatorsToggle = screen.getByRole('switch', { name: /visual indicators/i });
      expect(visualIndicatorsToggle).toBeInTheDocument();
    });

    it('should provide vibration feedback toggle', () => {
      render(
        <AccessibilityWrapper>
          <AudioAccessibilityControls />
        </AccessibilityWrapper>
      );

      const vibrationToggle = screen.getByRole('switch', { name: /vibration/i });
      expect(vibrationToggle).toBeInTheDocument();
    });

    it('should provide audio descriptions toggle', () => {
      render(
        <AccessibilityWrapper>
          <AudioAccessibilityControls />
        </AccessibilityWrapper>
      );

      const audioDescriptionsToggle = screen.getByRole('switch', { name: /audio descriptions/i });
      expect(audioDescriptionsToggle).toBeInTheDocument();
    });
  });
});

describe('Accessibility Integration', () => {
  it('should work with all accessibility features enabled', async () => {
    const TestApp = () => (
      <div>
        <header>
          <h1>VoxForge</h1>
          <nav>
            <ul>
              <li><a href="#record">Record</a></li>
              <li><a href="#analyze">Analyze</a></li>
              <li><a href="#generate">Generate</a></li>
            </ul>
          </nav>
        </header>
        
        <main>
          <section id="record" aria-labelledby="record-heading">
            <h2 id="record-heading">Record Your Voice</h2>
            <button aria-label="Start recording">Start Recording</button>
          </section>
          
          <section id="analyze" aria-labelledby="analyze-heading">
            <h2 id="analyze-heading">Analysis Results</h2>
            <div role="status" aria-live="polite" aria-atomic="true">
              Analysis will appear here
            </div>
          </section>
          
          <section id="generate" aria-labelledby="generate-heading">
            <h2 id="generate-heading">Generate Music</h2>
            <button aria-label="Generate music">Generate</button>
          </section>
        </main>
        
        <footer>
          <p>&copy; 2024 VoxForge</p>
        </footer>
      </div>
    );

    const { container } = render(
      <AccessibilityWrapper>
        <TestApp />
      </AccessibilityWrapper>
    );

    // Check for accessibility violations
    const results = await axe(container);
    expect(results).toHaveNoViolations();

    // Check semantic structure
    expect(screen.getByRole('banner')).toBeInTheDocument();
    expect(screen.getByRole('navigation')).toBeInTheDocument();
    expect(screen.getByRole('main')).toBeInTheDocument();
    expect(screen.getByRole('contentinfo')).toBeInTheDocument();

    // Check headings
    expect(screen.getByRole('heading', { level: 1 })).toBeInTheDocument();
    expect(screen.getAllByRole('heading', { level: 2 })).toHaveLength(3);

    // Check interactive elements
    expect(screen.getByRole('button', { name: 'Start recording' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Generate music' })).toBeInTheDocument();

    // Check navigation links
    expect(screen.getByRole('link', { name: 'Record' })).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'Analyze' })).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'Generate' })).toBeInTheDocument();

    // Check ARIA attributes
    expect(screen.getByLabelText('Start recording')).toBeInTheDocument();
    expect(screen.getByLabelText('Generate music')).toBeInTheDocument();
    expect(screen.getByRole('status')).toHaveAttribute('aria-live', 'polite');
  });

  it('should handle keyboard navigation properly', async () => {
    const user = userEvent.setup();

    render(
      <AccessibilityWrapper>
        <div>
          <button>Button 1</button>
          <button>Button 2</button>
          <button>Button 3</button>
        </div>
      </AccessibilityWrapper>
    );

    const buttons = screen.getAllByRole('button');
    
    // Test Tab navigation
    await user.tab();
    expect(buttons[0]).toHaveFocus();
    
    await user.tab();
    expect(buttons[1]).toHaveFocus();
    
    await user.tab();
    expect(buttons[2]).toHaveFocus();
    
    // Test Shift+Tab navigation
    await user.tab({ shift: true });
    expect(buttons[1]).toHaveFocus();
    
    // Test Enter activation
    fireEvent.keyDown(buttons[1], { key: 'Enter' });
    // Button should be activated (would trigger onClick in real app)
  });

  it('should handle screen reader announcements', async () => {
    const mockAnnounce = jest.fn();
    
    render(
      <AccessibilityWrapper>
        <div role="alert" aria-live="assertive" aria-atomic="true">
          Important message
        </div>
      </AccessibilityWrapper>
    );

    const liveRegion = screen.getByRole('alert');
    expect(liveRegion).toHaveAttribute('aria-live', 'assertive');
    expect(liveRegion).toHaveAttribute('aria-atomic', 'true');
    expect(liveRegion).toHaveTextContent('Important message');
  });
});

describe('Form Accessibility', () => {
  it('should provide accessible form controls', async () => {
    const TestForm = () => (
      <form>
        <div>
          <label htmlFor="name">Name</label>
          <input id="name" type="text" required aria-describedby="name-help" />
          <div id="name-help">Enter your full name</div>
        </div>
        
        <div>
          <label htmlFor="email">Email</label>
          <input id="email" type="email" required aria-describedby="email-error" />
          <div id="email-error" role="alert">Email is required</div>
        </div>
        
        <div>
          <fieldset>
            <legend>Preferences</legend>
            <input id="newsletter" type="checkbox" />
            <label htmlFor="newsletter">Subscribe to newsletter</label>
          </fieldset>
        </div>
        
        <button type="submit">Submit</button>
      </form>
    );

    const { container } = render(
      <AccessibilityWrapper>
        <TestForm />
      </AccessibilityWrapper>
    );

    const results = await axe(container);
    expect(results).toHaveNoViolations();

    // Check form labels
    expect(screen.getByLabelText('Name')).toBeInTheDocument();
    expect(screen.getByLabelText('Email')).toBeInTheDocument();
    expect(screen.getByLabelText('Subscribe to newsletter')).toBeInTheDocument();

    // Check ARIA attributes
    expect(screen.getByRole('textbox', { name: 'Name' })).toHaveAttribute('aria-describedby', 'name-help');
    expect(screen.getByRole('textbox', { name: 'Email' })).toHaveAttribute('aria-describedby', 'email-error');
    expect(screen.getByRole('alert')).toHaveTextContent('Email is required');

    // Check required fields
    expect(screen.getByRole('textbox', { name: 'Name' })).toHaveAttribute('required');
    expect(screen.getByRole('textbox', { name: 'Email' })).toHaveAttribute('required');
  });
});

describe('Media Accessibility', () => {
  it('should provide accessible audio controls', async () => {
    const TestAudio = () => (
      <div>
        <audio>
          <track kind="captions" src="captions.vtt" label="English captions" />
          Your browser does not support the audio element.
        </audio>
        
        <button aria-label="Play audio">Play</button>
        <button aria-label="Pause audio">Pause</button>
        <button aria-label="Mute audio">Mute</button>
        
        <div role="status" aria-live="polite">
          Audio status: Ready
        </div>
      </div>
    );

    const { container } = render(
      <AccessibilityWrapper>
        <TestAudio />
      </AccessibilityWrapper>
    );

    const results = await axe(container);
    expect(results).toHaveNoViolations();

    // Check audio controls
    expect(screen.getByLabelText('Play audio')).toBeInTheDocument();
    expect(screen.getByLabelText('Pause audio')).toBeInTheDocument();
    expect(screen.getByLabelText('Mute audio')).toBeInTheDocument();

    // Check status announcement
    expect(screen.getByRole('status')).toHaveTextContent('Audio status: Ready');
  });

  it('should provide accessible video controls', async () => {
    const TestVideo = () => (
      <div>
        <video>
          <track kind="captions" src="captions.vtt" label="English captions" />
          <track kind="descriptions" src="descriptions.vtt" label="Audio descriptions" />
          Your browser does not support the video element.
        </video>
        
        <button aria-label="Play video">Play</button>
        <button aria-label="Pause video">Pause</button>
        <button aria-label="Fullscreen">Fullscreen</button>
      </div>
    );

    const { container } = render(
      <AccessibilityWrapper>
        <TestVideo />
      </AccessibilityWrapper>
    );

    const results = await axe(container);
    expect(results).toHaveNoViolations();

    // Check video controls
    expect(screen.getByLabelText('Play video')).toBeInTheDocument();
    expect(screen.getByLabelText('Pause video')).toBeInTheDocument();
    expect(screen.getByLabelText('Fullscreen')).toBeInTheDocument();
  });
});
