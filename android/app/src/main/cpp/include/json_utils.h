#pragma once

#include <string>
#include <nlohmann/json.hpp>

namespace tinghook {

using json = nlohmann::json;

inline std::string createAuthMessage(const std::string& apiKey, const std::string& deviceUid) {
    json msg;
    msg["type"] = "AUTH";
    msg["data"]["api_key"] = apiKey;
    msg["data"]["device_uid"] = deviceUid;
    return msg.dump();
}

inline std::string createPingMessage(int battery, int signal) {
    json msg;
    msg["type"] = "PING";
    msg["data"]["battery"] = battery;
    msg["data"]["signal"] = signal;
    return msg.dump();
}

inline std::string createSMSReceivedMessage(const std::string& sender, const std::string& content, int simSlot) {
    json msg;
    msg["type"] = "SMS_RECEIVED";
    msg["data"]["sender"] = sender;
    msg["data"]["content"] = content;
    msg["data"]["sim_slot"] = simSlot;
    msg["data"]["timestamp"] = ""; // TODO: Add current timestamp
    return msg.dump();
}

} // namespace tinghook
