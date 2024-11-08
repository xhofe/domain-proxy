use std::env;

use axum::{
    body::Body,
    extract::State,
    http::{Request, Response, StatusCode},
    routing::any,
    Router,
};
use reqwest::Client;

#[tokio::main]
async fn main() -> Result<(), String> {
    let host = env::var("PROXY_HOST").map_err(|e| format!("Failed to get PROXY_HOST: {}", e))?;
    let scheme = env::var("PROXY_SCHEME").unwrap_or("https".to_string());
    let port = env::var("PROXY_PORT").unwrap_or("3000".to_string());

    let dst = format!("{}://{}", scheme, host);

    // Initialize reqwest client
    let client = Client::new();

    // Create the router
    let app = Router::new()
        .route("/", any(proxy_handler))
        .route("/*path", any(proxy_handler))
        .with_state((client, dst));

    let listener = tokio::net::TcpListener::bind(format!("0.0.0.0:{}", port))
        .await
        .unwrap();
    println!("listening on {}", listener.local_addr().unwrap());
    axum::serve(listener, app).await.unwrap();
    Ok(())
}

async fn proxy_handler(
    State((client, dst)): State<(Client, String)>,
    req: Request<Body>,
) -> Result<Response<Body>, (StatusCode, String)> {
    // Transform the incoming request to a reqwest request
    let path = req.uri().path();
    let path_query = req
        .uri()
        .path_and_query()
        .map(|v| v.as_str())
        .unwrap_or(path);
    let url = format!("{}{}", dst, path_query);

    let mut request_builder = client.request(req.method().clone(), &url);

    // Copy headers
    for (key, value) in req.headers().iter() {
        if key == axum::http::header::HOST {
            continue;
        }
        request_builder = request_builder.header(key, value);
    }

    // Send the request using reqwest
    match request_builder.send().await {
        Ok(response) => {
            let status = response.status();
            // println!("{} {} -> {}", req.method(), url, status);
            let headers = response.headers().clone();
            let body = response.bytes_stream();
            let mut response_builder = Response::builder().status(status);
            for (key, value) in headers {
                if let Some(key) = key {
                    response_builder = response_builder.header(key, value);
                }
            }
            response_builder
                .body(Body::from_stream(body))
                .map_err(|err| {
                    (
                        StatusCode::INTERNAL_SERVER_ERROR,
                        format!("Failed to build response: {}", err),
                    )
                })
        }
        Err(err) => Err((StatusCode::BAD_GATEWAY, format!("Request failed: {}", err))),
    }
}
