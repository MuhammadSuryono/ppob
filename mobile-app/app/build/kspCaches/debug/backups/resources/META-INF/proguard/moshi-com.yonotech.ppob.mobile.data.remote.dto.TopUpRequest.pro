-keepnames class com.yonotech.ppob.mobile.data.remote.dto.TopUpRequest
-if class com.yonotech.ppob.mobile.data.remote.dto.TopUpRequest
-keep class com.yonotech.ppob.mobile.data.remote.dto.TopUpRequestJsonAdapter {
    public <init>(com.squareup.moshi.Moshi);
}
