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
        if let Err(e) = Command::new("explorer.exe")
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
