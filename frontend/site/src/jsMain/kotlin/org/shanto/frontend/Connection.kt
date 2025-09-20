package org.shanto.frontend

import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.setValue
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.serialization.encodeToString
import kotlinx.serialization.json.Json
import kotlinx.serialization.json.decodeFromJsonElement
import org.shanto.frontend.model.AuthInfo
import org.shanto.frontend.model.AuthResponse
import org.shanto.frontend.model.Envelope
import org.w3c.dom.MessageEvent
import org.w3c.dom.WebSocket


object Connection {
    var socket: WebSocket? = null
    var isConnected by mutableStateOf(false)

    private val _status = MutableStateFlow<Boolean>(false)
    val status = _status.asStateFlow()

    fun connect(url: String) {
        if (socket != null) return
        try {
            socket = WebSocket(url)
            socket?.onopen = { isConnected = true }
            socket?.onclose = { socket!!.close(1000, "normal") }

            socket?.onmessage = { msg ->
                val rawData = msg.data.toString()
                try {
                    val packet = Json.decodeFromString<Envelope>(rawData)
                    when (packet.type) {
                        "auth" -> {
                            val authPacket = Json.decodeFromJsonElement<AuthResponse>(packet.payload)
                            _status.value = authPacket.status
                        }
                    }
                } catch (e: Exception) {
                    console.log(e.toString())
                }

            }
        } catch (e: Exception) {
            console.log(e.toString())
        }
    }

    fun sendAuthMessage(request: AuthInfo) {
        console.log(request)
        try {
            val jsonData = Json.encodeToString(request)
            socket?.send(jsonData)
        } catch (e: Exception) {
            console.log(e.toString())
        }
    }
}