
import React from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import { LogOut, Users, Shuffle, Settings, Home, HelpCircle } from 'lucide-react';
import { Button } from "@/components/ui/button";
import kugelblitzLogo from '../assets/logo.jpg';

const Layout = ({ children }) => {
  const { logout, user } = useAuth();
  const location = useLocation();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const navigation = [
    { name: 'Dashboard', href: '/', icon: Home },
    { name: 'Clients', href: '/clients', icon: Users },
    { name: 'Transform', href: '/transform', icon: Shuffle },
  ];

  return (
    <div className="min-h-screen bg-background">
      {/* Sidebar */}
      <div className="fixed inset-y-0 left-0 z-50 w-64 bg-card border-r">
        <div className="flex h-full flex-col">
          {/* Logo */}
          <div className="flex h-16 shrink-0 items-center px-6 border-b">
            <img src={kugelblitzLogo} alt="Kugelblitz" className="h-8 w-auto" />
          </div>
          
          {/* Navigation */}
          <nav className="flex flex-1 flex-col px-6 py-6">
            <ul className="space-y-1">
              {navigation.map((item) => {
                const isActive = location.pathname === item.href;
                const Icon = item.icon;
                
                return (
                  <li key={item.name}>
                    <Link
                      to={item.href}
                      className={`group flex gap-x-3 rounded-md p-2 text-sm leading-6 font-semibold transition-colors ${
                        isActive
                          ? 'bg-primary text-primary-foreground'
                          : 'text-muted-foreground hover:text-foreground hover:bg-muted'
                      }`}
                    >
                      <Icon className="h-5 w-5 shrink-0" />
                      {item.name}
                    </Link>
                  </li>
                );
              })}
            </ul>
            
            {/* Help Section */}
            <div className="mt-auto pt-6">
              <h3 className="text-xs font-semibold text-muted-foreground uppercase tracking-wider mb-2">
                Help & Resources
              </h3>
              <ul className="space-y-1">
                <li>
                  <a
                    href="/expression-help"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="group flex gap-x-3 rounded-md p-2 text-sm leading-6 font-semibold text-muted-foreground hover:text-foreground hover:bg-muted transition-colors"
                  >
                    <HelpCircle className="h-5 w-5 shrink-0" />
                    Expression Help
                  </a>
                </li>
              </ul>
            </div>
          </nav>
          
          {/* User info and logout */}
          <div className="border-t p-6">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <div className="h-8 w-8 rounded-full bg-primary flex items-center justify-center">
                  <span className="text-sm font-medium text-primary-foreground">
                    {user?.username?.charAt(0).toUpperCase()}
                  </span>
                </div>
              </div>
              <div className="ml-3 flex-1">
                <p className="text-sm font-medium">{user?.username}</p>
              </div>
              <Button
                variant="ghost"
                size="icon"
                onClick={handleLogout}
                title="Logout"
              >
                <LogOut className="h-4 w-4" />
              </Button>
            </div>
          </div>
        </div>
      </div>

      {/* Main content */}
      <div className="pl-64">
        <main className="py-6">
          <div className="mx-auto max-w-7xl px-6 lg:px-8">
            {children}
          </div>
        </main>
      </div>
    </div>
  );
};

export default Layout;
