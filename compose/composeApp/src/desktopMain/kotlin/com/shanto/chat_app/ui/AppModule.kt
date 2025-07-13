package com.shanto.chat_app.ui

import com.shanto.chat_app.ui.screen.ScreenViewModel
import org.koin.dsl.module

val appModule = module {
    factory{
        ScreenViewModel(get())
    }
}