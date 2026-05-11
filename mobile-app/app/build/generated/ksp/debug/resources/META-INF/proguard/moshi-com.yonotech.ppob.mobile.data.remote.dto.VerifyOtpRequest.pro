-keepnames class com.yonotech.ppob.mobile.data.remote.dto.VerifyOtpRequest
-if class com.yonotech.ppob.mobile.data.remote.dto.VerifyOtpRequest
-keep class com.yonotech.ppob.mobile.data.remote.dto.VerifyOtpRequestJsonAdapter {
    public <init>(com.squareup.moshi.Moshi);
}
