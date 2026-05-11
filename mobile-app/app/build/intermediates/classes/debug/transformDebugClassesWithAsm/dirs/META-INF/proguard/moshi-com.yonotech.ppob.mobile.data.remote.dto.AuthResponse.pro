-keepnames class com.yonotech.ppob.mobile.data.remote.dto.AuthResponse
-if class com.yonotech.ppob.mobile.data.remote.dto.AuthResponse
-keep class com.yonotech.ppob.mobile.data.remote.dto.AuthResponseJsonAdapter {
    public <init>(com.squareup.moshi.Moshi);
}
-if class com.yonotech.ppob.mobile.data.remote.dto.AuthResponse
-keepnames class kotlin.jvm.internal.DefaultConstructorMarker
-keepclassmembers class com.yonotech.ppob.mobile.data.remote.dto.AuthResponse {
    public synthetic <init>(java.lang.String,java.lang.String,boolean,com.yonotech.ppob.mobile.data.remote.dto.UserDto,int,kotlin.jvm.internal.DefaultConstructorMarker);
}
