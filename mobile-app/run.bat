@echo off
REM Jalankan gradle installDebug
echo Installing debug build...
call gradlew installDebug

REM Launch aplikasi setelah install
echo Launching app...
adb shell am start -n com.yonotech.ppob.mobile/.MainActivity

echo Done!
pause
