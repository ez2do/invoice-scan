import { useState, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  FileText,
  Camera,
  ChevronRight,
  ChevronLeft,
  CheckCircle2,
  XCircle,
  Loader2,
  Clock,
  Trash2
} from 'lucide-react';
import { apiClient, getImageUrl } from '@/lib/api';
import type { InvoiceStatus } from '@/types';

const PAGE_SIZE = 10;
const SWIPE_THRESHOLD = 80;

export default function ListInvoicesPage() {
  const navigate = useNavigate();
  const [page, setPage] = useState(1);
  const [deletingId, setDeletingId] = useState<string | null>(null);
  const [swipedId, setSwipedId] = useState<string | null>(null);
  const touchStartX = useRef<number>(0);
  const touchCurrentX = useRef<number>(0);
  const queryClient = useQueryClient();

  const { data, isLoading, error } = useQuery({
    queryKey: ['invoices', page],
    queryFn: () => apiClient.getInvoices(page, PAGE_SIZE),
    refetchInterval: (query) => {
      const invoices = query.state.data?.data || [];
      const hasActiveInvoices = invoices.some(
        (inv) => inv.status === 'pending' || inv.status === 'processing'
      );
      return hasActiveInvoices ? 10000 : false;
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => apiClient.deleteInvoice(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['invoices'] });
      setDeletingId(null);
      setSwipedId(null);
    },
    onError: () => {
      setDeletingId(null);
    },
  });

  const invoices = data?.data || [];
  const totalPages = data?.total_pages || 1;
  const total = data?.total || 0;

  const handleDelete = (invoiceId: string) => {
    setDeletingId(invoiceId);
    deleteMutation.mutate(invoiceId);
  };

  const handleTouchStart = (e: React.TouchEvent, invoiceId: string) => {
    touchStartX.current = e.touches[0].clientX;
    touchCurrentX.current = e.touches[0].clientX;
    // Close other swiped items
    if (swipedId && swipedId !== invoiceId) {
      setSwipedId(null);
    }
  };

  const handleTouchMove = (e: React.TouchEvent, invoiceId: string, element: HTMLDivElement | null) => {
    if (!element) return;
    touchCurrentX.current = e.touches[0].clientX;
    const diff = touchStartX.current - touchCurrentX.current;

    if (diff > 0) {
      // Swiping left - show delete
      const translateX = Math.min(diff, SWIPE_THRESHOLD);
      element.style.transform = `translateX(-${translateX}px)`;
    } else if (swipedId === invoiceId) {
      // Swiping right when already swiped - close
      const translateX = Math.max(SWIPE_THRESHOLD + diff, 0);
      element.style.transform = `translateX(-${translateX}px)`;
    }
  };

  const handleTouchEnd = (invoiceId: string, element: HTMLDivElement | null) => {
    if (!element) return;
    const diff = touchStartX.current - touchCurrentX.current;

    if (diff > SWIPE_THRESHOLD / 2) {
      // Swiped enough - keep open
      element.style.transform = `translateX(-${SWIPE_THRESHOLD}px)`;
      setSwipedId(invoiceId);
    } else if (swipedId === invoiceId && diff < -20) {
      // Swiped back - close
      element.style.transform = 'translateX(0)';
      setSwipedId(null);
    } else {
      // Not enough swipe - reset
      element.style.transform = swipedId === invoiceId ? `translateX(-${SWIPE_THRESHOLD}px)` : 'translateX(0)';
    }
  };

  const getStatusConfig = (status: InvoiceStatus) => {
    switch (status) {
      case 'completed':
        return {
          icon: <CheckCircle2 className="w-4 h-4" />,
          text: 'Hoàn thành',
          className: 'badge-success',
        };
      case 'failed':
        return {
          icon: <XCircle className="w-4 h-4" />,
          text: 'Thất bại',
          className: 'badge-error',
        };
      case 'processing':
        return {
          icon: <Loader2 className="w-4 h-4 animate-spin" />,
          text: 'Đang xử lý',
          className: 'badge-info',
        };
      case 'pending':
        return {
          icon: <Clock className="w-4 h-4" />,
          text: 'Chờ xử lý',
          className: 'badge-warning',
        };
      default:
        return {
          icon: <Clock className="w-4 h-4" />,
          text: status,
          className: 'badge',
        };
    }
  };

  const handleInvoiceClick = (invoiceId: string) => {
    // Prevent navigation if item is swiped
    if (swipedId === invoiceId) {
      setSwipedId(null);
      return;
    }
    navigate(`/extract-invoice-data/${invoiceId}`);
  };

  const handlePrevPage = () => {
    if (page > 1) setPage(page - 1);
  };

  const handleNextPage = () => {
    if (page < totalPages) setPage(page + 1);
  };

  return (
    <div className="page-container">
      <header className="page-header safe-top flex justify-center">
        <h1 className="page-title text-center">Hóa đơn</h1>
      </header>

      <main className="page-content p-4">
        {isLoading && (
          <div className="flex flex-col items-center justify-center h-64 animate-fade-in">
            <div className="w-12 h-12 rounded-2xl bg-primary-100 dark:bg-primary-900/30 flex items-center justify-center mb-4">
              <Loader2 className="w-6 h-6 text-primary-600 dark:text-primary-400 animate-spin" />
            </div>
            <p className="text-surface-500 dark:text-surface-400 text-sm">Đang tải hóa đơn...</p>
          </div>
        )}

        {error && (
          <div className="animate-fade-in p-4">
            <div className="card bg-error-50 dark:bg-error-500/10 border-error-200 dark:border-error-500/20">
              <div className="flex items-start gap-3">
                <XCircle className="w-5 h-5 text-error-500 shrink-0 mt-0.5" />
                <div>
                  <p className="font-medium text-error-700 dark:text-error-400">Không thể tải</p>
                  <p className="text-sm text-error-600 dark:text-error-400/80 mt-0.5">
                    {error instanceof Error ? error.message : 'Lỗi không xác định'}
                  </p>
                </div>
              </div>
            </div>
          </div>
        )}

        {!isLoading && !error && (
          <div className="flex flex-col gap-3 animate-fade-in">
            {invoices.length === 0 ? (
              <div className="flex flex-col items-center justify-center py-16 px-4">
                <div className="w-20 h-20 rounded-3xl bg-surface-100 dark:bg-surface-800 flex items-center justify-center mb-4">
                  <FileText className="w-10 h-10 text-surface-400 dark:text-surface-500" />
                </div>
                <h3 className="text-lg font-semibold text-surface-700 dark:text-surface-300 mb-1">
                  Chưa có hóa đơn nào
                </h3>
                <p className="text-sm text-surface-500 dark:text-surface-400 text-center max-w-[240px]">
                  Chụp ảnh hóa đơn để bắt đầu trích xuất dữ liệu
                </p>
              </div>
            ) : (
              <>
                {invoices.map((invoice, index) => {
                  const statusConfig = getStatusConfig(invoice.status);
                  let cardRef: HTMLDivElement | null = null;
                  return (
                    <div
                      key={invoice.id}
                      className="relative overflow-hidden rounded-2xl"
                      style={{ animationDelay: `${index * 50}ms` }}
                    >
                      {/* Delete action background */}
                      <div className="absolute inset-y-0 right-0 w-20 bg-error-500 flex items-center justify-center">
                        <button
                          onClick={() => handleDelete(invoice.id)}
                          disabled={deletingId === invoice.id}
                          className="w-full h-full flex items-center justify-center text-white"
                          aria-label="Xóa hóa đơn"
                        >
                          {deletingId === invoice.id ? (
                            <Loader2 className="w-6 h-6 animate-spin" />
                          ) : (
                            <Trash2 className="w-6 h-6" />
                          )}
                        </button>
                      </div>

                      {/* Swipeable card */}
                      <div
                        ref={(el) => { cardRef = el; }}
                        className="card-interactive flex items-center gap-4 text-left w-full relative bg-white dark:bg-surface-900 transition-transform duration-150"
                        onClick={() => handleInvoiceClick(invoice.id)}
                        onTouchStart={(e) => handleTouchStart(e, invoice.id)}
                        onTouchMove={(e) => handleTouchMove(e, invoice.id, cardRef)}
                        onTouchEnd={() => handleTouchEnd(invoice.id, cardRef)}
                      >
                        <div className="w-14 h-14 rounded-xl bg-surface-100 dark:bg-surface-800 overflow-hidden shrink-0 flex items-center justify-center">
                          {invoice.image_path ? (
                            <img
                              alt="Ảnh đại diện hóa đơn"
                              className="w-full h-full object-cover"
                              src={getImageUrl(invoice.image_path) || ''}
                              onError={(e) => {
                                const target = e.target as HTMLImageElement;
                                target.style.display = 'none';
                                target.nextElementSibling?.classList.remove('hidden');
                              }}
                            />
                          ) : null}
                          <FileText className={`w-6 h-6 text-surface-400 ${invoice.image_path ? 'hidden' : ''}`} />
                        </div>

                        <div className="flex-grow min-w-0">
                          <h2 className="font-semibold text-surface-900 dark:text-white truncate">
                            Hóa đơn #{invoice.id.slice(-8).toUpperCase()}
                          </h2>
                          <div className="flex items-center gap-2 mt-1.5">
                            <span className={statusConfig.className}>
                              {statusConfig.icon}
                              <span>{statusConfig.text}</span>
                            </span>
                          </div>
                          {invoice.error_message && (
                            <p className="text-xs text-error-500 mt-1 truncate">
                              {invoice.error_message}
                            </p>
                          )}
                        </div>

                        <ChevronRight className="w-5 h-5 text-surface-400 shrink-0" />
                      </div>
                    </div>
                  );
                })}

                {/* Pagination Controls */}
                {totalPages > 1 && (
                  <div className="flex items-center justify-center gap-3 mt-4 pb-20">
                    <button
                      onClick={handlePrevPage}
                      disabled={page <= 1}
                      className="pagination-btn"
                      aria-label="Trang trước"
                    >
                      <ChevronLeft className="w-5 h-5" />
                    </button>

                    <span className="text-sm text-surface-600 dark:text-surface-400 font-medium">
                      {page} / {totalPages}
                    </span>

                    <button
                      onClick={handleNextPage}
                      disabled={page >= totalPages}
                      className="pagination-btn"
                      aria-label="Trang sau"
                    >
                      <ChevronRight className="w-5 h-5" />
                    </button>
                  </div>
                )}

                {/* Total count indicator */}
                {total > 0 && (
                  <p className="text-center text-xs text-surface-400 dark:text-surface-500">
                    Tổng: {total} hóa đơn
                  </p>
                )}
              </>
            )}
          </div>
        )}
      </main>

      <div className="fixed bottom-6 right-6 safe-bottom z-50">
        <button
          className="fab"
          onClick={() => navigate('/take-picture')}
          aria-label="Quét hóa đơn mới"
        >
          <Camera className="w-6 h-6" />
        </button>
      </div>
    </div>
  );
}

