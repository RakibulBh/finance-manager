"use client";

import { DashboardLayout } from "@/components/layout/dashboard-layout";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { useTransactions } from "@/hooks/useTransactions";
import {
    ArrowDownLeft,
    ArrowUpRight,
    Calendar,
    Download,
    Filter,
    History as HistoryIcon,
    Search,
    Store,
    Tag
} from "lucide-react";

export default function TransactionsPage() {
  const { data: transactions = [], isLoading: loading } = useTransactions();

  return (
    <DashboardLayout>
      <div className="space-y-10">
        <div className="flex items-end justify-between">
          <div className="space-y-1">
            <h1 className="text-4xl font-bold tracking-tight text-white">Transactions</h1>
            <p className="text-zinc-500">A complete ledger of your activities</p>
          </div>
          <Button variant="outline" className="gap-2">
            <Download className="w-4 h-4" />
            Export CSV
          </Button>
        </div>

        {/* Filters Header */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
           <div className="md:col-span-2 bg-zinc-900/50 border border-zinc-900 rounded-2xl px-6 py-3 flex items-center gap-4 focus-within:ring-2 focus-within:ring-brand-lime/20 transition-all">
              <Search className="w-5 h-5 text-zinc-500" />
              <input placeholder="Search transactions..." className="bg-transparent border-none text-sm w-full focus:outline-none" />
           </div>

           <div className="bg-zinc-900/50 border border-zinc-900 rounded-2xl px-4 py-3 flex items-center gap-3 cursor-pointer hover:bg-zinc-900 transition-colors group">
              <Calendar className="w-4 h-4 text-zinc-500 group-hover:text-brand-lime transition-colors" />
              <span className="text-sm text-zinc-300">Last 30 Days</span>
           </div>

           <Button variant="outline" className="gap-2 h-full rounded-2xl">
             <Filter className="w-4 h-4" />
             Advanced Filters
           </Button>
        </div>

        {/* Transactions Table/List */}
        <Card className="p-0 overflow-hidden border-zinc-900/50 bg-zinc-900/10">
           <div className="overflow-x-auto">
             <table className="w-full text-left border-collapse">
               <thead>
                 <tr className="border-b border-zinc-900 bg-zinc-900/20">
                   <th className="px-8 py-5 text-xs font-bold uppercase tracking-widest text-zinc-500">Date</th>
                   <th className="px-8 py-5 text-xs font-bold uppercase tracking-widest text-zinc-500">Entity</th>
                   <th className="px-8 py-5 text-xs font-bold uppercase tracking-widest text-zinc-500">Category</th>
                   <th className="px-8 py-5 text-xs font-bold uppercase tracking-widest text-zinc-500">Account</th>
                   <th className="px-8 py-5 text-xs font-bold uppercase tracking-widest text-zinc-500 text-right">Amount</th>
                 </tr>
               </thead>
               <tbody className="divide-y divide-zinc-900/50">
                 {transactions.map((tx) => (
                   <tr
                key={tx.id}
                className="group border-b border-zinc-900/50 hover:bg-white/2 transition-colors"
              >
                     <td className="px-8 py-6">
                        <p className="text-sm font-medium text-zinc-400 group-hover:text-zinc-300 transition-colors">
                          {new Date(tx.date).toLocaleDateString(undefined, { month: 'short', day: 'numeric', year: 'numeric' })}
                        </p>
                     </td>
                     <td className="px-8 py-6">
                        <div className="flex items-center gap-4">
                           <div className="w-10 h-10 bg-zinc-900 border border-zinc-800 rounded-xl flex items-center justify-center">
                              {tx.merchant_name ? <Store className="w-5 h-5 text-zinc-500" /> : <Tag className="w-5 h-5 text-zinc-500" />}
                           </div>
                           <div>
                              <p className="font-semibold text-white">{tx.name}</p>
                              {tx.merchant_name && <p className="text-xs text-zinc-500">{tx.merchant_name}</p>}
                           </div>
                        </div>
                     </td>
                     <td className="px-8 py-6">
                        <Badge variant="neutral">{tx.category_name || "Uncategorized"}</Badge>
                     </td>
                     <td className="px-8 py-6 text-sm text-zinc-400">
                        {tx.account_name}
                     </td>
                     <td className="px-8 py-6 text-right">
                        <div className="flex items-center justify-end gap-2">
                           <span className={`text-lg font-bold tracking-tight ${tx.amount < 0 ? 'text-white' : 'text-brand-lime'}`}>
                             {tx.amount < 0 ? '-' : '+'}${Math.abs(tx.amount).toLocaleString()}
                           </span>
                           {tx.amount < 0 ? (
                             <ArrowDownLeft className="w-4 h-4 text-red-500/50" />
                           ) : (
                             <ArrowUpRight className="w-4 h-4 text-brand-lime/50" />
                           )}
                        </div>
                     </td>
                   </tr>
                 ))}
               </tbody>
             </table>

             {transactions.length === 0 && !loading && (
               <div className="py-32 text-center space-y-4 bg-black/20">
                  <div className="w-16 h-16 bg-zinc-900/50 rounded-2xl mx-auto flex items-center justify-center">
                    <HistoryIcon className="w-8 h-8 text-zinc-700" />
                  </div>
                  <div className="max-w-xs mx-auto">
                    <h3 className="text-xl font-bold">No transactions found</h3>
                    <p className="text-zinc-500 mt-1">Start by adding a manual transaction or connecting your bank.</p>
                  </div>
                  <Button className="mt-4">Add Transaction</Button>
               </div>
             )}

             {loading && (
               <div className="p-8 space-y-4">
                  {[...Array(5)].map((_, i) => (
                    <div key={i} className="h-16 bg-zinc-900/20 animate-pulse rounded-2xl w-full" />
                  ))}
               </div>
             )}
           </div>
        </Card>
      </div>
    </DashboardLayout>
  );
}
