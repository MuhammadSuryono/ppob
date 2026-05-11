-keepnames class com.yonotech.ppob.mobile.data.remote.dto.ErrorResponse
-if class com.yonotech.ppob.mobile.data.remote.dto.ErrorResponse
-keep class com.yonotech.ppob.mobile.data.remote.dto.ErrorResponseJsonAdapter {
    public <init>(com.squareup.moshi.Moshi);
}
