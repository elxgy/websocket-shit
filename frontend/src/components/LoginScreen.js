import React, { useState, useEffect } from 'react';
import { MessageCircle, LogIn, Loader2, Users } from 'lucide-react';
import config from '../utils/config';
import api from '../utils/api';

const LoginScreen = ({ onLogin }) => {
  const [formData, setFormData] = useState({ username: '', password: '' });
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const [serverStatus, setServerStatus] = useState(null);

  // Check server health on mount
  useEffect(() => {
    const checkHealth = async () => {
      try {
        const health = await api.healthCheck();
        setServerStatus(health);
      } catch (error) {
        console.error('Health check failed:', error);
        setServerStatus({ status: 'error', message: error.message });
      }
    };
    
    checkHealth();
  }, []);

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!formData.username.trim() || !formData.password) {
      setError('Please enter both username and password');
      return;
    }

    setIsLoading(true);
    setError('');

    try {
      await onLogin(formData.username.trim(), formData.password);
    } catch (error) {
      console.error('Login error:', error);
      setError(error.message || 'Login failed. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  const handleTestUser = (user) => {
    setFormData({ username: user.username, password: user.password });
    setError('');
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
    if (error) setError('');
  };

  return (
    <div className="min-h-screen flex items-center justify-center p-4">
      <div className="bg-white rounded-2xl shadow-2xl w-full max-w-md p-8">
        {/* Header */}
        <div className="text-center mb-8">
          <div className="flex justify-center mb-4">
            <div className="bg-gradient-to-br from-indigo-500 to-purple-600 p-3 rounded-full">
              <MessageCircle className="w-8 h-8 text-white" />
            </div>
          </div>
          <h1 className="text-3xl font-bold text-gray-900 mb-2">WebSocket Chat</h1>
          <p className="text-gray-600">Connect with up to 4 users in real-time</p>
        </div>

        {/* Server Status */}
        {serverStatus && (
          <div className={`mb-6 p-3 rounded-lg text-sm ${
            serverStatus.status === 'ok' 
              ? 'bg-green-50 text-green-700 border border-green-200' 
              : 'bg-red-50 text-red-700 border border-red-200'
          }`}>
            <div className="flex items-center gap-2">
              <div className={`w-2 h-2 rounded-full ${
                serverStatus.status === 'ok' ? 'bg-green-400 animate-pulse' : 'bg-red-400'
              }`} />
              {serverStatus.status === 'ok' ? (
                <span>Server Online • {serverStatus.clients || 0}/4 users connected</span>
              ) : (
                <span>Server Offline • {serverStatus.message}</span>
              )}
            </div>
          </div>
        )}

        {/* Login Form */}
        <form onSubmit={handleSubmit} className="space-y-6">
          <div>
            <label htmlFor="username" className="block text-sm font-medium text-gray-700 mb-2">
              Username
            </label>
            <input
              type="text"
              id="username"
              name="username"
              value={formData.username}
              onChange={handleChange}
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-colors"
              placeholder="Enter your username"
              maxLength={50}
              autoComplete="username"
              disabled={isLoading}
            />
          </div>

          <div>
            <label htmlFor="password" className="block text-sm font-medium text-gray-700 mb-2">
              Password
            </label>
            <input
              type="password"
              id="password"
              name="password"
              value={formData.password}
              onChange={handleChange}
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-colors"
              placeholder="Enter your password"
              autoComplete="current-password"
              disabled={isLoading}
            />
          </div>

          {error && (
            <div className="bg-red-50 text-red-700 p-3 rounded-lg text-sm border border-red-200">
              {error}
            </div>
          )}

          <button
            type="submit"
            disabled={isLoading}
            className="w-full bg-gradient-to-r from-indigo-500 to-purple-600 text-white py-3 px-4 rounded-lg hover:from-indigo-600 hover:to-purple-700 focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
          >
            {isLoading ? (
              <>
                <Loader2 className="w-4 h-4 animate-spin" />
                Signing In...
              </>
            ) : (
              <>
                <LogIn className="w-4 h-4" />
                Sign In
              </>
            )}
          </button>
        </form>

        {/* Test Users */}
        <div className="mt-8 pt-6 border-t border-gray-200">
          <div className="flex items-center gap-2 mb-4">
            <Users className="w-4 h-4 text-gray-500" />
            <span className="text-sm font-medium text-gray-700">Quick Login (Test Users)</span>
          </div>
          <div className="grid grid-cols-2 gap-2">
            {config.TEST_USERS.map((user) => (
              <button
                key={user.username}
                type="button"
                onClick={() => handleTestUser(user)}
                disabled={isLoading}
                className="text-xs px-3 py-2 bg-gray-100 text-gray-700 rounded-md hover:bg-gray-200 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {user.username}
              </button>
            ))}
          </div>
          <p className="text-xs text-gray-500 mt-2">
            Password for all test users: <code className="bg-gray-100 px-1 rounded">password123</code>
          </p>
        </div>
      </div>
    </div>
  );
};

export default LoginScreen;