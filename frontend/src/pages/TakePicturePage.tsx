import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { X, Zap, Camera, AlertCircle } from 'lucide-react';
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
      <div className="fixed inset-0 bg-surface-950 flex items-center justify-center z-50">
        <div className="text-center space-y-6 p-6 max-w-sm animate-fade-in">
          <div className="w-16 h-16 rounded-2xl bg-error-500/20 flex items-center justify-center mx-auto">
            <AlertCircle className="w-8 h-8 text-error-400" />
          </div>
          <div>
            <h2 className="text-xl font-semibold text-white mb-2">Không thể truy cập camera</h2>
            <p className="text-surface-400 text-sm">
              Thiết bị hoặc trình duyệt này không hỗ trợ truy cập camera
            </p>
          </div>
          <button
            onClick={handleClose}
            className="btn-primary w-full"
          >
            Quay lại
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="fixed inset-0 bg-black z-50">
      <video
        ref={videoRef}
        autoPlay
        playsInline
        muted
        className="absolute inset-0 w-full h-full object-cover"
        style={{ zIndex: 1 }}
      />

      <div
        className="absolute inset-0 bg-gradient-to-b from-black/40 via-transparent to-black/60"
        style={{ zIndex: 2 }}
      />

      <div
        className="absolute inset-0 flex flex-col safe-top safe-bottom"
        style={{ zIndex: 10 }}
      >
        <div className="flex items-center justify-between p-4">
          <button
            className="w-11 h-11 rounded-xl bg-white/10 backdrop-blur-sm flex items-center justify-center text-white hover:bg-white/20 active:scale-95 transition-all"
            onClick={handleClose}
            aria-label="Đóng camera"
          >
            <X className="w-5 h-5" />
          </button>
          <button 
            className="w-11 h-11 rounded-xl bg-white/10 backdrop-blur-sm flex items-center justify-center text-white hover:bg-white/20 active:scale-95 transition-all"
            aria-label="Đèn flash"
          >
            <Zap className="w-5 h-5" />
          </button>
        </div>

        <div className="flex-1 flex items-center justify-center px-6">
          <div className="w-full max-w-sm aspect-[0.707] relative">
            <div className="absolute inset-0 border-2 border-white/40 rounded-2xl" />
            <div className="absolute top-0 left-0 w-8 h-8 border-t-[3px] border-l-[3px] border-white rounded-tl-xl" />
            <div className="absolute top-0 right-0 w-8 h-8 border-t-[3px] border-r-[3px] border-white rounded-tr-xl" />
            <div className="absolute bottom-0 left-0 w-8 h-8 border-b-[3px] border-l-[3px] border-white rounded-bl-xl" />
            <div className="absolute bottom-0 right-0 w-8 h-8 border-b-[3px] border-r-[3px] border-white rounded-br-xl" />
          </div>
        </div>

        <div className="px-6 pb-4">
          <p className="text-white/80 text-center text-sm font-medium">
            Đặt hóa đơn trong khung hình
          </p>
        </div>

        <div className="flex items-center justify-center pb-8">
          <button
            className="w-[72px] h-[72px] rounded-full bg-white flex items-center justify-center shadow-soft-xl hover:scale-105 active:scale-95 transition-transform disabled:opacity-50 disabled:hover:scale-100"
            onClick={handleCapture}
            disabled={!isActive}
            aria-label="Chụp ảnh"
          >
            <div className="w-[60px] h-[60px] rounded-full bg-primary-600 flex items-center justify-center">
              <Camera className="w-7 h-7 text-white" />
            </div>
          </button>
        </div>
      </div>

      {error && (
        <div
          className="absolute top-20 left-4 right-4 animate-fade-in"
          style={{ zIndex: 20 }}
        >
          <div className="bg-error-500/90 backdrop-blur-sm text-white p-4 rounded-xl flex items-start gap-3">
            <AlertCircle className="w-5 h-5 shrink-0 mt-0.5" />
            <p className="text-sm">{error}</p>
          </div>
        </div>
      )}
    </div>
  );
}
