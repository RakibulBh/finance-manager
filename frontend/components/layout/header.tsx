"use client";

import { useAuthStore } from "@/store/useAuthStore";
import { Bell, ChevronDown, LogOut, Search, User } from "lucide-react";
import { useRouter } from "next/navigation";
import { useState } from "react";

export function Header() {
  const user = useAuthStore((state) => state.user);
  const logout = useAuthStore((state) => state.logout);
  const router = useRouter();
  const [isOpen, setIsOpen] = useState(false);

  // Derive display name from email if name property is missing in user model for now
  const displayName = user?.email?.split("@")[0] || "User";

  const handleLogout = () => {
    logout();
    router.push("/login");
  };

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

        <div className="relative">
          <button
            onClick={() => setIsOpen(!isOpen)}
            className="flex items-center gap-3 pl-6 border-l border-zinc-900 focus:outline-none group"
          >
            <div className="text-right hidden md:block">
              <p className="text-sm font-medium text-white capitalize group-hover:text-brand-lime transition-colors">
                {displayName}
              </p>
              <p className="text-xs text-zinc-500">Free Plan</p>
            </div>
            <div className="w-10 h-10 rounded-full bg-linear-to-br from-brand-lime to-emerald-500 border-2 border-zinc-900 flex items-center justify-center text-black font-bold text-lg">
              {displayName[0]?.toUpperCase()}
            </div>
            <ChevronDown className={`w-4 h-4 text-zinc-500 transition-transform duration-200 ${isOpen ? "rotate-180" : ""}`} />
          </button>

          {isOpen && (
            <>
              <div
                className="fixed inset-0 z-40"
                onClick={() => setIsOpen(false)}
              />
              <div className="absolute right-0 top-full mt-2 w-56 bg-zinc-900 border border-zinc-800 rounded-xl shadow-xl z-50 overflow-hidden py-1 animate-in fade-in zoom-in-95 duration-100">
                <div className="px-4 py-3 border-b border-zinc-800">
                  <p className="text-sm text-white font-medium">Signed in as</p>
                  <p className="text-xs text-zinc-400 truncate">{user?.email}</p>
                </div>

                <div className="p-1">
                  <button
                    onClick={() => router.push("/settings")}
                    className="flex w-full items-center gap-2 px-3 py-2 text-sm text-zinc-300 hover:text-white hover:bg-zinc-800 rounded-lg transition-colors"
                  >
                    <User className="w-4 h-4" />
                    Profile
                  </button>
                  <div className="h-px bg-zinc-800 my-1 mx-2" />
                  <button
                    onClick={handleLogout}
                    className="flex w-full items-center gap-2 px-3 py-2 text-sm text-red-400 hover:text-red-300 hover:bg-red-950/30 rounded-lg transition-colors"
                  >
                    <LogOut className="w-4 h-4" />
                    Log out
                  </button>
                </div>
              </div>
            </>
          )}
        </div>
      </div>
    </header>
  );
}
