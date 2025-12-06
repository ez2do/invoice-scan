import { forwardRef } from 'react';
import { cn } from '@/lib/utils';
import type { InputProps } from '@/types';

const Input = forwardRef<HTMLInputElement, InputProps>(
  ({ 
    label, 
    error, 
    className, 
    type = 'text',
    onChange,
    ...props 
  }, ref) => {
    return (
      <div className="w-full">
        {label && (
          <label className="block text-sm font-medium text-surface-600 dark:text-surface-400 mb-1.5">
            {label}
            {props.required && <span className="text-error-500 ml-1">*</span>}
          </label>
        )}
        <input
          ref={ref}
          type={type}
          className={cn(
            'w-full px-4 py-2.5 rounded-xl text-sm',
            'bg-surface-50 dark:bg-surface-800',
            'border border-surface-200 dark:border-surface-700',
            'text-surface-900 dark:text-white',
            'placeholder:text-surface-400 dark:placeholder:text-surface-500',
            'focus:outline-none focus:ring-2 focus:ring-primary-500/30 focus:border-primary-500',
            'transition-all duration-200',
            error && 'border-error-500 focus:border-error-500 focus:ring-error-500/30',
            className
          )}
          onChange={(e) => onChange?.(e.target.value)}
          {...props}
        />
        {error && (
          <p className="mt-1.5 text-sm text-error-500 flex items-center gap-1">
            <svg className="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <circle cx="12" cy="12" r="10" />
              <line x1="12" y1="8" x2="12" y2="12" />
              <line x1="12" y1="16" x2="12.01" y2="16" />
            </svg>
            {error}
          </p>
        )}
      </div>
    );
  }
);

Input.displayName = 'Input';

export { Input };
