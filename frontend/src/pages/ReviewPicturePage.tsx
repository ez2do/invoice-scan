import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useMutation } from '@tanstack/react-query';
import { ArrowLeft, RotateCcw, FileSearch, Loader2, CheckCircle2 } from 'lucide-react';
import { useAppStore } from '@/stores/app-store';
import { apiClient } from '@/lib/api';

export default function ReviewPicturePage() {
  const navigate = useNavigate();
  const { currentImage, clearData } = useAppStore();
  const [isUploading, setIsUploading] = useState(false);

  const uploadMutation = useMutation({
    mutationFn: apiClient.uploadInvoice.bind(apiClient),
    onSuccess: () => {
      setIsUploading(false);
      clearData();
      navigate('/list-invoices');
    },
    onError: (error: Error) => {
      setIsUploading(false);
      alert(`Không thể tải lên hóa đơn: ${error.message}`);
    },
  });

  const handleRetake = () => {
    navigate('/take-picture');
  };

  const handleExtractData = () => {
    if (!currentImage) return;
    setIsUploading(true);
    uploadMutation.mutate(currentImage);
  };

  const handleBack = () => {
    clearData();
    navigate('/list-invoices');
  };

  if (!currentImage) {
    navigate('/take-picture');
    return null;
  }

  return (
    <div className="page-container bg-surface-950">
      <header className="page-header bg-surface-900/80 border-surface-800/60 safe-top">
        <button
          className="icon-btn text-surface-400 hover:bg-surface-800"
          onClick={handleBack}
          aria-label="Quay lại"
        >
          <ArrowLeft className="w-5 h-5" />
        </button>
        <h1 className="page-title text-white">Xem lại ảnh</h1>
        <div className="w-10 h-10" />
      </header>

      <main className="flex-grow flex flex-col items-center justify-center px-4 py-6">
        <div className="w-full max-w-sm space-y-4 animate-fade-in">
          <div className="flex items-center justify-center gap-2 text-surface-400">
            <CheckCircle2 className="w-4 h-4 text-accent-500" />
            <p className="text-sm font-medium">Chụp ảnh thành công</p>
          </div>

          <div className="relative rounded-2xl overflow-hidden shadow-soft-xl bg-surface-800">
            <img
              alt="Ảnh hóa đơn đã chụp để xem lại"
              className="w-full h-auto object-contain max-h-[420px]"
              src={currentImage}
            />
            <div className="absolute inset-0 ring-1 ring-inset ring-white/10 rounded-2xl pointer-events-none" />
          </div>

          <p className="text-surface-500 text-center text-sm px-4">
            Đảm bảo hóa đơn hiển thị rõ ràng và tất cả văn bản đều đọc được
          </p>
        </div>
      </main>

      <footer className="shrink-0 p-4 pb-6 safe-bottom bg-surface-950 border-t border-surface-800/60">
        <div className="flex items-center gap-3 max-w-sm mx-auto">
          <button
            className="btn-secondary flex-1 h-14 bg-surface-800 hover:bg-surface-700 border-surface-700 text-white"
            onClick={handleRetake}
          >
            <RotateCcw className="w-5 h-5" />
            <span>Chụp lại</span>
          </button>
          <button
            className="btn-primary flex-1 h-14"
            onClick={handleExtractData}
            disabled={isUploading}
          >
            {isUploading ? (
              <>
                <Loader2 className="w-5 h-5 animate-spin" />
                <span>Đang tải lên...</span>
              </>
            ) : (
              <>
                <FileSearch className="w-5 h-5" />
                <span>Trích xuất dữ liệu</span>
              </>
            )}
          </button>
        </div>
      </footer>
    </div>
  );
}
