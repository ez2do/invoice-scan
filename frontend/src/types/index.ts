// Invoice Data Types
export interface KeyValuePair {
  key: string;
  value: string;
  confidence?: number;
}

export interface TableData {
  headers: string[];
  rows: string[][];
}

export interface ExtractedData {
  key_value_pairs: KeyValuePair[] | null;
  table: TableData;
  summary: KeyValuePair[];
  confidence?: number;
}

export interface InvoiceData {
  keyValuePairs: KeyValuePair[];
  table: TableData | null;
  summary: KeyValuePair[];
  confidence?: number;
}

// API Types
export interface ExtractRequest {
  image: string; // base64 encoded
}

export interface ExtractResponse {
  success: boolean;
  data?: InvoiceData;
  error?: string;
  processingTime?: number;
}

export type InvoiceStatus = 'pending' | 'processing' | 'completed' | 'failed';

export interface InvoiceListItem {
  id: string;
  status: InvoiceStatus;
  image_path: string;
  created_at: string;
  updated_at?: string;
  extracted_data?: ExtractedData;
  error_message?: string;
}

export interface UploadResponse {
  success: boolean;
  data?: InvoiceListItem;
  error?: string;
}

export interface InvoicesResponse {
  success: boolean;
  data?: InvoiceListItem[];
  error?: string;
}

export interface InvoiceResponse {
  success: boolean;
  data?: InvoiceListItem;
  error?: string;
}

// App State Types
export interface AppState {
  currentImage: string | null;
  extractedData: InvoiceData | null;
  isLoading: boolean;
  error: string | null;
}

// Camera Types
export interface CameraConfig {
  video: {
    facingMode: 'environment' | 'user';
    width: { ideal: number };
    height: { ideal: number };
  };
}

// UI Component Types
export interface ButtonProps {
  variant?: 'primary' | 'secondary' | 'ghost';
  size?: 'sm' | 'md' | 'lg';
  loading?: boolean;
  disabled?: boolean;
  children: React.ReactNode;
  onClick?: () => void;
  type?: 'button' | 'submit';
  className?: string;
}

export interface InputProps {
  label?: string;
  error?: string;
  placeholder?: string;
  value?: string;
  onChange?: (value: string) => void;
  type?: 'text' | 'number' | 'email';
  required?: boolean;
  className?: string;
}