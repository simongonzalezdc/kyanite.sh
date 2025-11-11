'use client';

import { useState, useEffect, useRef } from 'react';
import { cn } from '@/lib/utils';

interface ResponsiveGridProps {
  children: React.ReactNode[];
  className?: string;
  minItemWidth?: number;
  gap?: string;
  maxColumns?: number;
  adaptive?: boolean;
  onLayoutChange?: (columns: number) => void;
}

export default function ResponsiveGrid({
  children,
  className = '',
  minItemWidth = 300,
  gap = 'var(--mobile-gap)',
  maxColumns = 4,
  adaptive = true,
  onLayoutChange
}: ResponsiveGridProps) {
  const [columns, setColumns] = useState(1);
  const [containerWidth, setContainerWidth] = useState(0);
  const containerRef = useRef<HTMLDivElement>(null);

  // Calculate optimal columns based on container width
  const calculateColumns = (width: number) => {
    if (!adaptive) return 1;
    
    // Account for gap spacing
    const gapPixels = parseFloat(gap.replace(/[^\d.]/g, '')) || 16;
    const availableWidth = width - gapPixels;
    
    // Calculate how many columns fit
    const calculatedColumns = Math.floor(availableWidth / (minItemWidth + gapPixels));
    
    // Ensure we stay within bounds
    return Math.max(1, Math.min(calculatedColumns, maxColumns));
  };

  // Handle resize and layout calculations
  useEffect(() => {
    const updateLayout = () => {
      if (containerRef.current) {
        const width = containerRef.current.offsetWidth;
        setContainerWidth(width);
        
        const newColumns = calculateColumns(width);
        if (newColumns !== columns) {
          setColumns(newColumns);
          onLayoutChange?.(newColumns);
        }
      }
    };

    // Initial calculation
    updateLayout();

    // Set up resize observer for better performance
    const resizeObserver = new ResizeObserver(() => {
      updateLayout();
    });

    if (containerRef.current) {
      resizeObserver.observe(containerRef.current);
    }

    // Fallback to window resize for older browsers
    const handleResize = () => {
      updateLayout();
    };
    
    window.addEventListener('resize', handleResize);

    return () => {
      resizeObserver.disconnect();
      window.removeEventListener('resize', handleResize);
    };
  }, [minItemWidth, gap, maxColumns, adaptive, onLayoutChange]); // Remove columns from dependencies

  // Responsive breakpoints for static layouts
  const getResponsiveClasses = () => {
    if (!adaptive) {
      return cn(
        'grid grid-cols-1',
        'xs:grid-cols-1',
        'sm:grid-cols-1',
        'md:grid-cols-2',
        'lg:grid-cols-3',
        'xl:grid-cols-4',
        className
      );
    }

    return cn(
      'grid',
      `grid-cols-${Math.min(columns, maxColumns)}`,
      className
    );
  };

  // Get gap style
  const getGapStyle = () => {
    return { gap };
  };

  // Performance optimization: only render visible items on mobile
  const getVisibleChildren = () => {
    // On small screens, limit initial render for performance
    if (containerWidth < 480 && children.length > 6) {
      return children.slice(0, 6);
    }
    return children;
  };

  return (
    <div className="w-full">
      {/* Layout indicator for development */}
      {process.env.NODE_ENV === 'development' && (
        <div className="text-xs text-gray-500 mb-2 text-center">
          Grid: {columns} column{columns !== 1 ? 's' : ''} ({containerWidth}px)
        </div>
      )}
      
      {/* Main grid container */}
      <div
        ref={containerRef}
        className={getResponsiveClasses()}
        style={getGapStyle()}
      >
        {getVisibleChildren().map((child, index) => (
          <div
            key={index}
            className={cn(
              'w-full',
              // Performance optimizations
              'will-change-transform',
              // Touch-friendly sizing
              'min-h-[44px]',
              // Responsive item sizing
              adaptive && 'transition-all duration-300 ease-in-out'
            )}
            style={{
              // Adaptive item width based on columns
              ...(adaptive && {
                minWidth: `${minItemWidth}px`,
                maxWidth: columns === 1 ? '100%' : `${100 / columns}%`
              })
            }}
          >
            {child}
          </div>
        ))}
      </div>

      {/* Load more button for mobile when items are truncated */}
      {containerWidth < 480 && children.length > 6 && (
        <div className="mt-4 text-center">
          <button
            onClick={() => {
              // In a real implementation, this would load more items
              console.log('Load more items');
            }}
            className="px-4 py-2 bg-primary-500 hover:bg-primary-600 rounded-lg text-white font-medium transition-colors min-h-[44px]"
          >
            Load More ({children.length - 6} remaining)
          </button>
        </div>
      )}
    </div>
  );
}

// Helper component for grid items with consistent styling
export function ResponsiveGridItem({
  children,
  className = '',
  priority = false,
  ...props
}: {
  children: React.ReactNode;
  className?: string;
  priority?: boolean;
  [key: string]: any;
}) {
  return (
    <div
      className={cn(
        'bg-gray-900 rounded-xl border border-gray-800 p-4',
        'transition-all duration-200 ease-in-out',
        'hover:border-gray-700 hover:shadow-lg',
        // Touch-friendly interactions
        'active:scale-[0.98]',
        // Performance optimizations
        priority ? 'will-change-transform' : '',
        className
      )}
      {...props}
    >
      {children}
    </div>
  );
}

// Specialized grid for different content types
export function ResponsiveCardGrid({ children, ...props }: ResponsiveGridProps) {
  return (
    <ResponsiveGrid
      minItemWidth={280}
      maxColumns={3}
      gap="1rem"
      {...props}
    >
      {children.map((child, index) => (
        <ResponsiveGridItem key={index} priority={index < 4}>
          {child}
        </ResponsiveGridItem>
      ))}
    </ResponsiveGrid>
  );
}

export function ResponsiveControlGrid({ children, ...props }: ResponsiveGridProps) {
  return (
    <ResponsiveGrid
      minItemWidth={200}
      maxColumns={4}
      gap="0.75rem"
      adaptive={false}
      {...props}
    >
      {children.map((child, index) => (
        <div
          key={index}
          className="flex items-center justify-center min-h-[44px]"
        >
          {child}
        </div>
      ))}
    </ResponsiveGrid>
  );
}

export function ResponsiveStatsGrid({ children, ...props }: ResponsiveGridProps) {
  return (
    <ResponsiveGrid
      minItemWidth={150}
      maxColumns={4}
      gap="1rem"
      {...props}
    >
      {children.map((child, index) => (
        <ResponsiveGridItem key={index} className="text-center">
          {child}
        </ResponsiveGridItem>
      ))}
    </ResponsiveGrid>
  );
}