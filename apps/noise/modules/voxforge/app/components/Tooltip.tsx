'use client';

import { useState, useRef, useEffect } from 'react';
import { Info } from 'lucide-react';

interface TooltipProps {
  content: string;
  position?: 'top' | 'bottom' | 'left' | 'right';
  children?: React.ReactNode;
  icon?: boolean;
  className?: string;
}

export default function Tooltip({ 
  content, 
  position = 'top', 
  children, 
  icon = true,
  className = ''
}: TooltipProps) {
  const [isVisible, setIsVisible] = useState(false);
  const [tooltipPosition, setTooltipPosition] = useState({ top: 0, left: 0 });
  const triggerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (isVisible && triggerRef.current) {
      const rect = triggerRef.current.getBoundingClientRect();
      const scrollTop = window.pageYOffset || document.documentElement.scrollTop;
      const scrollLeft = window.pageXOffset || document.documentElement.scrollLeft;
      
      let top = 0;
      let left = 0;
      
      const tooltipWidth = 250;
      const tooltipHeight = 80;
      
      switch (position) {
        case 'top':
          top = rect.top + scrollTop - tooltipHeight - 10;
          left = rect.left + scrollLeft + rect.width / 2 - tooltipWidth / 2;
          break;
        case 'bottom':
          top = rect.bottom + scrollTop + 10;
          left = rect.left + scrollLeft + rect.width / 2 - tooltipWidth / 2;
          break;
        case 'left':
          top = rect.top + scrollTop + rect.height / 2 - tooltipHeight / 2;
          left = rect.left + scrollLeft - tooltipWidth - 10;
          break;
        case 'right':
          top = rect.top + scrollTop + rect.height / 2 - tooltipHeight / 2;
          left = rect.right + scrollLeft + 10;
          break;
      }
      
      // Ensure tooltip stays within viewport
      const viewportWidth = window.innerWidth;
      const viewportHeight = window.innerHeight;
      
      if (left < 10) left = 10;
      if (left + tooltipWidth > viewportWidth) left = viewportWidth - tooltipWidth - 10;
      if (top < 10) top = 10;
      if (top + tooltipHeight > viewportHeight) top = viewportHeight - tooltipHeight - 10;
      
      setTooltipPosition({ top, left });
    }
  }, [isVisible, position]);

  return (
    <>
      <div
        ref={triggerRef}
        className={`inline-flex items-center ${className}`}
        onMouseEnter={() => setIsVisible(true)}
        onMouseLeave={() => setIsVisible(false)}
      >
        {children || (
          icon && (
            <button
              className="p-1 text-gray-400 hover:text-primary-500 transition-colors rounded"
              aria-label="Show help"
            >
              <Info size={16} />
            </button>
          )
        )}
      </div>
      
      {isVisible && (
        <div
          className="fixed z-50 bg-gray-800 border border-gray-700 rounded-lg shadow-xl p-3 max-w-xs"
          style={{
            top: tooltipPosition.top,
            left: tooltipPosition.left,
            width: '250px'
          }}
        >
          <div className="text-sm text-gray-300 leading-relaxed">
            {content}
          </div>
          
          {/* Arrow */}
          <div 
            className={`absolute w-2 h-2 bg-gray-800 border-r border-b border-gray-700 transform rotate-45 ${
              position === 'top' ? 'bottom-[-5px] left-1/2 -translate-x-1/2' :
              position === 'bottom' ? 'top-[-5px] left-1/2 -translate-x-1/2' :
              position === 'left' ? 'right-[-5px] top-1/2 -translate-y-1/2' :
              'left-[-5px] top-1/2 -translate-y-1/2'
            }`}
          />
        </div>
      )}
    </>
  );
}