-keepnames class com.yonotech.ppob.mobile.data.remote.dto.StaffDto
-if class com.yonotech.ppob.mobile.data.remote.dto.StaffDto
-keep class com.yonotech.ppob.mobile.data.remote.dto.StaffDtoJsonAdapter {
    public <init>(com.squareup.moshi.Moshi);
}
