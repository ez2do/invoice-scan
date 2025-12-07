import { create } from 'zustand';
import type { InvoiceData } from '@/types';

interface AppStore {
  // State
  currentImage: string | null;
  extractedData: InvoiceData | null;
  isLoading: boolean;
  error: string | null;
  selectedInvoiceId: string | null;
  isDirty: boolean;

  // Actions
  setCurrentImage: (image: string | null) => void;
  setExtractedData: (data: InvoiceData | null) => void;
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  setSelectedInvoiceId: (id: string | null) => void;
  setDirty: (value: boolean) => void;
  resetDirty: () => void;
  clearData: () => void;
  updateKeyValue: (index: number, key: string, value: string) => void;
  updateTableCell: (rowIndex: number, colIndex: number, value: string) => void;
  updateSummary: (index: number, key: string, value: string) => void;
}

export const useAppStore = create<AppStore>((set, get) => ({
  // Initial state
  currentImage: null,
  extractedData: null,
  isLoading: false,
  error: null,
  selectedInvoiceId: null,
  isDirty: false,

  // Actions
  setCurrentImage: (image) => set({ currentImage: image }),

  setExtractedData: (data) => set({ extractedData: data, isDirty: false }),

  setLoading: (loading) => set({ isLoading: loading }),

  setError: (error) => set({ error }),

  setSelectedInvoiceId: (id) => set({ selectedInvoiceId: id }),

  setDirty: (value) => set({ isDirty: value }),

  resetDirty: () => set({ isDirty: false }),

  clearData: () => set({
    currentImage: null,
    extractedData: null,
    isLoading: false,
    error: null,
    selectedInvoiceId: null,
    isDirty: false,
  }),

  updateKeyValue: (index, key, value) => {
    const { extractedData } = get();
    if (!extractedData) return;

    const updated = { ...extractedData };
    updated.keyValuePairs = [...updated.keyValuePairs];
    updated.keyValuePairs[index] = { ...updated.keyValuePairs[index], key, value };

    set({ extractedData: updated, isDirty: true });
  },

  updateTableCell: (rowIndex, colIndex, value) => {
    const { extractedData } = get();
    if (!extractedData?.table) return;

    const updated = { ...extractedData };
    const table = updated.table;
    if (table) {
      updated.table = {
        headers: [...table.headers],
        rows: table.rows.map((row, idx) =>
          idx === rowIndex
            ? row.map((cell, cellIdx) => cellIdx === colIndex ? value : cell)
            : row
        )
      };
    }

    set({ extractedData: updated, isDirty: true });
  },

  updateSummary: (index, key, value) => {
    const { extractedData } = get();
    if (!extractedData) return;

    const updated = { ...extractedData };
    updated.summary = [...updated.summary];
    updated.summary[index] = { ...updated.summary[index], key, value };

    set({ extractedData: updated, isDirty: true });
  },
}));