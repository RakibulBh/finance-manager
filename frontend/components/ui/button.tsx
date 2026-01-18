import React from "react";

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: "primary" | "secondary" | "outline" | "ghost";
  size?: "sm" | "md" | "lg";
}

export const Button = ({
  children,
  className = "",
  variant = "primary",
  size = "md",
  ...props
}: ButtonProps) => {
  const baseStyles = "inline-flex items-center justify-center rounded-full font-medium transition-all active:scale-95 disabled:opacity-50 disabled:pointer-events-none";

  const variants = {
    primary: "bg-brand-lime text-black hover:brightness-110",
    secondary: "bg-brand-muted text-white hover:bg-zinc-800",
    outline: "border border-zinc-700 text-white hover:bg-zinc-900",
    ghost: "text-zinc-400 hover:text-white hover:bg-white/5",
  };

  const sizes = {
    sm: "px-4 py-1.5 text-sm",
    md: "px-6 py-2.5 text-base",
    lg: "px-8 py-3.5 text-lg",
  };

  const combinedClassName = `${baseStyles} ${variants[variant]} ${sizes[size]} ${className}`;

  return (
    <button className={combinedClassName} {...props}>
      {children}
    </button>
  );
};
