'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { signOut, useSession } from 'next-auth/react';
import { 
  LayoutDashboard, 
  Users, 
  Activity, 
  Eye, 
  AlertTriangle, 
  Zap,
  Database,
  Settings,
  Menu,
  X,
  LogOut
} from 'lucide-react';
import { cn } from '@/lib/utils';
import { Button } from './ui/button';
import { useState } from 'react';

interface NavItem {
  name: string;
  href: string;
  icon: React.ElementType;
  badge?: string;
}

const navItems: NavItem[] = [
  { name: 'Dashboard', href: '/', icon: LayoutDashboard },
  { name: 'All Logs', href: '/logs', icon: Database },
  { name: 'Behavioral', href: '/logs/behavioral', icon: Users },
  { name: 'Telemetry', href: '/logs/telemetry', icon: Activity },
  { name: 'Observability', href: '/logs/observability', icon: Eye },
  { name: 'Errors', href: '/logs/errors', icon: AlertTriangle },
  { name: 'Performance', href: '/logs/performance', icon: Zap },
];

export function Sidebar() {
  const pathname = usePathname();
  const { data: session } = useSession();
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  return (
    <>
      {/* Mobile menu button */}
      <Button
        variant="ghost"
        size="icon"
        className="fixed top-4 left-4 z-50 lg:hidden"
        onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
      >
        {isMobileMenuOpen ? <X className="h-6 w-6" /> : <Menu className="h-6 w-6" />}
      </Button>

      {/* Sidebar */}
      <aside
        className={cn(
          'fixed left-0 top-0 z-40 h-screen w-64 border-r border-zinc-800/50 bg-zinc-950/95 backdrop-blur-sm transition-transform duration-300',
          isMobileMenuOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'
        )}
      >
        <div className="flex h-full flex-col">
          {/* Logo/Header */}
          <div className="flex h-16 items-center border-b border-zinc-800/50 px-6 bg-black/50">
            <Activity className="h-6 w-6 text-zinc-500" />
            <h1 className="ml-3 text-lg font-bold text-gray-100">Analytics Hub</h1>
          </div>

          {/* Navigation */}
          <nav className="flex-1 space-y-1 overflow-y-auto px-3 py-4">
            {navItems.map((item) => {
              const Icon = item.icon;
              const isActive = pathname === item.href;

              return (
                <Link
                  key={item.href}
                  href={item.href}
                  onClick={() => setIsMobileMenuOpen(false)}
                  className={cn(
                    'flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-all duration-200',
                    isActive
                      ? 'bg-zinc-600 text-white shadow-lg shadow-zinc-600/30'
                      : 'text-gray-400 hover:bg-zinc-800/50 hover:text-gray-100'
                  )}
                >
                  <Icon className="h-5 w-5" />
                  <span>{item.name}</span>
                  {item.badge && (
                    <span className="ml-auto rounded-full bg-primary/10 px-2 py-0.5 text-xs">
                      {item.badge}
                    </span>
                  )}
                </Link>
              );
            })}
          </nav>

          {/* Footer */}
          <div className="border-t border-zinc-800/50 p-4 bg-black/50 space-y-3">
            <div className="rounded-lg bg-zinc-900/50 p-3 border border-zinc-800/50">
              <p className="text-xs font-medium text-gray-100">Log Server</p>
              <p className="text-xs text-gray-500">
                Connected & Active
              </p>
              <div className="mt-2 flex items-center gap-2">
                <div className="h-2 w-2 rounded-full bg-green-500 animate-pulse shadow-lg shadow-green-500/50" />
                <span className="text-xs text-gray-500">Online</span>
              </div>
            </div>
            
            {session && (
              <div className="space-y-2">
                <div className="rounded-lg bg-zinc-900/50 p-2 border border-zinc-800/50">
                  <p className="text-xs text-gray-500">Signed in as</p>
                  <p className="text-xs font-medium text-gray-100 truncate">{session.user?.name}</p>
                </div>
                <Button
                  onClick={() => signOut({ callbackUrl: '/login' })}
                  variant="outline"
                  size="sm"
                  className="w-full bg-zinc-900/50 border-zinc-800/50 hover:bg-zinc-800 text-gray-300 hover:text-gray-100"
                >
                  <LogOut className="h-4 w-4 mr-2" />
                  Sign Out
                </Button>
              </div>
            )}
          </div>
        </div>
      </aside>

      {/* Overlay for mobile */}
      {isMobileMenuOpen && (
        <div
          className="fixed inset-0 z-30 bg-black/50 lg:hidden"
          onClick={() => setIsMobileMenuOpen(false)}
        />
      )}
    </>
  );
}

