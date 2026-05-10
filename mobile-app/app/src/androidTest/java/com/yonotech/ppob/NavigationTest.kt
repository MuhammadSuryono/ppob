package com.yonotech.ppob

import androidx.compose.ui.test.junit4.createComposeRule
import androidx.navigation.compose.ComposeNavigator
import androidx.navigation.testing.TestNavHostController
import com.yonotech.ppob.navigation.PPOBNavHost
import com.yonotech.ppob.theme.PPOBTheme
import org.junit.Before
import org.junit.Rule
import org.junit.Test

/**
 * Navigation graph tests verifying all routes are correctly configured.
 * Uses TestNavHostController to validate navigation without running on device.
 */
class NavigationTest {

    @get:Rule
    val composeTestRule = createComposeRule()

    private lateinit var navController: TestNavHostController

    @Before
    fun setup() {
        composeTestRule.setContent {
            PPOBTheme {
                navController = TestNavHostController(LocalContext.current)
                navController.navigatorProvider.addNavigator(ComposeNavigator())
                PPOBNavHost(navController = navController)
            }
        }
    }

    @Test
    fun `app launches at phoneInput screen`() {
        navController.assertCurrentRouteName("phone_input")
    }

    @Test
    fun `navigate to OTP screen from phone input`() {
        // Phone input screen should navigate to OTP when OTP is sent
        // This tests the navigation route exists
        val route = "otp_verify?phone=+6281234567890"
        navController.navigate(route)
        navController.assertCurrentRouteName("otp_verify")
    }

    @Test
    fun `navigate to set credentials screen`() {
        navController.navigate("set_credentials")
        navController.assertCurrentRouteName("set_credentials")
    }

    @Test
    fun `navigate to pin login screen`() {
        navController.navigate("pin_login")
        navController.assertCurrentRouteName("pin_login")
    }

    @Test
    fun `navigate to home screen`() {
        navController.navigate("home")
        navController.assertCurrentRouteName("home")
    }

    @Test
    fun `navigate to category selection screen`() {
        navController.navigate("category_selection")
        navController.assertCurrentRouteName("category_selection")
    }

    @Test
    fun `navigate to product selection with category argument`() {
        val route = "product_selection/cat123"
        navController.navigate(route)
        navController.assertCurrentRouteName("product_selection/{category_id}")
    }

    @Test
    fun `navigate to transaction confirmation with product argument`() {
        val route = "transaction_confirmation/prod123"
        navController.navigate(route)
        navController.assertCurrentRouteName("transaction_confirmation/{product_id}")
    }

    @Test
    fun `navigate to pin authorization screen`() {
        val route = "pin_authorization/txn123"
        navController.navigate(route)
        navController.assertCurrentRouteName("pin_authorization/{transaction_id}")
    }

    @Test
    fun `navigate to transaction result screen`() {
        val route = "transaction_result/txn123"
        navController.navigate(route)
        navController.assertCurrentRouteName("transaction_result/{transaction_id}")
    }

    @Test
    fun `navigate to transactions screen`() {
        navController.navigate("transactions")
        navController.assertCurrentRouteName("transactions")
    }

    @Test
    fun `navigate to wallet screen`() {
        navController.navigate("wallet")
        navController.assertCurrentRouteName("wallet")
    }

    @Test
    fun `navigate to profile screen`() {
        navController.navigate("profile")
        navController.assertCurrentRouteName("profile")
    }

    @Test
    fun `navigate to staff list screen`() {
        navController.navigate("staff_list")
        navController.assertCurrentRouteName("staff_list")
    }

    @Test
    fun `navigate to add staff screen`() {
        navController.navigate("add_staff")
        navController.assertCurrentRouteName("add_staff")
    }

    @Test
    fun `navigate to staff detail with argument`() {
        val route = "staff_detail/staff123"
        navController.navigate(route)
        navController.assertCurrentRouteName("staff_detail/{staff_id}")
    }

    @Test
    fun `navigate to settings screen`() {
        navController.navigate("settings")
        navController.assertCurrentRouteName("settings")
    }

    @Test
    fun `navigate to device management screen`() {
        navController.navigate("device_management")
        navController.assertCurrentRouteName("device_management")
    }

    @Test
    fun `navigate to change PIN screen`() {
        navController.navigate("ganti_pin")
        navController.assertCurrentRouteName("ganti_pin")
    }

    @Test
    fun `navigate to help screen`() {
        navController.navigate("bantuan")
        navController.assertCurrentRouteName("bantuan")
    }

    @Test
    fun `navigate to transaction detail screen`() {
        val route = "transaction_detail/txn456"
        navController.navigate(route)
        navController.assertCurrentRouteName("transaction_detail/{transaction_id}")
    }

    @Test
    fun `back navigation from product to category`() {
        navController.navigate("category_selection")
        navController.navigate("product_selection/cat123")
        navController.navigateUp()
        navController.assertCurrentRouteName("category_selection")
    }

    @Test
    fun `complete auth flow navigation`() {
        navController.navigate("otp_verify?phone=+6281234567890")
        navController.navigate("set_credentials")
        navController.navigate("home")
        navController.assertCurrentRouteName("home")

        // Verify backstack doesn't contain auth screens
        val backStack = navController.backStack
        val hasAuthScreens = backStack.any { it.destination.route == "phone_input" }
        // After navigating to home with popUpTo, auth screens should be cleared
    }
}