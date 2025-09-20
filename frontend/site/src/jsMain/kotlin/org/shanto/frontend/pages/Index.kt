package org.shanto.frontend.pages

import androidx.compose.runtime.*
import com.varabyte.kobweb.compose.foundation.layout.Arrangement
import com.varabyte.kobweb.compose.foundation.layout.Box
import com.varabyte.kobweb.compose.foundation.layout.Column
import com.varabyte.kobweb.compose.ui.Alignment
import com.varabyte.kobweb.compose.ui.Modifier
import com.varabyte.kobweb.compose.ui.modifiers.fillMaxSize
import com.varabyte.kobweb.compose.ui.modifiers.padding
import com.varabyte.kobweb.core.Page
import com.varabyte.kobweb.core.rememberPageContext
import kotlinx.coroutines.channels.awaitClose
import kotlinx.coroutines.delay
import kotlinx.coroutines.flow.callbackFlow
import org.jetbrains.compose.web.css.px
import org.jetbrains.compose.web.css.vh
import org.jetbrains.compose.web.dom.P
import org.jetbrains.compose.web.dom.Text
import org.w3c.dom.Navigator
import org.w3c.dom.WebSocket
import com.varabyte.kobweb.navigation.Router
import com.varabyte.kobweb.silk.components.forms.Button
import org.jetbrains.compose.web.dom.Br
import org.jetbrains.compose.web.dom.H1
import org.shanto.frontend.Connection
import kotlin.js.Console

@Page
@Composable
fun HomePage() {
    val ctx = rememberPageContext()
    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(10.px),
        verticalArrangement = Arrangement.Center,
        horizontalAlignment = Alignment.CenterHorizontally
    ) {

        H1 { Text("THIS IS WELCOME PAGE, konnichiya ") }
        Button(onClick = {
            ctx.router.navigateTo("/login")
        }) {
            Text("Sign In")
        }
        Br()
        Button(onClick = {
            ctx.router.navigateTo("/signup")
        }) {
            Text("Sign Up")
        }
    }
}




