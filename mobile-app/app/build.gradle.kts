plugins {
    alias(libs.plugins.android.application)
    alias(libs.plugins.kotlin.android)
    alias(libs.plugins.kotlin.compose)
    alias(libs.plugins.hilt.android)
    alias(libs.plugins.google.services)
    alias(libs.plugins.google.devtools.ksp)
}

android {
    namespace = "com.yonotech.ppob"

    compileSdk = 36

    defaultConfig {
        applicationId = "com.yonotech.ppob"

        minSdk = 26
        targetSdk = 36

        versionCode = 1
        versionName = "1.0.0"

        testInstrumentationRunner = "androidx.test.runner.AndroidJUnitRunner"

        buildConfigField(
            "String",
            "BASE_URL",
            "\"https://fedora.sinauplatform.id/api/v1/\""
        )

        ksp {
            arg("room.schemaLocation", "$projectDir/schemas")
            arg("room.incremental", "true")
            arg("room.generateKotlin", "true")
        }
    }

    buildTypes {

        debug {
            isMinifyEnabled = false
            isDebuggable = true

            buildConfigField(
                "String",
                "BASE_URL",
                "\"https://fedora.sinauplatform.id/api/v1/\""
            )
        }

        release {
            isMinifyEnabled = true
            isShrinkResources = true
            isDebuggable = false

            buildConfigField(
                "String",
                "BASE_URL",
                "\"https://fedora.sinauplatform.id/api/v1/\""
            )

            proguardFiles(
                getDefaultProguardFile("proguard-android-optimize.txt"),
                "proguard-rules.pro"
            )
        }
    }

    compileOptions {
        sourceCompatibility = JavaVersion.VERSION_17
        targetCompatibility = JavaVersion.VERSION_17
    }

    kotlin {
        jvmToolchain(17)
    }

    buildFeatures {
        compose = true
        buildConfig = true
    }

    packaging {
        resources {
            excludes += "/META-INF/{AL2.0,LGPL2.1}"
            excludes += "META-INF/DEPENDENCIES"
        }
    }

    testOptions {
        unitTests {
            isReturnDefaultValues = true
        }
    }
}

kotlin {
    compilerOptions {
        jvmTarget.set(
            org.jetbrains.kotlin.gradle.dsl.JvmTarget.JVM_17
        )
    }
}

dependencies {

    // =========================================
    // Compose BOM
    // =========================================
    implementation(platform(libs.androidx.compose.bom))

    // =========================================
    // Core
    // =========================================
    implementation(libs.androidx.core.ktx)
    implementation(libs.androidx.core.splashscreen)

    // ========================================
    // Worker
    // ========================================
    implementation(libs.androidx.work.runtime.ktx)

    // =========================================
    // Coroutines
    // =========================================
    implementation(libs.kotlinx.coroutines.android)
    implementation(libs.kotlinx.coroutines.play.services)

    // =========================================
    // Activity Compose
    // =========================================
    implementation(libs.androidx.activity.compose)

    // =========================================
    // Compose UI
    // =========================================
    implementation(libs.androidx.compose.ui)
    implementation(libs.androidx.compose.ui.graphics)
    implementation(libs.androidx.compose.ui.tooling.preview)

    // =========================================
    // Compose Foundation
    // =========================================
    implementation(libs.androidx.compose.foundation)
    implementation(libs.androidx.compose.animation)

    // =========================================
    // Material 3
    // =========================================
    implementation(libs.androidx.compose.material3)

    // =========================================
    // Material Icons
    // =========================================
    implementation(libs.androidx.compose.material.icons.core)
    implementation(libs.androidx.compose.material.icons.extended)

    // =========================================
    // Navigation
    // =========================================
    implementation(libs.androidx.navigation.compose)

    // =========================================
    // Lifecycle
    // =========================================
    implementation(libs.androidx.lifecycle.runtime.ktx)
    implementation(libs.androidx.lifecycle.runtime.compose)
    implementation(libs.androidx.lifecycle.viewmodel.ktx)
    implementation(libs.androidx.lifecycle.viewmodel.compose)

    // =========================================
    // Hilt
    // =========================================
    implementation(libs.hilt.android)
    implementation(libs.androidx.hilt.navigation.compose)

    ksp(libs.hilt.compiler)

    // =========================================
    // Room
    // =========================================
    implementation(libs.androidx.room.runtime)
    implementation(libs.androidx.room.ktx)

    ksp(libs.androidx.room.compiler)

    // =========================================
    // DataStore
    // =========================================
    implementation(libs.androidx.datastore.preferences)

    // =========================================
    // Networking
    // =========================================
    implementation(libs.retrofit)
    implementation(libs.retrofit.converter.moshi)

    implementation(libs.okhttp)
    implementation(libs.okhttp.logging.interceptor)

    // =========================================
    // Moshi
    // =========================================
    implementation(libs.moshi)

    ksp(libs.moshi.codegen)

    // =========================================
    // Coil
    // =========================================
    implementation(libs.coil.compose)

    // =========================================
    // Paging
    // =========================================
    implementation(libs.androidx.paging.runtime.ktx)
    implementation(libs.androidx.paging.compose)

    // =========================================
    // Biometric
    // =========================================
    implementation(libs.androidx.biometric)

    // =========================================
    // Firebase
    // =========================================
    implementation(platform(libs.firebase.bom))

    implementation(libs.firebase.analytics)
    implementation(libs.firebase.messaging)
    implementation(libs.firebase.crashlytics)

    // =========================================
    // Security
    // =========================================
    implementation(libs.security.crypto)

    // =========================================
    // Logging
    // =========================================
    implementation(libs.timber)

    // =========================================
    // Debugging
    // =========================================
    debugImplementation(libs.chucker.debug)
    releaseImplementation(libs.chucker.release)

    debugImplementation(libs.leakcanary.android)

    // =========================================
    // Testing
    // =========================================
    testImplementation(libs.junit)
    testImplementation(libs.kotlinx.coroutines.test)
    testImplementation(libs.mockk)
    testImplementation(libs.turbine)

    androidTestImplementation(libs.androidx.junit)
    androidTestImplementation(libs.androidx.espresso.core)
    androidTestImplementation(libs.androidx.compose.ui.test.junit4)

    // =========================================
    // Compose Debug
    // =========================================
    debugImplementation(libs.androidx.compose.ui.tooling)
    debugImplementation(libs.androidx.compose.ui.test.manifest)
}