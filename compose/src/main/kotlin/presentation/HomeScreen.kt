import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.material.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color

@Composable
fun HomeScreen() {
    var toggle by remember { mutableStateOf(false) }

    Surface(modifier = Modifier.fillMaxSize()) {
        Row(modifier = Modifier.fillMaxSize()) {
            Column(modifier = Modifier.fillMaxHeight().fillMaxWidth(if (toggle) 0.3f else 1f)) {
                Button(onClick = { toggle = !toggle }) { Text("Chat Slide") }
            }

            if (toggle) {
                LazyColumn(
                        modifier = Modifier.fillMaxSize().background(Color.Green),
                        horizontalAlignment = Alignment.Start,
                        verticalArrangement = Arrangement.Bottom
                ) { items(200) { Text("NEW APP") } }
            }
        }
    }
}
