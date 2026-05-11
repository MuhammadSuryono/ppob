-keepnames class com.yonotech.ppob.mobile.data.remote.dto.LoginRequest
-if class com.yonotech.ppob.mobile.data.remote.dto.LoginRequest
-keep class com.yonotech.ppob.mobile.data.remote.dto.LoginRequestJsonAdapter {
    public <init>(com.squareup.moshi.Moshi);
}
