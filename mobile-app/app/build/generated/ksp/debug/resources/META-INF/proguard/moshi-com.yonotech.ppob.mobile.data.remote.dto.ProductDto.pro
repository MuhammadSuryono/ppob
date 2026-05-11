-keepnames class com.yonotech.ppob.mobile.data.remote.dto.ProductDto
-if class com.yonotech.ppob.mobile.data.remote.dto.ProductDto
-keep class com.yonotech.ppob.mobile.data.remote.dto.ProductDtoJsonAdapter {
    public <init>(com.squareup.moshi.Moshi);
}
-if class com.yonotech.ppob.mobile.data.remote.dto.ProductDto
-keepnames class kotlin.jvm.internal.DefaultConstructorMarker
-keepclassmembers class com.yonotech.ppob.mobile.data.remote.dto.ProductDto {
    public synthetic <init>(java.lang.String,java.lang.String,java.lang.String,java.lang.String,java.lang.String,double,java.lang.String,java.lang.String,int,kotlin.jvm.internal.DefaultConstructorMarker);
}
