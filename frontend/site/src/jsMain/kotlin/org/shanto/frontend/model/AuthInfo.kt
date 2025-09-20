package org.shanto.frontend.model
import kotlinx.serialization.Serializable
import kotlinx.serialization.json.JsonElement

@Serializable
data class AuthInfo(
    val method : String = "",
    val username: String = "",
    val password: String = ""
)


@Serializable
data class AuthResponse(
    val status: Boolean = false,
)

@Serializable
data class Envelope(
    val type: String,
    val payload: JsonElement
)