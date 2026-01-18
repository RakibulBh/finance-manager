"use client";

import { useAuthStore } from "@/store/useAuthStore";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { Header } from "./header";
import { Sidebar } from "./sidebar";

export function DashboardLayout({ children }: { children: React.ReactNode }) {
  const { isAuthenticated } = useAuthStore();
  const router = useRouter();

  useEffect(() => {
    // If we've determined the user is not authenticated after store hydration
    if (!isAuthenticated) {
      const token = localStorage.getItem("token");
      if (!token) {
        router.push("/login");
      }
    }
  }, [isAuthenticated, router]);

  // We can't easily wait for "hydration" without a hydration state,
  // but we can check if token exists as a fallback for the first render.
  const hasToken = typeof window !== 'undefined' ? !!localStorage.getItem("token") : false;

  if (!isAuthenticated && !hasToken) {
    return (
      <div className="min-h-screen bg-black flex items-center justify-center">
        <div className="w-8 h-8 border-4 border-brand-lime border-t-transparent rounded-full animate-spin" />
      </div>
    );
  }

  return (
    <div className="flex bg-black min-h-screen">
      <Sidebar />
      <div className="flex-1 ml-72">
        <Header />
        <main className="p-10 max-w-7xl mx-auto">
          {children}
        </main>
      </div>
    </div>
  );
}
