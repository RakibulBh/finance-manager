import React from "react";

interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
}

export const Input = ({ label, error, className = "", ...props }: InputProps) => {
  return (
    <div className="w-full">
      {label && (
        <label className="block text-sm font-medium text-zinc-400 mb-1.5 ml-1">
          {label}
        </label>
      )}
      <input
        className={`w-full bg-brand-muted border border-zinc-800 rounded-2xl px-4 py-3 text-white placeholder:text-zinc-600 focus:outline-none focus:ring-2 focus:ring-brand-lime/20 focus:border-brand-lime/50 transition-all ${className}`}
        {...props}
      />
      {error && <p className="text-sm text-red-500 mt-1.5 ml-1">{error}</p>}
    </div>
  );
};
