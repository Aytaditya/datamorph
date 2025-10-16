import React, { useState, useEffect } from 'react';
import { Users, Shuffle, CheckCircle, AlertCircle } from 'lucide-react';
import { clientsAPI } from '../services/api';
import toast from 'react-hot-toast';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";


const Dashboard = () => {
  const [stats, setStats] = useState({
    totalClients: 0,
    totalTransformations: 0,
    systemStatus: 'healthy'
  });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadStats();
  }, []);

  const loadStats = async () => {
    try {
      const clients = await clientsAPI.getAll();
      setStats(prev => ({
        ...prev,
        totalClients: clients.length
      }));
    } catch (error) {
      toast.error('Failed to load dashboard stats');
    } finally {
      setLoading(false);
    }
  };

  const statCards = [
    {
      name: 'Total Clients',
      value: stats.totalClients,
      icon: Users,
      change: '+2.5%',
      changeType: 'increase'
    },
    {
      name: 'Transformations',
      value: stats.totalTransformations,
      icon: Shuffle,
      change: '+54.02%',
      changeType: 'increase'
    },
    {
      name: 'System Status',
      value: stats.systemStatus,
      icon: stats.systemStatus === 'healthy' ? CheckCircle : AlertCircle,
      change: 'Operational',
      changeType: 'neutral'
    }
  ];

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="border-b pb-5">
        <h1 className="text-3xl font-bold tracking-tight">Dashboard</h1>
        <p className="text-muted-foreground">
          Welcome to the Data Mapping API dashboard. Monitor your transformations and manage clients.
        </p>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3">
        {statCards.map((stat) => {
          const Icon = stat.icon;
          return (
            <Card key={stat.name}>
              <CardContent className="p-6">
                <div className="flex items-center">
                  <div className="flex-shrink-0">
                    <Icon 
                      className={`h-8 w-8 ${
                        stat.name === 'System Status' 
                          ? stats.systemStatus === 'healthy' 
                            ? 'text-green-600' 
                            : 'text-red-600'
                          : 'text-primary'
                      }`} 
                    />
                  </div>
                  <div className="ml-5 w-0 flex-1">
                    <dl>
                      <dt className="text-sm font-medium text-muted-foreground truncate">
                        {stat.name}
                      </dt>
                      <dd className="flex items-baseline">
                        <div className="text-2xl font-semibold">
                          {stat.value}
                        </div>
                        <div className={`ml-2 flex items-baseline text-sm font-semibold ${
                          stat.changeType === 'increase' 
                            ? 'text-green-600' 
                            : stat.changeType === 'decrease'
                            ? 'text-red-600'
                            : 'text-muted-foreground'
                        }`}>
                          {stat.change}
                        </div>
                      </dd>
                    </dl>
                  </div>
                </div>
              </CardContent>
            </Card>
          );
        })}
      </div>

      {/* Quick Actions */}
      <div className="grid grid-cols-1 gap-5 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Quick Actions</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <a
              href="/clients"
              className="flex items-center p-3 rounded-lg border hover:border-primary/50 hover:bg-muted/50 transition-colors duration-200"
            >
              <Users className="h-5 w-5 text-primary mr-3" />
              <div>
                <p className="text-sm font-medium">Manage Clients</p>
                <p className="text-sm text-muted-foreground">Create and manage API clients</p>
              </div>
            </a>
            <a
              href="/transform"
              className="flex items-center p-3 rounded-lg border hover:border-primary/50 hover:bg-muted/50 transition-colors duration-200"
            >
              <Shuffle className="h-5 w-5 text-primary mr-3" />
              <div>
                <p className="text-sm font-medium">Test Transformations</p>
                <p className="text-sm text-muted-foreground">Test JSON data transformations</p>
              </div>
            </a>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Recent Activity</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-sm text-muted-foreground">
              No recent activity to display.
            </div>
          </CardContent>
        </Card>
      </div>

      {/* API Status */}
      <Card>
        <CardContent className="p-6">
          <div className="flex items-center justify-between">
            <div>
              <h3 className="text-lg font-medium">API Status</h3>
              <p className="text-sm text-muted-foreground">Backend API connectivity and health</p>
            </div>
            <div className="flex items-center">
              <CheckCircle className="h-5 w-5 text-green-500 mr-2" />
              <span className="text-sm font-medium text-green-600">Connected</span>
            </div>
          </div>
          <div className="mt-4 grid grid-cols-1 gap-4 sm:grid-cols-3">
            <div>
              <dt className="text-sm font-medium text-muted-foreground">Endpoint</dt>
              <dd className="mt-1 text-sm">https://localhost:8080</dd>
            </div>
            <div>
              <dt className="text-sm font-medium text-muted-foreground">Version</dt>
              <dd className="mt-1 text-sm">1.0.0</dd>
            </div>
            <div>
              <dt className="text-sm font-medium text-muted-foreground">Environment</dt>
              <dd className="mt-1 text-sm">Development</dd>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
};

export default Dashboard;
