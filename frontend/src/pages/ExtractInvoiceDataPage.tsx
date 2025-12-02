import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useMutation } from '@tanstack/react-query';
import { useAppStore } from '@/stores/app-store';
import { apiClient } from '@/lib/api';
import { Input } from '@/components/ui/Input';

export default function ExtractInvoiceDataPage() {
  const navigate = useNavigate();
  const { 
    currentImage, 
    extractedData, 
    isLoading,
    setExtractedData, 
    setLoading, 
    setError,
    updateKeyValue,
    updateTableCell,
    updateSummary,
    clearData
  } = useAppStore();

  const extractMutation = useMutation({
    mutationFn: apiClient.extractInvoice.bind(apiClient),
    onSuccess: (response: any) => {
      if (response.success && response.data) {
        setExtractedData(response.data);
      } else {
        setError(response.error || 'Failed to extract data');
      }
      setLoading(false);
    },
    onError: (error: Error) => {
      setError(error.message);
      setLoading(false);
    },
  });

  useEffect(() => {
    if (!currentImage) {
      navigate('/take-picture');
      return;
    }

    if (!extractedData && !isLoading) {
      setLoading(true);
      extractMutation.mutate(currentImage);
    }
  }, [currentImage, extractedData, isLoading]);

  const handleBack = () => {
    clearData();
    navigate('/list-invoices');
  };

  const handleComplete = () => {
    clearData();
    navigate('/list-invoices');
  };

  if (!currentImage) {
    return null;
  }

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
              <p className="text-gray-600 dark:text-gray-300">Extracting invoice data...</p>
            </div>
          </div>
        )}

        {/* Error State */}
        {extractMutation.error && (
          <div className="p-4">
            <div className="bg-red-50 border border-red-200 rounded-lg p-4">
              <p className="text-red-700">
                Failed to extract data: {extractMutation.error.message}
              </p>
            </div>
          </div>
        )}

        {/* Main Content */}
        {extractedData && (
          <div className="flex h-full">
            {/* Left Panel - Invoice Image */}
            <div className="w-1/2 border-r border-border-light dark:border-border-dark bg-gray-50 dark:bg-gray-800 p-4">
              <div className="space-y-2">
                <p className="text-xs font-medium text-slate-600 dark:text-slate-400">Viewing Invoice Image</p>
                <div className="rounded-lg overflow-hidden">
                  <img
                    src={currentImage}
                    alt="Invoice"
                    className="w-full h-auto object-contain max-h-96"
                  />
                </div>
              </div>
            </div>

            {/* Right Panel - Extracted Data */}
            <div className="w-1/2 flex flex-col">
              <div className="flex-1 overflow-y-auto p-4 space-y-6">
                {/* Key-Value Pairs */}
                {extractedData.keyValuePairs.length > 0 && (
                  <div className="space-y-3">
                    <h3 className="font-medium text-gray-900 dark:text-white text-sm">Invoice Information</h3>
                    <div className="space-y-3">
                      {extractedData.keyValuePairs.map((pair, index) => (
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
                {extractedData.table && (
                  <div className="space-y-3">
                    <h3 className="font-medium text-gray-900 dark:text-white text-sm">Line Items</h3>
                    <div className="overflow-x-auto">
                      <table className="min-w-full border border-gray-200 dark:border-gray-700 rounded-lg text-sm">
                        <thead className="bg-gray-50 dark:bg-gray-800">
                          <tr>
                            {extractedData.table.headers.map((header, index) => (
                              <th key={index} className="px-2 py-1 text-left text-xs font-medium text-gray-900 dark:text-white border-b">
                                {header}
                              </th>
                            ))}
                          </tr>
                        </thead>
                        <tbody>
                          {extractedData.table.rows.map((row, rowIndex) => (
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
                {extractedData.summary.length > 0 && (
                  <div className="space-y-3">
                    <h3 className="font-medium text-gray-900 dark:text-white text-sm">Summary</h3>
                    <div className="space-y-3">
                      {extractedData.summary.map((item, index) => (
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