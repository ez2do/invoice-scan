import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useCamera } from '@/hooks/useCamera';
import { useAppStore } from '@/stores/app-store';

export default function TakePicturePage() {
  const navigate = useNavigate();
  const { 
    videoRef, 
    captureImage, 
    error, 
    isSupported, 
    isActive, 
    startCamera, 
    stopCamera 
  } = useCamera();
  
  const { setCurrentImage } = useAppStore();

  useEffect(() => {
    if (isSupported) {
      startCamera();
    }

    return () => {
      stopCamera();
    };
  }, [isSupported]);

  const handleCapture = () => {
    const imageData = captureImage();
    if (imageData) {
      setCurrentImage(imageData);
      navigate('/review-picture');
    }
  };

  const handleClose = () => {
    stopCamera();
    navigate('/list-invoices');
  };

  if (!isSupported) {
    return (
      <div className="fixed inset-0 bg-black flex items-center justify-center z-50">
        <div className="text-center space-y-4">
          <p className="text-white text-lg">Camera not supported</p>
          <button 
            onClick={handleClose}
            className="px-4 py-2 bg-primary text-white rounded-lg"
          >
            Go Back
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="fixed inset-0 bg-black z-50">
      {/* Camera Video Background */}
      <video
        ref={videoRef}
        autoPlay
        playsInline
        muted
        className="absolute inset-0 w-full h-full object-cover"
        style={{ zIndex: 1 }}
      />
      
      {/* Dark overlay */}
      <div 
        className="absolute inset-0 bg-black/20"
        style={{ zIndex: 2 }}
      />
      
      {/* UI Controls */}
      <div 
        className="absolute inset-0 flex flex-col"
        style={{ zIndex: 10 }}
      >
        {/* Top Controls */}
        <div className="flex items-center justify-between p-4">
          <button 
            className="w-12 h-12 bg-black/30 rounded-full flex items-center justify-center text-white hover:bg-black/50 transition-colors"
            onClick={handleClose}
          >
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
          <button className="w-12 h-12 bg-black/30 rounded-full flex items-center justify-center text-white hover:bg-black/50 transition-colors">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
              <polygon points="13,2 3,14 12,14 11,22 21,10 12,10"></polygon>
            </svg>
          </button>
        </div>
        
        {/* Center Area - Frame Guide */}
        <div className="flex-1 flex items-center justify-center px-8">
          <div className="w-full max-w-sm aspect-[0.707] border-2 border-dashed border-white/60 rounded-lg bg-transparent"></div>
        </div>
        
        {/* Instruction Text */}
        <div className="px-6 pb-4">
          <p className="text-white text-center text-base">
            Position invoice within the frame. Ensure good lighting for best results.
          </p>
        </div>
        
        {/* Bottom Controls - Capture Button */}
        <div className="flex items-center justify-center pb-8">
          <button 
            className="w-20 h-20 bg-primary rounded-full border-4 border-white flex items-center justify-center shadow-lg hover:scale-105 active:scale-95 transition-transform disabled:opacity-50"
            onClick={handleCapture}
            disabled={!isActive}
          >
            <svg 
              width="32" 
              height="32" 
              viewBox="0 0 24 24" 
              fill="none" 
              stroke="white" 
              strokeWidth="2"
              strokeLinecap="round" 
              strokeLinejoin="round"
            >
              <path d="M14.5 4h-5L7 7H4a2 2 0 0 0-2 2v9a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2V9a2 2 0 0 0-2-2h-3l-2.5-3z"/>
              <circle cx="12" cy="13" r="3"/>
            </svg>
          </button>
        </div>
      </div>

      {/* Error Display */}
      {error && (
        <div 
          className="absolute top-20 left-4 right-4"
          style={{ zIndex: 20 }}
        >
          <div className="bg-red-500 text-white p-3 rounded-lg text-center">
            <p className="text-sm">{error}</p>
          </div>
        </div>
      )}
    </div>
  );
}