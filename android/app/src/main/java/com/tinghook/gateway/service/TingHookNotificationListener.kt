package com.tinghook.gateway.service

import android.service.notification.NotificationListenerService
import android.service.notification.StatusBarNotification
import android.util.Log
import com.tinghook.gateway.NativeEngine
import org.json.JSONObject

class TingHookNotificationListener : NotificationListenerService() {
    companion object {
        private const val TAG = "NotifListener"
        private val ALLOWED_PACKAGES = setOf<String>()
    }
    
    override fun onNotificationPosted(sbn: StatusBarNotification) {
        val packageName = sbn.packageName
        
        if (packageName == "com.tinghook.gateway") return
        if (packageName.startsWith("com.android")) return
        
        if (ALLOWED_PACKAGES.isNotEmpty() && packageName !in ALLOWED_PACKAGES) return
        
        val notification = sbn.notification
        val extras = notification.extras
        
        val title = extras.getString("android.title") ?: ""
        val content = extras.getCharSequence("android.text")?.toString() ?: ""
        val appName = try {
            packageManager.getApplicationLabel(
                packageManager.getApplicationInfo(packageName, 0)
            ).toString()
        } catch (e: Exception) {
            packageName
        }
        
        Log.d(TAG, "Notification from $appName: $title")
        
        if (NativeEngine.isConnected()) {
            val data = JSONObject().apply {
                put("app_package", packageName)
                put("app_name", appName)
                put("title", title)
                put("content", content)
                put("timestamp", System.currentTimeMillis())
            }
            NativeEngine.sendMessage("""{"type":"NOTIFICATION_RECEIVED","data":$data}""")
        }
    }
    
    override fun onNotificationRemoved(sbn: StatusBarNotification) {
        // Optional: track dismissed notifications
    }
}
