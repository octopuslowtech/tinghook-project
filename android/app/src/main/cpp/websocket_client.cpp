#include "include/websocket_client.h"
#include "include/tinghook_engine.h"
#include <android/log.h>
#include <chrono>

#define LOG_TAG "TingHookWS"
#define LOGI(...) __android_log_print(ANDROID_LOG_INFO, LOG_TAG, __VA_ARGS__)
#define LOGE(...) __android_log_print(ANDROID_LOG_ERROR, LOG_TAG, __VA_ARGS__)

namespace tinghook {

WebSocketClient::WebSocketClient() = default;

WebSocketClient::~WebSocketClient() {
    disconnect();
}

void WebSocketClient::connect(const std::string& url, const std::string& apiKey, const std::string& deviceUid) {
    if (m_state != ConnectionState::Disconnected) {
        return;
    }
    
    m_url = url;
    m_apiKey = apiKey;
    m_deviceUid = deviceUid;
    m_shouldRun = true;
    m_reconnectAttempts = 0;
    
    m_thread = std::thread(&WebSocketClient::run, this);
}

void WebSocketClient::disconnect() {
    m_shouldRun = false;
    m_state = ConnectionState::Disconnected;
    
    if (m_thread.joinable()) {
        m_thread.join();
    }
}

void WebSocketClient::send(const std::string& message) {
    std::lock_guard<std::mutex> lock(m_mutex);
    m_outgoingQueue.push(message);
}

void WebSocketClient::run() {
    while (m_shouldRun) {
        m_state = ConnectionState::Connecting;
        if (m_stateCallback) m_stateCallback(m_state.load());
        
        LOGI("Connecting to %s", m_url.c_str());
        
        // TODO: Implement actual WebSocket connection using libwebsockets or similar
        // For now, this is a placeholder that simulates connection
        
        // Simulate connection delay
        std::this_thread::sleep_for(std::chrono::milliseconds(1000));
        
        if (!m_shouldRun) break;
        
        m_state = ConnectionState::Connected;
        if (m_stateCallback) m_stateCallback(m_state.load());
        m_reconnectAttempts = 0;
        
        LOGI("Connected successfully");
        
        // Main loop - process messages
        while (m_shouldRun && m_state == ConnectionState::Connected) {
            processMessages();
            std::this_thread::sleep_for(std::chrono::milliseconds(100));
        }
        
        // Reconnect if needed
        if (m_shouldRun) {
            reconnect();
        }
    }
    
    m_state = ConnectionState::Disconnected;
    if (m_stateCallback) m_stateCallback(m_state.load());
}

void WebSocketClient::reconnect() {
    if (m_reconnectAttempts >= MAX_RECONNECT_ATTEMPTS) {
        LOGE("Max reconnect attempts reached");
        m_shouldRun = false;
        return;
    }
    
    m_state = ConnectionState::Reconnecting;
    if (m_stateCallback) m_stateCallback(m_state.load());
    
    m_reconnectAttempts++;
    int delay = RECONNECT_DELAY_MS * m_reconnectAttempts; // Exponential backoff
    
    LOGI("Reconnecting in %d ms (attempt %d/%d)", delay, m_reconnectAttempts, MAX_RECONNECT_ATTEMPTS);
    std::this_thread::sleep_for(std::chrono::milliseconds(delay));
}

void WebSocketClient::processMessages() {
    std::lock_guard<std::mutex> lock(m_mutex);
    while (!m_outgoingQueue.empty()) {
        auto msg = m_outgoingQueue.front();
        m_outgoingQueue.pop();
        // TODO: Send message via WebSocket
        LOGI("Sending: %s", msg.c_str());
    }
}

void WebSocketClient::setMessageCallback(MessageCallback callback) {
    m_messageCallback = std::move(callback);
}

void WebSocketClient::setStateCallback(StateCallback callback) {
    m_stateCallback = std::move(callback);
}

ConnectionState WebSocketClient::getState() const {
    return m_state.load();
}

bool WebSocketClient::isConnected() const {
    return m_state == ConnectionState::Connected;
}

} // namespace tinghook
