package com.yonotech.ppob.mobile.ui.staff

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.hilt.navigation.compose.hiltViewModel
import com.yonotech.ppob.mobile.data.remote.dto.CreateStaffRequest
import com.yonotech.ppob.mobile.data.remote.dto.StaffDto
import com.yonotech.ppob.mobile.ui.components.PpoButton
import com.yonotech.ppob.mobile.ui.components.PpoTextField
import com.yonotech.ppob.mobile.utils.Resource
import com.yonotech.ppob.mobile.viewmodels.staff.StaffViewModel
import java.text.NumberFormat
import java.util.Locale

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun StaffListScreen(
    onBackClick: () -> Unit,
    onTopUpClick: (StaffDto) -> Unit,
    viewModel: StaffViewModel = hiltViewModel()
) {
    val staffState by viewModel.staffListState.collectAsState()
    val currencyFormat = NumberFormat.getCurrencyInstance(Locale("in", "ID")).apply { maximumFractionDigits = 0 }

    LaunchedEffect(Unit) {
        viewModel.getStaffList()
    }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Kelola Staff") },
                navigationIcon = {
                    IconButton(onClick = onBackClick) {
                        Icon(Icons.Default.ArrowBack, contentDescription = "Back")
                    }
                },
                colors = TopAppBarDefaults.topAppBarColors(
                    containerColor = MaterialTheme.colorScheme.primary,
                    titleContentColor = MaterialTheme.colorScheme.onPrimary,
                    navigationIconContentColor = MaterialTheme.colorScheme.onPrimary
                )
            )
        },
        floatingActionButton = {
            FloatingActionButton(
                onClick = { /* Open Add Staff Dialog */ },
                containerColor = MaterialTheme.colorScheme.primary
            ) {
                Icon(Icons.Default.Add, contentDescription = "Tambah Staff")
            }
        }
    ) { padding ->
        when (staffState) {
            is Resource.Loading -> {
                Box(modifier = Modifier.fillMaxSize().padding(padding), contentAlignment = Alignment.Center) {
                    CircularProgressIndicator()
                }
            }
            is Resource.Success -> {
                val staffList = (staffState as Resource.Success).data
                if (staffList.isEmpty()) {
                    Box(modifier = Modifier.fillMaxSize().padding(padding), contentAlignment = Alignment.Center) {
                        Text("Belum ada staff")
                    }
                } else {
                    LazyColumn(
                        modifier = Modifier.fillMaxSize().padding(padding),
                        contentPadding = PaddingValues(16.dp),
                        verticalArrangement = Arrangement.spacedBy(12.dp)
                    ) {
                        items(staffList) { staff ->
                            StaffItem(staff = staff, onTopUpClick = onTopUpClick, currencyFormat = currencyFormat)
                        }
                    }
                }
            }
            is Resource.Error -> {
                Box(modifier = Modifier.fillMaxSize().padding(padding), contentAlignment = Alignment.Center) {
                    Text(text = (staffState as Resource.Error).message, color = MaterialTheme.colorScheme.error)
                }
            }
            else -> {}
        }
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun StaffItem(
    staff: StaffDto,
    onTopUpClick: (StaffDto) -> Unit,
    currencyFormat: NumberFormat
) {
    Card(
        modifier = Modifier.fillMaxWidth(),
        shape = MaterialTheme.shapes.medium,
        colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.surface)
    ) {
        Row(
            modifier = Modifier.padding(16.dp),
            verticalAlignment = Alignment.CenterVertically
        ) {
            Surface(
                shape = MaterialTheme.shapes.small,
                color = MaterialTheme.colorScheme.secondaryContainer,
                modifier = Modifier.size(48.dp)
            ) {
                Box(contentAlignment = Alignment.Center) {
                    Text(text = staff.name.take(1), fontWeight = FontWeight.Bold)
                }
            }

            Spacer(modifier = Modifier.width(16.dp))

            Column(modifier = Modifier.weight(1f)) {
                Text(text = staff.name, fontWeight = FontWeight.Bold, fontSize = 16.sp)
                Text(text = staff.phone, fontSize = 12.sp, color = MaterialTheme.colorScheme.onSurfaceVariant)
                Text(
                    text = "Limit: ${currencyFormat.format(staff.dailyLimit)}",
                    fontSize = 12.sp,
                    color = MaterialTheme.colorScheme.primary
                )
            }

            Column(horizontalAlignment = Alignment.End) {
                Text(
                    text = currencyFormat.format(staff.balance),
                    fontWeight = FontWeight.Bold,
                    color = MaterialTheme.colorScheme.primary
                )
                TextButton(onClick = { onTopUpClick(staff) }) {
                    Text("Top Up")
                }
            }
        }
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun AddStaffDialog(
    onDismiss: () -> Unit,
    viewModel: StaffViewModel = hiltViewModel()
) {
    var name by remember { mutableStateOf("") }
    var phone by remember { mutableStateOf("") }
    var email by remember { mutableStateOf("") }
    var marginValue by remember { mutableStateOf("") }
    var dailyLimit by remember { mutableStateOf("") }
    var marginScheme by remember { mutableStateOf("FixedAllowance") }
    val createState by viewModel.createStaffState.collectAsState()

    var expanded by remember { mutableStateOf(false) }
    val schemes = listOf("FixedAllowance", "MarginShare")

    LaunchedEffect(createState) {
        if (createState is Resource.Success) {
            onDismiss()
            viewModel.resetState()
        }
    }

    AlertDialog(
        onDismissRequest = onDismiss,
        title = { Text("Tambah Staff Baru") },
        text = {
            Column {
                PpoTextField(value = name, onValueChange = { name = it }, label = "Nama Lengkap")
                Spacer(modifier = Modifier.height(8.dp))
                PpoTextField(value = phone, onValueChange = { phone = it }, label = "Nomor Telepon")
                Spacer(modifier = Modifier.height(8.dp))
                PpoTextField(value = email, onValueChange = { email = it }, label = "Email")
                Spacer(modifier = Modifier.height(8.dp))

                ExposedDropdownMenuBox(expanded = expanded, onExpandedChange = { expanded = it }) {
                    OutlinedTextField(
                        value = marginScheme,
                        onValueChange = {},
                        readOnly = true,
                        label = { Text("Skema Margin") },
                        trailingIcon = { ExposedDropdownMenuDefaults.TrailingIcon(expanded) },
                        modifier = Modifier.menuAnchor()
                    )
                    ExposedDropdownMenu(expanded = expanded, onDismissRequest = { expanded = false }) {
                        schemes.forEach { scheme ->
                            DropdownMenuItem(text = { Text(scheme) }, onClick = {
                                marginScheme = scheme
                                expanded = false
                            })
                        }
                    }
                }

                Spacer(modifier = Modifier.height(8.dp))
                PpoTextField(value = marginValue, onValueChange = { marginValue = it }, label = "Nilai Margin")
                Spacer(modifier = Modifier.height(8.dp))
                PpoTextField(value = dailyLimit, onValueChange = { dailyLimit = it }, label = "Limit Harian (Rp)")
            }
        },
        confirmButton = {
            PpoButton(
                label = "Simpan",
                onClick = {
                    viewModel.createStaff(
                        CreateStaffRequest(
                            name = name,
                            phone = phone,
                            email = email,
                            marginScheme = marginScheme,
                            marginValue = marginValue.toDoubleOrNull() ?: 0.0,
                            dailyLimit = dailyLimit.toDoubleOrNull() ?: 0.0
                        )
                    )
                },
                isLoading = createState is Resource.Loading
            )
        },
        dismissButton = {
            TextButton(onClick = onDismiss) { Text("Batal") }
        }
    )
}