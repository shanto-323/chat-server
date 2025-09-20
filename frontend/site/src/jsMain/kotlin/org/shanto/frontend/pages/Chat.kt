package org.shanto.frontend.pages

import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import com.varabyte.kobweb.compose.css.MaxHeight
import com.varabyte.kobweb.compose.css.Overflow
import com.varabyte.kobweb.compose.css.OverflowWrap
import com.varabyte.kobweb.compose.css.ScrollBehavior
import com.varabyte.kobweb.compose.foundation.layout.Arrangement
import com.varabyte.kobweb.compose.foundation.layout.Box
import com.varabyte.kobweb.compose.foundation.layout.Column
import com.varabyte.kobweb.compose.foundation.layout.Row
import com.varabyte.kobweb.compose.foundation.layout.Spacer
import com.varabyte.kobweb.compose.ui.Alignment
import com.varabyte.kobweb.compose.ui.Modifier
import com.varabyte.kobweb.compose.ui.graphics.Colors
import com.varabyte.kobweb.compose.ui.modifiers.aspectRatio
import com.varabyte.kobweb.compose.ui.modifiers.background
import com.varabyte.kobweb.compose.ui.modifiers.backgroundColor
import com.varabyte.kobweb.compose.ui.modifiers.border
import com.varabyte.kobweb.compose.ui.modifiers.fillMaxHeight
import com.varabyte.kobweb.compose.ui.modifiers.fillMaxSize
import com.varabyte.kobweb.compose.ui.modifiers.fillMaxWidth
import com.varabyte.kobweb.compose.ui.modifiers.height
import com.varabyte.kobweb.compose.ui.modifiers.margin
import com.varabyte.kobweb.compose.ui.modifiers.onClick
import com.varabyte.kobweb.compose.ui.modifiers.overflow
import com.varabyte.kobweb.compose.ui.modifiers.overflowWrap
import com.varabyte.kobweb.compose.ui.modifiers.padding
import com.varabyte.kobweb.compose.ui.modifiers.scrollBehavior
import com.varabyte.kobweb.compose.ui.modifiers.size
import com.varabyte.kobweb.compose.ui.modifiers.width
import com.varabyte.kobweb.core.Page
import com.varabyte.kobweb.silk.components.forms.CheckboxKind
import com.varabyte.kobweb.silk.components.forms.Input
import com.varabyte.kobweb.silk.components.icons.fa.FaIcon
import com.varabyte.kobweb.silk.components.icons.fa.IconCategory
import com.varabyte.kobweb.silk.components.icons.fa.IconSize
import org.jetbrains.compose.web.attributes.InputType
import org.jetbrains.compose.web.css.LineStyle
import org.jetbrains.compose.web.css.dpi
import org.jetbrains.compose.web.css.percent
import org.jetbrains.compose.web.css.px
import org.jetbrains.compose.web.css.vh
import org.jetbrains.compose.web.css.vw
import org.jetbrains.compose.web.dom.Text
import org.shanto.frontend.model.MessageModel
import org.shanto.frontend.model.UserModel
import org.shanto.frontend.model.messageList

@Page("/chat")
@Composable
fun ChatPage() {
    var selected by remember { mutableStateOf(false) }
    var selectedUser by remember { mutableStateOf<UserModel?>(null) }

    // Chat
    Row(
        modifier = Modifier
            .height(100.percent)
            .width(100.percent)
            .margin(10.px)
    ) {
        Column(
            modifier = Modifier
                .width(if (selected) 30.percent else 100.percent)
        ) {
            Box(
                modifier = Modifier.fillMaxWidth()
                    .height(10.percent)
                    .backgroundColor(Colors.White)
                    .margin(bottom = 8.px)
                    .border(
                        width = 2.px,
                        color = Colors.Black,
                        style = LineStyle.Solid
                    ),
                contentAlignment = Alignment.Center
            ) {
                Text("CHAT-APP1")
            }


            Box(
                modifier = Modifier
                    .onClick { selected = !selected }
                    .fillMaxWidth()
                    .height(90.percent)
                    .backgroundColor(Colors.White)
                    .border(
                        width = 2.px,
                        color = Colors.Black,
                        style = LineStyle.Solid
                    )
            ) {
                UsersListScreen(
                    userList,
                    onclick = { user ->
                        selectedUser = user
                        console.log(user.name)
                    }
                )
            }
        }

        selectedUser?.let { user ->
            if (selected) {
                Box(
                    modifier = Modifier
                        .width(70.percent)
                        .height(100.percent)
                        .backgroundColor(Colors.White)
                        .margin(left = 10.px)
                        .border(
                            width = 2.px,
                            color = Colors.Black,
                            style = LineStyle.Solid
                        )
                ) {
                    ChatScreen(
                        user = user,
                        messages = messageList
                    )
                }
            }
        }

    }
}


// Active pool
@Composable
private fun UsersListScreen(list: List<UserModel>, onclick: (UserModel) -> Unit) {
    Column(
        modifier = Modifier.fillMaxWidth()
            .height(95.vh)
            .padding(8.px)
            .overflow(Overflow.Auto),
        verticalArrangement = Arrangement.Top,
        horizontalAlignment = Alignment.Start
    ) {
        list.forEach { user ->
            Box(
                modifier = Modifier
                    .height(24.px)
                    .fillMaxWidth()
                    .onClick { onclick(user) },
                contentAlignment = Alignment.CenterStart
            ) {
                Text(user.name)
            }
        }
    }
}
@Composable
private fun ChatScreen(user: UserModel, messages: List<MessageModel>) {
    var message by remember { mutableStateOf("") }
    Column(
        modifier = Modifier.fillMaxSize().padding(10.px),
        verticalArrangement = Arrangement.Top,
        horizontalAlignment = Alignment.Start
    ) {

        // User Profile
        Box(
            modifier = Modifier
                .fillMaxWidth()
                .weight(0.05f)
                .background(Colors.White)
                .border(
                    width = 1.px,
                    color = Colors.Black,
                    style = LineStyle.Ridge
                )
                .padding(8.px)
        ) {
            Text(user.name)
        }

        // Message field
        Box(
            modifier = Modifier
                .fillMaxWidth()
                .weight(0.85f)
                .background(Colors.White)
                .overflow(Overflow.Auto)
        ) {
            Column(
                modifier = Modifier.fillMaxWidth().padding(20.px),
                verticalArrangement = Arrangement.Top,
                horizontalAlignment = Alignment.CenterHorizontally
            ) {
                messages.forEach { msg ->
                    val isCurrentUser = msg.senderId == user.id
                    Box(
                        modifier = Modifier
                            .height(24.px)
                            .fillMaxWidth()
                            .background(if (isCurrentUser) Colors.Cyan else Colors.White)
                        ,
                        contentAlignment = if (isCurrentUser) Alignment.CenterEnd else Alignment.CenterStart
                    ) {
                        Text(
                            value = msg.message,
                        )
                    }
                }
            }
        }


        // Message Writing field
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .weight(0.05f)
                .background(Colors.White),
            verticalAlignment = Alignment.CenterVertically,
            horizontalArrangement = Arrangement.SpaceAround
        ) {
            Input(
                type = InputType.Text,
                value = message,
                onValueChange = { message = it },
                modifier = Modifier
                    .fillMaxHeight()
                    .weight(1f)
                    .padding(8.px)
                    .border(
                        width = 2.px,
                        color = Colors.Black,
                        style = LineStyle.Solid
                    )
                    .padding(8.px)
            )
            Box(
                modifier = Modifier.fillMaxHeight().aspectRatio(1f),
                contentAlignment = Alignment.Center
            ) {
                FaIcon(
                    name = "paper-plane",
                    style = IconCategory.REGULAR,
                    modifier = Modifier.padding(5.px),
                    size = IconSize.XL
                )
            }
        }
    }
}

val userList = listOf<UserModel>(
    UserModel(
        name = "User.1",
        id = "1"
    ),
    UserModel(
        name = "User.2",
        id = "2"
    ),
    UserModel(
        name = "User.1",
        id = "1"
    ),
    UserModel(
        name = "User.2",
        id = "2"
    ),
    UserModel(
        name = "User.1",
        id = "1"
    ),
    UserModel(
        name = "User.2",
        id = "2"
    ),
    UserModel(
        name = "User.1",
        id = "1"
    ),
    UserModel(
        name = "User.2",
        id = "2"
    ),
    UserModel(
        name = "User.1",
        id = "1"
    ),
    UserModel(
        name = "User.2",
        id = "2"
    ),
    UserModel(
        name = "User.1",
        id = "1"
    ),
    UserModel(
        name = "User.2",
        id = "2"
    ),
    UserModel(
        name = "User.1",
        id = "1"
    ),
    UserModel(
        name = "User.2",
        id = "2"
    ),
    UserModel(
        name = "User.1",
        id = "1"
    ),
    UserModel(
        name = "User.2",
        id = "2"
    ),
    UserModel(
        name = "User.1",
        id = "1"
    ),
    UserModel(
        name = "User.2",
        id = "2"
    ),
    UserModel(
        name = "User.1",
        id = "1"
    ),
    UserModel(
        name = "User.2",
        id = "2"
    ),
    UserModel(
        name = "User.1",
        id = "1"
    ),
    UserModel(
        name = "User.2",
        id = "2"
    ),
    UserModel(
        name = "User.1",
        id = "1"
    ),
    UserModel(
        name = "User.2",
        id = "2"
    ),
    UserModel(
        name = "User.1",
        id = "1"
    ),
    UserModel(
        name = "User.2",
        id = "2"
    ),
    UserModel(
        name = "User.1",
        id = "1"
    ),
    UserModel(
        name = "User.2",
        id = "2"
    ),
    UserModel(
        name = "User.1",
        id = "1"
    ),
    UserModel(
        name = "User.2",
        id = "2"
    ),
    UserModel(
        name = "User.1",
        id = "1"
    ),
    UserModel(
        name = "User.2",
        id = "2"
    ),
    UserModel(
        name = "User.1",
        id = "1"
    ),
    UserModel(
        name = "User.2",
        id = "2"
    ),
    UserModel(
        name = "User.1",
        id = "1"
    ),
    UserModel(
        name = "User.2",
        id = "2"
    ),
    UserModel(
        name = "User.1",
        id = "1"
    ),
    UserModel(
        name = "User.2",
        id = "2"
    ),
    UserModel(
        name = "User.1",
        id = "1"
    ),
    UserModel(
        name = "User.2",
        id = "2"
    ),
    UserModel(
        name = "User.1",
        id = "1"
    ),
    UserModel(
        name = "User.2",
        id = "2"
    ),
    UserModel(
        name = "User.1",
        id = "1"
    ),
    UserModel(
        name = "User.2",
        id = "2"
    ),
    UserModel(
        name = "User.1",
        id = "1"
    ),
    UserModel(
        name = "User.2",
        id = "2"
    ),
    UserModel(
        name = "User.1",
        id = "1"
    ),
    UserModel(
        name = "User.2",
        id = "2"
    ),
    UserModel(
        name = "User.1",
        id = "1"
    ),
    UserModel(
        name = "User.2",
        id = "2"
    ),
    UserModel(
        name = "User.1",
        id = "1"
    ),
    UserModel(
        name = "User.2",
        id = "2"
    ),
    UserModel(
        name = "User.1",
        id = "1"
    ),
    UserModel(
        name = "User.2",
        id = "2"
    ),
    UserModel(
        name = "User.1",
        id = "1"
    ),
    UserModel(
        name = "User.2",
        id = "2"
    ),
    UserModel(
        name = "User.1",
        id = "1"
    ),
    UserModel(
        name = "User.2",
        id = "2"
    ),
    UserModel(
        name = "User.1",
        id = "1"
    ),
    UserModel(
        name = "User.2",
        id = "2"
    ),
    UserModel(
        name = "User.1",
        id = "1"
    ),
    UserModel(
        name = "User.2",
        id = "2"
    ),
    UserModel(
        name = "User.1",
        id = "1"
    ),
    UserModel(
        name = "User.2",
        id = "2"
    ),
    UserModel(
        name = "User.1",
        id = "1"
    ),
    UserModel(
        name = "User.2",
        id = "2"
    ),
    UserModel(
        name = "User.1",
        id = "1"
    ),
    UserModel(
        name = "User.2",
        id = "2"
    ),
    UserModel(
        name = "User.1",
        id = "1"
    ),
    UserModel(
        name = "User.2",
        id = "2"
    ),
    UserModel(
        name = "User.1",
        id = "1"
    ),
    UserModel(
        name = "User.2",
        id = "2"
    ),
)

