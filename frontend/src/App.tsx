import { HashRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import ListInvoicesPage from '@/pages/ListInvoicesPage';
import TakePicturePage from '@/pages/TakePicturePage';
import ReviewPicturePage from '@/pages/ReviewPicturePage';
import ExtractInvoiceDataPage from '@/pages/ExtractInvoiceDataPage';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 5 * 60 * 1000, // 5 minutes
      retry: 2,
    },
  },
});

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <Router basename="/">
        <div className="min-h-screen bg-background-light dark:bg-background-dark">
          <Routes>
            <Route path="/" element={<Navigate to="/list-invoices" replace />} />
            <Route path="/list-invoices" element={<ListInvoicesPage />} />
            <Route path="/take-picture" element={<TakePicturePage />} />
            <Route path="/review-picture" element={<ReviewPicturePage />} />
            <Route path="/extract-invoice-data/:id?" element={<ExtractInvoiceDataPage />} />
            <Route path="*" element={<Navigate to="/list-invoices" replace />} />
          </Routes>
        </div>
      </Router>
    </QueryClientProvider>
  );
}

export default App;