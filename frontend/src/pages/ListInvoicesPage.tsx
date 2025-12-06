import { useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { useAppStore } from '@/stores/app-store';
import { apiClient, getImageUrl } from '@/lib/api';
import type { InvoiceStatus } from '@/types';

export default function ListInvoicesPage() {
  const navigate = useNavigate();
  const { setSelectedInvoiceId } = useAppStore();

  const { data, isLoading, error } = useQuery({
    queryKey: ['invoices'],
    queryFn: () => apiClient.getInvoices(),
    refetchInterval: 3000,
  });

  const invoices = data?.data || [];

  const getStatusIcon = (status: InvoiceStatus) => {
    switch (status) {
      case 'completed':
        return (
          <span className="text-base fill-current" style={{ fontVariationSettings: "'FILL' 1" }}>
            ✓
          </span>
        );
      case 'failed':
        return (
          <span className="text-base fill-current" style={{ fontVariationSettings: "'FILL' 1" }}>
            ✗
          </span>
        );
      case 'processing':
        return (
          <span className="text-base animate-spin">
            ↻
          </span>
        );
      case 'pending':
        return (
          <span className="text-base animate-pulse">
            ⏳
          </span>
        );
      default:
        return null;
    }
  };

  const getStatusText = (status: InvoiceStatus) => {
    switch (status) {
      case 'completed':
        return 'Extraction Done';
      case 'failed':
        return 'Extraction Failed';
      case 'processing':
        return 'Extracting...';
      case 'pending':
        return 'Waiting...';
      default:
        return status;
    }
  };

  const getStatusColor = (status: InvoiceStatus) => {
    switch (status) {
      case 'completed':
        return 'text-status-green';
      case 'failed':
        return 'text-red-500';
      case 'processing':
        return 'text-status-blue';
      case 'pending':
        return 'text-gray-500';
      default:
        return 'text-gray-500';
    }
  };

  const handleInvoiceClick = (invoiceId: string) => {
    setSelectedInvoiceId(invoiceId);
    navigate(`/extract-invoice-data/${invoiceId}`);
  };

  return (
    <div className="relative mx-auto flex h-screen max-h-[960px] w-full max-w-[480px] flex-col overflow-hidden bg-background-light dark:bg-background-dark font-display text-text-light dark:text-text-dark">
      {/* Header */}
      <header className="flex shrink-0 items-center justify-between border-b border-border-light dark:border-border-dark bg-surface-light dark:bg-surface-dark px-4 py-3">
        <button className="flex size-10 cursor-pointer items-center justify-center rounded-full">
          <span className="material-symbols-outlined text-2xl">←</span>
        </button>
        <h1 className="text-lg font-semibold">Scanned Invoices</h1>
        <button className="flex size-10 cursor-pointer items-center justify-center rounded-full">
          <span className="material-symbols-outlined text-2xl">🔍</span>
        </button>
      </header>

      {/* Main Content */}
      <main className="flex-grow overflow-y-auto p-4">
        {isLoading && (
          <div className="flex items-center justify-center h-64">
            <div className="text-center space-y-4">
              <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto"></div>
              <p className="text-gray-600 dark:text-gray-300">Loading invoices...</p>
            </div>
          </div>
        )}

        {error && (
          <div className="p-4">
            <div className="bg-red-50 border border-red-200 rounded-lg p-4">
              <p className="text-red-700">
                Failed to load invoices: {error instanceof Error ? error.message : 'Unknown error'}
              </p>
            </div>
          </div>
        )}

        {!isLoading && !error && (
          <div className="flex flex-col gap-4">
            {invoices.length === 0 ? (
              <div className="text-center py-12">
                <p className="text-gray-500 dark:text-gray-400">No invoices yet. Take a picture to get started!</p>
              </div>
            ) : (
              invoices.map((invoice) => (
                <button
                  key={invoice.id}
                  className="flex items-center gap-4 rounded-xl bg-surface-light p-3 shadow-sm transition-shadow hover:shadow-md dark:bg-surface-dark text-left w-full"
                  onClick={() => handleInvoiceClick(invoice.id)}
                >
                  <div className="h-16 w-16 flex-shrink-0 rounded-lg bg-gray-200 dark:bg-gray-700 overflow-hidden flex items-center justify-center">
                    {invoice.image_path ? (
                      <img 
                        alt="Invoice thumbnail" 
                        className="h-full w-full object-cover" 
                        src={getImageUrl(invoice.image_path) || ''}
                        onError={(e) => {
                          const target = e.target as HTMLImageElement;
                          target.style.display = 'none';
                        }}
                      />
                    ) : null}
                    {!invoice.image_path && (
                      <span className="text-gray-400 text-2xl">📄</span>
                    )}
                  </div>
                  <div className="flex-grow">
                    <h2 className="font-semibold text-text-light dark:text-text-dark">
                      Invoice {invoice.id.slice(-8)}
                    </h2>
                    <div className={`mt-1 flex items-center gap-1.5 text-sm ${getStatusColor(invoice.status)}`}>
                      {getStatusIcon(invoice.status)}
                      <p>{getStatusText(invoice.status)}</p>
                    </div>
                    {invoice.error_message && (
                      <p className="mt-1 text-xs text-red-500 truncate">{invoice.error_message}</p>
                    )}
                  </div>
                  <span className="material-symbols-outlined text-text-light-secondary dark:text-text-dark-secondary">›</span>
                </button>
              ))
            )}
          </div>
        )}
      </main>

      {/* Floating Action Button */}
      <div className="absolute bottom-6 right-6">
        <button 
          className="flex h-16 w-16 items-center justify-center rounded-full bg-primary text-white shadow-lg transition-transform hover:scale-105 active:scale-95"
          onClick={() => navigate('/take-picture')}
        >
          <span className="text-4xl" style={{ fontVariationSettings: "'FILL' 1" }}>📷</span>
        </button>
      </div>
    </div>
  );
}