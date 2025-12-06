import { useEffect, useMemo } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { useAppStore } from '@/stores/app-store';
import { apiClient, getImageUrl } from '@/lib/api';
import { Input } from '@/components/ui/Input';
import type { ExtractedData, InvoiceData } from '@/types';

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

  const { data: invoiceResponse, isLoading, error } = useQuery({
    queryKey: ['invoice', invoiceId],
    queryFn: () => {
      if (!invoiceId) throw new Error('Invoice ID is required');
      return apiClient.getInvoice(invoiceId);
    },
    enabled: !!invoiceId,
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
    if (invoiceData && !extractedData) {
      setExtractedData(invoiceData);
    }
  }, [invoiceData, extractedData, setExtractedData]);

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
  const displayData = extractedData || invoiceData;

  return (
    <div className="relative mx-auto flex h-screen max-h-[960px] w-full max-w-[480px] flex-col overflow-hidden bg-background-light dark:bg-background-dark">
      {/* Header */}
      <header className="flex items-center justify-between border-b border-border-light dark:border-border-dark bg-surface-light dark:bg-surface-dark px-4 py-3">
        <button 
          onClick={handleBack}
          className="flex h-10 w-10 items-center justify-center rounded-full"
        >
          <span className="text-2xl">←</span>
        </button>
        <h1 className="text-lg font-semibold">Invoice Data</h1>
        <div className="h-10 w-10"></div>
      </header>

      <div className="flex-1 overflow-hidden">
        {/* Loading State */}
        {isLoading && (
          <div className="flex items-center justify-center h-full">
            <div className="text-center space-y-4">
              <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto"></div>
              <p className="text-gray-600 dark:text-gray-300">
                {invoice?.status === 'processing' ? 'Extracting invoice data...' : 'Loading invoice...'}
              </p>
            </div>
          </div>
        )}

        {/* Error State */}
        {error && (
          <div className="p-4">
            <div className="bg-red-50 border border-red-200 rounded-lg p-4">
              <p className="text-red-700">
                Failed to load invoice: {error instanceof Error ? error.message : 'Unknown error'}
              </p>
            </div>
          </div>
        )}

        {invoice?.status === 'failed' && invoice.error_message && (
          <div className="p-4">
            <div className="bg-red-50 border border-red-200 rounded-lg p-4">
              <p className="text-red-700">Extraction failed: {invoice.error_message}</p>
            </div>
          </div>
        )}

        {/* Pending/Processing State - No Data Yet */}
        {!isLoading && invoice && invoice.status !== 'completed' && invoice.status !== 'failed' && (
          <div className="flex items-center justify-center h-full">
            <div className="text-center space-y-4">
              <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto"></div>
              <p className="text-gray-600 dark:text-gray-300">
                {invoice.status === 'processing' ? 'Extracting invoice data...' : 'Waiting to start extraction...'}
              </p>
            </div>
          </div>
        )}

        {/* Main Content */}
        {displayData && displayImage && invoice?.status === 'completed' && (
          <div className="flex flex-col h-full">
            {/* Top Panel - Invoice Image */}
            <div className="border-b border-border-light dark:border-border-dark bg-gray-50 dark:bg-gray-800 p-4">
              <div className="space-y-2">
                <p className="text-xs font-medium text-slate-600 dark:text-slate-400">Viewing Invoice Image</p>
                <div className="rounded-lg overflow-hidden flex justify-center">
                  <img
                    src={displayImage}
                    alt="Invoice"
                    className="w-full h-auto object-contain max-h-64"
                  />
                </div>
              </div>
            </div>

            {/* Bottom Panel - Extracted Data */}
            <div className="flex flex-col flex-1 min-h-0">
              <div className="flex-1 overflow-y-auto p-4 space-y-6">
                {/* Key-Value Pairs */}
                {displayData.keyValuePairs.length > 0 && (
                  <div className="space-y-3">
                    <h3 className="font-medium text-gray-900 dark:text-white text-sm">Invoice Information</h3>
                    <div className="space-y-3">
                      {displayData.keyValuePairs.map((pair, index) => (
                        <div key={index} className="space-y-1">
                          <label className="text-xs font-medium text-gray-600 dark:text-gray-400">
                            {pair.key}
                          </label>
                          <Input
                            value={pair.value}
                            onChange={(value) => updateKeyValue(index, pair.key, value)}
                            className="text-sm h-8"
                          />
                        </div>
                      ))}
                    </div>
                  </div>
                )}

                {/* Table Data */}
                {displayData.table && (
                  <div className="space-y-3">
                    <h3 className="font-medium text-gray-900 dark:text-white text-sm">Line Items</h3>
                    <div className="overflow-x-auto">
                      <table className="min-w-full border border-gray-200 dark:border-gray-700 rounded-lg text-sm">
                        <thead className="bg-gray-50 dark:bg-gray-800">
                          <tr>
                            {displayData.table.headers.map((header, index) => (
                              <th key={index} className="px-2 py-1 text-left text-xs font-medium text-gray-900 dark:text-white border-b">
                                {header}
                              </th>
                            ))}
                          </tr>
                        </thead>
                        <tbody>
                          {displayData.table.rows.map((row, rowIndex) => (
                            <tr key={rowIndex} className="border-b border-gray-100 dark:border-gray-700">
                              {row.map((cell, cellIndex) => (
                                <td key={cellIndex} className="px-2 py-1">
                                  <Input
                                    value={cell}
                                    onChange={(value) => updateTableCell(rowIndex, cellIndex, value)}
                                    className="border-0 shadow-none p-1 h-6 text-xs bg-transparent"
                                  />
                                </td>
                              ))}
                            </tr>
                          ))}
                        </tbody>
                      </table>
                    </div>
                  </div>
                )}

                {/* Summary */}
                {displayData.summary.length > 0 && (
                  <div className="space-y-3">
                    <h3 className="font-medium text-gray-900 dark:text-white text-sm">Summary</h3>
                    <div className="space-y-3">
                      {displayData.summary.map((item, index) => (
                        <div key={index} className="space-y-1">
                          <label className="text-xs font-medium text-gray-600 dark:text-gray-400">
                            {item.key}
                          </label>
                          <Input
                            value={item.value}
                            onChange={(value) => updateSummary(index, item.key, value)}
                            className="text-sm h-8"
                          />
                        </div>
                      ))}
                    </div>
                  </div>
                )}
              </div>

              {/* Footer Actions */}
              <div className="border-t border-border-light dark:border-border-dark p-4">
                <button 
                  className="w-full bg-primary text-white py-3 rounded-lg font-semibold hover:bg-primary/90 transition-colors"
                  onClick={handleComplete}
                >
                  Save & Complete
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}