#include <jni.h>
#include <string>
#include <memory>
#include "include/websocket_client.h"
#include "include/json_utils.h"
#include <android/log.h>

#define LOG_TAG "TingHookNative"
#define LOGI(...) __android_log_print(ANDROID_LOG_INFO, LOG_TAG, __VA_ARGS__)

namespace {
    std::unique_ptr<tinghook::WebSocketClient> g_wsClient;
    JavaVM* g_jvm = nullptr;
    jobject g_callbackObj = nullptr;
    jmethodID g_onMessageMethod = nullptr;
    jmethodID g_onStateChangeMethod = nullptr;
}

extern "C" {

JNIEXPORT jint JNICALL JNI_OnLoad(JavaVM* vm, void* reserved) {
    g_jvm = vm;
    return JNI_VERSION_1_6;
}

JNIEXPORT void JNICALL
Java_com_tinghook_gateway_NativeEngine_init(JNIEnv* env, jobject thiz, jobject callback) {
    LOGI("NativeEngine init");
    
    // Store callback reference
    g_callbackObj = env->NewGlobalRef(callback);
    
    // Get callback methods
    jclass callbackClass = env->GetObjectClass(callback);
    g_onMessageMethod = env->GetMethodID(callbackClass, "onMessage", "(Ljava/lang/String;Ljava/lang/String;)V");
    g_onStateChangeMethod = env->GetMethodID(callbackClass, "onStateChange", "(I)V");
    
    // Create WebSocket client
    g_wsClient = std::make_unique<tinghook::WebSocketClient>();
    
    // Set message callback
    g_wsClient->setMessageCallback([](const tinghook::Message& msg) {
        JNIEnv* env;
        if (g_jvm->AttachCurrentThread(&env, nullptr) != JNI_OK) return;
        
        jstring type = env->NewStringUTF(msg.type.c_str());
        jstring data = env->NewStringUTF(msg.data.c_str());
        env->CallVoidMethod(g_callbackObj, g_onMessageMethod, type, data);
        
        env->DeleteLocalRef(type);
        env->DeleteLocalRef(data);
    });
    
    // Set state callback
    g_wsClient->setStateCallback([](tinghook::ConnectionState state) {
        JNIEnv* env;
        if (g_jvm->AttachCurrentThread(&env, nullptr) != JNI_OK) return;
        
        env->CallVoidMethod(g_callbackObj, g_onStateChangeMethod, static_cast<int>(state));
    });
}

JNIEXPORT void JNICALL
Java_com_tinghook_gateway_NativeEngine_connect(JNIEnv* env, jobject thiz, 
                                                jstring serverUrl, jstring apiKey, jstring deviceUid) {
    if (!g_wsClient) return;
    
    const char* urlChars = env->GetStringUTFChars(serverUrl, nullptr);
    const char* keyChars = env->GetStringUTFChars(apiKey, nullptr);
    const char* uidChars = env->GetStringUTFChars(deviceUid, nullptr);
    
    std::string url(urlChars);
    std::string key(keyChars);
    std::string uid(uidChars);
    
    env->ReleaseStringUTFChars(serverUrl, urlChars);
    env->ReleaseStringUTFChars(apiKey, keyChars);
    env->ReleaseStringUTFChars(deviceUid, uidChars);
    
    LOGI("Connecting to %s", url.c_str());
    g_wsClient->connect(url, key, uid);
}

JNIEXPORT void JNICALL
Java_com_tinghook_gateway_NativeEngine_disconnect(JNIEnv* env, jobject thiz) {
    if (g_wsClient) {
        g_wsClient->disconnect();
    }
}

JNIEXPORT void JNICALL
Java_com_tinghook_gateway_NativeEngine_sendMessage(JNIEnv* env, jobject thiz, jstring message) {
    if (!g_wsClient) return;
    
    const char* msgChars = env->GetStringUTFChars(message, nullptr);
    std::string msg(msgChars);
    env->ReleaseStringUTFChars(message, msgChars);
    
    g_wsClient->send(msg);
}

JNIEXPORT void JNICALL
Java_com_tinghook_gateway_NativeEngine_sendPing(JNIEnv* env, jobject thiz, jint battery, jint signal) {
    if (!g_wsClient) return;
    
    auto msg = tinghook::createPingMessage(battery, signal);
    g_wsClient->send(msg);
}

JNIEXPORT void JNICALL
Java_com_tinghook_gateway_NativeEngine_sendSMSReceived(JNIEnv* env, jobject thiz,
                                                        jstring sender, jstring content, jint simSlot) {
    if (!g_wsClient) return;
    
    const char* senderChars = env->GetStringUTFChars(sender, nullptr);
    const char* contentChars = env->GetStringUTFChars(content, nullptr);
    
    auto msg = tinghook::createSMSReceivedMessage(senderChars, contentChars, simSlot);
    
    env->ReleaseStringUTFChars(sender, senderChars);
    env->ReleaseStringUTFChars(content, contentChars);
    
    g_wsClient->send(msg);
}

JNIEXPORT jboolean JNICALL
Java_com_tinghook_gateway_NativeEngine_isConnected(JNIEnv* env, jobject thiz) {
    return g_wsClient ? g_wsClient->isConnected() : false;
}

JNIEXPORT jint JNICALL
Java_com_tinghook_gateway_NativeEngine_getConnectionState(JNIEnv* env, jobject thiz) {
    return g_wsClient ? static_cast<int>(g_wsClient->getState()) : 0;
}

JNIEXPORT void JNICALL
Java_com_tinghook_gateway_NativeEngine_destroy(JNIEnv* env, jobject thiz) {
    if (g_wsClient) {
        g_wsClient->disconnect();
        g_wsClient.reset();
    }
    
    if (g_callbackObj) {
        env->DeleteGlobalRef(g_callbackObj);
        g_callbackObj = nullptr;
    }
}

} // extern "C"
