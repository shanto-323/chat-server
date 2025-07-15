package com.shanto.chat_app.ui.screen

import com.shanto.chat_app.api.ApiService
import com.shanto.chat_app.model.ActivePool
import com.shanto.chat_app.model.Client
import com.shanto.chat_app.model.Response
import com.shanto.chat_app.model.MessagePacket
import kotlinx.coroutines.*
import kotlinx.coroutines.channels.Channel
import kotlinx.coroutines.flow.*
import kotlinx.serialization.json.Json
import kotlinx.serialization.json.decodeFromJsonElement

class ScreenViewModel(
    private val apiService: ApiService
){
    private val _state = MutableStateFlow<State>(State())
    val state = _state.asStateFlow()
    private val _response = MutableStateFlow<Event>(Event.Loading)
    val response = _response.asStateFlow()

    private val scope = CoroutineScope(Dispatchers.IO + SupervisorJob())

    init {
       connect()
    }

    private fun connect(){
        scope.launch {
            _response.value = Event.Loading
            val conn = apiService.Connect()
            when(conn){
                is Response.Error -> {
                    //_response.value = Event.Loading
                    // TEST
                    _response.value = Event.Success
                }
                is Response.Success<Boolean> ->{
                    _response.value = Event.Success
                    _state.value = _state.value.copy(
                        Connected = true
                    )
                    collectPacket()
                }
            }
        }
    }

    private  fun collectPacket(){
        apiService.StartRead()
            .onEach {
                println("collecting data!!")
                scope.launch {
                   incommingMessage(it)
                }
            }
            .catch {
                println(it)
            }
            .launchIn(scope)
    }

    private fun incommingMessage(msg : MessagePacket){
        when(msg.msgType){
             "info" -> {
                 msg.payload?.let {
                     val client = Json.decodeFromJsonElement<Client>(it)
                     _state.value = state.value.copy(
                         UserId = client.id
                     )
                 }
             }
            "list" -> {
                msg.payload?.let {
                    val activePool = Json.decodeFromJsonElement<ActivePool>(it)
                    _state.value = state.value.copy(
                        ActivePool = activePool.aliveList
                    )
                }
                println(state.value.ActivePool)
            }
        }
    }

    fun sendList(){
        scope.launch {
            try {
                apiService.WriteMessage(
                    MessagePacket(
                        msgType = "list",
                        senderId = ""
                    )
                )
                println("Packet sent !!")
            }catch (e:Exception){
                println(e)
            }
            delay(10000)
        }
    }

    fun ReceiverIdChanged(id : String?){
        _state.value = _state.value.copy(
            ReceiverId = id
        )
    }
}

sealed class Event{
    data object Loading: Event()
    data object Success: Event()
    data object  Failure: Event()
}

data class State(
    val Connected : Boolean = false,
    val UserId : String = "2zuI55oWXISVs2MixhgyCnw9oDC",
    val ActivePool : List<String> = listOf(
        "2zuJw9oQmVrgFvDT4vMEJ30qfjs",
        "2zuI55oWXISVs2MixhgyCnw9oDC",
        "2zuJw9oQmVrgFvDT4vMEJ30qfjs",
        "2zuJw9oQmVrgFvDT4vMEJ30qfjs",
        "2zuI55oWXISVs2MixhgyCnw9oDC"
    ),
    val ReceiverId : String? = null
)


