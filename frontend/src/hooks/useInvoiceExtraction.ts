import { useMutation } from '@tanstack/react-query';
import { apiClient } from '@/lib/api';

export function useInvoiceExtraction() {
  return useMutation({
    mutationFn: (imageDataUrl: string) => apiClient.extractInvoice(imageDataUrl),
    onError: (error) => {
      console.error('Invoice extraction failed:', error);
    },
  });
}