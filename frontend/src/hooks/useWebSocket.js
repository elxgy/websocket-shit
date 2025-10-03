import { useState, useCallback, useRef, useEffect } from "react";
import config from "../utils/config";

export const useWebSocket = () => {
  const [connectionStatus, setConnectionStatus] = useState("idle");
  const [messages, setMessages] = useState([]);
  const [userCount, setUserCount] = useState(0);

  const wsRef = useRef(null);
  const reconnectAttemptsRef = useRef(0);
  const usernameRef = useRef(null);
  const reconnectTimeoutRef = useRef(null);

  const attemptReconnect = useCallback(() => {
    if (reconnectAttemptsRef.current >= config.RECONNECT_ATTEMPTS) {
      console.log("Max reconnection attempts reached");
      setConnectionStatus("failed");
      return;
    }

    reconnectAttemptsRef.current++;
    console.log(
      `Attempting to reconnect... (${reconnectAttemptsRef.current}/${config.RECONNECT_ATTEMPTS})`,
    );

    setConnectionStatus("reconnecting");

    reconnectTimeoutRef.current = setTimeout(() => {
      if (usernameRef.current) {
        const url = `${config.WS_URL}/ws?username=${encodeURIComponent(usernameRef.current)}`;
        console.log(`Reconnecting to WebSocket: ${url}`);

        try {
          wsRef.current = new WebSocket(url);

          wsRef.current.onopen = (event) => {
            console.log("WebSocket reconnected successfully");
            setConnectionStatus("connected");
            reconnectAttemptsRef.current = 0;
          };

          wsRef.current.onmessage = (event) => {
            try {
              const message = JSON.parse(event.data);
              console.log("WebSocket message received:", message);

              setMessages((prev) => {
                const newMessages = [...prev, message];
                return newMessages.slice(-config.MESSAGE_HISTORY_LIMIT);
              });

              if (
                message.type === "user_joined" ||
                message.type === "user_left"
              ) {
                setUserCount((prev) =>
                  message.type === "user_joined"
                    ? Math.min(prev + 1, 4)
                    : Math.max(prev - 1, 0),
                );
              }
            } catch (error) {
              console.error("Failed to parse WebSocket message:", error);
            }
          };

          wsRef.current.onclose = (event) => {
            console.log("WebSocket disconnected:", event.code, event.reason);
            setConnectionStatus("disconnected");

            if (
              event.code !== 1000 &&
              reconnectAttemptsRef.current < config.RECONNECT_ATTEMPTS
            ) {
              attemptReconnect();
            }
          };

          wsRef.current.onerror = (error) => {
            console.error("WebSocket error:", error);
            setConnectionStatus("error");

            if (
              usernameRef.current &&
              reconnectAttemptsRef.current < config.RECONNECT_ATTEMPTS
            ) {
              attemptReconnect();
            }
          };
        } catch (error) {
          console.error("WebSocket reconnection failed:", error);
          setConnectionStatus("error");
        }
      }
    }, config.RECONNECT_DELAY);
  }, []);

  const connect = useCallback(
    (username) => {
      if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
        console.log("WebSocket already connected");
        return;
      }

      // Clear any existing reconnection timeout
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
        reconnectTimeoutRef.current = null;
      }

      usernameRef.current = username;
      const url = `${config.WS_URL}/ws?username=${encodeURIComponent(username)}`;

      console.log(`Connecting to WebSocket: ${url}`);
      setConnectionStatus("connecting");

      try {
        wsRef.current = new WebSocket(url);

        wsRef.current.onopen = (event) => {
          console.log("WebSocket connected successfully");
          setConnectionStatus("connected");
          reconnectAttemptsRef.current = 0;
        };

        wsRef.current.onmessage = (event) => {
          try {
            const message = JSON.parse(event.data);
            console.log("WebSocket message received:", message);

            setMessages((prev) => {
              const newMessages = [...prev, message];
              return newMessages.slice(-config.MESSAGE_HISTORY_LIMIT);
            });

            if (
              message.type === "user_joined" ||
              message.type === "user_left"
            ) {
              setUserCount((prev) =>
                message.type === "user_joined"
                  ? Math.min(prev + 1, 4)
                  : Math.max(prev - 1, 0),
              );
            }
          } catch (error) {
            console.error("Failed to parse WebSocket message:", error);
          }
        };

        wsRef.current.onclose = (event) => {
          console.log("WebSocket disconnected:", event.code, event.reason);
          setConnectionStatus("disconnected");

          if (
            event.code !== 1000 &&
            reconnectAttemptsRef.current < config.RECONNECT_ATTEMPTS
          ) {
            attemptReconnect();
          }
        };

        wsRef.current.onerror = (error) => {
          console.error("WebSocket error:", error);
          setConnectionStatus("error");

          // Attempt reconnection on error if we have a username
          if (
            usernameRef.current &&
            reconnectAttemptsRef.current < config.RECONNECT_ATTEMPTS
          ) {
            attemptReconnect();
          }
        };
      } catch (error) {
        console.error("WebSocket connection failed:", error);
        setConnectionStatus("error");
      }
    },
    [attemptReconnect],
  );

  const sendMessage = useCallback((content) => {
    if (!wsRef.current || wsRef.current.readyState !== WebSocket.OPEN) {
      console.error("WebSocket not connected");
      throw new Error("WebSocket not connected");
    }

    const message = {
      type: "message",
      content: content.trim(),
    };

    console.log("Sending message:", message);
    wsRef.current.send(JSON.stringify(message));
  }, []);

  const disconnect = useCallback(() => {
    if (wsRef.current) {
      console.log("Manually disconnecting WebSocket");
      reconnectAttemptsRef.current = config.RECONNECT_ATTEMPTS;
      wsRef.current.close(1000, "Manual disconnect");
      wsRef.current = null;
    }
    usernameRef.current = null;
    setConnectionStatus("idle");
    setMessages([]);
    setUserCount(0);
  }, []);

  useEffect(() => {
    return () => {
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, []);

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
      }
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, []);

  return {
    connectionStatus,
    messages,
    userCount,
    sendMessage,
    connect,
    disconnect,
    isConnected: connectionStatus === "connected",
    isReconnecting: connectionStatus === "reconnecting",
  };
};
