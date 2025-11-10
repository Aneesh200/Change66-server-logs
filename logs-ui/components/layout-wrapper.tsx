"use client";

import { usePathname } from "next/navigation";
import { Sidebar } from "@/components/sidebar";

export function LayoutWrapper({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const isLoginPage = pathname === "/login";

  if (isLoginPage) {
    return <>{children}</>;
  }

  return (
    <div className="flex h-screen overflow-hidden bg-black">
      <Sidebar />
      <main className="flex-1 overflow-y-auto lg:ml-64 bg-black">
        {children}
      </main>
    </div>
  );
}

