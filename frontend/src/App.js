import React, { useState } from 'react';
import LoginScreen from './components/LoginScreen';
import ChatScreen from './components/ChatScreen';
import { useWebSocket } from './hooks/useWebSocket';
import api from './utils/api';

function App() {
  const [user, setUser] = useState(null);
  const [isLoading, setIsLoading] = useState(false);
  const webSocket = useWebSocket();

  const handleLogin = async (username, password) => {
    setIsLoading(true);
    
    try {
      console.log('Attempting login for:', username);
      const response = await api.login(username, password);
      
      if (response.success) {
        console.log('Login successful, connecting WebSocket');
        setUser({ username });
        
        webSocket.connect(username);
      } else {
        throw new Error(response.message || 'Login failed');
      }
    } catch (error) {
      console.error('Login failed:', error);
      throw error;
    } finally {
      setIsLoading(false);
    }
  };

  const handleLogout = () => {
    console.log('User logging out');
    webSocket.disconnect();
    setUser(null);
  };

  return (
    <div className="App">
      {!user ? (
        <LoginScreen onLogin={handleLogin} isLoading={isLoading} />
      ) : (
        <ChatScreen 
          username={user.username}
          webSocket={webSocket}
          onLogout={handleLogout}
        />
      )}
    </div>
  );
}

export default App;