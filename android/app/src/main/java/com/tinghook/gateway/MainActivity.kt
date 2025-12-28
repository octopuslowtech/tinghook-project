package com.tinghook.gateway

import android.content.Intent
import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.result.contract.ActivityResultContracts
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.unit.dp
import com.tinghook.gateway.service.GatewayService
import com.tinghook.gateway.ui.QrScannerActivity
import com.tinghook.gateway.util.PermissionHelper

class MainActivity : ComponentActivity(), NativeEngineCallback {
    
    private var connectionState by mutableStateOf(ConnectionState.DISCONNECTED)
    private var permissions by mutableStateOf(PermissionHelper.PermissionStatus(false, false, false, false))
    private var isServiceRunning by mutableStateOf(false)
    
    private var serverUrl: String = ""
    private var apiKey: String = ""
    private var deviceUid: String = ""
    
    private val smsPermissionLauncher = registerForActivityResult(
        ActivityResultContracts.RequestMultiplePermissions()
    ) { _ ->
        permissions = PermissionHelper.checkPermissions(this)
    }
    
    private val qrScannerLauncher = registerForActivityResult(
        ActivityResultContracts.StartActivityForResult()
    ) { result ->
        if (result.resultCode == RESULT_OK) {
            result.data?.let { data ->
                apiKey = data.getStringExtra(QrScannerActivity.EXTRA_API_KEY) ?: ""
                serverUrl = data.getStringExtra(QrScannerActivity.EXTRA_SERVER_URL) ?: ""
                deviceUid = generateDeviceUid()
                
                saveCredentials(serverUrl, apiKey, deviceUid)
                startGatewayService()
            }
        }
    }
    
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        
        loadCredentials()
        permissions = PermissionHelper.checkPermissions(this)
        
        setContent {
            MaterialTheme {
                MainScreen()
            }
        }
    }
    
    override fun onResume() {
        super.onResume()
        permissions = PermissionHelper.checkPermissions(this)
    }
    
    @Composable
    fun MainScreen() {
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(16.dp)
                .verticalScroll(rememberScrollState()),
            horizontalAlignment = Alignment.CenterHorizontally
        ) {
            Text("TingHook Gateway", style = MaterialTheme.typography.headlineLarge)
            Spacer(modifier = Modifier.height(24.dp))
            
            ConnectionStatusCard()
            
            Spacer(modifier = Modifier.height(16.dp))
            
            PermissionsCard()
            
            Spacer(modifier = Modifier.height(16.dp))
            
            ActionButtons()
        }
    }
    
    @Composable
    fun ConnectionStatusCard() {
        Card(modifier = Modifier.fillMaxWidth()) {
            Column(modifier = Modifier.padding(16.dp)) {
                Text("Connection Status", style = MaterialTheme.typography.titleMedium)
                Spacer(modifier = Modifier.height(8.dp))
                
                Row(verticalAlignment = Alignment.CenterVertically) {
                    val (color, text) = when (connectionState) {
                        ConnectionState.CONNECTED -> Color.Green to "Connected"
                        ConnectionState.CONNECTING -> Color.Yellow to "Connecting..."
                        ConnectionState.RECONNECTING -> Color.Yellow to "Reconnecting..."
                        ConnectionState.DISCONNECTED -> Color.Red to "Disconnected"
                    }
                    
                    Surface(
                        modifier = Modifier.size(12.dp),
                        shape = MaterialTheme.shapes.small,
                        color = color
                    ) {}
                    Spacer(modifier = Modifier.width(8.dp))
                    Text(text)
                }
                
                if (serverUrl.isNotEmpty()) {
                    Spacer(modifier = Modifier.height(8.dp))
                    Text("Server: $serverUrl", style = MaterialTheme.typography.bodySmall)
                }
            }
        }
    }
    
    @Composable
    fun PermissionsCard() {
        Card(modifier = Modifier.fillMaxWidth()) {
            Column(modifier = Modifier.padding(16.dp)) {
                Text("Permissions", style = MaterialTheme.typography.titleMedium)
                Spacer(modifier = Modifier.height(8.dp))
                
                PermissionRow("SMS Receive", permissions.smsReceive) {
                    smsPermissionLauncher.launch(PermissionHelper.getSmsPermissions())
                }
                PermissionRow("SMS Send", permissions.smsSend) {
                    smsPermissionLauncher.launch(PermissionHelper.getSmsPermissions())
                }
                PermissionRow("Notifications", permissions.notifications) {
                    PermissionHelper.openNotificationListenerSettings(this@MainActivity)
                }
                PermissionRow("Battery Optimization", permissions.batteryOptimization) {
                    PermissionHelper.requestBatteryOptimizationExemption(this@MainActivity)
                }
            }
        }
    }
    
    @Composable
    fun PermissionRow(name: String, granted: Boolean, onRequest: () -> Unit) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(vertical = 4.dp),
            horizontalArrangement = Arrangement.SpaceBetween,
            verticalAlignment = Alignment.CenterVertically
        ) {
            Text(name)
            if (granted) {
                Text("✓", color = Color.Green)
            } else {
                TextButton(onClick = onRequest) { Text("Grant") }
            }
        }
    }
    
    @Composable
    fun ActionButtons() {
        Column(modifier = Modifier.fillMaxWidth()) {
            if (apiKey.isEmpty()) {
                Button(
                    onClick = { qrScannerLauncher.launch(Intent(this@MainActivity, QrScannerActivity::class.java)) },
                    modifier = Modifier.fillMaxWidth()
                ) {
                    Text("Scan QR to Connect")
                }
            } else {
                if (isServiceRunning) {
                    Button(
                        onClick = { reconnect() },
                        modifier = Modifier.fillMaxWidth()
                    ) {
                        Text("Reconnect")
                    }
                    Spacer(modifier = Modifier.height(8.dp))
                    OutlinedButton(
                        onClick = { stopGatewayService() },
                        modifier = Modifier.fillMaxWidth()
                    ) {
                        Text("Stop Service")
                    }
                } else {
                    Button(
                        onClick = { startGatewayService() },
                        modifier = Modifier.fillMaxWidth()
                    ) {
                        Text("Start Service")
                    }
                }
            }
        }
    }
    
    private fun startGatewayService() {
        GatewayService.start(this, serverUrl, apiKey, deviceUid)
        isServiceRunning = true
    }
    
    private fun stopGatewayService() {
        GatewayService.stop(this)
        isServiceRunning = false
        connectionState = ConnectionState.DISCONNECTED
    }
    
    private fun reconnect() {
        val intent = Intent(this, GatewayService::class.java).apply {
            action = "com.tinghook.gateway.RECONNECT"
        }
        startService(intent)
    }
    
    private fun generateDeviceUid(): String {
        return java.util.UUID.randomUUID().toString()
    }
    
    private fun saveCredentials(url: String, key: String, uid: String) {
        getSharedPreferences("tinghook", MODE_PRIVATE).edit().apply {
            putString("server_url", url)
            putString("api_key", key)
            putString("device_uid", uid)
            apply()
        }
    }
    
    private fun loadCredentials() {
        val prefs = getSharedPreferences("tinghook", MODE_PRIVATE)
        serverUrl = prefs.getString("server_url", "") ?: ""
        apiKey = prefs.getString("api_key", "") ?: ""
        deviceUid = prefs.getString("device_uid", "") ?: ""
    }
    
    override fun onMessage(type: String, data: String) {
    }
    
    override fun onStateChange(state: Int) {
        connectionState = ConnectionState.fromInt(state)
    }
}
