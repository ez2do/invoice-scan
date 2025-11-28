import { useNavigate } from 'react-router-dom';
import { useAppStore } from '@/stores/app-store';

export default function ReviewPicturePage() {
  const navigate = useNavigate();
  const { currentImage, clearData } = useAppStore();

  const handleRetake = () => {
    navigate('/take-picture');
  };

  const handleExtractData = () => {
    navigate('/extract-invoice-data');
  };

  const handleBack = () => {
    clearData();
    navigate('/list-invoices');
  };

  if (!currentImage) {
    // Redirect if no image
    navigate('/take-picture');
    return null;
  }

  return (
    <div className="relative mx-auto flex h-screen max-h-[960px] w-full max-w-[480px] flex-col overflow-hidden bg-background-dark">
      {/* Header */}
      <header className="relative z-10 flex items-center justify-between p-4">
        <button 
          className="flex h-10 w-10 items-center justify-center rounded-full text-white"
          onClick={handleBack}
        >
          <span className="text-2xl">←</span>
        </button>
        <h1 className="text-lg font-semibold text-white">Review Invoice</h1>
        <div className="h-10 w-10"></div>
      </header>

      {/* Main Content */}
      <main className="flex flex-grow flex-col items-center justify-center px-6 pb-4 pt-2">
        <p className="mb-4 text-center text-base font-normal text-gray-300">
          Is the invoice clear and all details readable?
        </p>
        <div className="w-full overflow-hidden rounded-xl">
          <img 
            alt="A captured image of an invoice for review" 
            className="h-full w-full object-cover max-h-96" 
            src={currentImage}
          />
        </div>
      </main>

      {/* Footer */}
      <footer className="z-10 bg-background-dark p-6">
        <div className="flex items-center space-x-4">
          <button 
            className="flex h-14 flex-1 items-center justify-center gap-2 rounded-xl border border-white/50 bg-white/10 text-base font-semibold text-white transition-colors hover:bg-white/20"
            onClick={handleRetake}
          >
            <span className="text-xl">↻</span>
            Retake
          </button>
          <button 
            className="flex h-14 flex-1 items-center justify-center gap-2 rounded-xl bg-primary text-base font-semibold text-white transition-colors hover:bg-primary/90"
            onClick={handleExtractData}
          >
            <span className="text-xl">📄</span>
            Extract Data
          </button>
        </div>
      </footer>
    </div>
  );
}