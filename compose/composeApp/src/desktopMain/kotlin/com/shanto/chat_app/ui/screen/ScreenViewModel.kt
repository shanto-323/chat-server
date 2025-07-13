package com.shanto.chat_app.ui.screen

import com.shanto.chat_app.api.ApiService
import com.shanto.chat_app.model.Response
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.SupervisorJob
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch

class ScreenViewModel(
    private val apiService: ApiService
){
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
                    _response.value = Event.Loading
                }
                is Response.Success<Boolean> ->{
                    _response.value = Event.Success
                }
            }
        }
    }
}

sealed class Event{
    data object Loading: Event()
    data object Success: Event()
    data object  Failure: Event()
}