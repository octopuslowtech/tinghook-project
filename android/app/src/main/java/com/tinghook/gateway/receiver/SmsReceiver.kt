package com.tinghook.gateway.receiver

import android.content.BroadcastReceiver
import android.content.Context
import android.content.Intent
import android.provider.Telephony
import android.util.Log
import com.tinghook.gateway.NativeEngine

class SmsReceiver : BroadcastReceiver() {
    companion object {
        private const val TAG = "SmsReceiver"
    }
    
    override fun onReceive(context: Context, intent: Intent) {
        if (intent.action != Telephony.Sms.Intents.SMS_RECEIVED_ACTION) return
        
        val messages = Telephony.Sms.Intents.getMessagesFromIntent(intent)
        if (messages.isNullOrEmpty()) return
        
        for (message in messages) {
            val sender = message.originatingAddress ?: "Unknown"
            val content = message.messageBody ?: ""
            val simSlot = intent.getIntExtra("android.telephony.extra.SLOT_INDEX", 0)
            
            Log.i(TAG, "SMS received from $sender")
            
            if (NativeEngine.isConnected()) {
                NativeEngine.sendSMSReceivedEvent(sender, content, simSlot)
            }
        }
    }
}
