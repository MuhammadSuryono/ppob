-keepnames class com.yonotech.ppob.mobile.data.remote.dto.CategoryDto
-if class com.yonotech.ppob.mobile.data.remote.dto.CategoryDto
-keep class com.yonotech.ppob.mobile.data.remote.dto.CategoryDtoJsonAdapter {
    public <init>(com.squareup.moshi.Moshi);
}
-if class com.yonotech.ppob.mobile.data.remote.dto.CategoryDto
-keepnames class kotlin.jvm.internal.DefaultConstructorMarker
-keepclassmembers class com.yonotech.ppob.mobile.data.remote.dto.CategoryDto {
    public synthetic <init>(java.lang.String,java.lang.String,java.lang.String,int,kotlin.jvm.internal.DefaultConstructorMarker);
}
