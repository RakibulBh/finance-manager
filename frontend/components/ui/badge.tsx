import React from "react";

interface BadgeProps {
  children: React.ReactNode;
  variant?: "success" | "warning" | "danger" | "info" | "neutral";
}

export const Badge = ({ children, variant = "neutral" }: BadgeProps) => {
  const variants = {
    success: "bg-emerald-500/10 text-emerald-500 border-emerald-500/20",
    warning: "bg-amber-500/10 text-amber-500 border-amber-500/20",
    danger: "bg-red-500/10 text-red-500 border-red-500/20",
    info: "bg-blue-500/10 text-blue-500 border-blue-500/20",
    neutral: "bg-zinc-500/10 text-zinc-400 border-zinc-500/20",
  };

  return (
    <span className={`px-2 py-0.5 text-[10px] font-bold uppercase tracking-widest border rounded-md ${variants[variant]}`}>
      {children}
    </span>
  );
};
