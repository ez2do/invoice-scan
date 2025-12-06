import { useEffect, useMemo, useState, useRef, useCallback } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { 
  ArrowLeft, 
  Loader2, 
  CheckCircle2, 
  XCircle, 
  Image, 
  FileText, 
  Table2, 
  Calculator,
  Save,
  ChevronDown,
  Package,
  LayoutGrid,
  List
} from 'lucide-react';
import { useAppStore } from '@/stores/app-store';
import { apiClient, getImageUrl } from '@/lib/api';
import type { ExtractedData, InvoiceData } from '@/types';

interface AutoExpandTextareaProps {
  value: string;
  onChange: (value: string) => void;
  className?: string;
  placeholder?: string;
}

function AutoExpandTextarea({ value, onChange, className = '', placeholder }: AutoExpandTextareaProps) {
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  const adjustHeight = useCallback(() => {
    const textarea = textareaRef.current;
    if (textarea) {
      textarea.style.height = 'auto';
      textarea.style.height = `${Math.max(36, textarea.scrollHeight)}px`;
    }
  }, []);

  useEffect(() => {
    adjustHeight();
  }, [value, adjustHeight]);

  return (
    <textarea
      ref={textareaRef}
      value={value}
      onChange={(e) => onChange(e.target.value)}
      placeholder={placeholder}
      rows={1}
      className={`w-full resize-none overflow-hidden bg-surface-50 dark:bg-surface-800 
        border border-surface-200 dark:border-surface-700 rounded-lg px-3 py-2
        text-sm text-surface-900 dark:text-white
        placeholder:text-surface-400 dark:placeholder:text-surface-500
        focus:outline-none focus:ring-2 focus:ring-primary-500/30 focus:border-primary-500
        transition-all duration-200 ${className}`}
    />
  );
}

interface LineItemCardProps {
  row: string[];
  headers: string[];
  rowIndex: number;
  onCellChange: (rowIndex: number, cellIndex: number, value: string) => void;
}

function LineItemCard({ row, headers, rowIndex, onCellChange }: LineItemCardProps) {
  const [isExpanded, setIsExpanded] = useState(false);
  const primaryField = row[0] || `Mục ${rowIndex + 1}`;

  return (
    <div className="card overflow-hidden">
      <button
        type="button"
        onClick={() => setIsExpanded(!isExpanded)}
        className="w-full flex items-center justify-between gap-3 text-left"
      >
        <div className="flex items-center gap-3 min-w-0 flex-1">
          <div className="w-8 h-8 rounded-lg bg-primary-100 dark:bg-primary-900/30 flex items-center justify-center shrink-0">
            <Package className="w-4 h-4 text-primary-600 dark:text-primary-400" />
          </div>
          <div className="min-w-0 flex-1">
            <p className="text-sm font-medium text-surface-900 dark:text-white truncate">
              {primaryField}
            </p>
            {!isExpanded && row.length > 1 && (
              <p className="text-xs text-surface-500 dark:text-surface-400 truncate mt-0.5">
                {headers.slice(1).map((h, i) => `${h}: ${row[i + 1] || '-'}`).join(' • ')}
              </p>
            )}
          </div>
        </div>
        <ChevronDown 
          className={`w-5 h-5 text-surface-400 shrink-0 transition-transform duration-200 ${
            isExpanded ? 'rotate-180' : ''
          }`} 
        />
      </button>

      {isExpanded && (
        <div className="mt-4 pt-4 border-t border-surface-100 dark:border-surface-800 space-y-3 animate-fade-in">
          {headers.map((header, cellIndex) => (
            <div key={cellIndex}>
              <label className="block text-xs font-medium text-surface-500 dark:text-surface-400 mb-1.5">
                {header}
              </label>
              <AutoExpandTextarea
                value={row[cellIndex] || ''}
                onChange={(value) => onCellChange(rowIndex, cellIndex, value)}
              />
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

function convertExtractedDataToInvoiceData(extractedData: ExtractedData): InvoiceData {
  return {
    keyValuePairs: extractedData.key_value_pairs || [],
    table: extractedData.table?.headers?.length > 0 ? extractedData.table : null,
    summary: extractedData.summary || [],
    confidence: extractedData.confidence,
  };
}

export default function ExtractInvoiceDataPage() {
  const navigate = useNavigate();
  const { id } = useParams<{ id?: string }>();
  const [tableViewMode, setTableViewMode] = useState<'card' | 'table'>('card');
  const {
    currentImage,
    selectedInvoiceId,
    extractedData,
    setExtractedData,
    updateKeyValue,
    updateTableCell,
    updateSummary
  } = useAppStore();

  const invoiceId = id || selectedInvoiceId;

  useEffect(() => {
    if (invoiceId) {
      setExtractedData(null);
    }
  }, [invoiceId, setExtractedData]);

  const { data: invoiceResponse, isLoading, error } = useQuery({
    queryKey: ['invoice', invoiceId],
    queryFn: () => {
      if (!invoiceId) throw new Error('Yêu cầu ID hóa đơn');
      return apiClient.getInvoice(invoiceId);
    },
    enabled: !!invoiceId,
    refetchOnMount: true,
    staleTime: 0,
    refetchInterval: (query) => {
      const data = query.state.data;
      if (!data?.data) return false;
      const status = data.data.status;
      if (status === 'completed' || status === 'failed') {
        return false;
      }
      return 2000;
    },
  });

  const invoice = invoiceResponse?.data;
  const invoiceData = useMemo<InvoiceData | null>(() => {
    if (!invoice?.extracted_data) return null;
    return convertExtractedDataToInvoiceData(invoice.extracted_data);
  }, [invoice?.extracted_data]);

  useEffect(() => {
    if (invoiceData && invoice?.id === invoiceId) {
      setExtractedData(invoiceData);
    }
  }, [invoiceData, invoiceId, invoice?.id, setExtractedData]);

  const handleBack = () => {
    navigate('/list-invoices');
  };

  const handleComplete = () => {
    navigate('/list-invoices');
  };

  if (!invoiceId) {
    navigate('/list-invoices');
    return null;
  }

  const displayImage = getImageUrl(invoice?.image_path || currentImage);
  const displayData = invoiceData || extractedData;

  return (
    <div className="page-container">
      <header className="page-header safe-top">
        <button
          onClick={handleBack}
          className="icon-btn"
          aria-label="Quay lại"
        >
          <ArrowLeft className="w-5 h-5" />
        </button>
        <h1 className="page-title">Chi tiết hóa đơn</h1>
        <div className="w-10 h-10" />
      </header>

      <div className="flex-1 overflow-hidden">
        {isLoading && (
          <div className="flex flex-col items-center justify-center h-full animate-fade-in">
            <div className="w-16 h-16 rounded-2xl bg-primary-100 dark:bg-primary-900/30 flex items-center justify-center mb-4">
              <Loader2 className="w-8 h-8 text-primary-600 dark:text-primary-400 animate-spin" />
            </div>
            <p className="text-surface-600 dark:text-surface-400 font-medium">
              {invoice?.status === 'processing' ? 'Đang trích xuất dữ liệu...' : 'Đang tải hóa đơn...'}
            </p>
            <p className="text-surface-400 dark:text-surface-500 text-sm mt-1">
              Vui lòng đợi trong giây lát
            </p>
          </div>
        )}

        {error && (
          <div className="p-4 animate-fade-in">
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

        {invoice?.status === 'failed' && invoice.error_message && (
          <div className="p-4 animate-fade-in">
            <div className="card bg-error-50 dark:bg-error-500/10 border-error-200 dark:border-error-500/20">
              <div className="flex items-start gap-3">
                <XCircle className="w-5 h-5 text-error-500 shrink-0 mt-0.5" />
                <div>
                  <p className="font-medium text-error-700 dark:text-error-400">Trích xuất thất bại</p>
                  <p className="text-sm text-error-600 dark:text-error-400/80 mt-0.5">
                    {invoice.error_message}
                  </p>
                </div>
              </div>
            </div>
          </div>
        )}

        {!isLoading && invoice && invoice.status !== 'completed' && invoice.status !== 'failed' && (
          <div className="flex flex-col items-center justify-center h-full animate-fade-in">
            <div className="w-16 h-16 rounded-2xl bg-primary-100 dark:bg-primary-900/30 flex items-center justify-center mb-4">
              <Loader2 className="w-8 h-8 text-primary-600 dark:text-primary-400 animate-spin" />
            </div>
            <p className="text-surface-600 dark:text-surface-400 font-medium">
              {invoice.status === 'processing' ? 'Đang trích xuất dữ liệu hóa đơn...' : 'Đang chờ bắt đầu...'}
            </p>
            <p className="text-surface-400 dark:text-surface-500 text-sm mt-1">
              Vui lòng đợi trong khi chúng tôi phân tích hóa đơn của bạn
            </p>
          </div>
        )}

        {displayData && displayImage && invoice?.status === 'completed' && (
          <div className="flex flex-col h-full animate-fade-in">
            <div className="border-b border-surface-200 dark:border-surface-800 bg-surface-100 dark:bg-surface-900 p-4">
              <div className="flex items-center gap-2 mb-3">
                <Image className="w-4 h-4 text-surface-500" />
                <span className="text-xs font-medium text-surface-500 uppercase tracking-wide">Hình ảnh hóa đơn</span>
              </div>
              <div className="rounded-xl overflow-hidden bg-white dark:bg-surface-800 shadow-inner-soft">
                <img
                  src={displayImage}
                  alt="Hóa đơn"
                  className="w-full h-auto object-contain max-h-52"
                />
              </div>
            </div>

            <div className="flex flex-col flex-1 min-h-0">
              <div className="flex-1 overflow-y-auto p-4 space-y-6 scrollbar-hide">
                {displayData.keyValuePairs.length > 0 && (
                  <section className="space-y-3">
                    <div className="flex items-center gap-2">
                      <FileText className="w-4 h-4 text-primary-500" />
                      <h3 className="font-semibold text-surface-900 dark:text-white text-sm">
                        Thông tin hóa đơn
                      </h3>
                    </div>
                    <div className="space-y-3">
                      {displayData.keyValuePairs.map((pair, index) => (
                        <div key={index}>
                          <label className="block text-xs font-medium text-surface-500 dark:text-surface-400 mb-1.5">
                            {pair.key}
                          </label>
                          <AutoExpandTextarea
                            value={pair.value}
                            onChange={(value) => updateKeyValue(index, pair.key, value)}
                          />
                        </div>
                      ))}
                    </div>
                  </section>
                )}

                {displayData.table && (
                  <section className="space-y-3">
                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-2">
                        <Table2 className="w-4 h-4 text-primary-500" />
                        <h3 className="font-semibold text-surface-900 dark:text-white text-sm">
                          Chi tiết hàng hóa
                        </h3>
                        <span className="text-xs text-surface-500 bg-surface-100 dark:bg-surface-800 px-2 py-1 rounded-full">
                          {displayData.table.rows.length} mục
                        </span>
                      </div>
                      <div className="flex items-center gap-1 p-1 bg-surface-100 dark:bg-surface-800 rounded-lg">
                        <button
                          type="button"
                          onClick={() => setTableViewMode('card')}
                          className={`p-1.5 rounded-md transition-all duration-200 ${
                            tableViewMode === 'card'
                              ? 'bg-white dark:bg-surface-700 text-primary-600 dark:text-primary-400 shadow-sm'
                              : 'text-surface-500 hover:text-surface-700 dark:hover:text-surface-300'
                          }`}
                          aria-label="Xem dạng thẻ"
                        >
                          <LayoutGrid className="w-4 h-4" />
                        </button>
                        <button
                          type="button"
                          onClick={() => setTableViewMode('table')}
                          className={`p-1.5 rounded-md transition-all duration-200 ${
                            tableViewMode === 'table'
                              ? 'bg-white dark:bg-surface-700 text-primary-600 dark:text-primary-400 shadow-sm'
                              : 'text-surface-500 hover:text-surface-700 dark:hover:text-surface-300'
                          }`}
                          aria-label="Xem dạng bảng"
                        >
                          <List className="w-4 h-4" />
                        </button>
                      </div>
                    </div>

                    {tableViewMode === 'card' ? (
                      <div className="space-y-3">
                        {displayData.table.rows.map((row, rowIndex) => (
                          <LineItemCard
                            key={rowIndex}
                            row={row}
                            headers={displayData.table!.headers}
                            rowIndex={rowIndex}
                            onCellChange={updateTableCell}
                          />
                        ))}
                      </div>
                    ) : (
                      <div className="overflow-x-auto -mx-4 px-4">
                        <div className="inline-block min-w-full align-middle">
                          <div className="overflow-hidden rounded-xl border border-surface-200 dark:border-surface-700">
                            <table className="min-w-full divide-y divide-surface-200 dark:divide-surface-700">
                              <thead className="bg-surface-50 dark:bg-surface-800">
                                <tr>
                                  {displayData.table.headers.map((header, index) => (
                                    <th
                                      key={index}
                                      className="px-3 py-2.5 text-left text-xs font-semibold text-surface-600 dark:text-surface-300 uppercase tracking-wide whitespace-nowrap"
                                    >
                                      {header}
                                    </th>
                                  ))}
                                </tr>
                              </thead>
                              <tbody className="divide-y divide-surface-100 dark:divide-surface-800 bg-white dark:bg-surface-900">
                                {displayData.table.rows.map((row, rowIndex) => (
                                  <tr key={rowIndex}>
                                    {row.map((cell, cellIndex) => (
                                      <td key={cellIndex} className="px-3 py-2 align-top">
                                        <AutoExpandTextarea
                                          value={cell}
                                          onChange={(value) => updateTableCell(rowIndex, cellIndex, value)}
                                          className="border-0 bg-transparent p-0 min-h-[24px] text-xs focus:ring-0"
                                        />
                                      </td>
                                    ))}
                                  </tr>
                                ))}
                              </tbody>
                            </table>
                          </div>
                        </div>
                      </div>
                    )}
                  </section>
                )}

                {displayData.summary.length > 0 && (
                  <section className="space-y-3">
                    <div className="flex items-center gap-2">
                      <Calculator className="w-4 h-4 text-primary-500" />
                      <h3 className="font-semibold text-surface-900 dark:text-white text-sm">
                        Tổng kết
                      </h3>
                    </div>
                    <div className="card bg-primary-50 dark:bg-primary-900/20 border-primary-100 dark:border-primary-800/30">
                      <div className="space-y-3">
                        {displayData.summary.map((item, index) => (
                          <div key={index} className="flex items-start gap-4">
                            <label className="text-sm font-medium text-surface-600 dark:text-surface-400 w-[35%] shrink-0 pt-2">
                              {item.key}
                            </label>
                            <AutoExpandTextarea
                              value={item.value}
                              onChange={(value) => updateSummary(index, item.key, value)}
                              className="text-right font-semibold flex-1 bg-white dark:bg-surface-800"
                            />
                          </div>
                        ))}
                      </div>
                    </div>
                  </section>
                )}

                {displayData.confidence !== undefined && (
                  <div className="flex items-center justify-center gap-2 py-2">
                    <CheckCircle2 className="w-4 h-4 text-accent-500" />
                    <span className="text-xs text-surface-500">
                      Độ tin cậy trích xuất: {Math.round(displayData.confidence * 100)}%
                    </span>
                  </div>
                )}
              </div>

              <div className="shrink-0 p-4 border-t border-surface-200 dark:border-surface-800 bg-white dark:bg-surface-900">
                <button
                  className="btn-primary w-full h-12"
                  onClick={handleComplete}
                >
                  <Save className="w-5 h-5" />
                  <span>Lưu & Hoàn thành</span>
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
