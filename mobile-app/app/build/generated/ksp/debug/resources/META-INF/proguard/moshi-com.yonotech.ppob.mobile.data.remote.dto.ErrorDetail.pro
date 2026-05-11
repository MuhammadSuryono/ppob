-keepnames class com.yonotech.ppob.mobile.data.remote.dto.ErrorDetail
-if class com.yonotech.ppob.mobile.data.remote.dto.ErrorDetail
-keep class com.yonotech.ppob.mobile.data.remote.dto.ErrorDetailJsonAdapter {
    public <init>(com.squareup.moshi.Moshi);
}
