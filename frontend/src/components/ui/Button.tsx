import { forwardRef } from 'react';
import { cn } from '@/lib/utils';
import type { ButtonProps } from '@/types';

const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  ({ 
    variant = 'primary', 
    size = 'md', 
    loading = false, 
    disabled = false, 
    children, 
    className, 
    ...props 
  }, ref) => {
    const baseClasses = cn(
      'inline-flex items-center justify-center gap-2 font-semibold',
      'rounded-xl transition-all duration-200',
      'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-500/50 focus-visible:ring-offset-2',
      'disabled:opacity-50 disabled:cursor-not-allowed disabled:active:scale-100'
    );
    
    const variantClasses = {
      primary: cn(
        'bg-primary-600 text-white',
        'hover:bg-primary-700 active:bg-primary-800',
        'shadow-soft hover:shadow-soft-lg',
        'active:scale-[0.98]'
      ),
      secondary: cn(
        'bg-surface-100 dark:bg-surface-800 text-surface-700 dark:text-surface-200',
        'border border-surface-200 dark:border-surface-700',
        'hover:bg-surface-200 dark:hover:bg-surface-700',
        'active:bg-surface-300 dark:active:bg-surface-600',
        'active:scale-[0.98]'
      ),
      ghost: cn(
        'text-surface-600 dark:text-surface-400',
        'hover:bg-surface-100 dark:hover:bg-surface-800',
        'active:scale-[0.98]'
      ),
      danger: cn(
        'bg-error-500 text-white',
        'hover:bg-error-600 active:bg-error-600',
        'shadow-soft hover:shadow-soft-lg',
        'active:scale-[0.98]'
      ),
    };
    
    const sizeClasses = {
      sm: 'h-9 px-4 text-sm',
      md: 'h-11 px-5 text-sm',
      lg: 'h-12 px-6 text-base',
      xl: 'h-14 px-8 text-base',
    };

    return (
      <button
        ref={ref}
        className={cn(
          baseClasses,
          variantClasses[variant],
          sizeClasses[size],
          className
        )}
        disabled={disabled || loading}
        {...props}
      >
        {loading && (
          <span className="loading-spinner w-4 h-4" />
        )}
        {children}
      </button>
    );
  }
);

Button.displayName = 'Button';

export { Button };
