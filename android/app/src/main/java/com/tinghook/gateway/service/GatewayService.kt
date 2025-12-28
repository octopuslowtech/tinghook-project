package com.tinghook.gateway.service

import android.app.Notification
import android.app.NotificationChannel
import android.app.NotificationManager
import android.app.PendingIntent
import android.app.Service
import android.content.Context
import android.content.Intent
import android.os.BatteryManager
import android.os.Build
import android.os.IBinder
import android.os.PowerManager
import android.util.Log
import androidx.core.app.NotificationCompat
import com.tinghook.gateway.ConnectionState
import com.tinghook.gateway.MainActivity
import com.tinghook.gateway.NativeEngine
import com.tinghook.gateway.NativeEngineCallback
import com.tinghook.gateway.R
import kotlinx.coroutines.*

class GatewayService : Service(), NativeEngineCallback {
    
    companion object {
        private const val TAG = "GatewayService"
        private const val CHANNEL_ID = "tinghook_gateway_channel"
        private const val NOTIFICATION_ID = 1
        private const val PING_INTERVAL_MS = 30_000L
        
        private const val ACTION_STOP = "com.tinghook.gateway.STOP"
        private const val ACTION_RECONNECT = "com.tinghook.gateway.RECONNECT"
        
        fun start(context: Context, serverUrl: String, apiKey: String, deviceUid: String) {
            val intent = Intent(context, GatewayService::class.java).apply {
                putExtra("server_url", serverUrl)
                putExtra("api_key", apiKey)
                putExtra("device_uid", deviceUid)
            }
            if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
                context.startForegroundService(intent)
            } else {
                context.startService(intent)
            }
        }
        
        fun stop(context: Context) {
            context.stopService(Intent(context, GatewayService::class.java))
        }
    }
    
    private var wakeLock: PowerManager.WakeLock? = null
    private val serviceScope = CoroutineScope(Dispatchers.IO + SupervisorJob())
    private var pingJob: Job? = null
    
    private var serverUrl: String = ""
    private var apiKey: String = ""
    private var deviceUid: String = ""
    private var currentState: ConnectionState = ConnectionState.DISCONNECTED
    
    override fun onCreate() {
        super.onCreate()
        Log.i(TAG, "Service created")
        createNotificationChannel()
        acquireWakeLock()
    }
    
    override fun onStartCommand(intent: Intent?, flags: Int, startId: Int): Int {
        Log.i(TAG, "Service started")
        
        when (intent?.action) {
            ACTION_STOP -> {
                stopSelf()
                return START_NOT_STICKY
            }
            ACTION_RECONNECT -> {
                reconnect()
                return START_STICKY
            }
        }
        
        serverUrl = intent?.getStringExtra("server_url") ?: return START_NOT_STICKY
        apiKey = intent.getStringExtra("api_key") ?: return START_NOT_STICKY
        deviceUid = intent.getStringExtra("device_uid") ?: return START_NOT_STICKY
        
        startForeground(NOTIFICATION_ID, createNotification("Connecting..."))
        
        NativeEngine.initialize(this)
        NativeEngine.connectToServer(serverUrl, apiKey, deviceUid)
        
        startPingLoop()
        
        return START_STICKY
    }
    
    override fun onBind(intent: Intent?): IBinder? = null
    
    override fun onDestroy() {
        Log.i(TAG, "Service destroyed")
        pingJob?.cancel()
        serviceScope.cancel()
        NativeEngine.disconnect()
        NativeEngine.destroy()
        releaseWakeLock()
        super.onDestroy()
    }
    
    override fun onMessage(type: String, data: String) {
        Log.d(TAG, "Message received: $type")
        when (type) {
            "SEND_SMS" -> {
                // TODO: Parse data and send SMS via SmsManager
            }
        }
    }
    
    override fun onStateChange(state: Int) {
        currentState = ConnectionState.fromInt(state)
        Log.i(TAG, "Connection state: $currentState")
        updateNotification(getStateText(currentState))
    }
    
    private fun reconnect() {
        NativeEngine.disconnect()
        NativeEngine.connectToServer(serverUrl, apiKey, deviceUid)
    }
    
    private fun startPingLoop() {
        pingJob?.cancel()
        pingJob = serviceScope.launch {
            while (isActive) {
                delay(PING_INTERVAL_MS)
                if (NativeEngine.isConnected()) {
                    val battery = getBatteryLevel()
                    NativeEngine.sendPingStatus(battery, 4)
                }
            }
        }
    }
    
    private fun getBatteryLevel(): Int {
        val batteryManager = getSystemService(Context.BATTERY_SERVICE) as BatteryManager
        return batteryManager.getIntProperty(BatteryManager.BATTERY_PROPERTY_CAPACITY)
    }
    
    private fun getStateText(state: ConnectionState): String = when (state) {
        ConnectionState.DISCONNECTED -> "Disconnected"
        ConnectionState.CONNECTING -> "Connecting..."
        ConnectionState.CONNECTED -> "Connected"
        ConnectionState.RECONNECTING -> "Reconnecting..."
    }
    
    private fun createNotificationChannel() {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            val channel = NotificationChannel(
                CHANNEL_ID,
                "TingHook Gateway",
                NotificationManager.IMPORTANCE_LOW
            ).apply {
                description = "SMS Gateway service notification"
                setShowBadge(false)
            }
            val manager = getSystemService(NotificationManager::class.java)
            manager.createNotificationChannel(channel)
        }
    }
    
    private fun createNotification(status: String): Notification {
        val contentIntent = PendingIntent.getActivity(
            this,
            0,
            Intent(this, MainActivity::class.java),
            PendingIntent.FLAG_IMMUTABLE
        )
        
        val stopIntent = PendingIntent.getService(
            this,
            1,
            Intent(this, GatewayService::class.java).apply { action = ACTION_STOP },
            PendingIntent.FLAG_IMMUTABLE
        )
        
        val reconnectIntent = PendingIntent.getService(
            this,
            2,
            Intent(this, GatewayService::class.java).apply { action = ACTION_RECONNECT },
            PendingIntent.FLAG_IMMUTABLE
        )
        
        return NotificationCompat.Builder(this, CHANNEL_ID)
            .setContentTitle("TingHook Gateway")
            .setContentText(status)
            .setSmallIcon(R.drawable.ic_launcher_foreground)
            .setOngoing(true)
            .setContentIntent(contentIntent)
            .addAction(0, "Reconnect", reconnectIntent)
            .addAction(0, "Stop", stopIntent)
            .build()
    }
    
    private fun updateNotification(status: String) {
        val manager = getSystemService(NotificationManager::class.java)
        manager.notify(NOTIFICATION_ID, createNotification(status))
    }
    
    private fun acquireWakeLock() {
        val powerManager = getSystemService(Context.POWER_SERVICE) as PowerManager
        wakeLock = powerManager.newWakeLock(
            PowerManager.PARTIAL_WAKE_LOCK,
            "TingHook::GatewayWakeLock"
        ).apply {
            acquire(24 * 60 * 60 * 1000L)
        }
    }
    
    private fun releaseWakeLock() {
        wakeLock?.let {
            if (it.isHeld) it.release()
        }
        wakeLock = null
    }
}
