package com.tinghook.gateway

import android.app.Application
import android.util.Log

class TingHookApplication : Application() {

    companion object {
        private const val TAG = "TingHookApplication"

        @Volatile
        private var instance: TingHookApplication? = null

        fun getInstance(): TingHookApplication {
            return instance ?: throw IllegalStateException("Application not initialized")
        }
    }

    override fun onCreate() {
        super.onCreate()
        instance = this
        Log.d(TAG, "TingHook Application initialized")

        initializeNativeEngine()
    }

    private fun initializeNativeEngine() {
        try {
            NativeEngine.init()
            Log.d(TAG, "Native engine initialized successfully")
        } catch (e: Exception) {
            Log.e(TAG, "Failed to initialize native engine", e)
        }
    }

    override fun onTerminate() {
        super.onTerminate()
        NativeEngine.disconnect()
        Log.d(TAG, "TingHook Application terminated")
    }
}
