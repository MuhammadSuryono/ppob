-keepnames class com.yonotech.ppob.mobile.data.remote.dto.RegisterRequest
-if class com.yonotech.ppob.mobile.data.remote.dto.RegisterRequest
-keep class com.yonotech.ppob.mobile.data.remote.dto.RegisterRequestJsonAdapter {
    public <init>(com.squareup.moshi.Moshi);
}
