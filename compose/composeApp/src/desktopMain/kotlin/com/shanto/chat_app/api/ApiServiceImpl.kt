package com.shanto.chat_app.api
import com.shanto.chat_app.model.MessagePacket
import com.shanto.chat_app.model.Response
import io.ktor.client.*
import io.ktor.client.plugins.websocket.*
import io.ktor.websocket.*
import kotlinx.coroutines.flow.*
import kotlinx.serialization.encodeToString
import kotlinx.serialization.json.Json

class  ApiServiceImpl(
    private  val client : HttpClient
): ApiService{
    private val url = "ws://localhost:8080/ws"
    private var socketSession : WebSocketSession? = null
    override suspend fun Connect(): Response<Boolean> = try {
        socketSession = client.webSocketSession(url)
        Response.Success(true)
    }catch (e : Exception){
        Response.Error(e.toString())
    }

    override  fun StartRead(): Flow<MessagePacket> = flow{
        try {
           val session = socketSession ?: throw InternalError("session is nil")
            session .incoming.consumeAsFlow().filterIsInstance<Frame.Text>().mapNotNull {
                try {
                    Json.decodeFromString<MessagePacket>(it.readText())
                }catch (e:Exception){
                    throw e
                }
            }.collect{
                emit(it)
            }
        }catch (e : Exception){
            print(e.message)
        }finally {
            socketSession?.close()
            socketSession = null
        }
    }

    override suspend fun WriteMessage(msg: MessagePacket): Response<Boolean> = try {
        val message = Json.encodeToString(msg)
        socketSession?.outgoing?.send(
            Frame.Text(message)
        )
        Response.Success(true)
    } catch (e : Exception){
        Response.Error(e.toString())
    }
}