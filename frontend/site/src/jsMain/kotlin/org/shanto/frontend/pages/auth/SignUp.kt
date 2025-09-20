package org.shanto.frontend.pages.auth

import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import com.varabyte.kobweb.compose.foundation.layout.Arrangement
import com.varabyte.kobweb.compose.foundation.layout.Column
import com.varabyte.kobweb.compose.ui.Alignment
import com.varabyte.kobweb.compose.ui.Modifier
import com.varabyte.kobweb.compose.ui.modifiers.fillMaxSize
import com.varabyte.kobweb.compose.ui.modifiers.padding
import com.varabyte.kobweb.core.Page
import com.varabyte.kobweb.core.rememberPageContext
import com.varabyte.kobweb.silk.components.forms.Button
import com.varabyte.kobweb.silk.components.forms.Input
import org.jetbrains.compose.web.attributes.InputType
import org.jetbrains.compose.web.css.px
import org.jetbrains.compose.web.dom.Br
import org.jetbrains.compose.web.dom.H1
import org.jetbrains.compose.web.dom.P
import org.jetbrains.compose.web.dom.Text
import org.shanto.frontend.Connection
import org.shanto.frontend.model.AuthInfo

@Page("/signup")
@Composable
fun SignUpPage() {
    val ctx = rememberPageContext()

    var username by remember { mutableStateOf("") }
    var password by remember { mutableStateOf("") }

    val status = Connection.status.collectAsState().value
    LaunchedEffect(status) {
        if (status) {
            ctx.router.navigateTo("/chat")
        }
    }

    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(10.px),
        verticalArrangement = Arrangement.Center,
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        H1 { Text("Sign Up") }
        Input(
            type = InputType.Text,
            value = username,
            onValueChange = { username = it },
            placeholder = "Username"
        )
        Br()
        Input(
            type = InputType.Password,
            value = password,
            onValueChange = { password = it },
            placeholder = "Password"
        )
        Br()
        Button(
            onClick = {
                val auth = AuthInfo(
                    method = "signup",
                    username = username,
                    password = password
                )
                Connection.sendAuthMessage(auth)
            }
        ) {
            Text("Submit")
        }
    }
}