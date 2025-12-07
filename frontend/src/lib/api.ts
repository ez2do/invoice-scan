import type {
  ExtractResponse,
  UploadResponse,
  InvoiceResponse,
  PaginatedInvoicesResponse
} from '@/types';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'https://localhost:3001/api/v1';

function getImageUrl(imagePath: string | null | undefined): string | null {
  if (!imagePath) return null;

  if (imagePath.startsWith('data:')) {
    return imagePath;
  }

  const baseUrl = API_BASE_URL.replace('/api/v1', '');

  if (imagePath.startsWith('/')) {
    return `${baseUrl}${imagePath}`;
  }

  return `${baseUrl}/${imagePath}`;
}

export { getImageUrl };

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

  async uploadInvoice(imageDataUrl: string): Promise<UploadResponse> {
    try {
      const blob = await this.dataURLToBlob(imageDataUrl);
      const formData = new FormData();
      formData.append('image', blob, 'invoice.jpg');

      const response = await fetch(`${this.baseURL}/invoices/upload`, {
        method: 'POST',
        body: formData,
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.error || `HTTP error! status: ${response.status}`);
      }

      return await response.json();
    } catch (error) {
      console.error('Upload Error:', error);
      throw new Error(
        error instanceof Error ? error.message : 'Failed to upload invoice'
      );
    }
  }

  async getInvoices(page: number = 1, pageSize: number = 10): Promise<PaginatedInvoicesResponse> {
    try {
      const params = new URLSearchParams({
        page: page.toString(),
        page_size: pageSize.toString(),
      });
      const response = await fetch(`${this.baseURL}/invoices?${params}`);

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.error || `HTTP error! status: ${response.status}`);
      }

      return await response.json();
    } catch (error) {
      console.error('Get Invoices Error:', error);
      throw new Error(
        error instanceof Error ? error.message : 'Failed to fetch invoices'
      );
    }
  }

  async getInvoice(id: string): Promise<InvoiceResponse> {
    try {
      const response = await fetch(`${this.baseURL}/invoices/${id}`);

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.error || `HTTP error! status: ${response.status}`);
      }

      return await response.json();
    } catch (error) {
      console.error('Get Invoice Error:', error);
      throw new Error(
        error instanceof Error ? error.message : 'Failed to fetch invoice'
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