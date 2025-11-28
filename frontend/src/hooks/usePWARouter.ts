import { useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';

export function usePWARouter() {
  const navigate = useNavigate();
  const location = useLocation();

  useEffect(() => {
    // Handle PWA navigation state
    const handlePopState = (event: PopStateEvent) => {
      // Ensure proper navigation handling in PWA
      if (event.state?.path) {
        navigate(event.state.path, { replace: true });
      }
    };

    // Handle PWA app launch
    const handleAppInstalled = () => {
      // Redirect to main screen when app is launched
      if (location.pathname === '/') {
        navigate('/list-invoices', { replace: true });
      }
    };

    window.addEventListener('popstate', handlePopState);
    window.addEventListener('appinstalled', handleAppInstalled);

    // Handle initial route for PWA
    if (window.matchMedia('(display-mode: standalone)').matches) {
      // Running as PWA
      if (location.pathname === '/') {
        navigate('/list-invoices', { replace: true });
      }
    }

    return () => {
      window.removeEventListener('popstate', handlePopState);
      window.removeEventListener('appinstalled', handleAppInstalled);
    };
  }, [navigate, location]);

  const navigateWithState = (path: string, options?: any) => {
    // Enhanced navigation with proper state management
    window.history.pushState({ path }, '', path);
    navigate(path, options);
  };

  return { navigateWithState };
}