import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

// Image processing utilities
export function compressImage(file: File, maxWidth: number = 1024, quality: number = 0.8): Promise<string> {
  return new Promise((resolve) => {
    const canvas = document.createElement('canvas');
    const ctx = canvas.getContext('2d')!;
    const img = new Image();

    img.onload = () => {
      const ratio = Math.min(maxWidth / img.width, maxWidth / img.height);
      canvas.width = img.width * ratio;
      canvas.height = img.height * ratio;

      ctx.drawImage(img, 0, 0, canvas.width, canvas.height);
      resolve(canvas.toDataURL('image/jpeg', quality));
    };

    img.src = URL.createObjectURL(file);
  });
}

export function dataURLToFile(dataURL: string, filename: string = 'image.jpg'): File {
  const arr = dataURL.split(',');
  const mime = arr[0].match(/:(.*?);/)![1];
  const bstr = atob(arr[1]);
  let n = bstr.length;
  const u8arr = new Uint8Array(n);
  
  while (n--) {
    u8arr[n] = bstr.charCodeAt(n);
  }
  
  return new File([u8arr], filename, { type: mime });
}

// Format currency (Vietnamese Dong)
export function formatCurrency(value: string): string {
  const num = parseFloat(value.replace(/[^\d]/g, ''));
  if (isNaN(num)) return value;
  return new Intl.NumberFormat('vi-VN').format(num);
}

// Detect if value looks like currency
export function isCurrency(value: string): boolean {
  return /[\d,.]+(đ|vnd|dong|\₫|vnđ)/i.test(value) || /^\d{1,3}(,\d{3})*(\.\d{2})?$/.test(value);
}