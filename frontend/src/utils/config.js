const config = {
  API_URL: process.env.NODE_ENV === 'production' 
    ? 'https://websocket-shit-production.up.railway.app'
    : 'http://localhost:8080',
    
  WS_URL: process.env.NODE_ENV === 'production'
    ? 'wss://websocket-shit-production.up.railway.app'
    : 'ws://localhost:8080',
    
  RECONNECT_ATTEMPTS: 5,
  RECONNECT_DELAY: 3000,
  
  MAX_MESSAGE_LENGTH: 500,
  MESSAGE_HISTORY_LIMIT: 100,
  
  TEST_USERS: [
    { username: 'alice', password: 'password123' },
    { username: 'bob', password: 'password123' },
    { username: 'charlie', password: 'password123' },
    { username: 'diana', password: 'password123' }
  ]
};

export default config;