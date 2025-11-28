import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

/**
 * Hook to handle PWA-friendly navigation on iOS
 * Prevents address bar from showing during navigation
 */
export function usePWANavigation() {
  const navigate = useNavigate();

  useEffect(() => {
    // Detect if we're in standalone PWA mode
    const isStandalone = window.matchMedia('(display-mode: standalone)').matches;
    const isIOSStandalone = (window.navigator as any).standalone === true;
    const isPWA = isStandalone || isIOSStandalone;

    if (isPWA) {
      // Hide address bar on iOS by scrolling
      const hideAddressBar = () => {
        setTimeout(() => {
          window.scrollTo(0, 1);
        }, 100);
      };

      // Hide address bar initially
      hideAddressBar();

      // Hide address bar on orientation change
      window.addEventListener('orientationchange', hideAddressBar);
      window.addEventListener('resize', hideAddressBar);

      return () => {
        window.removeEventListener('orientationchange', hideAddressBar);
        window.removeEventListener('resize', hideAddressBar);
      };
    }
  }, []);

  // PWA-friendly navigate function
  const navigatePWA = (to: string, options?: { replace?: boolean }) => {
    // Use requestAnimationFrame to ensure smooth navigation
    requestAnimationFrame(() => {
      navigate(to, options);
      
      // Hide address bar after navigation
      setTimeout(() => {
        window.scrollTo(0, 1);
      }, 50);
    });
  };

  return { navigatePWA };
}