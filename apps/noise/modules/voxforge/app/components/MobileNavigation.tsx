'use client';

import { useState, useEffect } from 'react';
import { Mic, Activity, Music, Download, Menu, X } from 'lucide-react';
import { usePathname } from 'next/navigation';

interface MobileNavigationProps {
  activeSection?: string;
  onSectionChange?: (section: string) => void;
}

export default function MobileNavigation({ activeSection = 'record', onSectionChange }: MobileNavigationProps) {
  const [isVisible, setIsVisible] = useState(true);
  const [lastScrollY, setLastScrollY] = useState(0);
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const pathname = usePathname();

  const navigationItems = [
    {
      id: 'record',
      label: 'Record',
      icon: <Mic size={20} />,
      description: 'Record audio'
    },
    {
      id: 'analyze',
      label: 'Analyze',
      icon: <Activity size={20} />,
      description: 'View analysis'
    },
    {
      id: 'generate',
      label: 'Generate',
      icon: <Music size={20} />,
      description: 'Create music'
    },
    {
      id: 'export',
      label: 'Export',
      icon: <Download size={20} />,
      description: 'Export files'
    }
  ];

  // Handle scroll behavior to hide/show navigation
  useEffect(() => {
    const handleScroll = () => {
      const currentScrollY = window.scrollY;
      
      // Show navigation when scrolling up or at top
      if (currentScrollY < lastScrollY || currentScrollY < 100) {
        setIsVisible(true);
      } 
      // Hide navigation when scrolling down
      else if (currentScrollY > lastScrollY && currentScrollY > 100) {
        setIsVisible(false);
      }
      
      setLastScrollY(currentScrollY);
    };

    window.addEventListener('scroll', handleScroll, { passive: true });
    return () => window.removeEventListener('scroll', handleScroll);
  }, [lastScrollY]);

  // Handle keyboard avoidance
  useEffect(() => {
    const handleKeyboardShow = () => setIsVisible(false);
    const handleKeyboardHide = () => setIsVisible(true);

    window.addEventListener('focusin', handleKeyboardShow);
    window.addEventListener('focusout', handleKeyboardHide);
    
    return () => {
      window.removeEventListener('focusin', handleKeyboardShow);
      window.removeEventListener('focusout', handleKeyboardHide);
    };
  }, []);

  const handleSectionChange = (sectionId: string) => {
    onSectionChange?.(sectionId);
    setIsMenuOpen(false);
    
    // Smooth scroll to section if it exists
    const element = document.getElementById(sectionId);
    if (element) {
      element.scrollIntoView({ behavior: 'smooth', block: 'start' });
    }
  };

  // Only show on mobile devices
  if (typeof window !== 'undefined' && window.innerWidth >= 768) {
    return null;
  }

  return (
    <>
      {/* Mobile Navigation Bar */}
      <nav 
        className={`
          fixed bottom-0 left-0 right-0 z-50 bg-gray-900/95 backdrop-blur-lg border-t border-gray-800
          transition-transform duration-300 ease-in-out
          ${isVisible ? 'translate-y-0' : 'translate-y-full'}
          safe-b safe-l safe-r
        `}
        style={{ height: 'var(--mobile-nav-height)' }}
      >
        <div className="flex items-center justify-around h-full px-2">
          {navigationItems.map((item) => (
            <button
              key={item.id}
              onClick={() => handleSectionChange(item.id)}
              className={`
                flex flex-col items-center justify-center gap-1 py-2 px-3 rounded-lg transition-all duration-200
                min-h-[44px] min-w-[44px] touch-manipulation
                ${activeSection === item.id 
                  ? 'text-primary-500 bg-primary-500/10' 
                  : 'text-gray-400 hover:text-gray-300'
                }
              `}
              aria-label={item.label}
            >
              <div className={`
                transition-transform duration-200
                ${activeSection === item.id ? 'scale-110' : 'scale-100'}
              `}>
                {item.icon}
              </div>
              <span className="text-xs font-medium truncate max-w-[60px]">
                {item.label}
              </span>
            </button>
          ))}
          
          {/* Menu button for additional options */}
          <button
            onClick={() => setIsMenuOpen(!isMenuOpen)}
            className={`
              flex flex-col items-center justify-center gap-1 py-2 px-3 rounded-lg transition-all duration-200
              min-h-[44px] min-w-[44px] touch-manipulation
              ${isMenuOpen 
                ? 'text-primary-500 bg-primary-500/10' 
                : 'text-gray-400 hover:text-gray-300'
              }
            `}
            aria-label="Menu"
          >
            {isMenuOpen ? <X size={20} /> : <Menu size={20} />}
            <span className="text-xs font-medium">Menu</span>
          </button>
        </div>
      </nav>

      {/* Slide-up Menu Overlay */}
      {isMenuOpen && (
        <div 
          className="fixed inset-0 z-40 bg-black/50 backdrop-blur-sm animate-fade-in"
          onClick={() => setIsMenuOpen(false)}
        >
          <div 
            className="absolute bottom-0 left-0 right-0 bg-gray-900 border-t border-gray-800 animate-slide-up"
            style={{ paddingBottom: 'var(--safe-area-inset-bottom)' }}
          >
            <div className="p-4 space-y-2">
              <h3 className="text-lg font-semibold text-white mb-4">Quick Actions</h3>
              
              {/* Additional menu items */}
              <button
                onClick={() => {
                  handleSectionChange('record');
                  // Trigger new section action
                }}
                className="w-full flex items-center gap-3 p-3 rounded-lg bg-gray-800 hover:bg-gray-700 transition-colors text-left"
              >
                <Mic size={20} className="text-primary-500" />
                <div>
                  <p className="font-medium text-white">New Recording</p>
                  <p className="text-sm text-gray-400">Start a fresh recording</p>
                </div>
              </button>
              
              <button
                onClick={() => {
                  handleSectionChange('analyze');
                  // Trigger refresh action
                }}
                className="w-full flex items-center gap-3 p-3 rounded-lg bg-gray-800 hover:bg-gray-700 transition-colors text-left"
              >
                <Activity size={20} className="text-secondary-500" />
                <div>
                  <p className="font-medium text-white">Refresh Analysis</p>
                  <p className="text-sm text-gray-400">Re-analyze current audio</p>
                </div>
              </button>
              
              <button
                onClick={() => {
                  handleSectionChange('generate');
                  // Trigger regenerate action
                }}
                className="w-full flex items-center gap-3 p-3 rounded-lg bg-gray-800 hover:bg-gray-700 transition-colors text-left"
              >
                <Music size={20} className="text-green-500" />
                <div>
                  <p className="font-medium text-white">Regenerate Music</p>
                  <p className="text-sm text-gray-400">Create new arrangement</p>
                </div>
              </button>
              
              <button
                onClick={() => {
                  handleSectionChange('export');
                  // Trigger export all action
                }}
                className="w-full flex items-center gap-3 p-3 rounded-lg bg-gray-800 hover:bg-gray-700 transition-colors text-left"
              >
                <Download size={20} className="text-yellow-500" />
                <div>
                  <p className="font-medium text-white">Export All</p>
                  <p className="text-sm text-gray-400">Download all files</p>
                </div>
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Add padding to bottom of page to account for fixed navigation */}
      <div 
        className="md:hidden" 
        style={{ height: 'var(--mobile-nav-height)' }}
      />
    </>
  );
}