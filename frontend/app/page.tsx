"use client";

import { AddAccountModal } from "@/components/forms/add-account-modal";
import { AddTransactionModal } from "@/components/forms/add-transaction-modal";
import { DashboardLayout } from "@/components/layout/dashboard-layout";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { NetWorthChart } from "@/components/ui/net-worth-chart";
import {
  ArrowUpRight,
  Plus,
  TrendingUp,
  Wallet
} from "lucide-react";
import { useState } from "react";

import { useAccounts, useNetWorth } from "@/hooks/useAccounts";

export default function DashboardPage() {
  const { data: accounts = [], isLoading: accountsLoading, refetch: refetchAccounts } = useAccounts();
  const { data: netWorth = 0, isLoading: netWorthLoading, refetch: refetchNetWorth } = useNetWorth();

  const loading = accountsLoading || netWorthLoading;

  const [isTxModalOpen, setIsTxModalOpen] = useState(false);
  const [isAccModalOpen, setIsAccModalOpen] = useState(false);

  const refreshData = () => {
    refetchAccounts();
    refetchNetWorth();
  };

  if (loading) {
    return (
      <DashboardLayout>
        <div className="space-y-8 animate-pulse">
           <div className="h-40 bg-zinc-900 rounded-3xl w-full" />
           <div className="grid grid-cols-3 gap-8">
             <div className="h-32 bg-zinc-900 rounded-3xl" />
             <div className="h-32 bg-zinc-900 rounded-3xl" />
             <div className="h-32 bg-zinc-900 rounded-3xl" />
           </div>
        </div>
      </DashboardLayout>
    );
  }

  return (
    <DashboardLayout>
      <div className="space-y-10">
        {/* Welcome & Quick Actions */}
        <div className="flex items-end justify-between">
          <div className="space-y-1">
            <h1 className="text-4xl font-bold tracking-tight text-white">Dashboard</h1>
            <p className="text-zinc-500">Your financial health at a glance</p>
          </div>
          <div className="flex gap-4">
            <Button variant="outline" className="gap-2" onClick={() => setIsAccModalOpen(true)}>
              <Plus className="w-4 h-4" />
              Connect Account
            </Button>
            <Button className="gap-2" onClick={() => setIsTxModalOpen(true)}>
              <Plus className="w-4 h-4" />
              Add Transaction
            </Button>
          </div>
        </div>

        {/* Hero Cards */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          <Card className="col-span-2 bg-linear-to-br from-brand-lime/20 to-transparent border-brand-lime/10 relative overflow-hidden group">
            <div className="absolute top-0 right-0 p-8 opacity-10 group-hover:scale-110 transition-transform">
              <TrendingUp className="w-32 h-32" />
            </div>
            <div className="relative z-10 flex flex-col md:flex-row justify-between gap-10">
              <div className="space-y-6 shrink-0">
                <div>
                  <p className="text-zinc-400 font-medium">Net Worth</p>
                  <h2 className="text-6xl font-bold mt-2">${netWorth.toLocaleString()}</h2>
                </div>
                <div className="flex items-center gap-2 text-brand-lime font-medium">
                  <ArrowUpRight className="w-4 h-4" />
                  <span>+12.5% from last month</span>
                </div>
              </div>
              <div className="flex-1 min-w-0 h-48 md:h-auto">
                <NetWorthChart
                  data={[42000, 43500, 41000, 44000, 46000, 45500, netWorth]}
                  labels={['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul']}
                />
              </div>
            </div>
          </Card>

          <Card title="Quick Wallet" subtitle="Total Liquidity" className="bg-brand-muted/50">
            <div className="space-y-6">
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 bg-zinc-800 rounded-xl flex items-center justify-center">
                    <Wallet className="w-5 h-5 text-white" />
                  </div>
                  <div>
                    <p className="text-sm text-zinc-400">Cash Balance</p>
                    <p className="font-semibold">${accounts.filter(a => a.type === 'depository').reduce((sum, a) => sum + a.balance, 0).toLocaleString()}</p>
                  </div>
                </div>
              </div>
              <Button variant="secondary" className="w-full">Manage Wallet</Button>
            </div>
          </Card>
        </div>

        {/* Recent Accounts / Activity Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
           <Card title="Recent Accounts" className="lg:col-span-2">
              <div className="space-y-4">
                {accounts.slice(0, 5).map((account) => (
                  <div key={account.id} className="flex items-center justify-between p-4 rounded-2xl bg-zinc-900/40 border border-zinc-800/50 hover:border-brand-lime/20 transition-colors cursor-pointer group">
                    <div className="flex items-center gap-4">
                      <div className="w-10 h-10 bg-zinc-800 rounded-full flex items-center justify-center text-xs font-bold text-zinc-500">
                        {account.name.substring(0, 2).toUpperCase()}
                      </div>
                      <div>
                        <p className="font-medium text-white">{account.name}</p>
                        <p className="text-xs text-zinc-500 capitalize">{account.type} • {account.subtype}</p>
                      </div>
                    </div>
                    <div className="text-right">
                      <p className="font-semibold text-white">${account.balance.toLocaleString()}</p>
                      <p className="text-[10px] text-zinc-600 uppercase tracking-wider">{account.currency}</p>
                    </div>
                  </div>
                ))}
                {accounts.length === 0 && <p className="text-zinc-500 text-center py-10">No accounts connected yet.</p>}
              </div>
           </Card>

           <Card title="Insights">
              <div className="space-y-8">
                <div className="space-y-2">
                  <div className="flex justify-between text-sm">
                    <span className="text-zinc-400">Spending Goal</span>
                    <span className="text-white">72%</span>
                  </div>
                  <div className="h-2 bg-zinc-900 rounded-full overflow-hidden">
                    <div className="h-full bg-brand-lime w-[72%]" />
                  </div>
                </div>

                <div className="p-4 rounded-2xl bg-brand-lime/5 border border-brand-lime/10">
                   <p className="text-xs text-brand-lime font-bold uppercase tracking-widest mb-1">Tip of the day</p>
                   <p className="text-sm text-zinc-300">You saved <span className="text-brand-lime font-bold">$120</span> more than last month by cutting down on entertainment.</p>
                </div>
              </div>
           </Card>
        </div>
      </div>

      <AddTransactionModal
        isOpen={isTxModalOpen}
        onClose={() => setIsTxModalOpen(false)}
        onSuccess={refreshData}
      />
      <AddAccountModal
        isOpen={isAccModalOpen}
        onClose={() => setIsAccModalOpen(false)}
        onSuccess={refreshData}
      />
    </DashboardLayout>
  );
}
