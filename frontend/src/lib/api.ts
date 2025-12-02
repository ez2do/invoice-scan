import type { ExtractResponse } from '@/types';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:3001/api';

class APIClient {
  private baseURL: string;

  constructor(baseURL: string = API_BASE_URL) {
    this.baseURL = baseURL;
  }

  async extractInvoice(imageDataUrl: string): Promise<ExtractResponse> {
    try {
      const blob = await this.dataURLToBlob(imageDataUrl);
      const formData = new FormData();
      formData.append('image', blob, 'invoice.jpg');

      const response = await fetch(`${this.baseURL}/extract`, {
        method: 'POST',
        body: formData,
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.error || `HTTP error! status: ${response.status}`);
      }

      return await response.json();
    } catch (error) {
      console.error('API Error:', error);
      throw new Error(
        error instanceof Error ? error.message : 'Failed to extract invoice data'
      );
    }
  }

  private async dataURLToBlob(dataURL: string): Promise<Blob> {
    const response = await fetch(dataURL);
    return await response.blob();
  }

  async healthCheck(): Promise<{ status: string }> {
    try {
      const response = await fetch(`${this.baseURL}/health`);
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      return await response.json();
    } catch (error) {
      console.error('Health check failed:', error);
      throw new Error('API is not available');
    }
  }
}

export const apiClient = new APIClient();
export { APIClient };