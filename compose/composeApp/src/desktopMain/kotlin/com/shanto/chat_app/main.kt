package com.shanto.chat_app

import androidx.compose.ui.window.Window
import androidx.compose.ui.window.application
import com.shanto.chat_app.api.apiModule
import com.shanto.chat_app.ui.appModule
import com.shanto.chat_app.ui.screen.Screen
import com.shanto.chat_app.ui.screen.ScreenViewModel
import org.koin.compose.koinInject
import org.koin.core.context.startKoin

fun main() = application {
    startKoin {
        modules(
            apiModule,appModule
        )
    }
    val viewModel = koinInject<ScreenViewModel> ()
    Window(
        onCloseRequest = ::exitApplication,
        title = "Chat-Application",
    ) {
        Screen(viewModel)
    }
}