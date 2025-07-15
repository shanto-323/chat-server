package com.shanto.chat_app.api

import com.shanto.chat_app.model.MessagePacket
import com.shanto.chat_app.model.Response
import kotlinx.coroutines.flow.Flow

interface ApiService{
    suspend fun Connect() : Response<Boolean>
    fun StartRead(): Flow<MessagePacket>
    suspend fun WriteMessage(msg : MessagePacket) : Response<Boolean>
}