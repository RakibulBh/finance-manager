import React from "react";

interface CardProps {
  children: React.ReactNode;
  className?: string;
  title?: string;
  subtitle?: string;
  accent?: boolean;
}

export const Card = ({ children, className = "", title, subtitle, accent }: CardProps) => {
  return (
    <div className={`glass rounded-3xl p-6 ${accent ? 'ring-1 ring-brand-lime/20' : ''} ${className}`}>
      {(title || subtitle) && (
        <div className="mb-6">
          {title && <h3 className="text-lg font-semibold text-white">{title}</h3>}
          {subtitle && <p className="text-sm text-zinc-400 mt-1">{subtitle}</p>}
        </div>
      )}
      {children}
    </div>
  );
};
