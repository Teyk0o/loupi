"use client";

/**
 * Primary button component with variant support.
 * Follows the Loupi design system (soft & warm style).
 */

import { type ButtonHTMLAttributes, forwardRef } from "react";

type ButtonVariant = "primary" | "secondary" | "outline" | "danger" | "ghost";

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: ButtonVariant;
  isLoading?: boolean;
}

const variantStyles: Record<ButtonVariant, string> = {
  primary:
    "bg-primary text-white hover:bg-primary-hover active:scale-[0.98]",
  secondary:
    "bg-secondary text-foreground hover:bg-secondary-hover active:scale-[0.98]",
  outline:
    "border border-border bg-transparent text-foreground hover:bg-surface active:scale-[0.98]",
  danger:
    "bg-danger text-white hover:bg-danger-hover active:scale-[0.98]",
  ghost:
    "bg-transparent text-foreground-secondary hover:bg-surface active:scale-[0.98]",
};

export const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  ({ variant = "primary", isLoading, children, className = "", disabled, ...props }, ref) => {
    return (
      <button
        ref={ref}
        disabled={disabled || isLoading}
        className={`
          inline-flex items-center justify-center gap-2
          rounded-[--radius-md] px-5 py-3
          font-heading text-sm font-semibold
          transition-all duration-200
          disabled:opacity-50 disabled:cursor-not-allowed
          ${variantStyles[variant]}
          ${className}
        `}
        {...props}
      >
        {isLoading ? (
          <span className="h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent" />
        ) : null}
        {children}
      </button>
    );
  },
);

Button.displayName = "Button";
