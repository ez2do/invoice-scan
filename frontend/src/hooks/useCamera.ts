import { useState, useEffect, useRef } from 'react';
import type { CameraConfig } from '@/types';

interface UseCameraResult {
  stream: MediaStream | null;
  videoRef: React.RefObject<HTMLVideoElement | null>;
  captureImage: () => string | null;
  error: string | null;
  isSupported: boolean;
  isActive: boolean;
  startCamera: () => Promise<void>;
  stopCamera: () => Promise<void>;
  torch: boolean;
  supportsTorch: boolean;
  toggleTorch: () => Promise<void>;
}

const defaultConfig: CameraConfig = {
  video: {
    facingMode: 'environment',
    width: { ideal: 1024 },
    height: { ideal: 768 },
  },
};

export function useCamera(config: CameraConfig = defaultConfig): UseCameraResult {
  const [stream, setStream] = useState<MediaStream | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isActive, setIsActive] = useState(false);
  const [torch, setTorch] = useState(false);
  const [supportsTorch, setSupportsTorch] = useState(false);
  const videoRef = useRef<HTMLVideoElement>(null);
  const streamRef = useRef<MediaStream | null>(null);

  const isSupported = 'mediaDevices' in navigator && 'getUserMedia' in navigator.mediaDevices;

  const startCamera = async () => {
    if (!isSupported) {
      setError('Camera is not supported in this browser');
      return;
    }

    // Check if we're on HTTPS or localhost
    const isSecureContext = window.isSecureContext || location.protocol === 'https:' || location.hostname === 'localhost';
    if (!isSecureContext) {
      setError('Camera access requires HTTPS. Please use https:// instead of http://');
      return;
    }

    try {
      setError(null);

      // Try with mobile-optimized constraints first
      const mobileConfig = {
        video: {
          facingMode: { ideal: 'environment' },
          width: { ideal: 1280, max: 1920 },
          height: { ideal: 720, max: 1080 },
          aspectRatio: { ideal: 16 / 9 }
        }
      };

      let mediaStream;
      try {
        mediaStream = await navigator.mediaDevices.getUserMedia(mobileConfig);
      } catch (mobileErr) {
        // Fallback to basic constraints
        console.warn('Mobile-optimized camera failed, trying basic constraints:', mobileErr);
        mediaStream = await navigator.mediaDevices.getUserMedia(config);
      }

      setStream(mediaStream);
      streamRef.current = mediaStream;
      setIsActive(true);

      // Check for torch support
      const track = mediaStream.getVideoTracks()[0];
      const capabilities = track.getCapabilities() as any;
      if (capabilities.torch) {
        setSupportsTorch(true);
      }

      if (videoRef.current) {
        videoRef.current.srcObject = mediaStream;
      }
    } catch (err: any) {
      let errorMessage = 'Failed to access camera';

      if (err.name === 'NotAllowedError') {
        errorMessage = 'Camera permission denied. Please allow camera access and try again.';
      } else if (err.name === 'NotFoundError') {
        errorMessage = 'No camera found on this device.';
      } else if (err.name === 'NotSupportedError') {
        errorMessage = 'Camera not supported on this device.';
      } else if (err.name === 'NotReadableError') {
        errorMessage = 'Camera is being used by another application.';
      } else if (err.message) {
        errorMessage = err.message;
      }

      setError(errorMessage);
      console.error('Camera access error:', err);
    }
  };

  const stopCamera = async () => {
    const activeStream = streamRef.current;
    if (activeStream) {
      // First, explicitly turn off the torch if it's on
      if (torch) {
        const track = activeStream.getVideoTracks()[0];
        try {
          await track.applyConstraints({
            advanced: [{ torch: false } as any]
          });
        } catch (err) {
          console.error('Failed to turn off torch:', err);
        }
      }

      // Then stop all tracks
      activeStream.getTracks().forEach(track => {
        track.stop();
      });
      setStream(null);
      streamRef.current = null;
      setIsActive(false);
      setTorch(false);
      setSupportsTorch(false);
    }
  };

  const toggleTorch = async () => {
    const activeStream = streamRef.current;
    if (!activeStream) return;
    const track = activeStream.getVideoTracks()[0];
    const newTorchState = !torch;

    try {
      await track.applyConstraints({
        advanced: [{ torch: newTorchState } as any]
      });
      setTorch(newTorchState);
    } catch (err) {
      console.error('Failed to toggle torch:', err);
      // Don't set error state here as it's a non-critical error
    }
  };

  const captureImage = (): string | null => {
    if (!videoRef.current || !streamRef.current) {
      setError('Camera is not active');
      return null;
    }

    try {
      const video = videoRef.current;
      const canvas = document.createElement('canvas');
      const ctx = canvas.getContext('2d');

      if (!ctx) {
        setError('Canvas context not available');
        return null;
      }

      // Set canvas dimensions to video dimensions
      canvas.width = video.videoWidth;
      canvas.height = video.videoHeight;

      // Draw the current frame
      ctx.drawImage(video, 0, 0, canvas.width, canvas.height);

      // Return as data URL
      return canvas.toDataURL('image/jpeg', 0.9);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to capture image';
      setError(errorMessage);
      return null;
    }
  };

  useEffect(() => {
    // Cleanup on unmount
    return () => {
      stopCamera();
    };
  }, []);

  return {
    stream,
    videoRef,
    captureImage,
    error,
    isSupported,
    isActive,
    startCamera,
    stopCamera,
    torch,
    supportsTorch,
    toggleTorch,
  };
}