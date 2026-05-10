package com.yonotech.ppob

import org.junit.Test

/**
 * Example local unit test, which will execute on the development machine (host).
 *
 * See [testing documentation](http://d.android.com/tools/testing).
 */
class ExampleUnitTest {
    @Test
    fun addition_isCorrect() {
        assertEquals(4, 2 + 2)
    }

    @Test
    fun testPinValidation() {
        // Valid PINs
        assert(isValidPin("123456")) { "6-digit PIN should be valid" }
        assert(isValidPin("654321")) { "6-digit PIN should be valid" }

        // Invalid PINs
        assert(!isValidPin("12345")) { "5-digit PIN should be invalid" }
        assert(!isValidPin("1234567")) { "7-digit PIN should be invalid" }
        assert(!isValidPin("111111")) { "Repeated digits PIN should be invalid" }
        assert(!isValidPin("123456")) // Note: 123456 is sequential - blocked
        assert(!isValidPin("")) { "Empty PIN should be invalid" }
        assert(!isValidPin("abcdef")) { "Non-numeric PIN should be invalid" }
    }

    @Test
    fun testPhoneNumberValidation() {
        assert(isValidPhoneNumber("+6281234567890")) { "Valid Indonesian phone" }
        assert(!isValidPhoneNumber("081234567890")) { "Missing country code" }
        assert(!isValidPhoneNumber("+62812")) { "Too short" }
        assert(!isValidPhoneNumber("")) { "Empty phone" }
    }

    @Test
    fun testPasswordValidation() {
        assert(isValidPassword("StrongPass1!")) { "Valid password" }
        assert(!isValidPassword("weak")) { "Too short" }
        assert(!isValidPassword("alllowercase1!")) { "No uppercase" }
        assert(!isValidPassword("ALLUPPERCASE1!")) { "No lowercase" }
        assert(!isValidPassword("NoDigitsHere!!")) { "No digits" }
        assert(!isValidPassword("NoSpecial123")) { "No special char" }
    }

    @Test
    fun testIdempotencyKeyGeneration() {
        val key1 = generateIdempotencyKey()
        val key2 = generateIdempotencyKey()
        assert(key1.isNotEmpty()) { "Key should not be empty" }
        assert(key1 != key2) { "Keys should be unique" }
        assert(key1.length == 36) { "Should be UUID format" }
    }

    @Test
    fun testCurrencyFormatting() {
        assertEquals("Rp 1.000", formatCompact(1000.0))
        assertEquals("Rp 1.250.000", formatCompact(1250000.0))
        assertEquals("Rp 50.000", formatCompact(50000.0))
    }

    private fun isValidPin(pin: String): Boolean {
        if (pin.length != 6) return false
        if (!pin.all { it.isDigit() }) return false
        if (pin.all { it == pin[0] }) return false
        if (pin == "123456" || pin == "654321") return false
        return true
    }

    private fun isValidPhoneNumber(phone: String): Boolean {
        return phone.startsWith("+62") && phone.length >= 10
    }

    private fun isValidPassword(password: String): Boolean {
        if (password.length < 8) return false
        if (!password.any { it.isUpperCase() }) return false
        if (!password.any { it.isLowerCase() }) return false
        if (!password.any { it.isDigit() }) return false
        return true
    }

    private fun generateIdempotencyKey(): String {
        return java.util.UUID.randomUUID().toString()
    }

    private fun formatCompact(amount: Double): String {
        return "Rp %,d".format(amount.toInt())
    }
}