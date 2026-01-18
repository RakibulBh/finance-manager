"use client";

import { useAuthStore } from "@/store/useAuthStore";
import {
    BarChart3,
    CreditCard,
    History,
    Home,
    LogOut,
    Settings,
    Users
} from "lucide-react";
import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";

const navItems = [
  { icon: Home, label: "Overview", href: "/" },
  { icon: CreditCard, label: "Accounts", href: "/accounts" },
  { icon: History, label: "Transactions", href: "/transactions" },
  { icon: BarChart3, label: "Investments", href: "/investments" },
  { icon: Users, label: "Family", href: "/family" },
];

export function Sidebar() {
  const pathname = usePathname();
  const router = useRouter();
  const logout = useAuthStore((state) => state.logout);

  const handleLogout = () => {
    logout();
    router.push("/login");
  };

  return (
    <div className="w-72 h-screen border-r border-zinc-900 bg-black flex flex-col p-6 fixed left-0 top-0 overflow-y-auto">
      <div className="flex items-center gap-3 px-2 mb-10">
        <div className="w-10 h-10 bg-brand-lime rounded-xl flex items-center justify-center">
          <div className="w-6 h-6 bg-black rounded transform rotate-12" />
        </div>
        <span className="text-xl font-bold tracking-tight">maybe-clone</span>
      </div>

      <nav className="flex-1 space-y-2">
        {navItems.map((item) => {
          const isActive = pathname === item.href;
          return (
            <Link
              key={item.href}
              href={item.href}
              className={`flex items-center gap-3 px-4 py-3 rounded-2xl transition-all duration-200 group ${
                isActive
                  ? "bg-brand-lime text-black font-semibold"
                  : "text-zinc-500 hover:text-white hover:bg-zinc-900"
              }`}
            >
              <item.icon className={`w-5 h-5 ${isActive ? "text-black" : "group-hover:text-white"}`} />
              <span>{item.label}</span>
            </Link>
          );
        })}
      </nav>

      <div className="pt-6 border-t border-zinc-900 space-y-2">
        <Link href="/settings" className="flex items-center gap-3 px-4 py-3 text-zinc-500 hover:text-white w-full rounded-2xl hover:bg-zinc-900 transition-all group">
          <Settings className="w-5 h-5 group-hover:text-white" />
          <span>Settings</span>
        </Link>
        <button
          onClick={handleLogout}
          className="flex items-center gap-3 px-4 py-3 text-zinc-500 hover:text-red-400 w-full rounded-2xl hover:bg-red-950/10 transition-all group"
        >
          <LogOut className="w-5 h-5 group-hover:text-red-400" />
          <span>Log Out</span>
        </button>
      </div>
    </div>
  );
}
