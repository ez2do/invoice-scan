import { useMutation } from '@tanstack/react-query';
import { apiClient } from '@/lib/api';
import type { ExtractRequest } from '@/types';

export function useInvoiceExtraction() {
  return useMutation({
    mutationFn: (request: ExtractRequest) => apiClient.extractInvoice(request),
    onError: (error) => {
      console.error('Invoice extraction failed:', error);
    },
  });
}