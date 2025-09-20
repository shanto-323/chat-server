package org.shanto.frontend.model


import kotlin.time.Clock

data class MessageModel(
    var message: String,
    var senderId: String,
    var receiverId: String,
    var time: Double
)

val messageList = listOf(
    MessageModel("Hello there", "1", "2", 0.00),
    MessageModel("Hi, how are you?", "2", "1", 0.00),
    MessageModel("I'm doing great, thanks for asking", "1", "2", 0.00),
    MessageModel("What about you?", "1", "2", 0.00),
    MessageModel("I'm good too", "2", "1", 0.00),
    MessageModel("Any plans for today?", "2", "1", 0.00),
    MessageModel("Just working on some code", "1", "2", 0.00),
    MessageModel("Same here", "2", "1", 0.00),
    MessageModel("This is message 9", "1", "2", 0.00),
    MessageModel("And this is message 10", "2", "1", 0.00),
    MessageModel("Let's keep going", "1", "2", 0.00),
    MessageModel("Sure thing", "2", "1", 0.00),
    MessageModel("Message number 13", "1", "2", 0.00),
    MessageModel("Fourteen already", "2", "1", 0.00),
    MessageModel("Halfway to thirty", "1", "2", 0.00),
    MessageModel("Time flies", "2", "1", 0.00),
    MessageModel("Another message here", "1", "2", 0.00),
    MessageModel("Responding to that", "2", "1", 0.00),
    MessageModel("Almost there", "1", "2", 0.00),
    MessageModel("Just a few more", "2", "1", 0.00),
    MessageModel("Message twenty-one", "1", "2", 0.00),
    MessageModel("Twenty-two coming through", "2", "1", 0.00),
    MessageModel("Still going strong", "1", "2", 0.00),
    MessageModel("No stopping now", "2", "1", 0.00),
    MessageModel("Quarter century mark", "1", "2", 0.00),
    MessageModel("Twenty-sixth message", "2", "1", 0.00),
    MessageModel("Getting closer to fifty", "1", "2", 0.00),
    MessageModel("Twenty-eight and counting", "2", "1", 0.00),
    MessageModel("Almost thirty now", "1", "2", 0.00),
    MessageModel("Thirty messages done", "2", "1", 0.00),
    MessageModel("Thirty-one and still going", "1", "2", 0.00),
    MessageModel("Thirty-two messages sent", "2", "1", 0.00),
    MessageModel("This is getting long", "1", "2", 0.00),
    MessageModel("But we continue", "2", "1", 0.00),
    MessageModel("Thirty-five messages now", "1", "2", 0.00),
    MessageModel("Thirty-six completed", "2", "1", 0.00),
    MessageModel("Getting near the end", "1", "2", 0.00),
    MessageModel("Thirty-eight and counting", "2", "1", 0.00),
    MessageModel("Almost forty messages", "1", "2", 0.00),
    MessageModel("Forty exactly", "2", "1", 0.00),
    MessageModel("Forty-one messages done", "1", "2", 0.00),
    MessageModel("Forty-two completed", "2", "1", 0.00),
    MessageModel("Only a few left", "1", "2", 0.00),
    MessageModel("Forty-four messages", "2", "1", 0.00),
    MessageModel("Forty-five and counting", "1", "2", 0.00),
    MessageModel("Almost there now", "2", "1", 0.00),
    MessageModel("Forty-seven messages", "1", "2", 0.00),
    MessageModel("Forty-eight done", "2", "1", 0.00),
    MessageModel("Second to last message", "1", "2", 0.00),
    MessageModel("Fiftieth and final message", "2", "1", 0.00)
)
