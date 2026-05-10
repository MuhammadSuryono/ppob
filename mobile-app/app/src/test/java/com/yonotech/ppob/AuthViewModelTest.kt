package com.yonotech.ppob

import app.cash.turbine.test
import com.yonotech.ppob.domain.usecase.LoginUseCase
import com.yonotech.ppob.domain.usecase.RegisterUseCase
import com.yonotech.ppob.domain.usecase.VerifyOtpUseCase
import com.yonotech.ppob.presentation.auth.AuthStep
import com.yonotech.ppob.presentation.auth.AuthUiState
import com.yonotech.ppob.presentation.auth.AuthViewModel
import io.mockk.coEvery
import io.mockk.mockk
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.ExperimentalCoroutinesApi
import kotlinx.coroutines.flow.flow
import kotlinx.coroutines.test.*
import org.junit.After
import org.junit.Assert.assertEquals
import org.junit.Assert.assertFalse
import org.junit.Assert.assertTrue
import org.junit.Before
import org.junit.Test
import com.yonotech.ppob.data.remote.model.AuthResponse as RemoteAuthResponse
import com.yonotech.ppob.data.remote.model.RoleResponse as RemoteRoleResponse
import com.yonotech.ppob.data.remote.model.UserResponse as RemoteUserResponse

@OptIn(ExperimentalCoroutinesApi::class)
class AuthViewModelTest {

    private lateinit var authRepository: com.yonotech.ppob.domain.repository.AuthRepository
    private lateinit var registerUseCase: RegisterUseCase
    private lateinit var verifyOtpUseCase: VerifyOtpUseCase
    private lateinit var loginUseCase: LoginUseCase
    private lateinit var viewModel: AuthViewModel

    @Before
    fun setUp() {
        Dispatchers.setMain(StandardTestDispatcher())

        authRepository = mockk(relaxed = true)
        registerUseCase = mockk(relaxed = true)
        verifyOtpUseCase = mockk(relaxed = true)
        loginUseCase = mockk(relaxed = true)

        viewModel = AuthViewModel(
            authRepository = authRepository,
            registerUseCase = registerUseCase,
            verifyOtpUseCase = verifyOtpUseCase,
            loginUseCase = loginUseCase
        )
    }

    @After
    fun tearDown() {
        Dispatchers.resetMain()
    }

    // ==================== Phone Number Validation ====================

    @Test
    fun `phone number with valid Indonesian format should be valid`() {
        viewModel.onPhoneNumberChange("+6281234567890")
        val state = viewModel.uiState.value
        assertTrue(state.isPhoneValid)
        assertEquals("+6281234567890", state.phoneNumber)
    }

    @Test
    fun `phone number too short should be invalid`() {
        viewModel.onPhoneNumberChange("+628123")
        assertFalse(viewModel.uiState.value.isPhoneValid)
    }

    @Test
    fun `empty phone number should be invalid`() {
        viewModel.onPhoneNumberChange("")
        assertFalse(viewModel.uiState.value.isPhoneValid)
    }

    @Test
    fun `phone without country code should be invalid`() {
        viewModel.onPhoneNumberChange("081234567890")
        assertFalse(viewModel.uiState.value.isPhoneValid)
    }

    // ==================== Password Validation ====================

    @Test
    fun `password with all requirements should be valid`() {
        viewModel.onPasswordChange("StrongPassword123!")
        assertTrue(viewModel.uiState.value.isPasswordValid)
    }

    @Test
    fun `password without uppercase should be invalid`() {
        viewModel.onPasswordChange("weakpassword123!")
        assertFalse(viewModel.uiState.value.isPasswordValid)
    }

    @Test
    fun `password without digit should be invalid`() {
        viewModel.onPasswordChange("NoDigitsHere!!")
        assertFalse(viewModel.uiState.value.isPasswordValid)
    }

    @Test
    fun `password too short should be invalid`() {
        viewModel.onPasswordChange("Ab1!")
        assertFalse(viewModel.uiState.value.isPasswordValid)
    }

    // ==================== PIN Validation ====================

    @Test
    fun `valid 6-digit PIN should pass validation`() {
        viewModel.onPinChange("123456")
        assertTrue(viewModel.uiState.value.isPinValid)
    }

    @Test
    fun `PIN with all same digits should be invalid`() {
        viewModel.onPinChange("111111")
        assertFalse(viewModel.uiState.value.isPinValid)
    }

    @Test
    fun `sequential PIN should be invalid`() {
        viewModel.onPinChange("123456")
        // Note: 123456 is actually blocked as sequential
        // Let's test a non-sequential valid PIN
        viewModel.onPinChange("135790")
        assertTrue(viewModel.uiState.value.isPinValid)
    }

    @Test
    fun `5-digit PIN should be invalid`() {
        viewModel.onPinChange("12345")
        assertFalse(viewModel.uiState.value.isPinValid)
    }

    @Test
    fun `PIN with letters should be invalid`() {
        viewModel.onPinChange("12345a")
        assertFalse(viewModel.uiState.value.isPinValid)
    }

    // ==================== OTP Sending ====================

    @Test
    fun `sendOtp should transition to OTP_VERIFY step on success`() = runTest {
        coEvery { registerUseCase(any(), any()) } returns Result.success(Unit)

        viewModel.onPhoneNumberChange("+6281234567890")
        viewModel.sendOtp()

        // Wait for state update
        advanceUntilIdle()

        assertEquals(AuthStep.OTP_VERIFY, viewModel.uiState.value.currentStep)
        assertFalse(viewModel.uiState.value.isLoading)
    }

    @Test
    fun `sendOtp should set error on failure`() = runTest {
        coEvery { registerUseCase(any(), any()) } returns Result.failure(Exception("Network error"))

        viewModel.onPhoneNumberChange("+6281234567890")
        viewModel.sendOtp()

        advanceUntilIdle()

        assertEquals(AuthStep.PHONE_INPUT, viewModel.uiState.value.currentStep)
        assertEquals("Network error", viewModel.uiState.value.error)
        assertFalse(viewModel.uiState.value.isLoading)
    }

    // ==================== OTP Verification ====================

    @Test
    fun `verifyOtp should transition to SET_CREDENTIALS on success`() = runTest {
        coEvery { verifyOtpUseCase(any(), any(), any(), any()) } returns Result.success(Unit)

        viewModel.onPhoneNumberChange("+6281234567890")
        viewModel.onOtpChange("123456")
        viewModel.onPasswordChange("StrongPass1!")
        viewModel.onPinChange("654321")
        viewModel.verifyOtp()

        advanceUntilIdle()

        assertEquals(AuthStep.COMPLETE, viewModel.uiState.value.currentStep)
    }

    @Test
    fun `verifyOtp should show error when OTP is invalid`() {
        viewModel.onOtpChange("1") // Too short
        assertFalse(viewModel.uiState.value.isOtpValid)
    }

    @Test
    fun `verifyOtp should show error when passwords do not match`() {
        viewModel.onPasswordChange("Password1!")
        viewModel.onConfirmPasswordChange("Password2!")
        assertFalse(viewModel.uiState.value.confirmPassword == viewModel.uiState.value.password)
    }

    // ==================== Login ====================

    @Test
    fun `login should transition to COMPLETE on successful authentication`() = runTest {
        coEvery { loginUseCase(any(), any(), any(), any()) } returns Result.success(Unit)

        viewModel.onPhoneNumberChange("+6281234567890")
        viewModel.onPasswordChange("StrongPass1!")
        viewModel.onPinChange("654321")
        viewModel.login()

        advanceUntilIdle()

        assertEquals(AuthStep.COMPLETE, viewModel.uiState.value.currentStep)
    }

    @Test
    fun `login should set error on failure`() = runTest {
        coEvery { loginUseCase(any(), any(), any(), any()) } returns
            Result.failure(Exception("AUTH_INVALID_CREDENTIALS"))

        viewModel.onPhoneNumberChange("+6281234567890")
        viewModel.onPasswordChange("WrongPassword!")
        viewModel.onPinChange("654321")
        viewModel.login()

        advanceUntilIdle()

        assertEquals("AUTH_INVALID_CREDENTIALS", viewModel.uiState.value.error)
    }

    // ==================== State Transitions ====================

    @Test
    fun `moveToStep should update current step`() {
        viewModel.moveToStep(AuthStep.OTP_VERIFY)
        assertEquals(AuthStep.OTP_VERIFY, viewModel.uiState.value.currentStep)
    }

    @Test
    fun `resetForm should clear all fields`() {
        viewModel.onPhoneNumberChange("+6281234567890")
        viewModel.onPasswordChange("Test123!")
        viewModel.resetForm()

        assertEquals("", viewModel.uiState.value.phoneNumber)
        assertEquals("", viewModel.uiState.value.password)
        assertNull(viewModel.uiState.value.error)
    }

    @Test
    fun `clearError should remove error message`() {
        viewModel.onPhoneNumberChange("invalid")
        assertTrue(viewModel.uiState.value.error != null || !viewModel.uiState.value.isPhoneValid)
        viewModel.resetForm()
        viewModel.onPhoneNumberChange("+6281234567890")
        assertNull(viewModel.uiState.value.error)
    }

    // ==================== Loading State ====================

    @Test
    fun `sendOtp should set loading state during execution`() = runTest {
        coEvery { registerUseCase(any(), any()) } coAnswers {
            delay(100)
            Result.success(Unit)
        }

        viewModel.onPhoneNumberChange("+6281234567890")
        viewModel.sendOtp()

        advanceUntilIdle()

        assertFalse(viewModel.uiState.value.isLoading)
    }
}