import { useRef, useCallback, TouchEvent } from "react";

const SWIPE_THRESHOLD_DEFAULT = 50;
const SWIPE_MAX_VERTICAL_MOVEMENT_DEFAULT = 100;

interface SwipeConfig {
  /** Minimum distance (px) to trigger swipe */
  threshold?: number;
  /** Maximum vertical movement allowed */
  maxVerticalMovement?: number;
  /** Callback for left swipe */
  onSwipeLeft?: () => void;
  /** Callback for right swipe */
  onSwipeRight?: () => void;
  /** Callback for up swipe */
  onSwipeUp?: () => void;
  /** Callback for down swipe */
  onSwipeDown?: () => void;
}

interface SwipeHandlers {
  onTouchStart: (e: TouchEvent) => void;
  onTouchMove: (e: TouchEvent) => void;
  onTouchEnd: () => void;
}

interface TouchPosition {
  x: number;
  y: number;
}

/**
 * Hook for handling swipe gestures on touch devices
 *
 * @example
 * const swipeHandlers = useSwipe({
 *   onSwipeLeft: () => nextTab(),
 *   onSwipeRight: () => prevTab(),
 *   threshold: 50,
 * });
 *
 * return <div {...swipeHandlers}>Swipeable content</div>
 */
export function useSwipe(config: SwipeConfig): SwipeHandlers {
  const {
    threshold = SWIPE_THRESHOLD_DEFAULT,
    maxVerticalMovement = SWIPE_MAX_VERTICAL_MOVEMENT_DEFAULT,
    onSwipeLeft,
    onSwipeRight,
    onSwipeUp,
    onSwipeDown,
  } = config;

  const startPos = useRef<TouchPosition | null>(null);
  const currentPos = useRef<TouchPosition | null>(null);

  const onTouchStart = useCallback((e: TouchEvent) => {
    const touch = e.touches[0];
    startPos.current = { x: touch.clientX, y: touch.clientY };
    currentPos.current = { x: touch.clientX, y: touch.clientY };
  }, []);

  const onTouchMove = useCallback((e: TouchEvent) => {
    if (!startPos.current) return;

    const touch = e.touches[0];
    currentPos.current = { x: touch.clientX, y: touch.clientY };
  }, []);

  const onTouchEnd = useCallback(() => {
    if (!startPos.current || !currentPos.current) {
      startPos.current = null;
      currentPos.current = null;
      return;
    }

    const deltaX = currentPos.current.x - startPos.current.x;
    const deltaY = currentPos.current.y - startPos.current.y;
    const absDeltaX = Math.abs(deltaX);
    const absDeltaY = Math.abs(deltaY);

    // Determine if this is a horizontal or vertical swipe
    const isHorizontalSwipe = absDeltaX > absDeltaY;

    if (isHorizontalSwipe && absDeltaX >= threshold && absDeltaY <= maxVerticalMovement) {
      // Horizontal swipe
      if (deltaX < 0 && onSwipeLeft) {
        onSwipeLeft();
      } else if (deltaX > 0 && onSwipeRight) {
        onSwipeRight();
      }
    } else if (!isHorizontalSwipe && absDeltaY >= threshold && absDeltaX <= maxVerticalMovement) {
      // Vertical swipe
      if (deltaY < 0 && onSwipeUp) {
        onSwipeUp();
      } else if (deltaY > 0 && onSwipeDown) {
        onSwipeDown();
      }
    }

    // Reset
    startPos.current = null;
    currentPos.current = null;
  }, [threshold, maxVerticalMovement, onSwipeLeft, onSwipeRight, onSwipeUp, onSwipeDown]);

  return {
    onTouchStart,
    onTouchMove,
    onTouchEnd,
  };
}
