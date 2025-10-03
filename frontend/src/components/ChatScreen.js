import React, { useState, useRef, useEffect } from 'react';
import { Send, LogOut, MessageCircle, Users, Wifi, WifiOff, RotateCcw } from 'lucide-react';
import config from '../utils/config';

const ChatScreen = ({ username, webSocket, onLogout }) => {
  const [message, setMessage] = useState('');
  const [notification, setNotification] = useState(null);
  const messagesEndRef = useRef(null);
  const inputRef = useRef(null);

  const { connectionStatus, messages, userCount, sendMessage, isConnected } = webSocket;

  // Auto-scroll to bottom when new messages arrive
  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  // Focus input when connected
  useEffect(() => {
    if (isConnected && inputRef.current) {
      inputRef.current.focus();
    }
  }, [isConnected]);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    
    if (!message.trim() || !isConnected) return;

    if (message.length > config.MAX_MESSAGE_LENGTH) {
      showNotification(`Message too long (max ${config.MAX_MESSAGE_LENGTH} characters)`, 'error');
      return;
    }

    try {
      sendMessage(message);
      setMessage('');
    } catch (error) {
      console.error('Failed to send message:', error);
      showNotification('Failed to send message. Please check your connection.', 'error');
    }
  };

  const showNotification = (text, type = 'info') => {
    setNotification({ text, type });
    setTimeout(() => setNotification(null), 5000);
  };

  const handleLogout = () => {
    if (window.confirm('Are you sure you want to sign out?')) {
      onLogout();
    }
  };

  const getConnectionIcon = () => {
    switch (connectionStatus) {
      case 'connected':
        return <Wifi className="w-4 h-4 text-green-400" />;
      case 'connecting':
      case 'reconnecting':
        return <RotateCcw className="w-4 h-4 text-yellow-400 animate-spin" />;
      default:
        return <WifiOff className="w-4 h-4 text-red-400" />;
    }
  };

  const getConnectionText = () => {
    switch (connectionStatus) {
      case 'connected': return 'Connected';
      case 'connecting': return 'Connecting...';
      case 'reconnecting': return 'Reconnecting...';
      case 'error': return 'Connection Error';
      default: return 'Disconnected';
    }
  };

  const formatTime = (timestamp) => {
    return new Date(timestamp).toLocaleTimeString([], { 
      hour: '2-digit', 
      minute: '2-digit' 
    });
  };

  const renderMessage = (msg, index) => {
    if (msg.type === 'message') {
      const isOwn = msg.username === username;
      
      return (
        <div key={index} className={`mb-4 ${isOwn ? 'text-right' : 'text-left'}`}>
          <div className={`inline-block max-w-xs lg:max-w-md ${isOwn ? 'ml-auto' : 'mr-auto'}`}>
            <div className="text-xs text-gray-500 mb-1 px-1">
              {isOwn ? 'You' : msg.username} • {formatTime(msg.timestamp)}
            </div>
            <div className={`message-bubble ${isOwn ? 'own' : 'other'}`}>
              <span className="chat-message">{msg.content}</span>
            </div>
          </div>
        </div>
      );
    } else {
      // System message
      return (
        <div key={index} className="mb-3 text-center">
          <div className="inline-block bg-gray-100 text-gray-600 text-sm px-3 py-1 rounded-full">
            {msg.content} • {formatTime(msg.timestamp)}
          </div>
        </div>
      );
    }
  };

  return (
    <div className="h-screen flex flex-col bg-white">
      {/* Header */}
      <div className="bg-gradient-to-r from-indigo-500 to-purple-600 text-white p-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <MessageCircle className="w-6 h-6" />
            <h1 className="text-xl font-semibold">Chat Room</h1>
            <div className="flex items-center gap-2 bg-white/20 px-3 py-1 rounded-full">
              <Users className="w-4 h-4" />
              <span className="text-sm">{userCount}/4</span>
            </div>
          </div>
          
          <div className="flex items-center gap-4">
            <div className="connection-indicator text-white/90">
              {getConnectionIcon()}
              <span className="text-sm">{getConnectionText()}</span>
            </div>
            <button
              onClick={handleLogout}
              className="bg-white/20 hover:bg-white/30 px-3 py-1 rounded-lg transition-colors flex items-center gap-2"
            >
              <LogOut className="w-4 h-4" />
              <span className="text-sm">Sign Out</span>
            </button>
          </div>
        </div>
      </div>

      {/* Messages Area */}
      <div className="flex-1 overflow-y-auto p-4 bg-gray-50">
        {messages.length === 0 ? (
          <div className="text-center py-12">
            <MessageCircle className="w-16 h-16 text-gray-300 mx-auto mb-4" />
            <h3 className="text-lg font-medium text-gray-900 mb-2">
              Welcome to the chat room, {username}!
            </h3>
            <p className="text-gray-500">
              You can chat with up to 3 other users in real-time.
            </p>
          </div>
        ) : (
          <div className="max-w-4xl mx-auto">
            {messages.map(renderMessage)}
            <div ref={messagesEndRef} />
          </div>
        )}
      </div>

      {/* Message Input */}
      <div className="border-t bg-white p-4">
        <form onSubmit={handleSubmit} className="max-w-4xl mx-auto">
          <div className="flex gap-3">
            <input
              ref={inputRef}
              type="text"
              value={message}
              onChange={(e) => setMessage(e.target.value)}
              placeholder={isConnected ? "Type your message..." : "Connecting..."}
              className="flex-1 px-4 py-3 border border-gray-300 rounded-full focus:ring-2 focus:ring-indigo-500 focus:border-transparent disabled:bg-gray-100 disabled:cursor-not-allowed"
              maxLength={config.MAX_MESSAGE_LENGTH}
              disabled={!isConnected}
            />
            <button
              type="submit"
              disabled={!isConnected || !message.trim()}
              className="bg-gradient-to-r from-indigo-500 to-purple-600 text-white px-6 py-3 rounded-full hover:from-indigo-600 hover:to-purple-700 focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
            >
              <Send className="w-4 h-4" />
              <span className="hidden sm:inline">Send</span>
            </button>
          </div>
          <div className="flex justify-between mt-2 px-1">
            <span className="text-xs text-gray-500">
              {message.length}/{config.MAX_MESSAGE_LENGTH}
            </span>
            {!isConnected && (
              <span className="text-xs text-red-500">
                Reconnecting...
              </span>
            )}
          </div>
        </form>
      </div>

      {/* Notification */}
      {notification && (
        <div className={`fixed top-4 right-4 z-50 p-3 rounded-lg shadow-lg animate-slide-in ${
          notification.type === 'error' 
            ? 'bg-red-500 text-white' 
            : 'bg-blue-500 text-white'
        }`}>
          {notification.text}
        </div>
      )}
    </div>
  );
};

export default ChatScreen;