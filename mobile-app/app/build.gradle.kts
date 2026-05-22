plugins {
    alias(libs.plugins.android.application)

    alias(libs.plugins.kotlin.android)
    alias(libs.plugins.compose.compiler)

    alias(libs.plugins.hilt.android)

    alias(libs.plugins.ksp)

    alias(libs.plugins.google.gms.google.services)
    alias(libs.plugins.firebase.crashlytics)
}

android {

    namespace = "com.yonotech.ppob.mobile"

    compileSdk = 35

    defaultConfig {

        applicationId = "com.yonotech.ppob.mobile"

        minSdk = 26
        targetSdk = 35

        versionCode = 1
        versionName = "1.0.0"

        testInstrumentationRunner =
            "androidx.test.runner.AndroidJUnitRunner"

        vectorDrawables {
            useSupportLibrary = true
        }
    }

    buildTypes {

        debug {

            isDebuggable = true

            // applicationIdSuffix = ".debug"
            // versionNameSuffix = "-DEBUG"
        }

        release {

            isMinifyEnabled = true
            isShrinkResources = true

            proguardFiles(
                getDefaultProguardFile(
                    "proguard-android-optimize.txt"
                ),
                "proguard-rules.pro"
            )
        }
    }

    compileOptions {

    }

    kotlinOptions {

        jvmTarget = "17"

        freeCompilerArgs += listOf(
            "-Xjvm-default=all"
        )
    }

    buildFeatures {

        compose = true

        buildConfig = true
    }

    packaging {

        resources {

            excludes.addAll(
                listOf(
                    "/META-INF/{AL2.0,LGPL2.1}",
                    "META-INF/LICENSE.md",
                    "META-INF/LICENSE-notice.md"
                )
            )
        }

        jniLibs {
            excludes += setOf(
                "**/libandroidx.graphics.path.so",
                "**/libdatastore_shared_counter.so"
            )
        }
    }

    lint {

        abortOnError = false
        checkReleaseBuilds = false
    }
}

kotlin {

    jvmToolchain(17)
}

dependencies {

    // =====================================================
    // CORE
    // =====================================================
    implementation(libs.androidx.core.ktx)

    // =====================================================
    // LIFECYCLE
    // =====================================================
    implementation(libs.androidx.lifecycle.runtime.ktx)

    implementation(
        libs.androidx.lifecycle.viewmodel.compose
    )

    // =====================================================
    // COMPOSE
    // =====================================================
    implementation(
        libs.androidx.activity.compose
    )

    implementation(
        platform(libs.androidx.compose.bom)
    )

    implementation(libs.androidx.ui)

    implementation(libs.androidx.ui.graphics)

    implementation(
        libs.androidx.ui.tooling.preview
    )

    implementation(libs.androidx.material3)

    implementation(libs.androidx.material.icons.extended)

    // =====================================================
    // NAVIGATION
    // =====================================================
    implementation(
        libs.androidx.navigation.compose
    )

    // =====================================================
    // HILT
    // =====================================================
    implementation(libs.hilt.android)

    implementation(
        libs.androidx.hilt.navigation.compose
    )

    implementation(
        libs.androidx.hilt.work
    )

    ksp(libs.hilt.compiler)

    ksp(libs.androidx.hilt.compiler)

    // =====================================================
    // RETROFIT
    // =====================================================
    implementation(libs.retrofit)

    implementation(
        libs.retrofit.converter.moshi
    )

    // =====================================================
    // OKHTTP
    // =====================================================
    implementation(libs.okhttp)

    implementation(
        libs.okhttp.logging.interceptor
    )

    // =====================================================
    // MOSHI
    // =====================================================
    implementation(libs.moshi)

    implementation(libs.moshi.kotlin)

    ksp(libs.moshi.kotlin.codegen)

    // =====================================================
    // ROOM
    // =====================================================
    implementation(libs.room.runtime)

    implementation(libs.room.ktx)

    ksp(libs.room.compiler)

    // =====================================================
    // WORK MANAGER
    // =====================================================
    implementation(
        libs.androidx.work.runtime.ktx
    )

    // =====================================================
    // DATASTORE
    // =====================================================
    implementation(
        libs.androidx.datastore.preferences
    )

    // =====================================================
    // COIL
    // =====================================================
    implementation(libs.coil.compose)

    // =====================================================
    // BIOMETRIC
    // =====================================================
    implementation(libs.androidx.biometric)

    // =====================================================
    // FIREBASE
    // =====================================================
    implementation(
        platform(libs.firebase.bom)
    )

    implementation(libs.firebase.analytics)

    implementation(libs.firebase.messaging)

    implementation(libs.firebase.crashlytics)

    // =====================================================
    // LOGGER
    // =====================================================
    implementation(libs.timber)

    // =====================================================
    // CHUCKER
    // =====================================================
    debugImplementation(libs.chucker.debug)

    releaseImplementation(libs.chucker.release)

    // =====================================================
    // TEST
    // =====================================================
    testImplementation(libs.junit)

    androidTestImplementation(
        libs.androidx.junit
    )

    androidTestImplementation(
        libs.androidx.espresso.core
    )

    androidTestImplementation(
        platform(libs.androidx.compose.bom)
    )

    androidTestImplementation(
        libs.androidx.ui.test.junit4
    )

    debugImplementation(
        libs.androidx.ui.tooling
    )

    debugImplementation(
        libs.androidx.ui.test.manifest
    )
}
