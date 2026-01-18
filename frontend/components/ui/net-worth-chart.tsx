"use client";


interface ChartProps {
  data: number[];
  labels: string[];
}

export const NetWorthChart = ({ data, labels }: ChartProps) => {
  const max = Math.max(...data, 1);
  const min = Math.min(...data, 0);
  const range = max - min;

  const points = data.map((val, i) => {
    const x = (i / (data.length - 1)) * 100;
    const y = 100 - ((val - min) / range) * 100;
    return `${x},${y}`;
  }).join(" ");

  return (
    <div className="w-full h-64 relative group">
      <svg
        viewBox="0 0 100 100"
        preserveAspectRatio="none"
        className="w-full h-full overflow-visible"
      >
        <defs>
          <linearGradient id="gradient" x1="0%" y1="0%" x2="0%" y2="100%">
            <stop offset="0%" stopColor="var(--color-brand-lime)" stopOpacity="0.2" />
            <stop offset="100%" stopColor="var(--color-brand-lime)" stopOpacity="0" />
          </linearGradient>
        </defs>

        {/* Area fill */}
        <polyline
          fill="url(#gradient)"
          stroke="none"
          points={`0,100 ${points} 100,100`}
          className="transition-all duration-1000 ease-in-out"
        />

        {/* Line */}
        <polyline
          fill="none"
          stroke="var(--color-brand-lime)"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
          points={points}
          className="transition-all duration-1000 ease-in-out drop-shadow-[0_0_8px_rgba(159,232,112,0.5)]"
        />
      </svg>

      {/* Grid lines & Labels */}
      <div className="absolute inset-0 flex flex-col justify-between pointer-events-none pb-2">
         {[...Array(4)].map((_, i) => (
           <div key={i} className="border-t border-zinc-900/50 w-full" />
         ))}
      </div>

      <div className="absolute bottom-0 left-0 right-0 flex justify-between px-2 pt-4 border-t border-zinc-800/30">
         {labels.map((l, i) => (
           <span key={i} className="text-[10px] font-bold text-zinc-600 uppercase tracking-tighter">{l}</span>
         ))}
      </div>
    </div>
  );
};

// CSS variable references in SVG (Tailwind compatibility)
// Note: In global.css I defined --brand-lime, I'll use that.
