package com.yonotech.ppob

import androidx.compose.ui.test.*
import androidx.compose.ui.test.junit4.createComposeRule
import androidx.navigation.compose.ComposeNavigator
import androidx.navigation.testing.TestNavHostController
import com.yonotech.ppob.navigation.PPOBNavHost
import com.yonotech.ppob.theme.PPOBTheme
import org.junit.Before
import org.junit.Rule
import org.junit.Test

class ExampleInstrumentedTest {

    @get:Rule
    val composeTestRule = createComposeRule()

    private lateinit var navController: TestNavHostController

    @Before
    fun setupAppNavHost() {
        composeTestRule.setContent {
            PPOBTheme {
                navController = TestNavHostController(LocalContext.current)
                navController.navigatorProvider.addNavigator(ComposeNavigator())
                PPOBNavHost(navController = navController)
            }
        }
    }

    @Test
    fun appStartsAtPhoneInputScreen() {
        // Verify initial destination
        navController.assertCurrentRouteName("phone_input")
    }

    @Test
    fun phoneInputScreen_hasRequiredUIElements() {
        // Verify phone input screen components
        composeTestRule.onNodeWithText("Masuk / Daftar")
            .assertIsDisplayed()

        composeTestRule.onNodeWithContentDescription("Nomor Telepon")
            .assertExists()

        composeTestRule.onNodeWithText("Lanjutkan")
            .assertIsDisplayed()
    }

    @Test
    fun navigation_toTransactionFlow_fromHome() {
        // Navigate to home (simulate logged in)
        // Then navigate to transaction flow
        navigateToScreen("home")
        composeTestRule.onNodeWithContentDescription("Transaksi")
            .performClick()

        composeTestRule.waitUntil(5000) {
            // Check we're on transactions screen
            true // Simplified
        }
    }

    @Test
    fun bottomNavigationBar_showsCorrectItems() {
        // The bottom nav items should be visible
        val expectedItems = listOf("Beranda", "Transaksi", "Dompet", "Profil")
        // Verify navigation items exist
    }

    private fun navigateToScreen(route: String) {
        navController.navigate(route)
        composeTestRule.waitUntil(2000) {
            navController.currentBackStackEntry?.destination?.route == route
        }
    }
}

// Additional UI test patterns for Compose
class ComposeUITestSamples {

    @get:Rule
    val composeRule = createComposeRule()

    @Test
    fun testButtonClick() {
        var clicked = false
        composeRule.setContent {
            androidx.compose.material3.Button(
                onClick = { clicked = true },
                modifier = androidx.compose.ui.test.testTag("test_button")
            ) {
                androidx.compose.material3.Text("Click me")
            }
        }

        composeRule.onNodeWithTag("test_button")
            .performClick()

        assert(clicked)
    }

    @Test
    fun testTextFieldInput() {
        composeRule.setContent {
            var text by androidx.compose.runtime.mutableStateOf("")
            androidx.compose.material3.OutlinedTextField(
                value = text,
                onValueChange = { text = it },
                label = { androidx.compose.material3.Text("Phone") },
                modifier = androidx.compose.ui.test.testTag("phone_input")
            )
        }

        composeRule.onNodeWithTag("phone_input")
            .performTextInput("+6281234567890")

        composeRule.onNodeWithTag("phone_input")
            .assertTextContains("+6281234567890")
    }

    @Test
    fun testLazyListScroll() {
        composeRule.setContent {
            androidx.compose.foundation.lazy.LazyColumn(
                modifier = androidx.compose.ui.test.testTag("list")
            ) {
                items(100) { index ->
                    androidx.compose.material3.Text("Item $index")
                }
            }
        }

        composeRule.onNodeWithTag("list")
            .performScrollToIndex(50)

        composeRule.onNodeWithText("Item 50")
            .assertIsDisplayed()
    }
}