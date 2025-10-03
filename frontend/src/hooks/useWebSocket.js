import { useState, useCallback, useRef, useEffect } from 'react';
import config from '../utils/config';

export const useWebSocket = () => {
  const [connectionStatus, setConnectionStatus] = useState('disconnected');
  const [messages, setMessages] = useState([]);
  const [userCount, setUserCount] = useState(0);
  
  const wsRef = useRef(null);
  const reconnectAttemptsRef = useRef(0);
  const usernameRef = useRef(null);

  const connect = useCallback((username) => {
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      console.log('WebSocket already connected');
      return;
    }

    usernameRef.current = username;
    const url = `${config.WS_URL}/ws?username=${encodeURIComponent(username)}`;
    
    console.log(`Connecting to WebSocket: ${url}`);
    setConnectionStatus('connecting');

    try {
      wsRef.current = new WebSocket(url);
      
      wsRef.current.onopen = (event) => {
        console.log('WebSocket connected successfully');
        setConnectionStatus('connected');
        reconnectAttemptsRef.current = 0;
      };

      wsRef.current.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data);
          console.log('WebSocket message received:', message);
          
          setMessages(prev => {
            const newMessages = [...prev, message];
            return newMessages.slice(-config.MESSAGE_HISTORY_LIMIT);
          });

          if (message.type === 'user_joined' || message.type === 'user_left') {
            setUserCount(prev => 
              message.type === 'user_joined' ? Math.min(prev + 1, 4) : Math.max(prev - 1, 0)
            );
          }
        } catch (error) {
          console.error('Failed to parse WebSocket message:', error);
        }
      };

      wsRef.current.onclose = (event) => {
        console.log('WebSocket disconnected:', event.code, event.reason);
        setConnectionStatus('disconnected');
        
        if (event.code !== 1000 && reconnectAttemptsRef.current < config.RECONNECT_ATTEMPTS) {
          attemptReconnect();
        }
      };

      wsRef.current.onerror = (error) => {
        console.error('WebSocket error:', error);
        setConnectionStatus('error');
      };

    } catch (error) {
      console.error('WebSocket connection failed:', error);
      setConnectionStatus('error');
    }
  }, []);

  const attemptReconnect = useCallback(() => {
    if (reconnectAttemptsRef.current >= config.RECONNECT_ATTEMPTS) {
      console.log('Max reconnection attempts reached');
      return;
    }

    reconnectAttemptsRef.current++;
    console.log(`Attempting reconnection ${reconnectAttemptsRef.current}/${config.RECONNECT_ATTEMPTS}`);
    
    setConnectionStatus('reconnecting');

    setTimeout(() => {
      if (usernameRef.current) {
        connect(usernameRef.current);
      }
    }, config.RECONNECT_DELAY);
  }, [connect]);

  const sendMessage = useCallback((content) => {
    if (!wsRef.current || wsRef.current.readyState !== WebSocket.OPEN) {
      console.error('WebSocket not connected');
      throw new Error('WebSocket not connected');
    }

    const message = {
      type: 'message',
      content: content.trim()
    };

    console.log('Sending message:', message);
    wsRef.current.send(JSON.stringify(message));
  }, []);

  const disconnect = useCallback(() => {
    if (wsRef.current) {
      console.log('Manually disconnecting WebSocket');
      reconnectAttemptsRef.current = config.RECONNECT_ATTEMPTS;
      wsRef.current.close(1000, 'Manual disconnect');
      wsRef.current = null;
    }
    usernameRef.current = null;
    setConnectionStatus('disconnected');
    setMessages([]);
    setUserCount(0);
  }, []);

  const isConnected = connectionStatus === 'connected';

  useEffect(() => {
    return () => {
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, []);

  return {
    connectionStatus,
    messages,
    userCount,
    connect,
    sendMessage,
    disconnect,
    isConnected
  };
};