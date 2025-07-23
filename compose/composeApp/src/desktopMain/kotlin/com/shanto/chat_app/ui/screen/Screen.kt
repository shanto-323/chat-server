package com.shanto.chat_app.ui.screen

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.TextStyle
import androidx.compose.ui.text.font.FontStyle
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.shanto.chat_app.model.ActivePool
import com.shanto.chat_app.model.Client
import kotlinx.coroutines.delay
import kotlinx.serialization.json.Json
import kotlinx.serialization.json.decodeFromJsonElement
import org.jetbrains.compose.ui.tooling.preview.Preview

@Composable
fun Screen(
    viewModel: ScreenViewModel
) {
    val state = viewModel.state.collectAsState().value
    val response = viewModel.response.collectAsState().value

    var open by remember { mutableStateOf(false) }
    var message by remember { mutableStateOf("") }

    LaunchedEffect(Unit) {
        while (true) {
            viewModel.sendList()
            delay(10000)
        }
    }


    when (response) {
        Event.Failure -> {
            print("Error")
        }

        Event.Loading -> {
            Column(
                modifier = Modifier
                    .fillMaxSize()
                    .background(Color.White),
                horizontalAlignment = Alignment.CenterHorizontally,
                verticalArrangement = Arrangement.Center
            ) {
                CircularProgressIndicator(
                    color = Color.Black,
                    strokeWidth = 2.dp
                )
            }
        }

        Event.Success -> {
            Column(
                modifier = Modifier
                    .fillMaxSize(),
                verticalArrangement = Arrangement.Top
            ) {

                Row(
                    modifier = Modifier
                        .fillMaxWidth()
                        .fillMaxHeight(0.1F)
                        .padding(2.dp)
                        .background(Color.Gray),
                    verticalAlignment = Alignment.CenterVertically,
                    horizontalArrangement = Arrangement.spacedBy(2.dp)
                ) {
                    Box(
                        modifier = Modifier
                            .fillMaxHeight()
                            .weight(1f),
                        contentAlignment = Alignment.Center
                    ) {
                        Text("CHAT APP")
                    }
                    Text(
                        text = "Icon",
                        modifier = Modifier
                            .fillMaxHeight()
                            .aspectRatio(1f),
                        style = TextStyle(
                            textAlign = TextAlign.Center
                        )
                    )
                    Text(
                        text = "Icon",
                        modifier = Modifier
                            .fillMaxHeight()
                            .aspectRatio(1f),
                        style = TextStyle(
                            textAlign = TextAlign.Center
                        )
                    )
                }

                Row(
                    modifier = Modifier
                        .fillMaxWidth()
                        .weight(1f)
                        .padding(2.dp)
                        .background(Color.White),
                    verticalAlignment = Alignment.Top,
                ) {
                    // Active List...
                    Column(
                        modifier = Modifier
                            .fillMaxHeight()
                            .fillMaxWidth(0.3F)
                    ) {
                        Box(
                            modifier = Modifier
                                .fillMaxWidth()
                                .fillMaxHeight(0.05F)
                                .background(Color.Gray),
                            contentAlignment = Alignment.Center
                        ) {
                            Text("ACTIVE LIST")
                        }
                        LazyColumn(
                            modifier = Modifier
                                .fillMaxWidth()
                                .weight(1f)
                                .padding(horizontal = 5.dp, vertical = 0.dp),
                            horizontalAlignment = Alignment.Start,
                            verticalArrangement = Arrangement.Top
                        ) {
                            items(state.ActivePool) {
                                UserCard(
                                    onClick = {
                                        viewModel.ReceiverIdChanged(
                                            id = it
                                        )
                                        open = true
                                    },
                                    id = it,
                                )
                            }
                        }
                    }

                    // Chat Box
                    ChatCard(
                        id = state.ReceiverId,
                        modifier = Modifier
                            .fillMaxHeight()
                            .weight(1f),
                        onClick = {
                            viewModel.ReceiverIdChanged(
                                id = null
                            )
                            open = false
                        },
                        open = open,
                        onValueChange = {
                            message = it
                        },
                        message = message
                    )

                }

            }
        }
    }
}


@Composable
fun UserCard(
    onClick: () -> Unit = {},
    id: String = "Id"
) {
    Button(
        onClick = onClick,
        modifier = Modifier
            .fillMaxWidth()
            .aspectRatio(4f)
            .padding(5.dp),
        colors = ButtonDefaults.buttonColors(
            backgroundColor = Color.Black,
            contentColor = Color.White
        )
    ) {
        Text(
            text = id,
            textAlign = TextAlign.Center,
            fontSize = 10.sp,
            fontWeight = FontWeight.Bold,
            maxLines = 1
        )
    }
}

@Composable
fun ChatCard(
    modifier: Modifier = Modifier,
    id: String?,
    onClick: () -> Unit = {},
    open: Boolean,
    onValueChange: (String) -> Unit = {},
    message : String
) {
    Column(
        modifier = Modifier
            .then(modifier)
            .padding(10.dp)
            .background(Color.White),
    ) {
        if (open && id != null) {
            Row(
                modifier = Modifier
                    .fillMaxWidth()
                    .height(32.dp),
                verticalAlignment = Alignment.CenterVertically,
                horizontalArrangement = Arrangement.Center
            ) {
                Text(
                    text = id,
                    modifier = Modifier.padding(end = 2.dp).fillMaxHeight().weight(1f)
                )
                Button(
                    onClick = { onClick.invoke() },
                    modifier = Modifier.padding(start = 2.dp).fillMaxHeight().aspectRatio(1f),
                    colors = ButtonDefaults.buttonColors(
                        backgroundColor = Color.Black,
                        contentColor = Color.White
                    )
                ) {
                    Text("Close")
                }
            }

            LazyColumn(
                modifier = Modifier
                    .fillMaxWidth()
                    .weight(1f),
                horizontalAlignment = Alignment.Start,
                verticalArrangement = Arrangement.Bottom
            ) {
                items(20) {
                    MessageCard("Message!!")
                }
            }

            Row(
                modifier = Modifier
                    .fillMaxWidth()
                    .height(50.dp),
                verticalAlignment = Alignment.CenterVertically,
                horizontalArrangement = Arrangement.Center
            ) {
                OutlinedTextField(
                    value = message,
                    onValueChange = onValueChange,
                    modifier = Modifier.weight(1f),
                    placeholder = {
                        Text(
                            text = "write something....",
                            modifier = Modifier
                                .fillMaxSize(),
                            style = TextStyle(
                                textAlign = TextAlign.Start,
                                fontSize = 12.sp
                            )
                        )
                    },
                )
                Button(
                    onClick = {},
                    modifier = Modifier.padding(start = 2.dp).fillMaxHeight().aspectRatio(1f),
                    colors = ButtonDefaults.buttonColors(
                        backgroundColor = Color.Black,
                        contentColor = Color.White
                    )
                ) {
                    Text(
                        text = "Sned",
                        modifier = Modifier
                            .fillMaxSize(),
                        style = TextStyle(
                            textAlign = TextAlign.Center,
                            fontSize = 12.sp
                        )
                    )
                }
            }
        }
    }
}

@Composable
fun MessageCard(
    text: String,
    color: Color = Color.Gray
) {
    Box(
        modifier = Modifier.padding(2.dp).clip(RoundedCornerShape(4.dp)).padding(2.dp).background(color),
        contentAlignment = Alignment.Center
    ) {
        Text(
            text = text,
            textAlign = TextAlign.Center,
            color = Color.Black,
            fontSize = 12.sp,
            fontWeight = FontWeight.SemiBold,
            modifier = Modifier.padding(5.dp)
        )
    }
}