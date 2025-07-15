package com.shanto.chat_app.ui.screen

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material.Button
import androidx.compose.material.CircularProgressIndicator
import androidx.compose.material.Icon
import androidx.compose.material.Text
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.unit.dp
import com.shanto.chat_app.model.ActivePool
import com.shanto.chat_app.model.Client
import kotlinx.coroutines.delay
import kotlinx.serialization.json.Json
import kotlinx.serialization.json.decodeFromJsonElement

@Composable
fun Screen(
    viewModel: ScreenViewModel
) {
    val state = viewModel.state.collectAsState().value
    val response = viewModel.response.collectAsState().value

    var open by remember { mutableStateOf(false) }

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
            Row(
                modifier = Modifier
                    .fillMaxSize()
                    .background(Color.White),
                verticalAlignment = Alignment.Top,
            ) {
                LazyColumn(
                    modifier = Modifier
                        .fillMaxHeight()
                        .fillMaxWidth(0.3F)
                        .padding(10.dp),
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
                            id = it
                        )
                    }
                }

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
                    open = open
                )

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
            .aspectRatio(2f)
            .padding(5.dp)
    ) {
        Text(id)
    }
}

@Composable
fun ChatCard(
    modifier: Modifier = Modifier,
    id: String?,
    onClick: () -> Unit = {},
    open :Boolean
) {
    LazyColumn(
        modifier = Modifier
            .then(modifier)
            .padding(5.dp)
            .background(Color.Gray),
        horizontalAlignment = Alignment.Start,
        verticalArrangement = Arrangement.Top
    ) {
        if (open && id != null) {
            item {
                Row(
                    modifier = Modifier
                        .fillMaxSize()
                        .background(Color.White),
                    verticalAlignment = Alignment.CenterVertically,
                    horizontalArrangement = Arrangement.End
                ) {
                    Text(id)
                    Button(
                        onClick = {
                            onClick.invoke()
                        }
                    ) {
                        Text("close")
                    }
                }
            }
            items(20) {
                Text("message")
            }
        }
    }
}