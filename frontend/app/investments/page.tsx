"use client";

import { DashboardLayout } from "@/components/layout/dashboard-layout";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import {
    Activity,
    ArrowUpRight,
    Briefcase,
    Layers,
    PieChart,
    TrendingUp,
    Zap
} from "lucide-react";

import { useInvestments } from "@/hooks/useInvestments";

export default function InvestmentsPage() {
  const { data: securities = [], isLoading: loading } = useInvestments();

  const totalValue = securities.reduce((sum, s) => sum + s.balance, 0);

  return (
    <DashboardLayout>
      <div className="space-y-10">
        <div className="flex items-end justify-between">
          <div className="space-y-1">
            <h1 className="text-4xl font-bold tracking-tight text-white">Investments</h1>
            <p className="text-zinc-500">Track your portfolio performance</p>
          </div>
          <div className="flex gap-3">
             <Button variant="outline" className="gap-2">
                <PieChart className="w-4 h-4" />
                Allocation
             </Button>
             <Button className="gap-2">
                <Zap className="w-4 h-4" />
                Trade
             </Button>
          </div>
        </div>

        {/* Stats Grid */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
           <Card className="p-6 bg-zinc-900/40 space-y-2">
              <p className="text-xs font-bold text-zinc-500 uppercase tracking-widest">Total Value</p>
              <p className="text-3xl font-bold text-white">${totalValue.toLocaleString()}</p>
              <div className="flex items-center gap-1 text-brand-lime text-xs font-bold">
                 <ArrowUpRight className="w-3 h-3" />
                 <span>+4.2%</span>
              </div>
           </Card>

           <Card className="p-6 bg-zinc-900/40 space-y-2">
              <p className="text-xs font-bold text-zinc-500 uppercase tracking-widest">Day Change</p>
              <p className="text-3xl font-bold text-emerald-500">+$842.10</p>
              <p className="text-[10px] text-zinc-600 font-bold uppercase">Updated 2m ago</p>
           </Card>

           <Card className="p-6 bg-zinc-900/40 space-y-2">
              <p className="text-xs font-bold text-zinc-500 uppercase tracking-widest">Positions</p>
              <p className="text-3xl font-bold text-white">{securities.length}</p>
              <p className="text-[10px] text-zinc-600 font-bold uppercase">Active Securities</p>
           </Card>

           <Card className="p-6 bg-zinc-900/40 space-y-2 border-brand-lime/20">
              <p className="text-xs font-bold text-zinc-500 uppercase tracking-widest">Risk Level</p>
              <p className="text-3xl font-bold text-brand-lime">Moderate</p>
              <div className="w-full bg-zinc-800 h-1 rounded-full mt-2 overflow-hidden">
                 <div className="bg-brand-lime w-2/3 h-full" />
              </div>
           </Card>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
           {/* Holdings Table */}
           <Card className="lg:col-span-2 p-0 bg-transparent overflow-hidden">
              <div className="px-8 py-6 border-b border-zinc-900 bg-zinc-900/20 flex items-center justify-between">
                 <h3 className="font-bold text-lg flex items-center gap-2">
                    <Briefcase className="w-5 h-5 text-brand-lime" />
                    Your Holdings
                 </h3>
                 <Badge variant="neutral">{securities.length} Assets</Badge>
              </div>
              <div className="overflow-x-auto">
                 <table className="w-full text-left border-collapse">
                    <thead>
                       <tr className="bg-zinc-900/10">
                          <th className="px-8 py-4 text-[10px] font-bold uppercase tracking-widest text-zinc-500 border-b border-zinc-800">Security</th>
                          <th className="px-8 py-4 text-[10px] font-bold uppercase tracking-widest text-zinc-500 border-b border-zinc-800 text-right">Price</th>
                          <th className="px-8 py-4 text-[10px] font-bold uppercase tracking-widest text-zinc-500 border-b border-zinc-800 text-right">Value</th>
                          <th className="px-8 py-4 text-[10px] font-bold uppercase tracking-widest text-zinc-500 border-b border-zinc-800 text-right">Return</th>
                       </tr>
                    </thead>
                    <tbody className="divide-y divide-zinc-900">
                       {securities.map(s => (
                         <tr key={s.id} className="hover:bg-brand-lime/[0.02] transition-colors group">
                           <td className="px-8 py-6">
                              <div className="flex items-center gap-4">
                                 <div className="w-10 h-10 rounded-xl bg-zinc-900 flex items-center justify-center font-bold text-brand-lime">
                                    {s.name.substring(0, 1)}
                                 </div>
                                 <div>
                                    <p className="font-bold text-white group-hover:text-brand-lime transition-colors">{s.name}</p>
                                    <p className="text-xs text-zinc-500 uppercase tracking-tighter">Ticker: {s.name.substring(0, 4).toUpperCase()}</p>
                                 </div>
                              </div>
                           </td>
                           <td className="px-8 py-6 text-right font-medium text-zinc-400">$142.20</td>
                           <td className="px-8 py-6 text-right font-bold text-white">${s.balance.toLocaleString()}</td>
                           <td className="px-8 py-6 text-right font-bold text-brand-lime">+12.40%</td>
                         </tr>
                       ))}
                       {securities.length === 0 && !loading && (
                         <tr>
                            <td colSpan={4} className="px-8 py-20 text-center">
                               <Layers className="w-12 h-12 text-zinc-800 mx-auto mb-4" />
                               <p className="text-zinc-500">No investment accounts found.</p>
                            </td>
                         </tr>
                       )}
                    </tbody>
                 </table>
              </div>
           </Card>

           {/* Insights/Activity */}
           <div className="space-y-8">
              <Card className="bg-zinc-900/20 p-8 space-y-6">
                 <h3 className="font-bold text-lg flex items-center gap-2">
                    <TrendingUp className="w-5 h-5 text-brand-lime" />
                    Portfolio Health
                 </h3>
                 <div className="space-y-4">
                    <div className="p-4 bg-zinc-900/50 rounded-2xl border border-zinc-800/50">
                       <p className="text-xs font-bold text-zinc-500 uppercase mb-2">Diversification</p>
                       <div className="flex items-center gap-4">
                          <div className="flex-1 bg-zinc-800 h-2 rounded-full overflow-hidden">
                             <div className="bg-brand-lime w-3/4 h-full" />
                          </div>
                          <span className="text-sm font-bold text-white">75%</span>
                       </div>
                    </div>
                    <div className="p-4 bg-zinc-900/50 rounded-2xl border border-zinc-800/50">
                       <p className="text-xs font-bold text-zinc-500 uppercase mb-2">Benchmarks</p>
                       <p className="text-sm text-zinc-300">You are outperforming the S&P 500 by <span className="text-brand-lime font-bold">2.4%</span></p>
                    </div>
                 </div>
              </Card>

              <Card className="bg-linear-to-b from-brand-lime/10 to-transparent p-8 space-y-4 border-brand-lime/20">
                 <div className="w-12 h-12 bg-brand-lime rounded-2xl flex items-center justify-center">
                    <Activity className="w-6 h-6 text-brand-dark" />
                 </div>
                 <h3 className="text-xl font-bold">AI Insight</h3>
                 <p className="text-zinc-400 text-sm leading-relaxed">
                    Based on your portfolio risk, consider rebalancing your tech holdings. NVDA has seen significant growth recently.
                 </p>
                 <Button variant="outline" className="w-full border-brand-lime/30 text-brand-lime hover:bg-brand-lime hover:text-brand-dark">
                    View Recommendation
                 </Button>
              </Card>
           </div>
        </div>
      </div>
    </DashboardLayout>
  );
}
