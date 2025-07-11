plugins {
    kotlin("jvm") version "1.9.22" // ✅ safe Kotlin version for Compose Desktop
    id("org.jetbrains.compose") version "1.6.10"
}

repositories {
    google()
    mavenCentral()
    maven("https://maven.pkg.jetbrains.space/public/p/compose/dev")
}

kotlin { jvmToolchain(17) }

dependencies {
    implementation(compose.desktop.currentOs)
    implementation("io.insert-koin:koin-core:3.5.3")
    implementation("io.insert-koin:koin-compose:1.0.3")
}

compose.desktop { application { mainClass = "MainKt" } }
