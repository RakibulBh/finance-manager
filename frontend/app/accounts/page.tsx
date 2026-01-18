"use client";

import { DashboardLayout } from "@/components/layout/dashboard-layout";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import {
    ArrowRight,
    Briefcase,
    CreditCard,
    Filter,
    Landmark,
    Plus,
    Search
} from "lucide-react";

import { useAccounts } from "@/hooks/useAccounts";

export default function AccountsPage() {
  const { data: accounts = [], isLoading: loading } = useAccounts();

  const getIcon = (type: string) => {
    switch (type) {
      case 'depository': return Landmark;
      case 'credit': return CreditCard;
      case 'investment': return Briefcase;
      default: return Landmark;
    }
  };

  return (
    <DashboardLayout>
      <div className="space-y-10">
        <div className="flex items-end justify-between">
          <div className="space-y-1">
            <h1 className="text-4xl font-bold tracking-tight text-white">Accounts</h1>
            <p className="text-zinc-500">Manage your assets and liabilities</p>
          </div>
          <Button className="gap-2">
            <Plus className="w-4 h-4" />
            Connect Account
          </Button>
        </div>

        <div className="flex items-center gap-4">
           <div className="flex-1 bg-zinc-900/50 border border-zinc-900 rounded-2xl px-6 py-3 flex items-center gap-4 focus-within:ring-2 focus-within:ring-brand-lime/20 transition-all">
              <Search className="w-5 h-5 text-zinc-500" />
              <input placeholder="Search by name, type or institution..." className="bg-transparent border-none text-sm w-full focus:outline-none" />
           </div>
           <Button variant="outline" className="gap-2 h-12">
             <Filter className="w-4 h-4" />
             Filters
           </Button>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
           {['depository', 'credit', 'investment', 'loan'].map(type => {
             const typeAccounts = accounts.filter(a => a.type === type);
             if (typeAccounts.length === 0) return null;

             const total = typeAccounts.reduce((sum, a) => sum + a.balance, 0);
             const Icon = getIcon(type);

             return (
               <Card key={type} className="bg-zinc-900/20">
                 <div className="flex items-center justify-between mb-8">
                    <div className="flex items-center gap-3">
                      <div className="w-12 h-12 bg-zinc-800 rounded-2xl flex items-center justify-center">
                        <Icon className="w-6 h-6 text-brand-lime" />
                      </div>
                      <div>
                        <h2 className="text-xl font-bold capitalize">{type}</h2>
                        <p className="text-xs text-zinc-500">{typeAccounts.length} Connected</p>
                      </div>
                    </div>
                    <div className="text-right">
                       <p className="text-2xl font-bold">${Math.abs(total).toLocaleString()}</p>
                       <p className={`text-[10px] font-bold uppercase tracking-widest ${total >= 0 ? 'text-emerald-500' : 'text-red-500'}`}>
                         {total >= 0 ? 'Asset' : 'Liability'}
                       </p>
                    </div>
                 </div>

                 <div className="space-y-3">
                    {typeAccounts.map(account => (
                      <div key={account.id} className="flex items-center justify-between p-4 rounded-xl bg-black/40 border border-zinc-800/50 hover:border-brand-lime/30 transition-all group cursor-pointer">
                        <div className="flex items-center gap-4">
                           <div className="w-1 h-8 bg-zinc-800 rounded-full group-hover:bg-brand-lime transition-colors" />
                           <div>
                             <p className="font-medium text-zinc-200 group-hover:text-white transition-colors">{account.name}</p>
                             <p className="text-xs text-zinc-500">{account.institution_name || 'Manual'}</p>
                           </div>
                        </div>
                        <div className="flex items-center gap-6">
                           <div className="text-right">
                              <p className="font-semibold text-white tracking-tight">${account.balance.toLocaleString()}</p>
                              <Badge variant={account.balance >= 0 ? "success" : "danger"}>
                                {account.currency}
                              </Badge>
                           </div>
                           <ArrowRight className="w-4 h-4 text-zinc-700 group-hover:text-brand-lime transition-colors" />
                        </div>
                      </div>
                    ))}
                 </div>
               </Card>
             );
           })}
        </div>

        {accounts.length === 0 && !loading && (
          <div className="py-20 text-center space-y-4">
             <div className="w-20 h-20 bg-zinc-900 rounded-full mx-auto flex items-center justify-center">
                <Landmark className="w-10 h-10 text-zinc-700" />
             </div>
             <div>
                <h3 className="text-xl font-bold">No accounts found</h3>
                <p className="text-zinc-500">Connect your bank or add an account manually to get started.</p>
             </div>
             <Button className="mt-4">Connect First Account</Button>
          </div>
        )}
      </div>
    </DashboardLayout>
  );
}
