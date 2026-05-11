-keepnames class com.yonotech.ppob.mobile.data.remote.dto.WalletResponse
-if class com.yonotech.ppob.mobile.data.remote.dto.WalletResponse
-keep class com.yonotech.ppob.mobile.data.remote.dto.WalletResponseJsonAdapter {
    public <init>(com.squareup.moshi.Moshi);
}
-if class com.yonotech.ppob.mobile.data.remote.dto.WalletResponse
-keepnames class kotlin.jvm.internal.DefaultConstructorMarker
-keepclassmembers class com.yonotech.ppob.mobile.data.remote.dto.WalletResponse {
    public synthetic <init>(java.lang.String,double,double,double,java.lang.String,int,kotlin.jvm.internal.DefaultConstructorMarker);
}
