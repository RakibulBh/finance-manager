"use client";

import { useAuthStore } from "@/store/useAuthStore";
import { Bell, Search } from "lucide-react";

export function Header() {
  const user = useAuthStore((state) => state.user);
  const userName = user?.email ? user.email.split("@")[0] : "User";

  return (
    <header className="h-20 border-b border-zinc-900 bg-black/50 backdrop-blur-xl flex items-center justify-between px-10 sticky top-0 z-40">
      <div className="flex items-center gap-4 bg-zinc-900/50 border border-zinc-800 rounded-2xl px-4 py-2 w-96">
        <Search className="w-4 h-4 text-zinc-500" />
        <input
          placeholder="Search transactions, accounts..."
          className="bg-transparent border-none text-sm focus:outline-none w-full text-zinc-300"
        />
        <span className="text-[10px] bg-zinc-800 text-zinc-500 px-1.5 py-0.5 rounded border border-zinc-700 font-mono">⌘K</span>
      </div>

      <div className="flex items-center gap-6">
        <button
          className="relative p-2 text-zinc-400 hover:text-white transition-colors"
          title="Notifications"
        >
          <Bell className="w-5 h-5" />
          <span className="absolute top-2 right-2 w-2 h-2 bg-brand-lime rounded-full border-2 border-black" />
        </button>

        <div className="flex items-center gap-3 pl-6 border-l border-zinc-900">
          <div className="text-right">
            <p className="text-sm font-medium text-white capitalize">{userName}</p>
            <p className="text-xs text-zinc-500">Premium Member</p>
          </div>
          <div className="w-10 h-10 rounded-full bg-linear-to-br from-brand-lime to-emerald-500 border-2 border-zinc-900" />
        </div>
      </div>
    </header>
  );
}
