use std::env;

#[derive(Debug, Clone)]
pub struct AppConfig {
    db_host: String,
    db_port: String,
    db_name: String,
    db_user: String,
    db_password: String,
    db_ssl_mode: String,
    http_host: String,
    http_port: String,
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
            http_host: read_env("HTTP_HOST", "0.0.0.0"),
            http_port: read_env("HTTP_PORT", "8091"),
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

    pub fn http_address(&self) -> String {
        format!("{}:{}", self.http_host, self.http_port)
    }
}

fn read_env(key: &str, fallback: &str) -> String {
    env::var(key).unwrap_or_else(|_| fallback.to_string())
}
