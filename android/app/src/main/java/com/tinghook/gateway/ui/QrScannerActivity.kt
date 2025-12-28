package com.tinghook.gateway.ui

import android.Manifest
import android.app.Activity
import android.content.Intent
import android.content.pm.PackageManager
import android.os.Bundle
import android.util.Log
import android.view.SurfaceHolder
import android.view.SurfaceView
import android.widget.Toast
import androidx.activity.result.contract.ActivityResultContracts
import androidx.appcompat.app.AppCompatActivity
import androidx.camera.core.CameraSelector
import androidx.camera.core.ImageAnalysis
import androidx.camera.core.Preview
import androidx.camera.lifecycle.ProcessCameraProvider
import androidx.camera.view.PreviewView
import androidx.core.content.ContextCompat
import com.google.mlkit.vision.barcode.BarcodeScanning
import com.google.mlkit.vision.barcode.common.Barcode
import com.google.mlkit.vision.common.InputImage
import org.json.JSONObject
import java.util.concurrent.ExecutorService
import java.util.concurrent.Executors

class QrScannerActivity : AppCompatActivity() {
    
    companion object {
        private const val TAG = "QrScannerActivity"
        const val EXTRA_API_KEY = "api_key"
        const val EXTRA_SERVER_URL = "server_url"
    }
    
    private lateinit var cameraExecutor: ExecutorService
    private lateinit var previewView: PreviewView
    private var isProcessing = false
    
    private val cameraPermissionLauncher = registerForActivityResult(
        ActivityResultContracts.RequestPermission()
    ) { isGranted ->
        if (isGranted) {
            startCamera()
        } else {
            Toast.makeText(this, "Camera permission required for QR scanning", Toast.LENGTH_LONG).show()
            finish()
        }
    }
    
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        
        previewView = PreviewView(this)
        setContentView(previewView)
        
        cameraExecutor = Executors.newSingleThreadExecutor()
        
        if (ContextCompat.checkSelfPermission(this, Manifest.permission.CAMERA) == PackageManager.PERMISSION_GRANTED) {
            startCamera()
        } else {
            cameraPermissionLauncher.launch(Manifest.permission.CAMERA)
        }
    }
    
    private fun startCamera() {
        val cameraProviderFuture = ProcessCameraProvider.getInstance(this)
        
        cameraProviderFuture.addListener({
            val cameraProvider = cameraProviderFuture.get()
            
            val preview = Preview.Builder().build().also {
                it.setSurfaceProvider(previewView.surfaceProvider)
            }
            
            val imageAnalysis = ImageAnalysis.Builder()
                .setBackpressureStrategy(ImageAnalysis.STRATEGY_KEEP_ONLY_LATEST)
                .build()
                .also {
                    it.setAnalyzer(cameraExecutor) { imageProxy ->
                        if (!isProcessing) {
                            isProcessing = true
                            processImage(imageProxy)
                        } else {
                            imageProxy.close()
                        }
                    }
                }
            
            val cameraSelector = CameraSelector.DEFAULT_BACK_CAMERA
            
            try {
                cameraProvider.unbindAll()
                cameraProvider.bindToLifecycle(this, cameraSelector, preview, imageAnalysis)
            } catch (e: Exception) {
                Log.e(TAG, "Camera binding failed", e)
            }
        }, ContextCompat.getMainExecutor(this))
    }
    
    @androidx.annotation.OptIn(androidx.camera.core.ExperimentalGetImage::class)
    private fun processImage(imageProxy: ImageAnalysis.ImageProxy) {
        val mediaImage = imageProxy.image
        if (mediaImage != null) {
            val image = InputImage.fromMediaImage(mediaImage, imageProxy.imageInfo.rotationDegrees)
            val scanner = BarcodeScanning.getClient()
            
            scanner.process(image)
                .addOnSuccessListener { barcodes ->
                    for (barcode in barcodes) {
                        if (barcode.valueType == Barcode.TYPE_TEXT || barcode.valueType == Barcode.TYPE_URL) {
                            barcode.rawValue?.let { handleQrCode(it) }
                            return@addOnSuccessListener
                        }
                    }
                    isProcessing = false
                }
                .addOnFailureListener {
                    Log.e(TAG, "Barcode scanning failed", it)
                    isProcessing = false
                }
                .addOnCompleteListener {
                    imageProxy.close()
                }
        } else {
            imageProxy.close()
            isProcessing = false
        }
    }
    
    private fun handleQrCode(rawValue: String) {
        try {
            val json = JSONObject(rawValue)
            val apiKey = json.optString("apiKey", json.optString("api_key", ""))
            val serverUrl = json.optString("serverUrl", json.optString("server_url", ""))
            
            if (apiKey.isNotEmpty() && serverUrl.isNotEmpty()) {
                val resultIntent = Intent().apply {
                    putExtra(EXTRA_API_KEY, apiKey)
                    putExtra(EXTRA_SERVER_URL, serverUrl)
                }
                setResult(Activity.RESULT_OK, resultIntent)
                finish()
            } else {
                runOnUiThread {
                    Toast.makeText(this, "Invalid QR code format", Toast.LENGTH_SHORT).show()
                }
                isProcessing = false
            }
        } catch (e: Exception) {
            Log.e(TAG, "Failed to parse QR code", e)
            runOnUiThread {
                Toast.makeText(this, "Invalid QR code format", Toast.LENGTH_SHORT).show()
            }
            isProcessing = false
        }
    }
    
    override fun onDestroy() {
        super.onDestroy()
        cameraExecutor.shutdown()
    }
}
