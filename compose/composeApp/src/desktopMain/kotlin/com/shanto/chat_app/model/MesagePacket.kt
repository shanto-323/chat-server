package com.shanto.chat_app.model
import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable
import kotlinx.serialization.json.JsonElement

@Serializable
data class MessagePacket(
        @SerialName("type") val msgType: String,
        @SerialName("sender_id") val senderId: String? = null,
        @SerialName("receiver_id") val receiverId: String? = null,
        @SerialName("payload") val payload: JsonElement? = null, // Raw JSON
        @SerialName("timestamp") val timestamp: Long? = null
)

@Serializable data class ActivePool(@SerialName("alive_list") val aliveList: List<String>)

@Serializable data class Client(@SerialName("id") val id: String)