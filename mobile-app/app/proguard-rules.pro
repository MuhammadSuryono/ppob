/*
 * ProGuard rules for PPOB Mobile Application
 */

# Keep our application class and entry points
-keep class com.yonotech.ppob.** { *; }
-keep class com.yonotech.ppob.PPOBApp { *; }

# Keep ViewModels (Hilt / Android lifecycle)
-keep class * extends androidx.lifecycle.ViewModel { *; }
-keepclassmembers class * {
    @androidx.lifecycle.ViewModelProvider.NewInstanceFactory *;
}

# Keep Room entities
-keep @androidx.room.Entity class * { *; }
-keep @androidx.room.Dao class * { *; }
-keep @androidx.room.Database class * { *; }
-keep @androidx.room.TypeConverter class * { *; }

# Keep Moshi models
-keepattributes *Annotation*
-keep class com.yonotech.ppob.data.remote.model.** { *; }
-keep class com.yonotech.ppob.domain.model.** { *; }
-keep class com.squareup.moshi.** { *; }

# Keep Retrofit interfaces
-keepattributes Signature
-keepattributes Exceptions
-keepclasseswithmembers class * {
    @retrofit2.http.* <methods>;
}

# Keep Hilt generated code
-keep class dagger.hilt.** { *; }
-keep class **_Hilt_*.class { *; }
-keep class **$$HiltGeneratedPredicateProvider { *; }
-keepclasseswithmembers,includedescriptorclasses class * {
    public <init>(...);
}
-keepnames @javax.inject.Inject class *
-keepclassmembers class * {
    @javax.inject.Inject *;
    @dagger.Provides *;
    @dagger.hilt.InstallIn *;
}

# Keep Firebase
-keep class com.google.firebase.** { *; }
-dontwarn com.google.firebase.**

# Keep WorkManager
-keep class androidx.work.** { *; }
-dontwarn androidx.work.**

# Keep Coil (image loading)
-keep class coil.** { *; }
-keep class coil3.** { *; }

# Keep Compose runtime
-keep class androidx.compose.** { *; }
-keep class kotlinx.coroutines.** { *; }

# Keep Kotlin metadata
-keepattributes RuntimeVisibleAnnotations
-keep class kotlin.Metadata { *; }

# Keep Parcelable
-keep class * implements android.os.Parcelable {
    public static final android.os.Parcelable$Creator *;
}

# Keep Serializable
-keepclassmembers class * implements java.io.Serializable {
    static final long serialVersionUID;
    private static final java.io.ObjectStreamField[] serialPersistentFields;
    private void writeObject(java.io.ObjectOutputStream);
    private void readObject(java.io.ObjectInputStream);
    java.lang.Object writeReplace();
    java.lang.Object readResolve();
}

# Keep enum classes
-keepclassmembers enum * {
    public static **[] values();
    public static ** valueOf(java.lang.String);
}

# SQLCipher
-keep class net.sqlcipher.** { *; }
-keep class net.zetetic.** { *; }

# DataStore
-keep class androidx.datastore.** { *; }

# Biometric
-keep class androidx.biometric.** { *; }

# OkHttp / Interceptors
-keep class okhttp3.** { *; }
-keep interface okhttp3.** { *; }
-dontwarn okhttp3.**

# Retrofit
-keep class retrofit2.** { *; }
-dontwarn retrofit2.**

# Logging interceptor - keep in debug, can be stripped in release
-assumenosideeffects class android.util.Log {
    public static boolean isLoggable(java.lang.String, int);
    public static int v(...);
    public static int i(...);
    public static int w(...);
    public static int d(...);
    public static int e(...);
}