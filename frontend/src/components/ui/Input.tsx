"use client";

/**
 * Text input component with label and error display.
 * Follows the Loupi design system (soft & warm style).
 */

import { type InputHTMLAttributes, forwardRef } from "react";

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
}

export const Input = forwardRef<HTMLInputElement, InputProps>(
  ({ label, error, className = "", id, ...props }, ref) => {
    return (
      <div className="flex flex-col gap-1.5">
        {label ? (
          <label
            htmlFor={id}
            className="font-heading text-sm font-medium text-foreground"
          >
            {label}
          </label>
        ) : null}
        <input
          ref={ref}
          id={id}
          className={`
            w-full rounded-[--radius-md] border border-border
            bg-surface px-4 py-3
            font-body text-sm text-foreground
            placeholder:text-foreground-secondary
            outline-none transition-all duration-200
            focus:border-primary focus:ring-2 focus:ring-primary-light
            disabled:opacity-50 disabled:cursor-not-allowed
            ${error ? "border-danger focus:border-danger focus:ring-danger/20" : ""}
            ${className}
          `}
          {...props}
        />
        {error ? (
          <p className="text-xs text-danger">{error}</p>
        ) : null}
      </div>
    );
  },
);

Input.displayName = "Input";
