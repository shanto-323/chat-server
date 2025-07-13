package com.shanto.chat_app.api
import io.ktor.client.*
import io.ktor.client.engine.cio.*
import io.ktor.client.plugins.websocket.*
import org.koin.dsl.module

val apiModule = module {
    single {
        HttpClient(CIO){
            install(WebSockets){
                pingInterval = 15000
                maxFrameSize = 1024
            }
        }
    }

    single<ApiService> {
        ApiServiceImpl(get())
    }
}