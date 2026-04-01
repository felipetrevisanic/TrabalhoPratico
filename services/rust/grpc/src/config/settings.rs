use std::env;

#[derive(Debug, Clone)]
pub struct AppConfig {
    db_host: String,
    db_port: String,
    db_name: String,
    db_user: String,
    db_password: String,
    db_ssl_mode: String,
    grpc_host: String,
    grpc_port: String,
}

impl AppConfig {
    pub fn from_env() -> Self {
        Self {
            db_host: read_env("DB_HOST", "localhost"),
            db_port: read_env("DB_PORT", "5432"),
            db_name: read_env("DB_NAME", "tcc_banco"),
            db_user: read_env("DB_USER", "postgres"),
            db_password: read_env("DB_PASSWORD", "postgres"),
            db_ssl_mode: read_env("DB_SSLMODE", "disable"),
            grpc_host: read_env("GRPC_HOST", "0.0.0.0"),
            grpc_port: read_env("GRPC_PORT", "50052"),
        }
    }

    pub fn database_url(&self) -> String {
        format!(
            "postgres://{}:{}@{}:{}/{}?sslmode={}",
            self.db_user,
            self.db_password,
            self.db_host,
            self.db_port,
            self.db_name,
            self.db_ssl_mode
        )
    }

    pub fn grpc_address(&self) -> String {
        format!("{}:{}", self.grpc_host, self.grpc_port)
    }
}

fn read_env(key: &str, fallback: &str) -> String {
    env::var(key).unwrap_or_else(|_| fallback.to_string())
}
