#pragma once

#include <string>
#include <functional>
#include <thread>
#include <atomic>
#include <mutex>
#include <queue>

namespace tinghook {

enum class ConnectionState {
    Disconnected,
    Connecting,
    Connected,
    Reconnecting
};

struct Message {
    std::string type;
    std::string data;
};

using MessageCallback = std::function<void(const Message&)>;
using StateCallback = std::function<void(ConnectionState)>;

class WebSocketClient {
public:
    WebSocketClient();
    ~WebSocketClient();
    
    void connect(const std::string& url, const std::string& apiKey, const std::string& deviceUid);
    void disconnect();
    void send(const std::string& message);
    
    void setMessageCallback(MessageCallback callback);
    void setStateCallback(StateCallback callback);
    
    ConnectionState getState() const;
    bool isConnected() const;
    
private:
    void run();
    void reconnect();
    void processMessages();
    
    std::string m_url;
    std::string m_apiKey;
    std::string m_deviceUid;
    
    std::atomic<ConnectionState> m_state{ConnectionState::Disconnected};
    std::atomic<bool> m_shouldRun{false};
    
    std::thread m_thread;
    std::mutex m_mutex;
    std::queue<std::string> m_outgoingQueue;
    
    MessageCallback m_messageCallback;
    StateCallback m_stateCallback;
    
    int m_reconnectAttempts{0};
    static constexpr int MAX_RECONNECT_ATTEMPTS = 10;
    static constexpr int RECONNECT_DELAY_MS = 5000;
};

} // namespace tinghook
