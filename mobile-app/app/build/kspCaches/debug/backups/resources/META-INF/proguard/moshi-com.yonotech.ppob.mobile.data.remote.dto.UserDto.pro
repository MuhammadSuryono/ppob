-keepnames class com.yonotech.ppob.mobile.data.remote.dto.UserDto
-if class com.yonotech.ppob.mobile.data.remote.dto.UserDto
-keep class com.yonotech.ppob.mobile.data.remote.dto.UserDtoJsonAdapter {
    public <init>(com.squareup.moshi.Moshi);
}
