use std::{env, ffi::OsString, process::Command};

use axum::{Json, extract::State, response::IntoResponse};
use problem_details::ProblemDetails;
use reqwest::StatusCode;

use crate::app::{api::AppState, runner::AppRunner};

/// Restart Application
///
/// Shuts the current SISR instance down and relaunches it with the same arguments.
#[utoipa::path(
    post,
    path = "/api/v1/restart_sisr",
    tag = "ui",
    responses(
        (status = 200),
        (status = 424, description = "Failed to schedule restart"),
        (status = 501, description = "Restart is not implemented on this platform"),
    )
)]
pub async fn restart_sisr(State(_state): State<AppState>) -> impl IntoResponse {
    tracing::debug!("Received request to restart SISR");

    let current_exe = match env::current_exe() {
        Ok(path) => path,
        Err(e) => {
            tracing::error!("Failed to resolve current executable: {}", e);
            return ProblemDetails::from_status_code(StatusCode::FAILED_DEPENDENCY)
                .with_detail(format!("Failed to resolve current executable: {}", e))
                .into_response();
        }
    };
    let current_args: Vec<OsString> = env::args_os().skip(1).collect();

    AppRunner::run_cleanup_handler();

    #[cfg(windows)]
    {
        if launched_from_terminal() {
            let ps_quote = |value: &str| format!("'{}'", value.replace('\'', "''"));

            let exe = ps_quote(&current_exe.as_os_str().to_string_lossy());
            let arg_list = if current_args.is_empty() {
                String::from("@()")
            } else {
                format!(
                    "@({})",
                    current_args
                        .iter()
                        .map(|arg| ps_quote(&arg.to_string_lossy()))
                        .collect::<Vec<_>>()
                        .join(", ")
                )
            };
            let current_pid = std::process::id();
            let command = format!(
                "& {{ while (Get-Process -Id {current_pid} -ErrorAction SilentlyContinue) {{ Start-Sleep -Milliseconds 200 }}; Start-Process -FilePath {exe} -ArgumentList {arg_list} }}"
            );

            if let Err(e) = Command::new("powershell")
                .args([
                    "-NoProfile",
                    "-NonInteractive",
                    "-WindowStyle",
                    "Hidden",
                    "-Command",
                    &command,
                ])
                .spawn()
            {
                tracing::error!("Failed to schedule SISR restart: {}", e);
                return ProblemDetails::from_status_code(StatusCode::FAILED_DEPENDENCY)
                    .with_detail(format!("Failed to schedule SISR restart: {}", e))
                    .into_response();
            }
        } else if let Err(e) = Command::new("explorer.exe")
            .arg(&current_exe)
            .args(&current_args)
            .spawn()
        {
            tracing::error!("Failed to schedule SISR restart: {}", e);
            return ProblemDetails::from_status_code(StatusCode::FAILED_DEPENDENCY)
                .with_detail(format!("Failed to schedule SISR restart: {}", e))
                .into_response();
        }

        tokio::spawn(async move {
            tokio::time::sleep(std::time::Duration::from_millis(150)).await;
            std::process::exit(0);
        });
    }

    #[cfg(target_os = "linux")]
    {
        let current_pid = std::process::id();
        let pid = current_pid.to_string();

        if let Err(e) = Command::new("sh")
            .arg("-lc")
            .arg("while kill -0 \"$1\" 2>/dev/null; do sleep 0.2; done; shift; exec \"$@\"")
            .arg("restart_sisr")
            .arg(&pid)
            .arg(&current_exe)
            .args(&current_args)
            .spawn()
        {
            tracing::error!("Failed to schedule SISR restart: {}", e);
            return ProblemDetails::from_status_code(StatusCode::FAILED_DEPENDENCY)
                .with_detail(format!("Failed to schedule SISR restart: {}", e))
                .into_response();
        }

        tokio::spawn(async move {
            tokio::time::sleep(std::time::Duration::from_millis(150)).await;
            AppRunner::shutdown_without_cleanup();
        });
    }

    (StatusCode::OK, Json(serde_json::json!({}))).into_response()
}

#[cfg(windows)]
fn launched_from_terminal() -> bool {
    use windows_sys::Win32::System::Diagnostics::ToolHelp::{
        CreateToolhelp32Snapshot, Process32FirstW, Process32NextW, PROCESSENTRY32W,
        TH32CS_SNAPPROCESS,
    };

    unsafe {
        let snapshot = CreateToolhelp32Snapshot(TH32CS_SNAPPROCESS, 0);
        if snapshot == windows_sys::Win32::Foundation::INVALID_HANDLE_VALUE {
            return false;
        }

        let current_pid = std::process::id();
        let mut entry: PROCESSENTRY32W = std::mem::zeroed();
        entry.dwSize = std::mem::size_of::<PROCESSENTRY32W>() as u32;

        let mut parent_pid = 0;
        if Process32FirstW(snapshot, &mut entry) != 0 {
            loop {
                if entry.th32ProcessID == current_pid {
                    parent_pid = entry.th32ParentProcessID;
                    break;
                }

                if Process32NextW(snapshot, &mut entry) == 0 {
                    break;
                }
            }
        }

        windows_sys::Win32::Foundation::CloseHandle(snapshot);

        if parent_pid == 0 {
            return false;
        }

        let snapshot = CreateToolhelp32Snapshot(TH32CS_SNAPPROCESS, 0);
        if snapshot == windows_sys::Win32::Foundation::INVALID_HANDLE_VALUE {
            return false;
        }

        let mut entry: PROCESSENTRY32W = std::mem::zeroed();
        entry.dwSize = std::mem::size_of::<PROCESSENTRY32W>() as u32;

        let mut terminal_parent = false;
        if Process32FirstW(snapshot, &mut entry) != 0 {
            loop {
                if entry.th32ProcessID == parent_pid {
                    let len = entry
                        .szExeFile
                        .iter()
                        .position(|&ch| ch == 0)
                        .unwrap_or(entry.szExeFile.len());
                    let exe_name = String::from_utf16_lossy(&entry.szExeFile[..len]);
                    terminal_parent = matches!(
                        exe_name.to_ascii_lowercase().as_str(),
                        "cmd.exe"
                            | "powershell.exe"
                            | "pwsh.exe"
                            | "windowsterminal.exe"
                            | "conhost.exe"
                    );
                    break;
                }

                if Process32NextW(snapshot, &mut entry) == 0 {
                    break;
                }
            }
        }

        windows_sys::Win32::Foundation::CloseHandle(snapshot);
        terminal_parent
    }
}
