package com.tinghook.gateway

import android.util.Log

enum class ConnectionState(val value: Int) {
    DISCONNECTED(0),
    CONNECTING(1),
    CONNECTED(2),
    RECONNECTING(3);
    
    companion object {
        fun fromInt(value: Int) = entries.find { it.value == value } ?: DISCONNECTED
    }
}

interface NativeEngineCallback {
    fun onMessage(type: String, data: String)
    fun onStateChange(state: Int)
}

object NativeEngine {
    private const val TAG = "NativeEngine"
    
    private var callback: NativeEngineCallback? = null
    private var isInitialized = false
    
    init {
        System.loadLibrary("tinghook-engine")
    }
    
    fun initialize(callback: NativeEngineCallback) {
        if (isInitialized) {
            Log.w(TAG, "Already initialized")
            return
        }
        
        this.callback = callback
        init(object : NativeEngineCallback {
            override fun onMessage(type: String, data: String) {
                callback.onMessage(type, data)
            }
            
            override fun onStateChange(state: Int) {
                callback.onStateChange(state)
            }
        })
        isInitialized = true
        Log.i(TAG, "Initialized")
    }
    
    fun connectToServer(serverUrl: String, apiKey: String, deviceUid: String) {
        if (!isInitialized) {
            Log.e(TAG, "Not initialized")
            return
        }
        connect(serverUrl, apiKey, deviceUid)
    }
    
    fun sendPingStatus(battery: Int, signal: Int) {
        sendPing(battery, signal)
    }
    
    fun sendSMSReceivedEvent(sender: String, content: String, simSlot: Int) {
        sendSMSReceived(sender, content, simSlot)
    }
    
    // Native methods
    private external fun init(callback: NativeEngineCallback)
    private external fun connect(serverUrl: String, apiKey: String, deviceUid: String)
    external fun disconnect()
    external fun sendMessage(message: String)
    private external fun sendPing(battery: Int, signal: Int)
    private external fun sendSMSReceived(sender: String, content: String, simSlot: Int)
    external fun isConnected(): Boolean
    external fun getConnectionState(): Int
    external fun destroy()
}
