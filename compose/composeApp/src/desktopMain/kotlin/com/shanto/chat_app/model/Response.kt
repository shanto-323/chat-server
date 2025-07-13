package com.shanto.chat_app.model

sealed class Response<out T>{
    data class Success<out T>(val data : T): Response<T>()
    data class Error(val msg : String) : Response<Nothing>()
}